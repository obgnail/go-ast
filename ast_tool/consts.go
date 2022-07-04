package ast_tool

type ImportType int

const (
	dotImport ImportType = iota
	normalImport
	aliasImport
	underscoreImport
)
