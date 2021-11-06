package zapredactor

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	redactedString = "[REDACTED]"
	nilString      = "[NIL]"
)

type TagRedactor struct {
	val interface{}
}

func (r TagRedactor) Redact(encoder zapcore.ObjectEncoder) error {
	val := reflect.ValueOf(r.val)
	t := reflect.TypeOf(r.val)

	if val.Kind() == reflect.Invalid || (val.Kind() == reflect.Ptr && val.IsNil()) {
		encoder.AddString("", nilString)
		return nil
	}

	for t.Kind() == reflect.Ptr {
		val = val.Elem()
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		logRedacted(val, t, "", encoder)
		return nil
	}
	logObject(val, t, "", encoder)
	return nil
}

func fieldPrefix(prefix string, field reflect.StructField) string {
	name := field.Name
	if tag, ok := field.Tag.Lookup("json"); ok {
		name = strings.Split(tag, ",")[0]
	} else if tag, ok := field.Tag.Lookup("yaml"); ok {
		name = strings.Split(tag, ",")[0]
	}
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

func logObject(vt reflect.Value, t reflect.Type, prefix string, encoder zapcore.ObjectEncoder) {
	fieldsNum := vt.NumField()
	for i := 0; i < fieldsNum; i++ {
		field := vt.Field(i)
		fieldType := t.Field(i)

		name := fieldPrefix(prefix, fieldType)

		v, ok := fieldType.Tag.Lookup("redact")
		if v == "-" {
			continue
		}
		if !ok || v != "allow" {
			logRedacted(field, fieldType.Type, name, encoder)
			continue
		}

		logValue(field, fieldType.Type, name, encoder)
	}
}

func logRedacted(val reflect.Value, t reflect.Type, name string, encoder zapcore.ObjectEncoder) {
	if val.Kind() == reflect.Invalid || (val.Kind() == reflect.Ptr && val.IsNil()) {
		encoder.AddString(name, redactedString)
		return
	}
	switch t.Kind() {
	case reflect.Ptr:
		logRedacted(val.Elem(), val.Elem().Type(), name, encoder)
	case reflect.Array:
		encoder.AddString(name, fmt.Sprintf("[REDACTED!len=%d]", val.Len()))
	case reflect.Slice:
		encoder.AddString(name, fmt.Sprintf("[REDACTED!len=%d]", val.Len()))
	case reflect.Map:
		encoder.AddString(name, fmt.Sprintf("[REDACTED!len=%d]", val.Len()))
	default:
		encoder.AddString(name, redactedString)
	}
}

func logValue(val reflect.Value, t reflect.Type, name string, encoder zapcore.ObjectEncoder) {
	if val.Kind() == reflect.Invalid || (val.Kind() == reflect.Ptr && val.IsNil()) {
		encoder.AddString(name, "[NIL]")
		return
	}

	switch t.Kind() {
	case reflect.Ptr:
		logValue(val.Elem(), val.Elem().Type(), name, encoder)
	case reflect.Bool:
		encoder.AddBool(name, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		encoder.AddInt64(name, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		encoder.AddUint64(name, val.Uint())
	case reflect.Uintptr:
		encoder.AddString(name, "[unsupported:uintptr]")
	case reflect.Float32, reflect.Float64:
		encoder.AddFloat64(name, val.Float())
	case reflect.Complex64, reflect.Complex128:
		encoder.AddComplex128(name, val.Complex())
	case reflect.Array:
		encoder.AddString(name, "[unsupported:array]")
	case reflect.Map:
		encoder.AddString(name, "[unsupported:map]")
	case reflect.Slice:
		encoder.AddString(name, "[unsupported:slice]")
	case reflect.Interface:
		logValue(val.Elem(), t.Elem(), name, encoder)
	case reflect.Struct:
		logObject(val, t, name, encoder)
	case reflect.String:
		encoder.AddString(name, val.String())
	}
}
