package ast_tool

import (
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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

func NewPackage(path string, pkg *ast.Package) *Package {
	p := &Package{
		path:           path,
		pkg:            pkg,
		PackageTopDecl: make(map[string]*Decl),
		ImportTopDecl:  make(map[ImportType]map[string]*Decl),
	}
	return p
}

type DeclWalker struct {
	pkgs   map[string]*Package // map[pckPath]*Package
	filter func(info fs.FileInfo) bool
}

func NewDeclWalker(filter func(info fs.FileInfo) bool) *DeclWalker {
	d := &DeclWalker{
		pkgs:   make(map[string]*Package),
		filter: filter,
	}
	return d
}

func (w *DeclWalker) Visit(node ast.Node) ast.Visitor {
	return w
}

func (w *DeclWalker) ParsePackage(pkgPath string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, w.filter, parser.ParseComments)
	if err != nil {
		return errors.Trace(err)
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("expected len(pkgs) == 1, got %d", len(pkgs))
	}
	for _, pkg := range pkgs {
		w.SelectPackageTopDecl(pkgPath, pkg)
	}

	return nil
}

func (w *DeclWalker) SelectPackageTopDecl(pkgPath string, pkg *ast.Package) {
	w.pkgs[pkgPath] = NewPackage(pkgPath, pkg)
	for _, fl := range pkg.Files {
		for _, d := range fl.Decls {
			switch specDecl := d.(type) {
			case *ast.FuncDecl:
				// 方法不处理
				if specDecl.Recv == nil {
					declName := specDecl.Name.String()
					w.pkgs[pkgPath].PackageTopDecl[declName] = NewDecl(declName, pkgPath, specDecl)
				}
			case *ast.GenDecl:
				switch specDecl.Tok {
				case token.CONST, token.VAR:
					for _, valueSpec := range specDecl.Specs {
						declName := valueSpec.(*ast.ValueSpec).Names[0].Name
						w.pkgs[pkgPath].PackageTopDecl[declName] = NewDecl(declName, pkgPath, specDecl)
					}
				case token.TYPE:
					for _, typeSpec := range specDecl.Specs {
						declName := typeSpec.(*ast.TypeSpec).Name.Name
						w.pkgs[pkgPath].PackageTopDecl[declName] = NewDecl(declName, pkgPath, specDecl)
					}
				}
			}
		}
	}
}
