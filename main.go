package main

import (
	"fmt"
	"github.com/obgnail/go-ast/ast_tool"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	d := ast_tool.GetDeclWalker()
	d.ParsePackage(`D:\golang\src\github.com\obgnail\net-link`)
}

func main2() {
	arg := `D:\golang\src\github.com\obgnail\net-link\link.go`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, arg, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		for _, v := range f.Decls {
			if s, ok := v.(*ast.GenDecl); ok && s.Tok == token.IMPORT {
				for _, v := range s.Specs {
					fmt.Println("import:", v.(*ast.ImportSpec).Path.Value)
				}
			}
		}
		return false
	})
}
