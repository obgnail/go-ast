package decl

import (
	"github.com/obgnail/go-ast/namespace"
	"go/ast"
)

type Package struct {
	path           string
	pkg            *ast.Package
	PackageTopDecl map[string]*Decl    // map[declName]*Decl
	Imports        *namespace.Tertiary // map[filePath, ImportType, importName]importPath
}
