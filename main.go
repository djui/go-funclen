package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// TODO(uwe): Sort list
// TODO(uwe): Recursive
// TODO(uwe): Flexible output width
// TODO(uwe): Make cli
// TODO(uwe): Allow specifying directory
// TODO(uwe): Format printer.Fprint with width 4

func main() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch fn := n.(type) {
				case *ast.FuncDecl:
					funcLine := fset.Position(fn.Pos()).String()
					funcSignature := funcSignature(fn, fset)
					funcLen := funcLen(fn, fset)
					fmt.Printf("%-25s % 4d %s\n", funcLine, funcLen, funcSignature)
				}
				return true
			})
		}
	}
}

func funcSignature(f *ast.FuncDecl, fset *token.FileSet) string {
	name := f.Name.Name

	var recv string
	if f.Recv.NumFields() > 0 {
		var recvName string
		r := f.Recv.List[0]
		if len(r.Names) > 0 {
			recvName = fmt.Sprintf("%s ", r.Names[0].Name)
		}

		recvType, _ := formatExpr(r.Type, fset)

		recv = fmt.Sprintf("(%s%s) ", recvName, recvType)
	}

	ftype, _ := formatExpr(f.Type, fset)
	paramsAndResults := ftype[4:] // Shortcut: "^func"

	return fmt.Sprintf("%s%s%s", recv, name, paramsAndResults)
}

func funcLen(f *ast.FuncDecl, fset *token.FileSet) int {
	start := fset.Position(f.Body.Pos()).Line
	end := fset.Position(f.Body.End()).Line
	return end - start
}

func formatExpr(node ast.Expr, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
