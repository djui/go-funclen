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
	root := "."

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [PATH]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Print a list of all Go functions and their body length recursively in the given directory PATH.\n")
	}

	flag.Parse()
	if flag.NArg() == 1 {
		root = flag.Arg(0)
	} else if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	run(root)
}

func run(root string) {
	//var funcs []*FuncSig
	//storeFound := func(sig *FuncSig) { funcs = append(funcs, sig) }
	// :
	// for _, fn := range funcs {
	// 	fmt.Println(fn)
	// }

	printFound := func(sig *FuncSig) { fmt.Println(sig) }
	f := NewFuncFinder(printFound)
	err := filepath.Walk(root, f.VisitDirs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

type FoundFunc func(sig *FuncSig)

type FuncFinder struct {
	fset  *token.FileSet
	found FoundFunc
}

type FuncToken struct {
	Pos  token.Pos
	Body *ast.BlockStmt
	Name string
}

type FuncSig struct {
	loc string
	sig string
	len int
}

func (f *FuncSig) String() string {
	return fmt.Sprintf("%s: %d %s", f.loc, f.len, f.sig)
}

func NewFuncFinder(found FoundFunc) *FuncFinder {
	return &FuncFinder{token.NewFileSet(), found}
}

func (f *FuncFinder) VisitDirs(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		if err := f.parseDir(path); err != nil {
			return err
		}
	}

	return nil
}

func (f *FuncFinder) parseDir(path string) error {
	pkgs, err := parser.ParseDir(f.fset, path, nil, 0)
	if err != nil {
		return err
	}

	return f.parsePackages(pkgs)
}

func (f *FuncFinder) parsePackages(pkgs map[string]*ast.Package) error {
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if err := f.parseFile(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *FuncFinder) parseFile(file ast.Node) error {
	var err error

	inspector := func(n ast.Node) bool {
		var funcToken *FuncToken
		switch fn := n.(type) {
		default:
			return true
		case *ast.FuncDecl:
			funcToken = &FuncToken{Pos: fn.Pos(), Body: fn.Body, Name: fn.Name.String()}
		case *ast.FuncLit:
			funcToken = &FuncToken{Pos: fn.Pos(), Body: fn.Body, Name: "(anonymous)"}
		}
		f.found(f.parseFunc(funcToken))
		return true
	}

	ast.Inspect(file, inspector)
	return err
}

func (f *FuncFinder) parseFunc(fn *FuncToken) *FuncSig {
	loc := f.fset.Position(fn.Pos).String()
	len := f.funcLen(fn)
	sig := fn.Name
	return &FuncSig{loc, sig, len}
}

func (f *FuncFinder) funcLen(fn *FuncToken) int {
	if fn.Body == nil {
		return 0 // forward declaration
	}

	sLine := f.fset.Position(fn.Body.Pos()).Line
	eLine := f.fset.Position(fn.Body.End()).Line
	bodyLen := eLine - sLine
	if bodyLen == 0 {
		// Assuming at least one statement was on same line: "func() { stmt }"
		// which is incorrect for e.g. "func() {}".
		return 1
	}

	return bodyLen
}
