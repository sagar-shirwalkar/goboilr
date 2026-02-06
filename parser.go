package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
	"unicode"
)

func ParseFile(filePath string) (*FileData, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	data := &FileData{
		PackageName: node.Name.Name,
		Imports:     node.Imports,
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// Check doc comments for generation triggers
		hasConstructorTag := false
		hasBuilderTag := false

		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				text := comment.Text
				if strings.Contains(text, "gen:new") {
					hasConstructorTag = true
				}
				if strings.Contains(text, "gen:builder") {
					hasBuilderTag = true
				}
			}
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			sData := StructData{
				StructName:          typeSpec.Name.Name,
				GenerateConstructor: hasConstructorTag,
				GenerateBuilder:     hasBuilderTag,
			}

			for _, field := range structType.Fields.List {
				typeBuf := &bytes.Buffer{}
				printer.Fprint(typeBuf, fset, field.Type)
				typeName := typeBuf.String()

				var genVal string
				if field.Tag != nil {
					tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
					genVal = tag.Get("gen")
				}

				if len(field.Names) > 0 {
					for _, name := range field.Names {
						info := createFieldInfo(name.Name, typeName, genVal)
						sData.AllFields = append(sData.AllFields, info)
						if genVal != "" {
							sData.Accessors = append(sData.Accessors, info)
						}
					}
				} else {
					fieldName := extractEmbeddedName(typeName)
					info := createFieldInfo(fieldName, typeName, genVal)
					sData.AllFields = append(sData.AllFields, info)
				}
			}

			// Add struct if we have fields (and triggers usually imply fields exist)
			if len(sData.AllFields) > 0 {
				data.Structs = append(data.Structs, sData)
			}
		}
	}

	return data, nil
}

func createFieldInfo(name, typeName, genVal string) FieldInfo {
	info := FieldInfo{
		Name:    name,
		ArgName: lowerFirst(name),
		Type:    typeName,
	}
	if genVal != "" {
		info.MethodName = capitalize(name)
		info.HasGetter = strings.Contains(genVal, "get")
		info.HasSetter = strings.Contains(genVal, "set")
		info.HasValidator = strings.Contains(genVal, "val")
	}
	return info
}

func extractEmbeddedName(typeName string) string {
	s := strings.TrimLeft(typeName, "*")
	if idx := strings.LastIndex(s, "."); idx != -1 {
		return s[idx+1:]
	}
	return s
}

func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	r := []rune(str)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func lowerFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	r := []rune(str)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
