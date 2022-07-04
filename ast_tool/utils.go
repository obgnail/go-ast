package ast_tool

import (
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
)

func ParseDir(dir string, filter func(fs.FileInfo) bool) (pkgs map[string]*ast.Package, err error) {
	fset := token.NewFileSet()
	pkgs, err = parser.ParseDir(fset, dir, filter, parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected len(pkgs) == 1, got %d", len(pkgs))
	}
	return
}

func DeclListToMap(delcs []*Decl) map[string]*Decl {
	res := make(map[string]*Decl)
	for _, decl := range delcs {
		res[decl.name] = decl
	}
	return res
}

func FileMapToList(m map[string]*ast.File) (files []*ast.File) {
	for _, f := range m {
		files = append(files, f)
	}
	return files
}
