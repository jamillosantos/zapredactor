package redactreflection

import (
	"reflect"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor"
	"github.com/jamillosantos/zapredactor/redactors"
)

type redactInfo struct {
	allow    bool
	hidden   bool
	redactor redactors.Redactor
}

type redactByReflection struct {
	val                interface{}
	info               map[string]redactInfo
	fieldNameExtractor func(s reflect.StructField) string
}

func (r redactByReflection) Redact(encoder zapcore.ObjectEncoder) error {
	walkFields(r.val, "", r.fieldNameExtractor, func(field string, value interface{}) {
		info, _ := r.info[field]
		if info.hidden {
			return
		}
		if info.allow {
			addField(encoder, field, value)
			return
		}
		encoder.AddString(field, zapredactor.RedactValue(value, info.redactor))
	})
	return nil
}

// addField adds a field to the encoder considering the type of the value.
func addField(encoder zapcore.ObjectEncoder, field string, value interface{}) {
	switch v := value.(type) {
	case string:
		encoder.AddString(field, v)
	case bool:
		encoder.AddBool(field, v)
	case int:
		encoder.AddInt(field, v)
	case int8:
		encoder.AddInt8(field, v)
	case int16:
		encoder.AddInt16(field, v)
	case int32:
		encoder.AddInt32(field, v)
	case int64:
		encoder.AddInt64(field, v)
	case uint:
		encoder.AddUint(field, v)
	case uint8:
		encoder.AddUint8(field, v)
	case uint16:
		encoder.AddUint16(field, v)
	case uint32:
		encoder.AddUint32(field, v)
	case uint64:
		encoder.AddUint64(field, v)
	case uintptr:
		encoder.AddUintptr(field, v)
	case float32:
		encoder.AddFloat32(field, v)
	case float64:
		encoder.AddFloat64(field, v)
	case complex64:
		encoder.AddComplex64(field, v)
	case complex128:
		encoder.AddComplex128(field, v)
	case []byte:
		encoder.AddBinary(field, v)
	case time.Time:
		encoder.AddTime(field, v)
	case time.Duration:
		encoder.AddDuration(field, v)
	case error:
		encoder.AddString(field, v.Error())
	case zapcore.ObjectMarshaler:
		_ = encoder.AddObject(field, v)
	case zapcore.ArrayMarshaler:
		_ = encoder.AddArray(field, v)
	default:
		_ = encoder.AddReflected(field, v)
	}
}

// walkFields walks the fields of the given object using reflection. It returns the list of fields.
// Whenever a field is a struct, it is recursively walked.
func walkFields(obj interface{}, prefix string, fieldNameFormatter func(p reflect.StructField) string, f func(field string, value interface{})) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	fields := make([]string, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldName := prefix
		if prefix != "" {
			fieldName += "."
		}
		fieldName += fieldNameFormatter(field)
		if field.Type.Kind() == reflect.Struct {
			walkFields(val.Field(i).Interface(), fieldName, fieldNameFormatter, f)
			continue
		}
		f(fieldName, val.Field(i).Interface())
		fields = append(fields, fieldName)
	}
}
