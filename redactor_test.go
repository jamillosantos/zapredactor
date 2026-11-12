package zapredactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
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
