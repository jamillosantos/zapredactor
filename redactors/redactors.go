package redactors

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	notCompatible = "[redactor not compatible]"
)

type Redactor string

const (
	Default  Redactor = ""
	PAN64    Redactor = "pan64"
	BIN      Redactor = "bin"
	Star     Redactor = "star"
	Asterisk Redactor = "*"
	Len      Redactor = "len"
)

var (
	DefaultRedactedString = "[redacted]"

	// DefaultRedactor is the default redactor used whenever a redactor is not defined or the given redactor is not found.
	// The default implementation returns the DefaultRedactedString always.
	DefaultRedactor = func(data interface{}) string {
		return DefaultRedactedString
	}
)

// PAN64Redactor is a redactor for PANs it ouputs the first 6 and the last 4. If the PAN is less than 12 digits the
// DefaultRedactor is used.
func PAN64Redactor(data interface{}) string {
	if data == nil {
		return ""
	}
	var pan string
	switch p := data.(type) {
	case string:
		pan = p
	case *string:
		pan = *p
	default:
		return notCompatible
	}
	if len(pan) < 6+4+2 {
		return DefaultRedactor(data)
	}
	return pan[:6] + strings.Repeat("*", len(pan)-6-4) + pan[len(pan)-4:]
}

// BINRedactor is a redactor for PANs. It outputs the first 6 digits. If the PAN is less than 6 digits the DefaultRedactor is used.
func BINRedactor(data interface{}) string {
	if data == nil {
		return ""
	}
	var pan string
	switch p := data.(type) {
	case string:
		pan = p
	case *string:
		pan = *p
	default:
		return notCompatible
	}
	if len(pan) <= 6 {
		return DefaultRedactor(data)
	}
	return pan[:6]
}

// StarRedactor is a redactor for any string or *string. It returns a string with the same length but masked with "*".
func StarRedactor(data interface{}) string {
	if data == nil {
		return ""
	}
	l := 0
	switch p := data.(type) {
	case string:
		l = len(p)
	case *string:
		l = len(*p)
	default:
		return notCompatible
	}
	return strings.Repeat("*", l)
}

// LenRedactor is a redactor for arrays, string and *string. It will output the number of elements formatted as "[len:X]". If the
// given data is not supported it will try using reflection. If the given data is not an array it will return not compatible.
func LenRedactor(data interface{}) string {
	var l int
	switch p := data.(type) {
	case nil:
		l = 0
	case []string:
		l = len(p)
	case []int:
		l = len(p)
	case []int8:
		l = len(p)
	case []int16:
		l = len(p)
	case []int32:
		l = len(p)
	case []int64:
		l = len(p)
	case []uint:
		l = len(p)
	case []uint8:
		l = len(p)
	case []uint16:
		l = len(p)
	case []uint32:
		l = len(p)
	case []uint64:
		l = len(p)
	case []float32:
		l = len(p)
	case []float64:
		l = len(p)
	case []bool:
		l = len(p)
	case []interface{}:
		l = len(p)
	case string:
		l = len(p)
	case *string:
		l = len(*p)
	default:
		val := reflect.ValueOf(p)
		if val.Kind() != reflect.Slice {
			return notCompatible
		}
		l = val.Len()
	}
	return "[len:" + strconv.Itoa(l) + "]"
}
