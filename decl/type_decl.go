package decl

import (
	"fmt"
	"github.com/pingcap/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
)

type StructCategory = int
type FieldCategory = int

const (
	StructValid        StructCategory = 0
	StructFieldInvalid StructCategory = 1

	FieldValid              FieldCategory = 0
	FieldTypeAnonymous      FieldCategory = 1  // 已解决
	FieldTypeVoidInterface  FieldCategory = 2  // 已解决
	FieldTypeNonBangPackage FieldCategory = 4  // 可能为系统包或环境包
	FieldTypeStruct         FieldCategory = 8  // 未解决 匿名结构体做类型
	FieldTypeUnknown        FieldCategory = 16 // 未解决 函数做类型
	FieldTypeMultiSelector  FieldCategory = 32 // 未解决 a.b.c类型
)

type Interface struct {
	Name    string      `json:"name"`
	File    string      `json:"file"`
	Package string      `json:"package"`
	Funcs   []*Function `json:"functions"`
}

type Struct struct {
	Name           string         `json:"name"`
	File           string         `json:"file"`
	Package        string         `json:"package"`
	FieldList      []*Field       `json:"field_list"`
	StructCategory StructCategory `json:"struct_category"`
}

type Alias struct {
	Name    string      `json:"name"`
	Package string      `json:"package"`
	AliasOf interface{} `json:"alias_of"` // Function / string / Decl
}

type Field struct {
	FieldName     string        `json:"field_name"`
	FieldType     string        `json:"field_type"`
	FieldTag      string        `json:"field_tag"`
	FieldCategory FieldCategory `json:"field_category"`
}

func NewField(fieldNode *ast.Field, fset *token.FileSet) *Field {
	field := &Field{
		FieldName:     GetFieldName(fieldNode),
		FieldType:     ExprTokenToStr(fieldNode.Type, fset),
		FieldTag:      GetFieldTag(fieldNode),
		FieldCategory: 0,
	}
	return field
}

type TypeWalker struct {
	File       string
	fset       *token.FileSet
	Structs    []*Struct
	Interfaces []*Interface
	Aliases    []*Alias
}

func (w *TypeWalker) Visit(node ast.Node) ast.Visitor {
	switch currNode := node.(type) {
	case *ast.GenDecl:
		if currNode.Tok == token.TYPE {
			for _, spec := range currNode.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					w.CollectStruct(typeSpec)
					w.CollectInterface(typeSpec)
					w.CollectAlias(typeSpec)
				}
			}
			return nil
		}
		return nil
	}
	return w
}

func (w *TypeWalker) CollectStruct(typeSpec *ast.TypeSpec) {
	if structNode, ok := typeSpec.Type.(*ast.StructType); ok {
		_struct := &Struct{
			Name:    typeSpec.Name.Name,
			File:    w.File,
			Package: filepath.Dir(w.File),
		}

		for _, fieldNode := range structNode.Fields.List {
			field := NewField(fieldNode, w.fset)
			_struct.FieldList = append(_struct.FieldList, field)
		}

		w.Structs = append(w.Structs, _struct)
	}
}

/*
分类:
	- FuncType    : type RunBatchTaskFunc func(tx *gorp.Transaction, task *batch.BatchTask) error
	- nil         : type MyInt int
	- SelectorExpr: type tx gorp.Transaction
	- Ident       : type tx Interface
*/
func (w *TypeWalker) CollectAlias(typeSpec *ast.TypeSpec) {
	if _, ok := typeSpec.Type.(*ast.StructType); ok {
		return
	}
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return
	}

	//if selectorExpr, ok := typeSpec.Type.(*ast.SelectorExpr); ok {
	//	name := selectorExpr.X.(*ast.Ident).Name
	//	sel := selectorExpr.Sel.Name
	//	fmt.Println(name, sel)
	//}
	//
	//if funcNode, ok := typeSpec.Type.(*ast.FuncType); ok {
	//	fmt.Println(funcNode)
	//}

	alias := &Alias{
		Name:    typeSpec.Name.Name,
		Package: filepath.Dir(w.File),
		AliasOf: ExprTokenToStr(typeSpec.Type, w.fset),
	}
	w.Aliases = append(w.Aliases, alias)
}

func (w *TypeWalker) CollectInterface(typeSpec *ast.TypeSpec) {
	if interfaceNode, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		fmt.Println(interfaceNode)
	}
}

func ParseType(filePath string) (*TypeWalker, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	v := &TypeWalker{File: filePath, fset: fset}
	ast.Walk(v, f)
	return v, nil
}

func GetFieldName(fieldNode *ast.Field) (fieldName string) {
	if len(fieldNode.Names) > 0 {
		fieldName = fieldNode.Names[0].Name
	}
	return
}

func GetFieldTag(fieldNode *ast.Field) (fieldTag string) {
	if fieldNode.Tag != nil {
		fieldTag, _ = strconv.Unquote(fieldNode.Tag.Value)
	}
	return
}
