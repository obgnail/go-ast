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

func NewInterface(typeSpec *ast.TypeSpec, file string) *Struct {
	return nil
}

type Struct struct {
	ASTStruct *ast.TypeSpec
	Name      string   `json:"name"`
	File      string   `json:"file"`
	Package   string   `json:"package"`
	FieldList []*Field `json:"field_list"`
	//StructCategory StructCategory `json:"struct_category"`
}

func NewStruct(typeSpec *ast.TypeSpec, file string, fieldList []*Field) *Struct {
	_struct := &Struct{
		ASTStruct: typeSpec,
		Name:      typeSpec.Name.Name,
		File:      file,
		Package:   filepath.Dir(file),
		FieldList: fieldList,
	}
	return _struct
}

type Alias struct {
	ASTAlias *ast.TypeSpec
	Name     string      `json:"name"`
	Package  string      `json:"package"`
	AliasOf  interface{} `json:"alias_of"` // Function / string / Decl
}

func NewAlias(typeSpec *ast.TypeSpec, fset *token.FileSet, file string) *Alias {
	alias := &Alias{
		ASTAlias: typeSpec,
		Name:     typeSpec.Name.Name,
		Package:  filepath.Dir(file),
		AliasOf:  ExprTokenToStr(typeSpec.Type, fset),
	}
	return alias
}

type Field struct {
	ASTField  *ast.Field
	FieldName string `json:"field_name"`
	FieldType string `json:"field_type"`
	FieldTag  string `json:"field_tag"`
	//FieldCategory FieldCategory `json:"field_category"`
}

func NewField(fieldNode *ast.Field, fset *token.FileSet) *Field {
	field := &Field{
		ASTField:  fieldNode,
		FieldName: GetFieldName(fieldNode),
		FieldType: ExprTokenToStr(fieldNode.Type, fset),
		FieldTag:  GetFieldTag(fieldNode),
		//FieldCategory: 0,
	}
	return field
}

type TypeWalker struct {
	File       string
	Fset       *token.FileSet
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
		var fieldList []*Field
		for _, fieldNode := range structNode.Fields.List {
			field := NewField(fieldNode, w.Fset)
			fieldList = append(fieldList, field)
		}
		_struct := NewStruct(typeSpec, w.File, fieldList)
		w.Structs = append(w.Structs, _struct)
	}
}

/*
分类:
	- FuncType    : type RunBatchTaskFunc func(tx *gorp.Transaction, task *batch.BatchTask) error
	- nil         : type MyInt int 或者 type MyInt = int
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

	alias := NewAlias(typeSpec, w.Fset, filepath.Dir(w.File))
	w.Aliases = append(w.Aliases, alias)
}

func (w *TypeWalker) CollectInterface(typeSpec *ast.TypeSpec) {
	if interfaceNode, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		name := typeSpec.Name.Name

		for _, field := range interfaceNode.Methods.List {
			if function, ok := field.Type.(*ast.FuncType); ok {
				fmt.Println(function)
			} else if selector, ok := field.Type.(*ast.SelectorExpr); ok {
				fmt.Println(selector)
			} else {
				panic("error: not such type")
			}
		}
		fmt.Println(interfaceNode, name)
	}
}

func ParseType(filePath string) (*TypeWalker, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	v := &TypeWalker{File: filePath, Fset: fset}
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
