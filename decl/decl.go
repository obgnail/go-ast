package decl

import (
	"github.com/obgnail/go-ast/namespace"
	"github.com/pingcap/errors"
	"go/ast"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Decls struct {
	*namespace.Secondary // map[pkgPath, declName]*Decl
}

func NewDecls() *Decls {
	return &Decls{namespace.NewSecondary()}
}

func FixBaseDecl(decls *Decls, imports *Imports) {

}

func FindDeclFromSelectorExpr(expr *ast.SelectorExpr, filePath string, decls *namespace.Secondary) *Decl {
	return nil
}

func (d *Decls) CollectPkg(pkgPath string, fileFilter func(info fs.FileInfo) bool) (err error) {
	if d.Secondary == nil {
		d.Secondary = namespace.NewSecondary()
	}
	if pkgPath == "" {
		return
	}
	if _, ok := d.Get(pkgPath); ok {
		return
	}
	m, err := CollectTopDeclFromPkgPath(pkgPath, fileFilter)
	if err != nil {
		return errors.Trace(err)
	}
	for key, value := range m {
		d.Set(pkgPath, key, value)
	}
	return
}

type Decl struct {
	name     string
	belong   string // belong package
	decl     ast.Decl
	baseDecl ast.Decl // 专门针对type of、alias
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

// 获取某个package下的所有顶级decls
// return: map[declName]*Decl
func CollectTopDeclFromPkgPath(pkgPath string, fileFilter func(info fs.FileInfo) bool) (decls map[string]*Decl, err error) {
	pkgPath, pkg, err := NewASTPackage(pkgPath, fileFilter)
	if err != nil {
		return nil, errors.Trace(err)
	}
	res, err := collectPkgTopDecl(pkgPath, pkg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func collectPkgTopDecl(pkgPath string, pkg *ast.Package) (decls map[string]*Decl, err error) {
	decls = make(map[string]*Decl)

	for _, fl := range pkg.Files {
		for _, d := range fl.Decls {

			switch specDecl := d.(type) {
			case *ast.FuncDecl:
				// 只处理函数,方法不处理
				if specDecl.Recv == nil {
					declName := specDecl.Name.String()
					decls[declName] = NewDecl(declName, pkgPath, specDecl)
				}

			case *ast.GenDecl:
				switch specDecl.Tok {
				case token.CONST, token.VAR:
					for _, valueSpec := range specDecl.Specs {
						declName := valueSpec.(*ast.ValueSpec).Names[0].Name
						decls[declName] = NewDecl(declName, pkgPath, specDecl)
					}
				case token.TYPE:
					for _, typeSpec := range specDecl.Specs {
						declName := typeSpec.(*ast.TypeSpec).Name.Name
						decls[declName] = NewDecl(declName, pkgPath, specDecl)
					}
				}
			}
		}
	}
	return
}

// 获取某个dir下的所有顶级decls
// return: map[pkgPath, declName]*Decl
func CollectTopDeclFromDir(dir string, fileFilter func(info fs.FileInfo) bool) (decls *namespace.Secondary, err error) {
	decls = namespace.NewSecondary()

	err = filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		relP, err := filepath.Rel(dir, p)
		if err != nil {
			return errors.Trace(err)
		}
		base := path.Base(relP)
		if strings.HasPrefix(base, "test") {
			return filepath.SkipDir
		}

		list, err := os.ReadDir(p)
		if err != nil {
			return errors.Trace(err)
		}
		for _, d := range list {
			// 是目录,并且底下有go文件 则视为package
			if strings.HasSuffix(d.Name(), ".go") {
				m, err := CollectTopDeclFromPkgPath(p, fileFilter)
				if err != nil {
					return errors.Trace(err)
				}

				for declName, decl := range m {
					decls.Set(p, declName, decl)
				}
				break
			}
		}

		return nil
	})

	return
}

// 获取某个dir下的所有import package的decls
// return: map[pkgPath, declName]*Decl
func CollectDirTopDeclFromImports(
	dir string,
	fileFilter func(info fs.FileInfo) bool,
	importFilter func(path string) bool,
	fixImportPath func(path string) string,
) (imports *namespace.Secondary, err error) {

	imports = namespace.NewSecondary()

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		if fileFilter != nil && !fileFilter(info) {
			return nil
		}

		f, err := ParseFile(path)
		if err != nil {
			return errors.Trace(err)
		}

		for _, Import := range f.Imports {

			imp, err := strconv.Unquote(Import.Path.Value)
			if err != nil {
				return errors.Trace(err)
			}

			if importFilter != nil && !importFilter(imp) {
				continue
			}

			if fixImportPath != nil {
				imp = fixImportPath(imp)
			}

			if _, ok := imports.Get(imp); ok {
				continue
			}

			decls, err := CollectTopDeclFromPkgPath(imp, fileFilter)
			if err != nil {
				return errors.Trace(err)
			}
			for name, decl := range decls {
				imports.Set(imp, name, decl)
			}

		}
		return nil
	})

	return
}

func FixImportPath(path string) string {
	// TODO
	// 本地
	if strings.HasPrefix(path, "github.com/bangwork/bang-api") {
		path = filepath.Join("/Users/heyingliang/go/src/", path)
		// vendor
	} else {
		path = filepath.Join("/Users/heyingliang/go/src/github.com/bangwork/bang-api/vendor/", path)
	}
	return path
}

func FilterImport(path string) bool {
	if !strings.HasPrefix(path, "github.com") {
		return false
	}
	return true
}
