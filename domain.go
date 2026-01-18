package main

import "go/ast"

// FileData represents the result of parsing a file
type FileData struct {
	PackageName string
	Imports     []*ast.ImportSpec
	Structs     []StructData
}

// FieldInfo represents a single field in a struct
type FieldInfo struct {
	Name    string
	ArgName string
	Type    string

	// Accessor details
	MethodName                         string
	HasGetter, HasSetter, HasValidator bool
}

// StructData represents a struct and its metadata
type StructData struct {
	StructName          string
	GenerateConstructor bool
	AllFields           []FieldInfo
	Accessors           []FieldInfo
}
