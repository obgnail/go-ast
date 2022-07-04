package ast_tool

import (
	"github.com/pingcap/errors"
	"io/fs"
	"sync"
)

var once sync.Once

var DeclWalker *declWalker

func GetDeclWalker() *declWalker {
	once.Do(func() {
		DeclWalker = NewDeclWalker(nil, nil)
	})
	return DeclWalker
}

type declWalker struct {
	pkgs              map[string]*Package // map[pckPath]*Package
	parsePkgFilter    func(info fs.FileInfo) bool
	parseImportFilter func(path string) bool
}

func NewDeclWalker(
	parsePkgFilter func(info fs.FileInfo) (bool, ),
	parseImportFilter func(path string) (bool, ),
) *declWalker {
	d := &declWalker{
		pkgs:              make(map[string]*Package),
		parsePkgFilter:    parsePkgFilter,
		parseImportFilter: parseImportFilter,
	}
	return d
}

func (w *declWalker) ParsePackage(pkgPath string) error {
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

//func (w *declWalker) Visit(node ast.Node) ast.Visitor {
//	return w
//}
