package enumify

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	g "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

type Enum interface {
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
	sql.Scanner
	driver.Valuer
}

// Options defines how the code generation should be performed. These values are set by
// the CLI flags passed via the go generate directive.
type Options struct {
	File           string // The file that kicked off the generation (from $GOFILE)
	Pkg            string // The package that the enum is in (from $GOPACKAGE)
	NameVar        string // The variable name that contains the names of the enum values
	CaseSensitive  bool   // Whether to make the enum case sensitive (defaults to false)
	SpaceSensitive bool   // Whether to make the enum space sensitive (defaults to false)
	NoTests        bool   // Whether to skip the testing code generation (defaults to false)
	NoParser       bool   // Whether to skip the parser code generation (defaults to false)
	NoStringer     bool   // Whether to skip the Stringer interface code generation (defaults to false)
	NoText         bool   // Whether to skip the text interfaces code generation (defaults to false)
	NoBinary       bool   // Whether to skip the binary interfaces code generation (defaults to false)
	NoJSON         bool   // Whether to skip the JSON interfaces code generation (defaults to false)
	NoYAML         bool   // Whether to skip the YAML interfaces code generation (defaults to false)
	NoSQL          bool   // Whether to skip the SQL interfaces code generation (defaults to false)
}

// Discover is the entry point for all code generation. It will parse the file specified
// by the go generate directive into an AST, then loop through all of the data in the
// file to discover the Enum types (anything that extends uint8), the constant values
// for the enum types, and the names variable based on the name pattern, the passed in
// name, and the allowed name var types.
func Discover(o Options) (etypes EnumTypes, err error) {
	// Build tool package discovery configuration.
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}

	// NOTE: do not use the driver query file={os.File} here because it will load the
	// entire package instead of just the contents of the file. As a result, the
	// types discovered by packages.Load will have the package "command-line-arguments"
	// scope rather than the scope of the package being inspected.
	//
	// We prefer this so we can isolate the specific files that have a go generate
	// directive and ignore the other files including other files that may also have
	// go generate directives.
	//
	// TODO: what if multiple enums are defined in the same file?
	var pkgs []*packages.Package
	if pkgs, err = packages.Load(cfg, o.File); err != nil {
		return nil, fmt.Errorf("failed to load package %q for inspection: %w", o.File, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for inspection")
	}

	if len(pkgs) > 1 {
		return nil, fmt.Errorf("multiple packages found for inspection: %v", pkgs)
	}

	gopkg := pkgs[0]
	if len(gopkg.Errors) > 0 {
		return nil, fmt.Errorf("package errors: %v", pkgs[0].Errors)
	}

	// Get the predeclared uint8 type for comparison
	uint8Type := types.Typ[types.Uint8]

	// First pass: discover all the enum types in the package.
	// An enum type is a type whose underlying type is uint8.
	etypes = make(EnumTypes, 0, 1)
	scope := gopkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if typ, ok := obj.(*types.TypeName); ok {
			if types.Identical(typ.Type().Underlying(), uint8Type) {
				etypes = append(etypes, &EnumType{
					Name:  name,
					Type:  obj.Type(),
					gopkg: gopkg,
					scope: scope,
				})
			}
		}
	}

	// Second pass: populate the enum types with consts and the name variable.
	for _, etype := range etypes {
		// Discover the consts and names variable for the enum type.
		etype.discover()

		// If the names variable is not set, but one was passed in from the command
		// line, then attempt to set it on the enum type.
		if etype.NamesVar == nil && o.NameVar != "" {
			if err = etype.setNamesVar(o.NameVar); err != nil {
				return nil, fmt.Errorf("failed to set names variable for enum type %q: %w", etype.Name, err)
			}
		}

		// The enum type must be valid before we can generate code for it.
		if err = etype.validate(); err != nil {
			return nil, err
		}
	}
	return etypes, nil
}

//============================================================================
// Code Generation Entry Points
//============================================================================

// The generation comment is considered best practice when writing golang code generators.
// In addition to the standard warning, it also includes the version of the generator that
// was used to generate the code for debugging purposes.
func GenerateComment() string {
	return fmt.Sprintf("Code generated by enumify v%s. DO NOT EDIT.", Version())
}

// Generate the code file for the enum types.
func Generate(opts Options) (err error) {
	f := g.NewFile(opts.Pkg)
	f.PackageComment(GenerateComment())

	// Manage import names
	f.ImportName("go.rtnl.ai/enumify", "enumify")
	f.ImportName("encoding/json", "json")
	f.ImportName("database/sql/driver", "driver")

	var etypes EnumTypes
	if etypes, err = Discover(opts); err != nil {
		return err
	}

	// Write the code for the enum types.
	for _, etype := range etypes {
		evar := g.Id("err")                       // err variable name
		s := etype.ReceiverId()                   // e if Enum or s if Status etc.
		ptrS := g.Op("*").Add(etype.ReceiverId()) // *e if Enum or *s if Status etc.
		typeID := etype.Id()                      // Enum
		namesVar := etype.NamesVarId()            // enumNames
		namesVarSlice := etype.NamesVarSliceId()  // enumNames or enumNames[0]

		f.Comment("//============================================================================")
		f.Comment("// " + etype.Name + " Enum Type: Generated Functions and Methods")
		f.Comment("//============================================================================")
		f.Line()

		// Write the parser code for the enum type
		if !opts.NoParser {
			// TODO: this uses the factory function to create the parser function. This
			// is not ideal because it creates a dependency on the go.rtnl.ai/enumify
			// package. See https://github.com/rotationalio/enumify/issues/7 for
			// ideas on how to improve this.
			parserVar := g.Id(LowerFirst(etype.Name) + "Parser")
			parserType := g.Func().Params(g.Any()).Params(typeID, g.Error())
			f.Var().Add(parserVar, parserType).Op("=").Qual("go.rtnl.ai/enumify", "ParseFactory").Types(typeID).Call(namesVar)
			f.Line()

			f.Commentf("Parse%s parses the given value into a %s.", etype.Name, etype.Name)
			parserFunc := f.Func().Id("Parse" + UpperFirst(etype.Name))
			parserFunc.Params(s.Clone().Any()).Params(typeID, g.Error()).Block(
				g.Return(parserVar.Clone().Call(s.Clone())),
			)
			f.Line()
		}

		methodSig := g.Func().Params(s.Clone().Add(typeID))
		methodPtrSig := g.Func().Params(s.Clone().Op("*").Add(typeID))

		// Write the stringer code for the enum type
		if !opts.NoStringer {
			f.Commentf("Ensure %s implements fmt.Stringer.", etype.Name)

			method := methodSig.Clone().Id("String").Call().String()
			method.Block(
				g.If(s.Clone().Op(">=").Add(typeID).Call(g.Len(namesVarSlice))).Block(
					g.Return(etype.IndexNames(etype.ZeroConstId())),
				),
				g.Return(etype.IndexNames(s)),
			)
			f.Add(method)
			f.Line()
		}

		// JSON Marshal and Unmarshal code
		if !opts.NoJSON {
			f.Commentf("Ensure %s implements json.Marshaler.", etype.Name)
			method := methodSig.Clone().Id("MarshalJSON").Call().Params(g.Id("[]byte"), g.Error())
			method.Block(
				g.Return(g.Qual("encoding/json", "Marshal").Call(s.Clone().Dot("String").Call())),
			)
			f.Add(method)
			f.Line()

			f.Commentf("Ensure %s implements json.Unmarshaler.", etype.Name)

			sv := g.Id("sv")
			method = methodPtrSig.Clone().Id("UnmarshalJSON").Call(g.Id("data").Id("[]byte")).Params(evar.Clone().Error())
			method.Block(
				g.Var().Add(sv).Any(),
				g.If(
					evar.Clone().Op("=").Qual("encoding/json", "Unmarshal").Call(g.Id("data"), g.Op("&").Add(sv)).Op(";").Id("err").Op("!=").Nil()).
					Block(
						g.Return(evar),
					),
				g.Line(),
				g.If(
					ptrS.Clone().Op(",").Id("err").Op("=").Id("Parse"+etype.Name).Call(sv).Op(";").Id("err").Op("!=").Nil().
						Block(
							g.Return(evar),
						),
				),
				g.Return().Nil(),
			)
			f.Add(method)
			f.Line()
		}

		// YAML Marshal and Unmarshal code
		if !opts.NoYAML {
			f.Commentf("Ensure %s implements yaml.Marshaler.", etype.Name)
			method := methodSig.Clone().Id("MarshalYAML").Call().Params(g.Any(), g.Error())
			method.Block(
				g.Return(s.Clone().Dot("String").Call(), g.Nil()),
			)
			f.Add(method)
			f.Line()

			f.Commentf("Ensure %s implements yaml.Unmarshaler.", etype.Name)

			sv := g.Id("sv")
			method = methodPtrSig.Clone().Id("UnmarshalYAML").Call(g.Id("unmarshal").Func().Params(g.Any()).Params(g.Error())).Params(evar.Clone().Error())
			method.Block(
				g.Var().Add(sv).String(),
				g.If(evar.Clone().Op("=").Id("unmarshal").Call(g.Op("&").Add(sv)).Op(";").Id("err").Op("!=").Nil()).Block(
					g.Return(evar),
				),
				g.Line(),
				g.If(g.Add(ptrS.Clone()).Op(",").Id("err").Op("=").Id("Parse"+etype.Name).Call(sv).Op(";").Id("err").Op("!=").Nil().
					Block(
						g.Return(evar),
					),
				),
				g.Return(g.Nil()),
			)
			f.Add(method)
			f.Line()
		}

		// SQL Scanner and Valuer code
		if !opts.NoSQL {
			f.Commentf("Ensure %s implements sql.Scanner.", etype.Name)
			val := g.Id("val")
			method := methodPtrSig.Clone().Id("Scan").Call(g.Id("src").Any()).Params(evar.Clone().Error())
			method.Block(
				g.Switch(val.Clone().Op(":=").Id("src").Assert(g.Type())).Block(
					g.Case(g.Nil()).Block(g.Return(g.Nil())),
					g.Case(g.String()).Block(
						g.List(ptrS.Clone(), evar.Clone().Op("=").Id("Parse"+etype.Name).Call(val)),
						g.Return(evar),
					),
					g.Case(g.Id("[]byte")).Block(
						g.List(ptrS.Clone(), evar.Clone().Op("=").Id("Parse"+etype.Name).Call(g.String().Call(val))),
						g.Return(evar),
					),
					g.Default().Block(
						g.Return(g.Qual("fmt", "Errorf").Call(g.Lit("cannot scan %T into "+etype.Name), val)),
					),
				),
			)

			f.Add(method)
			f.Line()

			f.Commentf("Ensure %s implements driver.Valuer.", etype.Name)
			method = methodSig.Clone().Id("Value").Call().Params(g.Qual("database/sql/driver", "Value"), g.Error())
			method.Block(
				g.Return(s.Clone().Dot("String").Call(), g.Nil()),
			)
			f.Add(method)
			f.Line()
		}

	}

	if err = f.Save(opts.fileName()); err != nil {
		return err
	}
	return nil
}

// Generate the test file for the enum types.
func GenerateTests(opts Options) (err error) {
	f := g.NewFile(opts.Pkg)
	f.PackageComment(GenerateComment())

	// Manage import names
	f.ImportName("testing", "testing")
	f.ImportName("go.rtnl.ai/enumify", "enumify")

	// Discover the enum types in the package.
	var etypes EnumTypes
	if etypes, err = Discover(opts); err != nil {
		return err
	}

	// For each enum type, write a test function that creates and executes an enumify
	// test suite for the enum type.
	for _, etype := range etypes {
		f.Func().Id("Test"+UpperFirst(etype.Name)).Add(TestingT).Block(
			g.Id("suite").Op(":=").Qual("go.rtnl.ai/enumify", "TestSuite").Types(etype.Id(), etype.NamesVarTypeId()).Block(
				g.Id("Values").Op(":").Add(etype.ConstLiteral()).Op(","),
				g.Id("Names").Op(":").Add(etype.NamesVarId()).Op(","),
				g.Id("ICase").Op(":").Lit(!opts.CaseSensitive).Op(","),
				g.Id("ISpace").Op(":").Lit(!opts.SpaceSensitive).Op(","),
			),
			g.Line(),
			g.Id("suite").Dot("Run").Call(g.Id("t")),
		).Line()
	}

	if err = f.Save(opts.testFileName()); err != nil {
		return err
	}
	return nil
}

//============================================================================
// Options Methods
//============================================================================

func (o Options) fileName() string {
	ext := filepath.Ext(o.File)
	return strings.TrimSuffix(filepath.Base(o.File), ext) + "_gen" + ext
}

func (o Options) testFileName() string {
	ext := filepath.Ext(o.File)
	return strings.TrimSuffix(filepath.Base(o.File), ext) + "_gen_test" + ext
}
