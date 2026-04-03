# Enumify

**Easily create and manage typed enumerations with code generation and automated enum tests.**

We found ourselves generating a lot of boilerplate code for our Enums: particularly for parsing, serialization/deserialization, database storage, and tests. Really we just want to be able to describe an Enum and get all of this code for free! Enumify is the combination of a code generator for the boilerplate code as well as a package for reducing the boilerplate and ensuring that everything works as expected.

## Getting Started

Install enumify:

```sh
go install go.rtnl.ai/enumify/cmd/enumify@latest
```

This will add the `enumify` command to your `$PATH`.

Next steps:

1. Define enum schema file
2. Create generate command
3. Run go generate ./...
4. Run go mod tidy

Boom - your package is ready to go with enums!

## Theory

For us, an `enum` is a set of specific values that have string representations. In code we want our `enum` values to be `uint8` (if you have more than 255 values you probably want a database table instead of an `enum`) for compact memory usage and ease of comparisons. However for marshaling/unmarshaling the value to json/yaml or saving the value in a database, we'd like to use the string representation instead.

However, the boilerplate code for implementing `fmt.Stringer`, `json.Marshaler`, `json.Unmarshaler`, etc. is pretty verbose -- not to mention the testing. This library unifies all of that functionality into generic functions and also provides code generation to make it simple to define enumerations. Modifications require a change to the enum schema file and running `go generate ./...` and presto chango, a minimal set of code is added to your project, but code that is fully tested and organized.

Key theory point: the value behind the enum shouldn't matter, just its ordering. The first value (0) should always be the "unknown" or "default" value. You can compare enums but you cannot have them be specific runes or other values.

## License

This project is licensed under the BSD 3-Clause License. See [`LICENSE`](./LICENSE) for details. Please feel free to use Enumify in your own projects and applications.

## About Rotational Labs

Enumify is developed by [Rotational Labs](https://rotational.io), a team of engineers and scientists building AI infrastructure for serious work.
