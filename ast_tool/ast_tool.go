package ast_tool

import (
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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

type DeclWalker struct {
	pkgs              map[string]*Package // map[pckPath]*Package
	parsePkgFilter    func(info fs.FileInfo) bool
	parseImportFilter func(path string) bool
}

func NewDeclWalker(
	parsePkgFilter func(info fs.FileInfo) (bool, ),
	parseImportFilter func(path string) (bool, ),
) *DeclWalker {
	d := &DeclWalker{
		pkgs:              make(map[string]*Package),
		parsePkgFilter:    parsePkgFilter,
		parseImportFilter: parseImportFilter,
	}
	return d
}

//func (w *DeclWalker) Visit(node ast.Node) ast.Visitor {
//	return w
//}

func (w *DeclWalker) ParsePackage(pkgPath string) error {
	pkgs, err := ParseDir(pkgPath, w.parsePkgFilter)
	if err != nil {
		return errors.Trace(err)
	}
	for _, pkg := range pkgs {
		files := FileMapToList(pkg.Files)
		pkgDecls := CollectPackageTopDecl(pkgPath, files)
		importDecls := CollectPackageImportTopDecl(pkgPath, pkg, w.parseImportFilter)

		w.pkgs[pkgPath] = &Package{
			path:           pkgPath,
			pkg:            pkg,
			PackageTopDecl: DeclListToMap(pkgDecls),
			ImportTopDecl:  importDecls,
		}
	}

	return nil
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

				// dotImport
			} else if importIdent.Name == "." {

				// underscoreImport
			} else if importIdent.Name == "_" {

				// aliasImport
				// 注意alias一样
			} else {

			}
		}
	}
	return res
}

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
