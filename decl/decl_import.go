package decl

import (
	"fmt"
	"github.com/obgnail/go-ast/namespace"
	"github.com/pingcap/errors"
	"go/ast"
	"os"
	"path"
	"strconv"
	"strings"
)

type Imports struct {
	*namespace.Tertiary // map[pkgPath, filePath, importName]importPath
}

func NewImports() *Imports {
	return &Imports{namespace.NewTertiary()}
}

func (i *Imports) CollectPkg(pkgPath string) (err error) {
	if i.Tertiary == nil {
		i.Tertiary = namespace.NewTertiary()
	}
	if pkgPath == "" {
		return
	}
	if _, ok := i.Get(pkgPath); ok {
		return
	}
	m, err := CollectImportFromPkgPath(pkgPath)
	if err != nil {
		return errors.Trace(err)
	}
	i.AppendSecondary(pkgPath, m)
	return
}

// map[filePath, importName]importPath
func CollectImportFromPkgPath(pkgPath string) (*namespace.Secondary, error) {
	res := namespace.NewSecondary()

	list, err := os.ReadDir(pkgPath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, d := range list {
		name := path.Join(pkgPath, d.Name())
		if strings.HasSuffix(name, ".go") {
			s, err := CollectImportFromFilePath(name)
			if err != nil {
				return nil, errors.Trace(err)
			}

			for ns2, value := range s {
				res.Set(name, ns2, value)
			}
		}
	}
	return res, nil
}

// return: map[importName]importPath
func CollectImportFromFilePath(filePath string) (map[string]string, error) {
	f, err := ParseFile(filePath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	res, err := CollectImportFromFile(f)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

// return: map[importName]importPath
func CollectImportFromFile(file *ast.File) (map[string]string, error) {
	res := make(map[string]string)

	for _, Import := range file.Imports {
		imp, err := strconv.Unquote(Import.Path.Value)
		if err != nil {
			return nil, errors.Trace(err)
		}

		importIdent := Import.Name

		// normalImport
		if importIdent == nil {
			res[path.Base(imp)] = imp
			// dotImport
		} else if importIdent.Name == "." {
			res["."] = imp
			// underscoreImport
		} else if importIdent.Name == "_" {
			res["_"] = imp
			// aliasImport
		} else if len(importIdent.Name) != 0 {
			res[importIdent.Name] = imp
			// error
		} else {
			return nil, fmt.Errorf("no such import type: %v", importIdent)
		}
	}

	return res, nil
}
