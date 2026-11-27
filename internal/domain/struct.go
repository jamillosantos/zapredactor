package domain

import (
	"slices"

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
	OmitEmpty    bool
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
			OmitEmpty:    getOmitEmpty(field.Tag),
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

// getRedactor returns the first part after the name that is not one of the allow and omitempty flags.
func getRedactor(tag myasthurts.Tag) string {
	t := tag.TagParamByName("redact")
	if t == nil {
		return ""
	}
	for _, option := range t.Options {
		if option != "" && option != "allow" && option != "omitempty" {
			return option
		}
	}
	return ""
}

func getSkip(tag myasthurts.Tag) bool {
	t := tag.TagParamByName("redact")
	return t == nil || t.Value == "-"
}

// getOmitEmpty reports whether one of the parts after the name is the literal "omitempty".
func getOmitEmpty(tag myasthurts.Tag) bool {
	t := tag.TagParamByName("redact")
	return t != nil && slices.Contains(t.Options, "omitempty")
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

// getPackageName returns the field type with every pointer, slice and channel level in it, so that a type the
// templates do not know, such as "[]*builtin.string", is reported as unsupported instead of rendered as "[]string".
func getPackageName(refType myasthurts.RefType) string {
	switch t := refType.(type) {
	case *myasthurts.ArrayRefType:
		// shortcut: the parser drops the length, so [4]T is seen as []T; needs the parser to keep it.
		return "[]" + getPackageName(t.RefType)
	case *myasthurts.StarRefType:
		return "*" + getPackageName(t.RefType)
	case *myasthurts.ChanRefType:
		return "chan " + getPackageName(t.RefType)
	}
	return refType.Pkg().Name + "." + refType.Name()
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
