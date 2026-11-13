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

// builtinDefaultRedactor identifies the DefaultRedactor of this package, to tell when it was replaced.
var builtinDefaultRedactor = reflect.ValueOf(DefaultRedactor).Pointer()

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
		if p == nil {
			return ""
		}
		pan = *p
	default:
		return notCompatible
	}
	if len(pan) < 6+4+2 {
		return DefaultRedactor(data)
	}
	return pan[:6] + stars(len(pan)-6-4) + pan[len(pan)-4:]
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
		if p == nil {
			return ""
		}
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
		if p == nil {
			return ""
		}
		l = len(*p)
	default:
		return notCompatible
	}
	return stars(l)
}

// stars returns n asterisks. It only allocates when n is above the length of the constant.
func stars(n int) string {
	const s = "****************************************************************"
	if n <= len(s) {
		return s[:n]
	}
	return strings.Repeat("*", n)
}

// AppendString appends s, redacted by one of the built-in redactors, to dst. The text is the same the redactor
// function returns, and nothing is allocated. It appends nothing and reports false when it cannot do that: the
// redactor is not a built-in one, or DefaultRedactor was replaced.
func AppendString(dst []byte, s string, redactor Redactor) ([]byte, bool) {
	if reflect.ValueOf(DefaultRedactor).Pointer() != builtinDefaultRedactor {
		return dst, false
	}
	switch redactor {
	case Default:
	case PAN64:
		if len(s) >= 6+4+2 {
			dst = append(dst, s[:6]...)
			dst = append(dst, stars(len(s)-6-4)...)
			return append(dst, s[len(s)-4:]...), true
		}
	case BIN:
		if len(s) > 6 {
			return append(dst, s[:6]...), true
		}
	case Star, Asterisk:
		return append(dst, stars(len(s))...), true
	case Len:
		dst = append(dst, "[len:"...)
		dst = strconv.AppendInt(dst, int64(len(s)), 10)
		return append(dst, ']'), true
	default:
		return dst, false
	}
	return append(dst, DefaultRedactedString...), true
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
		if p != nil {
			l = len(*p)
		}
	default:
		val := reflect.ValueOf(p)
		if val.Kind() != reflect.Slice {
			return notCompatible
		}
		l = val.Len()
	}
	return "[len:" + strconv.Itoa(l) + "]"
}
