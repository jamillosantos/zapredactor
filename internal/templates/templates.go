//go:generate go run github.com/valyala/quicktemplate/qtc -dir=.

package templates

import (
	"github.com/valyala/quicktemplate"

	"github.com/jamillosantos/zapredactor/internal/domain"
)

var (
	// valueTypes maps a field type to the zapcore.ObjectEncoder method (without the "Add" prefix) that logs it.
	valueTypes = map[string]string{
		"builtin.string":      "String",
		"*builtin.string":     "String",
		"builtin.int":         "Int",
		"*builtin.int":        "Int",
		"builtin.int8":        "Int8",
		"*builtin.int8":       "Int8",
		"builtin.int16":       "Int16",
		"*builtin.int16":      "Int16",
		"builtin.int32":       "Int32",
		"*builtin.int32":      "Int32",
		"builtin.int64":       "Int64",
		"*builtin.int64":      "Int64",
		"builtin.uint":        "Uint",
		"*builtin.uint":       "Uint",
		"builtin.uint8":       "Uint8",
		"*builtin.uint8":      "Uint8",
		"builtin.uint16":      "Uint16",
		"*builtin.uint16":     "Uint16",
		"builtin.uint32":      "Uint32",
		"*builtin.uint32":     "Uint32",
		"builtin.uint64":      "Uint64",
		"*builtin.uint64":     "Uint64",
		"builtin.float32":     "Float32",
		"*builtin.float32":    "Float32",
		"builtin.float64":     "Float64",
		"*builtin.float64":    "Float64",
		"builtin.complex64":   "Complex64",
		"*builtin.complex64":  "Complex64",
		"builtin.complex128":  "Complex128",
		"*builtin.complex128": "Complex128",
		"builtin.bool":        "Bool",
		"*builtin.bool":       "Bool",
		"time.Time":           "Time",
		"*time.Time":          "Time",
		"time.Duration":       "Duration",
		"*time.Duration":      "Duration",
		"[]builtin.byte":      "Binary",
	}

	// arrayTypes maps a field type to the zap field constructor that logs it.
	arrayTypes = map[string]string{
		"[]builtin.string":     "Strings",
		"[]builtin.int":        "Ints",
		"[]builtin.int8":       "Int8s",
		"[]builtin.int16":      "Int16s",
		"[]builtin.int32":      "Int32s",
		"[]builtin.int64":      "Int64s",
		"[]builtin.uint":       "Uints",
		"[]builtin.uint8":      "Uint8s",
		"[]builtin.uint16":     "Uint16s",
		"[]builtin.uint32":     "Uint32s",
		"[]builtin.uint64":     "Uint64s",
		"[]builtin.float32":    "Float32s",
		"[]builtin.float64":    "Float64s",
		"[]builtin.complex64":  "Complex64s",
		"[]builtin.complex128": "Complex128s",
		"[]builtin.bool":       "Bools",
		"[]time.Time":          "Times",
		"[]time.Duration":      "Durations",
	}
)

// isAllowedValue reports whether the field is logged with its real value, which is the only case rendered through
// valueTypes and arrayTypes.
func isAllowedValue(field domain.RedactedField) bool {
	return !field.Skip && !field.IsStruct && !field.IsRedacted
}

// UsesZap reports whether the code rendered for the field references the zap package.
func UsesZap(field domain.RedactedField) bool {
	_, ok := arrayTypes[field.Type]
	return ok && isAllowedValue(field)
}

func StreamRenderField(writer *quicktemplate.Writer, field domain.RedactedField) {
	switch {
	case field.Skip:
		return
	case field.IsStruct:
		StreamRenderObjectRedacted(writer, field)
		return
	case field.IsRedacted:
		StreamRenderValueRedacted(writer, field, field.Redactor)
		return
	}
	if t, ok := arrayTypes[field.Type]; ok {
		StreamRenderArray(writer, field, t)
		return
	}
	if t, ok := valueTypes[field.Type]; ok {
		StreamRenderValueType(writer, field, t)
		return
	}
	panic("unsupported field type: " + field.Type)
}
