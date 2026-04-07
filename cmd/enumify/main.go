package main

import (
	"flag"
	"fmt"
	"os"

	"go.rtnl.ai/enumify"
)

func main() {
	// Set up the flags and options for go generate processing.
	versionFlag := flag.Bool("version", false, "print the version and exit")
	opts := enumify.Options{
		File: os.Getenv("GOFILE"),
		Pkg:  os.Getenv("GOPACKAGE"),
	}

	// Bind flag variables directly to the opts struct
	flag.StringVar(&opts.NameVar, "names", "", "variable name that contains the string reprs of the enum values")
	flag.BoolVar(&opts.CaseSensitive, "case-sensitive", false, "make the enum case sensitive")
	flag.BoolVar(&opts.SpaceSensitive, "space-sensitive", false, "make the enum space sensitive")
	flag.BoolVar(&opts.NoTests, "no-tests", false, "skip testing code generation")
	flag.BoolVar(&opts.NoParser, "no-parser", false, "skip parser code generation")
	flag.BoolVar(&opts.NoStringer, "no-stringer", false, "skip Stringer interface code generation")
	flag.BoolVar(&opts.NoText, "no-text", false, "skip text interfaces code generation")
	flag.BoolVar(&opts.NoBinary, "no-binary", false, "skip binary interfaces code generation")
	flag.BoolVar(&opts.NoJSON, "no-json", false, "skip JSON interfaces code generation")
	flag.BoolVar(&opts.NoYAML, "no-yaml", false, "skip YAML interfaces code generation")
	flag.BoolVar(&opts.NoSQL, "no-sql", false, "skip SQL interfaces code generation")

	// Set the usage function and parse command line arguments
	flag.Usage = usage
	flag.Parse()

	// Print the version and exit if the version flag is set
	if *versionFlag {
		fmt.Printf("enumify %s\n", enumify.Version())
		os.Exit(0)
	}

	// Perform code generation!
	if err := enumify.Generate(opts); err != nil {
		fmt.Fprintf(os.Stderr, "enumify: %v\n", err)
		os.Exit(1)
	}

	// Generate tests if not skipped
	if !opts.NoTests {
		if err := enumify.GenerateTests(opts); err != nil {
			fmt.Fprintf(os.Stderr, "enumify: %v\n", err)
			os.Exit(1)
		}
	}

	// Signal that we succeeded with the code generation
	os.Exit(0)
}

func usage() {
	// Get the default output (usually stderr)
	w := flag.CommandLine.Output()

	fmt.Fprintln(w, "Enumify is a code generation tool for easily creating and managing")
	fmt.Fprintln(w, "typed enumerations with code generation and automated enum tests.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage: add a go directive to a go file file for generation:")
	fmt.Fprintln(w, "\n  //go:generate enumify [-version] [-help] [opts]")
	fmt.Fprintln(w, "\nThen run run the go tool code generation command:")
	fmt.Fprintln(w, "\n  go generate ./...")
	fmt.Fprintln(w, "\nOptions:")
	flag.PrintDefaults()
}
