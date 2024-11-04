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

type fieldVisitor func(encoder zapcore.ObjectEncoder, fieldPath, fieldName string, value interface{})

func (r redactByReflection) Redact(encoder zapcore.ObjectEncoder) error {
	walkFields(encoder, r.val, "", r.fieldNameExtractor, func(encoder zapcore.ObjectEncoder, fieldPath, fieldName string, value interface{}) {
		info, _ := r.info[fieldPath]
		if info.hidden {
			return
		}
		if info.allow {
			addField(encoder, fieldName, value)
			return
		}
		encoder.AddString(fieldName, zapredactor.RedactValue(value, info.redactor))
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
	case *time.Time:
		encoder.AddTime(field, *v)
	case time.Duration:
		encoder.AddDuration(field, v)
	case *time.Duration:
		encoder.AddDuration(field, *v)
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
func walkFields(encoder zapcore.ObjectEncoder, obj interface{}, prefix string, fieldNameFormatter func(p reflect.StructField) string, f fieldVisitor) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		runStruct(encoder, val, prefix, fieldNameFormatter, f)
	case reflect.Map:
		runMap(encoder, val, prefix, fieldNameFormatter, f)
	default:
	}
}

func runMap(encoder zapcore.ObjectEncoder, val reflect.Value, prefix string, fieldNameFormatter func(p reflect.StructField) string, f fieldVisitor) {
	for _, key := range val.MapKeys() {
		fieldName := key.String()
		fieldPath := fieldName
		if prefix != "" {
			fieldPath = prefix + "." + fieldName
		}
		fieldValue := val.MapIndex(key).Interface()
		switch fv := fieldValue.(type) {
		case time.Time:
			f(encoder, fieldPath, fieldName, fv)
			continue
		}
		mapVal := reflect.ValueOf(fieldValue)
		switch {
		case mapVal.Kind() == reflect.Struct:
			_ = encoder.AddObject(fieldName, zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
				walkFields(encoder, mapVal.Interface(), fieldPath, fieldNameFormatter, f)
				return nil
			}))
		case mapVal.Kind() == reflect.Map:
			_ = encoder.AddObject(fieldName, zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
				runMap(encoder, mapVal, fieldPath, fieldNameFormatter, f)
				return nil
			}))
		default:
			f(encoder, fieldPath, fieldName, fieldValue)
		}
	}
}

func runStruct(encoder zapcore.ObjectEncoder, val reflect.Value, prefix string, fieldNameFormatter func(p reflect.StructField) string, f fieldVisitor) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldPath := prefix
		if prefix != "" {
			fieldPath += "."
		}
		fieldNameKey := fieldNameFormatter(field)
		fieldPath += fieldNameKey
		switch field.Type.Kind() {
		case reflect.Struct:
			_ = encoder.AddObject(fieldNameKey, zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
				runStruct(encoder, val.Field(i), fieldPath, fieldNameFormatter, f)
				return nil
			}))
			continue
		case reflect.Map:
			_ = encoder.AddObject(fieldNameKey, zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
				runMap(encoder, val.Field(i), fieldPath, fieldNameFormatter, f)
				return nil
			}))
			continue
		default:
			f(encoder, fieldPath, fieldNameKey, val.Field(i).Interface())
		}
	}
}
