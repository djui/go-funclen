package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [PATH]
Print a list of unused global tokens (constants, variables, functions, types,
interfaces) and function parameters found recursively in packages in the given
directory PATH.
`, filepath.Base(os.Args[0]))
	}

	flag.Parse()

	switch flag.NArg() {
	case 0:
		run(".")
	case 1:
		run(flag.Arg(0))
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func run(rootDir string) {
	findUnusedConstants(rootDir)
	findUnusedVariables(rootDir)
	findUnusedTypes(rootDir)
	findUnusedInterfaces(rootDir)
	findUnusedFunctions(rootDir)
	findUnusedParameters(rootDir)
}

var fset = token.NewFileSet()

func findUnusedConstants(rootDir string) {
	err := filepath.Walk(rootDir, visitDirs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func findUnusedVariables(rootDir string) {
	// TODO
}

func findUnusedTypes(rootDir string) {
	// TODO
}

func findUnusedInterfaces(rootDir string) {
	// TODO
}

func findUnusedFunctions(rootDir string) {
	// TODO
}

func findUnusedParameters(rootDir string) {
	// TODO: Find all funtions

	// TODO: Get all parameter names per function

	// TODO: Check if parameter is actually used in the function.

	// TODO: Print function name and paramter name

	// Don't think about shadowed variables for now.
}

func visitDirs(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		if err := parsePath(path); err != nil {
			return err
		}
	}

	return nil
}

func parsePath(path string) error {
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return err
	}

	return parsePackages(pkgs)
}

func parsePackages(pkgs map[string]*ast.Package) error {
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if err := parseFile(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseFile(root ast.Node) error {
	var err error

	inspector := func(n ast.Node) bool {
		consts := collectConsts(n)
		for _, c := range consts {
			fmt.Printf("%s\n", c)
		}
		return true
	}

	ast.Inspect(root, inspector)
	return err
}

func collectConsts(n ast.Node) []string {
	// TODO: Return package name as well: {pkg, name}
	var consts []string
	if d, ok := n.(*ast.GenDecl); ok && d.Tok == token.CONST {
		for _, s := range d.Specs {
			consts = append(consts, s.(*ast.ValueSpec).Names[0].Name)
		}
	}
	return consts
}
