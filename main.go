package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "."

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [PATH]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Print length of all functions in PATH.\n")
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

type FuncSig struct {
	loc string
	sig string
	len int
}

func (f *FuncSig) String() string {
	return fmt.Sprintf("% 4d %s %s", f.len, f.loc, f.sig)
}

type FoundFunc func(sig *FuncSig)

type FuncFinder struct {
	fset  *token.FileSet
	found FoundFunc
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

func (f *FuncFinder) parseFile(file *ast.File) error {
	var err error

	inspector := func(n ast.Node) bool {
		var funcToken *FuncToken
		switch fn := n.(type) {
		default:
			return true
		case *ast.FuncDecl:
			funcToken = &FuncToken{Pos: fn.Pos(), Body: fn.Body, Identity: fn}
		case *ast.FuncLit:
			funcToken = &FuncToken{Pos: fn.Pos(), Body: fn.Body, Identity: fn}
		}
		sig, err := f.parseFunc(funcToken)
		if err != nil {
			return false
		}
		f.found(sig)
		return true
	}

	ast.Inspect(file, inspector)
	return err
}

type FuncToken struct {
	Pos      token.Pos
	Body     *ast.BlockStmt
	Identity interface{}
}

func (f *FuncFinder) parseFunc(fn *FuncToken) (*FuncSig, error) {
	loc := f.fset.Position(fn.Pos).String()

	len, err := f.funcLen(fn)
	if err != nil {
		return nil, err
	}

	sig, err := f.funcSignature(fn)
	if err != nil {
		return nil, err
	}

	return &FuncSig{loc, sig, len}, nil
}

func (f *FuncFinder) funcLen(fn *FuncToken) (int, error) {
	if fn.Body == nil { // forward declaration
		return 0, nil
	}

	sLine := f.fset.Position(fn.Body.Pos()).Line
	eLine := f.fset.Position(fn.Body.End()).Line
	bodyLen := eLine - sLine
	if bodyLen == 0 {
		// Assuming at least one statement was on same line: "func() { stmt }"
		// which is incorrect for e.g. "func() {}".
		return 1, nil
	}
	return bodyLen, nil
}

const funcPrefixLen = len("func")

func (f *FuncFinder) funcSignature(fn *FuncToken) (string, error) {
	funcString, err := sprintNode(fn.Identity, f.fset)
	if err != nil {
		return "", err
	}

	if fn.Body == nil { // forward declaration
		return funcString[funcPrefixLen:], nil
	}

	// FIXME(uwe): I can't explain this, but somehow the sig vs body position is
	// influenced by tabs. And I assume replacing it with 1x8+1 spaces makes it
	// work?!
	funcString = strings.Replace(funcString, "\t", "         ", -1)

	sigLen := int(fn.Body.Pos() - fn.Pos)
	// FIXME(uwe): This should also not happen...
	if sigLen > len(funcString) {
		sigLen = len(funcString)
	}

	sig := strings.TrimSpace(funcString[funcPrefixLen:sigLen])
	return sig, nil
}

func sprintNode(node interface{}, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	p := &printer.Config{Tabwidth: 4}
	err := p.Fprint(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
