package decl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ParseFile(file string) (f *ast.File, err error) {
	fset := token.NewFileSet()
	f, err = parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return
}

// 是目录,并且底下有go文件 则视为package
func IsGolangPkg(dir string) (bool, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return false, errors.Trace(err)
	}
	for _, d := range list {
		if strings.HasSuffix(d.Name(), ".go") {
			return true, nil
		}
	}
	return false, nil
}

func getStdLibs(goPath string) (libs []string, err error) {
	err = filepath.Walk(goPath, func(p string, info fs.FileInfo, err error) error {
		relP, err := filepath.Rel(goPath, p)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		base := path.Base(relP)
		if strings.HasPrefix(base, "test") || strings.HasPrefix(base, "vendor") || base == "internal" {
			return filepath.SkipDir
		}

		list, err := os.ReadDir(p)
		if err != nil {
			return err
		}
		for _, d := range list {
			if strings.HasSuffix(d.Name(), ".go") {
				libs = append(libs, relP)
				break
			}
		}
		return nil
	})
	return
}

func ReadFirstLine(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("read first line err")
}

// 某些结构体调用直接在ast语法树上处理不方便，所以我们把他转换成字符串的形式来处理
func ExprTokenToStr(node ast.Expr, Fset *token.FileSet) string {
	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := format.Node(buffer, Fset, node); err != nil {
		panic(err)
	}
	return buffer.String()
}
