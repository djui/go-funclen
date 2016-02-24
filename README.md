go-funclen
==========

Prints a list of all Go functions and their body length recursively in a given
directory.


# Usage

    $ go-funclen -h
    Usage: main [PATH]
    Print a list of all Go functions and their body length recursively in the given directory PATH.

    $ go-funclen .
    17 ../../../go-funclen/main.go:16:1 main()
     3 ../../../go-funclen/main.go:19:15 ()
    14 ../../../go-funclen/main.go:35:1 run(root string)
     1 ../../../go-funclen/main.go:43:16 (sig *FuncSig)
     2 ../../../go-funclen/main.go:57:1 (f *FuncSig) String() string
     2 ../../../go-funclen/main.go:68:1 NewFuncFinder(found FoundFunc) *FuncFinder
    12 ../../../go-funclen/main.go:72:1 (f *FuncFinder) VisitDirs(path string, fi os.FileInfo, err     error) error
     7 ../../../go-funclen/main.go:86:1 (f *FuncFinder) parseDir(path string) error
    10 ../../../go-funclen/main.go:95:1 (f *FuncFinder) parsePackages(pkgs map[string]*ast.Package)     error
    23 ../../../go-funclen/main.go:107:1 (f *FuncFinder) parseFile(file *ast.File) error
    16 ../../../go-funclen/main.go:110:15 (n ast.Node) bool
    14 ../../../go-funclen/main.go:138:1 (f *FuncFinder) parseFunc(fn *FuncToken) (*FuncSig, error)
    14 ../../../go-funclen/main.go:154:1 (f *FuncFinder) funcLen(fn *FuncToken) (int, error)
    23 ../../../go-funclen/main.go:172:1 (f *FuncFinder) funcSignature(fn *FuncToken) (string, error)
     8 ../../../go-funclen/main.go:197:1 sprintNode(node interface{}, fset *token.FileSet) (string, error)

# Installation

    $ go get -v -u github.com/djui/go-funclen
