//go:generate go run github.com/valyala/quicktemplate/qtc -dir=.

package templates

import (
	"github.com/valyala/quicktemplate"

	"github.com/jamillosantos/zapredactor/internal/domain"
)

type fieldRendererFnc func(writer *quicktemplate.Writer, field domain.RedactedField)

var (
	renderFieldMap = map[string]fieldRendererFnc{
		"builtin.string":       streamRenderValueType("String"),
		"*builtin.string":      streamRenderValueType("String"),
		"[]builtin.string":     streamRenderArrayType("StringArray"),
		"builtin.int":          streamRenderValueType("Int"),
		"[]builtin.int":        streamRenderArrayType("Ints"),
		"*builtin.int":         streamRenderValueType("Int"),
		"builtin.int8":         streamRenderValueType("Int8"),
		"[]builtin.int8":       streamRenderArrayType("Int8s"),
		"*builtin.int8":        streamRenderValueType("Int8"),
		"builtin.int16":        streamRenderValueType("Int16"),
		"[]builtin.int16":      streamRenderArrayType("Int16s"),
		"*builtin.int16":       streamRenderValueType("Int16"),
		"builtin.int32":        streamRenderValueType("Int32"),
		"[]builtin.int32":      streamRenderArrayType("Int32s"),
		"*builtin.int32":       streamRenderValueType("Int32"),
		"builtin.int64":        streamRenderValueType("Int64"),
		"[]builtin.int64":      streamRenderArrayType("Int64s"),
		"*builtin.int64":       streamRenderValueType("Int64"),
		"builtin.uint":         streamRenderValueType("Uint"),
		"[]builtin.uint":       streamRenderArrayType("Uints"),
		"*builtin.uint":        streamRenderValueType("Uint"),
		"builtin.uint8":        streamRenderValueType("Uint8"),
		"[]builtin.uint8":      streamRenderArrayType("Uint8s"),
		"*builtin.uint8":       streamRenderValueType("Uint8"),
		"builtin.uint16":       streamRenderValueType("Uint16"),
		"[]builtin.uint16":     streamRenderArrayType("Uint16s"),
		"*builtin.uint16":      streamRenderValueType("Uint16"),
		"builtin.uint32":       streamRenderValueType("Uint32"),
		"[]builtin.uint32":     streamRenderArrayType("Uint32s"),
		"*builtin.uint32":      streamRenderValueType("Uint32"),
		"builtin.uint64":       streamRenderValueType("Uint64"),
		"[]builtin.uint64":     streamRenderArrayType("Uint64s"),
		"*builtin.uint64":      streamRenderValueType("Uint64"),
		"builtin.float32":      streamRenderValueType("Float32"),
		"[]builtin.float32":    streamRenderArrayType("Float32s"),
		"*builtin.float32":     streamRenderValueType("Float32"),
		"builtin.float64":      streamRenderValueType("Float64"),
		"[]builtin.float64":    streamRenderArrayType("Float64s"),
		"*builtin.float64":     streamRenderValueType("Float64"),
		"builtin.complex64":    streamRenderValueType("Complex64"),
		"[]builtin.complex64":  streamRenderArrayType("Complex64s"),
		"*builtin.complex64":   streamRenderValueType("Complex64"),
		"builtin.complex128":   streamRenderValueType("Complex128"),
		"[]builtin.complex128": streamRenderArrayType("Complex128s"),
		"*builtin.complex128":  streamRenderValueType("Complex128"),
		"builtin.bool":         streamRenderValueType("Bool"),
		"[]builtin.bool":       streamRenderArrayType("Bools"),
		"*builtin.bool":        streamRenderValueType("Bool"),
		"time.Time":            streamRenderValueType("Time"),
		"[]time.Time":          streamRenderArrayType("Times"),
		"*time.Time":           streamRenderValueType("Time"),
		"time.Duration":        streamRenderValueType("Duration"),
		"[]time.Duration":      streamRenderArrayType("Durations"),
		"*time.Duration":       streamRenderValueType("Duration"),
		"[]builtin.byte":       streamRenderValueType("Binary"),
	}
)

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
	fr, ok := renderFieldMap[field.Type]
	if !ok {
		panic("unsupported field type: " + field.Type)
	}
	fr(writer, field)
}

func streamRenderValueType(t string) fieldRendererFnc {
	return func(writer *quicktemplate.Writer, field domain.RedactedField) {
		StreamRenderValueType(writer, field, t)
	}
}

func streamRenderArrayType(t string) fieldRendererFnc {
	return func(writer *quicktemplate.Writer, field domain.RedactedField) {
		StreamRenderArray(writer, field, t)
	}
}
