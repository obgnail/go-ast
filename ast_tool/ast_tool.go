package ast_tool

import (
	"go/ast"
	"go/token"
	"strings"
)

type Decl struct {
	name     string
	belong   string // belong package
	decl     ast.Decl
	baseDecl ast.Decl
}

func NewDecl(name, belong string, decl ast.Decl) *Decl {
	d := &Decl{
		name:     name,
		belong:   belong,
		decl:     decl,
		baseDecl: decl,
	}
	return d
}

type Package struct {
	path string
	pkg  *ast.Package

	PackageTopDecl map[string]*Decl                // map[declName]*Decl
	ImportTopDecl  map[ImportType]map[string]*Decl // map[ImportType]map[packagePath]*Decl
}

func CollectPackageTopDecl(pkgPath string, files []*ast.File) (decls []*Decl) {
	for _, fl := range files {
		for _, d := range fl.Decls {
			switch specDecl := d.(type) {
			case *ast.FuncDecl:
				// 方法不处理
				if specDecl.Recv == nil {
					declName := specDecl.Name.String()
					decls = append(decls, NewDecl(declName, pkgPath, specDecl))
				}
			case *ast.GenDecl:
				switch specDecl.Tok {
				case token.CONST, token.VAR:
					for _, valueSpec := range specDecl.Specs {
						declName := valueSpec.(*ast.ValueSpec).Names[0].Name
						decls = append(decls, NewDecl(declName, pkgPath, specDecl))
					}
				case token.TYPE:
					for _, typeSpec := range specDecl.Specs {
						declName := typeSpec.(*ast.TypeSpec).Name.Name
						decls = append(decls, NewDecl(declName, pkgPath, specDecl))
					}
				}
			}
		}
	}
	return
}

func CollectPackageImportTopDecl(pkgPath string, pkg *ast.Package, filter func(path string) bool) map[ImportType]map[string]*Decl {
	res := make(map[ImportType]map[string]*Decl)

	for _, fl := range pkg.Files {
		for _, imp := range fl.Imports {
			im := strings.Trim(imp.Path.Value, `"`)
			if filter != nil && !filter(im) {
				continue
			}

			importIdent := imp.Name

			// normalImport
			if importIdent == nil {
				CollectPackageTopDecl(imp.Path.Value)

				// dotImport
			} else if importIdent.Name == "." {

				// underscoreImport
			} else if importIdent.Name == "_" {

				// aliasImport
			} else if len(importIdent.Name) != 0 {

				// error
			} else {

			}
		}
	}
	return res
}
