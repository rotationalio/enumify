package enumify

import (
	"fmt"
	"go/constant"
	"go/types"
	"slices"
	"unicode"
	"unicode/utf8"

	g "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

type EnumTypes []*EnumType

type EnumType struct {
	Name     string
	Type     types.Type
	Consts   []types.Object
	NamesVar types.Object
	gopkg    *packages.Package
	scope    *types.Scope
}

//============================================================================
// EnumType Code Generation Utilities
//============================================================================

func (e *EnumType) Id() *g.Statement {
	return g.Id(e.Name)
}

func (e *EnumType) ReceiverId() *g.Statement {
	if e.Name == "" {
		return g.Id("e")
	}

	r, _ := utf8.DecodeRuneInString(e.Name)
	if r == utf8.RuneError {
		return g.Id("e")
	}

	return g.Id(string(unicode.ToLower(r)))
}

func (e *EnumType) NamesVarId() *g.Statement {
	return g.Id(e.NamesVar.Name())
}

func (e *EnumType) NamesVarTypeId() *g.Statement {
	switch {
	case isStringTable(e.NamesVar.Type()):
		return g.Id("[][]string")
	case isStringSlice(e.NamesVar.Type()):
		return g.Id("[]string")
	default:
		panic(fmt.Errorf("unsupported names variable type: %T", e.NamesVar.Type()))
	}
}

func (e *EnumType) NamesVarSliceId() *g.Statement {
	switch {
	case isStringTable(e.NamesVar.Type()):
		return e.NamesVarId().Index(g.Lit(0))
	case isStringSlice(e.NamesVar.Type()):
		return e.NamesVarId()
	default:
		panic(fmt.Errorf("unsupported names variable type: %T", e.NamesVar.Type()))
	}
}

func (e *EnumType) IndexNames(i g.Code) *g.Statement {
	switch {
	case isStringTable(e.NamesVar.Type()):
		return e.NamesVarId().Index(g.Lit(0)).Index(i)
	case isStringSlice(e.NamesVar.Type()):
		return e.NamesVarId().Index(i)
	default:
		panic(fmt.Errorf("unsupported names variable type: %T", e.NamesVar.Type()))

	}
}

func (e *EnumType) ConstLiteral() *g.Statement {
	constLits := make([]g.Code, 0, len(e.Consts))
	for _, constObj := range e.Consts {
		constLit := g.Id(constObj.Name())
		constLits = append(constLits, constLit)
	}
	return g.Id("[]" + e.Name).Add(g.Values(constLits...))
}

func (e *EnumType) ZeroConstId() *g.Statement {
	if zeroConst := e.zeroConst(); zeroConst != nil {
		return g.Id(zeroConst.Name())
	}
	return g.Lit(0)
}

//============================================================================
// EnumType AST Validation and Parsing
//============================================================================

func (e *EnumType) validate() error {
	// Not sure how this would happen ... but we're going to LBYL because types are hard.
	if e.Name == "" || e.Type == nil {
		return fmt.Errorf("enum type %q has no name or type", e.Name)
	}

	// Need to have at least one constant declared for the enum type.
	nConsts := len(e.Consts)
	if nConsts == 0 {
		return fmt.Errorf("enum type %q has no constants", e.Name)
	}

	// Ensure all discovered constants are actually constants.
	for _, obj := range e.Consts {
		if _, ok := obj.(*types.Const); !ok {
			return fmt.Errorf("enum type %q has a non-constant value: %T", e.Name, obj)
		}
	}

	// Need to have a names variable declared for the enum type.
	if e.NamesVar == nil {
		return fmt.Errorf("enum type %q has no names variable", e.Name)
	}

	// There must be a zero-valued constant declared for the enum type.
	if e.zeroConst() == nil {
		return fmt.Errorf("enum type %q has no zero-valued constant", e.Name)
	}

	// The number of constants must match the number of names in the names variable.
	if namesVar, ok := e.NamesVar.(*types.Var); ok {
		if namesVar.Kind() != types.PackageVar {
			return fmt.Errorf("enum type %q has a non-package variable names: %T", e.Name, e.NamesVar)
		}

		if !isNamesType(namesVar.Type()) {
			return fmt.Errorf("enum type %q has a names variable that is not a string slice or string table: %T", e.Name, e.NamesVar.Type().Underlying())
		}

		// TODO: use the packages.NeedsSyntax mode to get the ast and find the
		// CompositeLit that represents the names variable. You can then use its Elts
		// property to get the number of elements and compare that to the number of
		// constants. Unfortunately, there is no mapping from a types.Object to an
		// ast.Node to do checking if it is an ast.CompositeLit so ast traversal is
		// required, which is a bit too far for me to implement right now.
	}

	return nil
}

func (e *EnumType) zeroConst() types.Object {
	for _, obj := range e.Consts {
		if v, err := constValue(obj); err == nil && v == 0 {
			return obj
		}
	}
	return nil
}

// Attempts automatic discovery of the enum constants and 1D or 2D string slice variable
// that contains the string representations of the enum values.
func (e *EnumType) discover() {
	e.Consts = make([]types.Object, 0, 1)
	e.NamesVar = nil

	for _, name := range e.scope.Names() {
		// Constants are the same type as the enum type.
		obj := e.scope.Lookup(name)
		if types.Identical(obj.Type(), e.Type) {
			// Only add constants (skipping at the very least, the declared type)
			if _, ok := obj.(*types.Const); !ok {
				continue
			}

			e.Consts = append(e.Consts, obj)
			continue
		}

		// Name match means that the variable name is [name]Names.
		// If the name matches this pattern and it is a string slice or string table,
		// then it is set as the names variable.
		if e.namesMatch(name) && isNamesType(obj.Type()) {
			e.NamesVar = obj
			continue
		}
	}

	// Sort the consts by their values.
	slices.SortFunc(e.Consts, func(i, j types.Object) int {
		valI, _ := constValue(i)
		valJ, _ := constValue(j)
		return int(valI - valJ)
	})
}

func (e *EnumType) setNamesVar(name string) error {
	if e.NamesVar != nil {
		if e.NamesVar.Name() != name {
			return fmt.Errorf("already set or discovered names variable %q, cannot set to %q", e.NamesVar.Name(), name)
		}
		return nil
	}

	var obj types.Object
	if obj = e.scope.Lookup(name); obj == nil {
		return fmt.Errorf("names variable %q not found", name)
	}

	if !isNamesType(obj.Type()) {
		return fmt.Errorf("names variable %q is not a string slice, string array, or string table", name)
	}

	// Success! Set the names variable and return.
	e.NamesVar = obj
	return nil
}

func (e *EnumType) namesMatch(name string) bool {
	target := e.Name + "Names"
	ucc := UpperFirst(target)
	lcc := LowerFirst(target)
	return name == ucc || name == lcc
}

//============================================================================
// Helper Functions
//============================================================================

func isNamesType(typ types.Type) bool {
	return isStringArray(typ) || isStringSlice(typ) || isStringTable(typ) || isStrings2DArray(typ)
}

func isStringArray(typ types.Type) bool {
	array, ok := typ.Underlying().(*types.Array)
	if !ok {
		return false
	}
	basic, ok := array.Elem().Underlying().(*types.Basic)
	return ok && basic.Kind() == types.String
}

func isStringSlice(typ types.Type) bool {
	slice, ok := typ.Underlying().(*types.Slice)
	if !ok {
		return false
	}

	basic, ok := slice.Elem().Underlying().(*types.Basic)
	return ok && basic.Kind() == types.String
}

func isStringTable(typ types.Type) bool {
	outer, ok := typ.Underlying().(*types.Slice)
	if !ok {
		return false
	}
	return isStringSlice(outer.Elem())
}

func isStrings2DArray(typ types.Type) bool {
	outer, ok := typ.Underlying().(*types.Array)
	if !ok {
		return false
	}
	return isStringArray(outer.Elem())
}

func constValue(obj types.Object) (uint64, error) {
	if c, ok := obj.(*types.Const); ok {
		if v, ok := constant.Uint64Val(c.Val()); ok {
			return v, nil
		}
		return 0, fmt.Errorf("const %q has no uint64 value", obj.Name())
	}
	return 0, fmt.Errorf("%q is not a constant", obj.Name())
}
