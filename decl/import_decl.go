package decl

import (
	"fmt"
	"github.com/obgnail/go-ast/namespace"
	"github.com/pingcap/errors"
	"go/ast"
	"strconv"
)

// return: map[ImportType, importName]importPath
func CollectImportFromFilePath(filePath string) (*namespace.Secondary, error) {
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

// return: map[ImportType, importName]importPath
func CollectImportFromFile(file *ast.File) (*namespace.Secondary, error) {
	ns := namespace.NewSecondary()

	for _, Import := range file.Imports {
		imp, err := strconv.Unquote(Import.Path.Value)
		if err != nil {
			return nil, errors.Trace(err)
		}

		importIdent := Import.Name

		// normalImport
		if importIdent == nil {
			ns.Set(normalImport, imp, imp)
			// dotImport
		} else if importIdent.Name == "." {
			ns.Set(dotImport, imp, imp)
			// underscoreImport
		} else if importIdent.Name == "_" {
			ns.Set(underscoreImport, imp, imp)
			// aliasImport
		} else if len(importIdent.Name) != 0 {
			ns.Set(underscoreImport, importIdent.Name, imp)
			// error
		} else {
			return nil, fmt.Errorf("no such import type: %v", importIdent)
		}
	}

	return ns, nil
}
