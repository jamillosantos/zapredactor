package zapredactor

import "reflect"

// IsEmpty reports whether v is empty the way the omitempty of encoding/json sees it: false, 0, a nil pointer or
// interface, and an empty string, slice, array or map. A struct is never empty. The generated code calls it for the
// fields tagged with omitempty.
func IsEmpty(v any) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Struct:
		return false
	}
	return rv.IsZero()
}
