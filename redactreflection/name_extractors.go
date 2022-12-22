package redactreflection

import (
	"reflect"
	"strings"
)

// DefaultNameExtractor receives a reflect.StructField and returns the name of the field. It prioritizes the json tag
// before the actual name of the field.
func DefaultNameExtractor(s reflect.StructField) string {
	jsonT := s.Tag.Get("json")
	if jsonT != "" {
		jsonT = strings.Split(jsonT, ",")[0]
	}
	if jsonT == "" {
		return s.Name
	}
	return jsonT
}
