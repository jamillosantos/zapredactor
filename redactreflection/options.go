package redactreflection

import (
	"reflect"

	"github.com/jamillosantos/zapredactor/redactors"
)

type ReflectionOption = func(*redactByReflection)

// WithAllowFields allows the fields to be logged without redaction.
// - If the field is already defined, it will be overwritten.
func WithAllowFields(fields ...string) ReflectionOption {
	return func(r *redactByReflection) {
		for _, field := range fields {
			r.info[field] = redactInfo{allow: true}
		}
	}
}

// WithRedactor defines a redactor for a given field.
// - If the field is already defined, it will be overwritten.
// - If the field is already allowed, it will be overwritten with allowed = false.
func WithRedactor(field string, redactor redactors.Redactor) ReflectionOption {
	return func(r *redactByReflection) {
		r.info[field] = redactInfo{redactor: redactor}
	}
}

// FieldNameExtractor defines a function to extract the field name from a struct field.
// - Calling this will override the previous value set.
func FieldNameExtractor(extractor func(s reflect.StructField) string) ReflectionOption {
	return func(r *redactByReflection) {
		r.fieldNameExtractor = extractor
	}
}

// WithHiddenField defines a field that will be hidden.
// - If the field is already defined, it will be overwritten.
// - If the field is already allowed, it will be overwritten with allowed = false.
// - You can call this function multiple times to define multiple fields.
func WithHiddenField(fields ...string) ReflectionOption {
	return func(r *redactByReflection) {
		for _, field := range fields {
			r.info[field] = redactInfo{hidden: true}
		}
	}
}
