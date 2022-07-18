package main

import (
	"fmt"
	"github.com/obgnail/go-ast/decl"
	"path/filepath"
	"strings"
)

func FixImportPath(path string) string {
	// TODO
	// 本地
	if strings.HasPrefix(path, "github.com/") {
		path = filepath.Join("/Users/heyingliang/go/src/", path)
		// vendor
	} else {
		//path = filepath.Join("/Users/heyingliang/go/src/github.com/bangwork/bang-api/vendor/", path)
	}
	return path
}

func main() {
	p := decl.NewPackages()
	err := p.Collect("/Users/heyingliang/go/src/github.com/obgnail/go-ast/example/testdata/net-link", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)

	//decls, err := decl.CollectTopDeclFromDir(
	//	"/Users/heyingliang/go/src/github.com/obgnail/go-ast/example/testdata/net-link", nil,
	//)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(decls)

	res, err := decl.ParseType("/Users/heyingliang/go/src/github.com/obgnail/go-ast/example/testdata/net-link/codec/codec.go")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
