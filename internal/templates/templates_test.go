package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jamillosantos/zapredactor/internal/domain"
)

func TestUsesZap(t *testing.T) {
	tests := []struct {
		name  string
		field domain.RedactedField
		want  bool
	}{
		{"allowed slice", domain.RedactedField{Type: "[]builtin.string", IsArray: true}, true},
		{"redacted slice", domain.RedactedField{Type: "[]builtin.string", IsArray: true, IsRedacted: true}, false},
		{"skipped slice", domain.RedactedField{Type: "[]builtin.string", IsArray: true, Skip: true}, false},
		{"allowed []byte", domain.RedactedField{Type: "[]builtin.byte", IsArray: true}, false},
		{"allowed value", domain.RedactedField{Type: "builtin.string"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, UsesZap(tt.field))
		})
	}
}

func TestNotEmpty(t *testing.T) {
	for name, tt := range map[string]struct {
		field domain.RedactedField
		want  string
	}{
		"pointer":      {domain.RedactedField{Type: "*builtin.string", IsPointer: true}, "s.F != nil"},
		"slice":        {domain.RedactedField{Type: "[]builtin.int", IsArray: true}, "len(s.F) != 0"},
		"string":       {domain.RedactedField{Type: "builtin.string"}, `s.F != ""`},
		"bool":         {domain.RedactedField{Type: "builtin.bool"}, "s.F"},
		"number":       {domain.RedactedField{Type: "builtin.float64"}, "s.F != 0"},
		"duration":     {domain.RedactedField{Type: "time.Duration"}, "s.F != 0"},
		"time":         {domain.RedactedField{Type: "time.Time"}, ""},
		"struct value": {domain.RedactedField{Type: "main.Inner", IsStruct: true}, ""},
		"named type":   {domain.RedactedField{Type: "main.Status"}, "!zapredactor.IsEmpty(s.F)"},
	} {
		t.Run(name, func(t *testing.T) {
			tt.field.Name = "F"
			assert.Equal(t, tt.want, notEmpty(tt.field))
		})
	}
}

func TestRenderValueType_nilPointerUsesExportedName(t *testing.T) {
	got := RenderValueType(domain.RedactedField{Name: "Email", ExportedName: "email", IsPointer: true}, "String")
	assert.Contains(t, got, `encoder.AddReflected("email", nil)`)
}

func TestRenderValueRedacted_doesNotHTMLEscapeTheKey(t *testing.T) {
	got := RenderValueRedacted(domain.RedactedField{Name: "QA", ExportedName: "q&a"}, "")
	assert.Contains(t, got, `"q&a"`)
}
