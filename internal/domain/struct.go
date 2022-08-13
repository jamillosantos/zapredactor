package domain

import (
	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
)

type RedactedField struct {
	Name         string
	ExportedName string
	Type         string
	IsRedacted   bool
	Skip         bool
	Redactor     string
	IsStruct     bool
	IsPointer    bool
	IsArray      bool
}

type RedactedStruct struct {
	Name   string
	Fields []RedactedField
}

func ToRedactedStruct(redactedStruct *myasthurts.Struct) RedactedStruct {
	redactedFields := make([]RedactedField, len(redactedStruct.Fields))
	for i, field := range redactedStruct.Fields {
		redactedFields[i] = RedactedField{
			Name:         field.Name,
			ExportedName: getExportedName(field),
			Type:         getPackageName(field.RefType),
			IsRedacted:   getIsRedacted(field.Tag),
			IsPointer:    getIsPointer(field.RefType),
			Skip:         getSkip(field.Tag),
			Redactor:     getRedactor(field.Tag),
			IsStruct:     getIsStruct(field.RefType),
			IsArray:      getIsArray(field.RefType),
		}
	}
	return RedactedStruct{
		Name:   redactedStruct.Name(),
		Fields: redactedFields,
	}
}

func getIsStruct(refType myasthurts.RefType) bool {
	switch refType.Type().(type) {
	case *myasthurts.Struct:
		return true
	}
	return false
}

func getRedactor(tag myasthurts.Tag) string {
	t := tag.TagParamByName("redact")
	if t == nil || len(t.Options) < 3 {
		return ""
	}
	return t.Options[2]
}

func getSkip(tag myasthurts.Tag) bool {
	t := tag.TagParamByName("redact")
	return t != nil && t.Value == "-"
}

func getIsRedacted(tag myasthurts.Tag) bool {
	t := tag.TagParamByName("redact")
	if t == nil {
		return true

	}
	for _, option := range t.Options {
		if option == "allow" {
			return false
		}
	}
	return true
}

func getPackageName(refType myasthurts.RefType) string {
	prefix := ""
	if getIsArray(refType) {
		prefix = "[]"
	}
	if getIsPointer(refType) {
		prefix += "*"
	}
	return prefix + refType.Pkg().Name + "." + refType.Name()
}

func getIsArray(refType myasthurts.RefType) bool {
	switch refType.(type) {
	case *myasthurts.ArrayRefType:
		return true
	}
	return false
}

func getIsPointer(refType myasthurts.RefType) bool {
	switch refType.(type) {
	case *myasthurts.StarRefType:
		return true
	}
	return false
}

func getExportedName(field *myasthurts.Field) string {
	if t := field.Tag.TagParamByName("redact"); t != nil && t.Value != "" && t.Value != "-" {
		return t.Value
	}
	if t := field.Tag.TagParamByName("json"); t != nil && t.Value != "-" && t.Value != "" {
		return t.Value
	}
	return field.Name
}
