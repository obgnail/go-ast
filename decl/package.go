package decl

import (
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewPackage(pkgPath string, filter func(fs.FileInfo) bool) (path string, pkg *ast.Package, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, filter, parser.ParseComments)
	if err != nil {
		return "", nil, errors.Trace(err)
	}
	// TODO
	// 当存在 交叉编译标签 或者 //+build ignore 标签时, 有可能出现 len(pkgs) != 1
	if len(pkgs) != 1 {
		for name, pkg := range pkgs {
			for file := range pkg.Files {
				first, err := ReadFirstLine(file)
				if err != nil {
					return "", nil, errors.Trace(err)
				}
				if strings.Contains(first, "//+build") {
					delete(pkgs, name)
				}
				break
			}
		}
	}

	if len(pkgs) != 1 {
		return "", nil, fmt.Errorf("expected len(pkgs) == 1, got %d", len(pkgs))
	}

	for path, pkg := range pkgs {
		return path, pkg, nil
	}
	return "", nil, nil
}

func CollectPkgFromDir(dir string, filter func(fs.FileInfo) bool) (pkgs map[string]*ast.Package, err error) {
	pkgs = make(map[string]*ast.Package)

	err = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		ok, err := IsGolangPkg(p)
		if err != nil {
			return errors.Trace(err)
		}
		if ok {
			pkgPath, pkg, err := NewPackage(p, filter)
			if err != nil {
				return errors.Trace(err)
			}
			pkgs[pkgPath] = pkg
		}

		return nil
	})

	return
}
