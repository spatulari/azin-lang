package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azin-lang/Azin/internal/compiler"
	"github.com/azin-lang/Azin/internal/fs"
)

const Version = "0.2.3-dev"

var (
	debug           = flag.Bool("debug", false, "Enable debug output")
	printTokens     = flag.Bool("tokens", false, "Print lexer tokens")
	output          = flag.String("o", "", "Output file")
	ignoreExtension = flag.Bool("ignore-extension", false, "Ignore source file extension")
	version         = flag.Bool("version", false, "Print compiler version")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <file>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Azin compiler %s\n", Version)
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		printDebug()
	}

	source := mustReadSource(flag.Arg(0))

	if err := compiler.Compile(source, *printTokens); err != nil {
		fatal(err)
	}

	if *debug {
		fmt.Printf("Compiled %.2f KiB\n", float64(len(source))/1024)
	}
}

func printDebug() {
	fmt.Printf("Debug: %t\n", *debug)
	fmt.Printf("Print tokens: %t\n", *printTokens)
	fmt.Printf("Output: %q\n", *output)
}

func mustReadSource(filename string) []byte {
	source, err := fs.ReadSourceFile(filename, *ignoreExtension)
	if err != nil {
		fatal(err)
	}
	return source
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
