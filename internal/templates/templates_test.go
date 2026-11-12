package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jamillosantos/zapredactor/internal/domain"
)

func TestUsesZapArray(t *testing.T) {
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
			assert.Equal(t, tt.want, UsesZapArray(tt.field))
		})
	}
}

func TestRenderValueType_nilPointerUsesExportedName(t *testing.T) {
	got := RenderValueType(domain.RedactedField{Name: "Email", ExportedName: "email", IsPointer: true}, "String")
	assert.Contains(t, got, `encoder.AddReflected("email", nil)`)
}
