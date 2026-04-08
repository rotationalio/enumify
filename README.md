# Enumify

[![CI Tests](https://github.com/rotationalio/enumify/actions/workflows/tests.yaml/badge.svg)](https://github.com/rotationalio/enumify/actions/workflows/tests.yaml)

**Easily create and manage typed enumerations with code generation and automated enum tests.**

We found ourselves generating a lot of boilerplate code for our Enums: particularly for parsing, serialization/deserialization, database storage, and tests. Really we just want to be able to describe an Enum and get all of this code for free! Enumify is the combination of a code generator for the boilerplate code as well as a package for reducing the boilerplate and ensuring that everything works as expected.

## Getting Started

Install enumify:

```sh
go install go.rtnl.ai/enumify/cmd/enumify@latest
```

This will add the `enumify` command to your `$PATH`.

Next, in the package that you'd like to create an enum for, define a `.go` file as follows:

```go
//go:generate enumify
type Status uint8

// Constants for the Status enum values.
// These values should be discovered by the enumify generator since they use the same
// type as the Status enum, which is the type being generated.
const (
	StatusUnknown Status = iota
	StatusPending
	StatusRunning
	StatusFailed
	StatusSuccess
	StatusCancelled
)

// String representations for the Status enum values; the enumify generator will
// discover this since it has the pattern [enum]Names and a []string or [][]string type.
var statusNames []string{
    "unknown", "pending", "running", "failed", "success", "canceled",
}
```

Run `go generate ./...` and two files `[file]_gen.go` and `file_gen_test.go` will be generated with a Parser function and all of the interface methods for the enum as well as a test suite for testing the enum wrt the enumify package.

Enumify uses static parsing to discover the type definitions in your `.go` file. It looks for any `type [Enum] uint8` definition, and determines that is the enum being used. Any constants that are of type `[Enum]` will be extracted as well as the string reprs in a `[enum]Names` variable.

You can also use a `[][]string` table to capture other enum information. For example:

```go
//go:generate enumify -names colorTable
type Color uint8

const (
    NoColor Color = iota
    Red
    Blue
    Green
)

var colorTable [][]string{
    {"none", "red", "blue", "green"},
    {"", "#FF0000", "#0000FF", "#00FF00"}
}

func (c Color) Hex() string {
    return colorNames[1][c]
}
```

You can call the string representations variable anything you'd like so long as you pass in the `-names` flag with the variable name; otherwise the `[enum]Names` pattern works for both `[]string` and `[][]string` variables.

> **NOTE**: The generator only works on one file at a time, if you have multiple files with enums in each file; you'll need a `go:generate` directive in each file. However, if you have multiple enums in a single file, you'll only need one `go:generate` directive.

## CLI Options

You can pass options to the `go:generate` directive as you would any CLI program. The latest help text is as follows:

```text
Enumify is a code generation tool for easily creating and managing
typed enumerations with code generation and automated enum tests.

Usage: add a go directive to a go file file for generation:

  //go:generate enumify [-version] [-help] [opts]

Then run run the go tool code generation command:

  go generate ./...

Options:
  -case-sensitive
    	make the enum case sensitive
  -names string
    	variable name that contains the string reprs of the enum values
  -no-binary
    	skip binary interfaces code generation
  -no-json
    	skip JSON interfaces code generation
  -no-parser
    	skip parser code generation
  -no-sql
    	skip SQL interfaces code generation
  -no-stringer
    	skip Stringer interface code generation
  -no-tests
    	skip testing code generation
  -no-text
    	skip text interfaces code generation
  -no-yaml
    	skip YAML interfaces code generation
  -space-sensitive
    	make the enum space sensitive
  -version
    	print the version and exit
```

## Theory

For us, an `enum` is a set of specific values that have string representations. In code we want our `enum` values to be `uint8` (if you have more than 255 values you probably want a database table instead of an `enum`) for compact memory usage and ease of comparisons. However for marshaling/unmarshaling the value to json/yaml or saving the value in a database, we'd like to use the string representation instead.

However, the boilerplate code for implementing `fmt.Stringer`, `json.Marshaler`, `json.Unmarshaler`, etc. is pretty verbose -- not to mention the testing. This library unifies all of that functionality into generic functions and also provides code generation to make it simple to define enumerations. Modifications require a change to the enum schema file and running `go generate ./...` and presto chango, a minimal set of code is added to your project, but code that is fully tested and organized.

Key theory point: the value behind the enum shouldn't matter, just its ordering. The first value (0) should always be the "unknown" or "default" value. You can compare enums but you cannot have them be specific runes or other values.

## References

The following blog posts and package documentation was invaluable to the implementation of this package:

- [A taste of Go code generator magic: a quick guide to getting started](https://evilmartians.com/chronicles/a-taste-of-go-code-generator-magic-a-quick-guide-to-getting-started)
- [Jennifer is a code generator for Go.](https://github.com/dave/jennifer)

## License

This project is licensed under the BSD 3-Clause License. See [`LICENSE`](./LICENSE) for details. Please feel free to use Enumify in your own projects and applications.

## About Rotational Labs

Enumify is developed by [Rotational Labs](https://rotational.io), a team of engineers and scientists building AI infrastructure for serious work.
