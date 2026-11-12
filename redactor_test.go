package zapredactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor/redactors"
)

type testRedactable struct{}

func (testRedactable) Redact(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("field", "value")
	return nil
}

func TestRedact_nil(t *testing.T) {
	for name, val := range map[string]Redactable{
		"nil pointer":   (*testRedactable)(nil),
		"nil interface": nil,
	} {
		t.Run(name, func(t *testing.T) {
			enc := zapcore.NewMapObjectEncoder()
			Redact("key", val).AddTo(enc)
			assert.Equal(t, map[string]interface{}{"key": nil}, enc.Fields)
		})
	}
}

func TestRedactValue(t *testing.T) {
	RegisterRedactor("test_custom", func(interface{}) string { return "[custom]" })

	assert.Equal(t, "[custom]", RedactValue("secret", "test_custom"))
	assert.Equal(t, "123456******4321", RedactValue("1234567890124321", redactors.PAN64))
	assert.Equal(t, redactors.DefaultRedactedString, RedactValue("secret", redactors.Default))
	assert.Equal(t, redactors.DefaultRedactedString, RedactValue("secret", "unknown"), "an unknown redactor must not leak the value")
}
