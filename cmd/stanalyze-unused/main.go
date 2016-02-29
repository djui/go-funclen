package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
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
	// findUnusedVariables(rootDir)
	// findUnusedTypes(rootDir)
	// findUnusedInterfaces(rootDir)
	// findUnusedFunctions(rootDir)
	// findUnusedParameters(rootDir)
}

type constRef struct {
	pkg  string
	name string
	pos  token.Position
	nInt int
	nExt int
}

func findUnusedConstants(rootDir string) {
	consts, err := findConstants(rootDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	var fset = token.NewFileSet()
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".go") {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		inspector := func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.Ident: // locally referenced
				pkg := f.Name.Name
				name := n.Name
				fullname := fmt.Sprintf("%s.%s", pkg, name)
				if _, ok := consts[fullname]; ok {
					consts[fullname].nInt++
				}
			case *ast.SelectorExpr: // globally referenced
				if x, ok := n.X.(*ast.Ident); ok {
					pkg := x.Name
					name := n.Sel.Name
					fullname := fmt.Sprintf("%s.%s", pkg, name)
					if _, ok := consts[fullname]; ok {
						consts[fullname].nExt++
					}
				}
			}
			return true
		}

		//ast.Print(fset, f)
		ast.Inspect(f, inspector)

		return nil
	}

	if err := filepath.Walk(rootDir, walker); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

	for _, c := range consts {
		if c.nExt == 0 {
			if c.nInt == 1 {
				fmt.Printf("%v: const %s unused\n", c.pos, c.name)
			} else if ast.IsExported(c.name) {
				fmt.Printf("%v: const %s exported unnecessarily\n", c.pos, c.name)
			}
		}
	}
}

func findConstants(rootDir string) (map[string]*constRef, error) {
	var fset = token.NewFileSet()
	var consts = map[string]*constRef{}

	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".go") {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			// return err
			fmt.Fprintln(os.Stderr, "Error:", err)
		}

		inspector := func(n ast.Node) bool {
			if d, ok := n.(*ast.GenDecl); ok && d.Tok == token.CONST {
				for _, s := range d.Specs {
					ident := s.(*ast.ValueSpec).Names[0]
					if ident.IsExported() && f.Scope.Lookup(ident.Name) == nil {
						// Exported (uppercased) but has no global scope
						continue
					}

					pkg := f.Name.Name
					name := ident.Name
					fullname := fmt.Sprintf("%s.%s", pkg, name)

					consts[fullname] = &constRef{
						pkg:  pkg,
						name: name,
						pos:  fset.Position(s.Pos()),
					}
				}
			}
			return true
		}

		ast.Inspect(f, inspector)

		return nil
	}

	if err := filepath.Walk(rootDir, walker); err != nil {
		return nil, err
	}

	return consts, nil
}

// func findUnusedVariables(rootDir string) {
// 	// TODO
// }
//
// func findUnusedTypes(rootDir string) {
// 	// TODO
// }
//
// func findUnusedInterfaces(rootDir string) {
// 	// TODO
// }
//
// func findUnusedFunctions(rootDir string) {
// 	// TODO
// }
//
// func findUnusedParameters(rootDir string) {
// 	// TODO: Find all funtions
//
// 	// TODO: Get all parameter names per function
//
// 	// TODO: Check if parameter is actually used in the function.
//
// 	// TODO: Print function name and paramter name
//
// 	// Don't think about shadowed variables for now.
// }
