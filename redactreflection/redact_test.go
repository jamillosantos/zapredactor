package redactreflection

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRedactor(t *testing.T) {
	c, ob := observer.New(zapcore.DebugLevel)
	logger := zap.New(c, zap.ErrorOutput(zapcore.AddSync(io.Discard)))

	redactor := Redactor()

	data := defaultTestingValue
	logger.Info("test", redactor("data", data))

	entries := ob.All()
	assert.Len(t, entries, 1)
	assert.Len(t, entries[0].Context, 1)
}

func TestRedactor_keepsOptionsAndValuesApart(t *testing.T) {
	redactor := Redactor(WithAllowFields("Name"))

	fields := map[string]zap.Field{
		"first":  redactor("data", struct{ Name string }{"first"}),
		"second": redactor("data", struct{ Name string }{"second"}),
	}
	for want, field := range fields {
		enc := zapcore.NewMapObjectEncoder()
		field.AddTo(enc)
		assert.Equal(t, map[string]interface{}{"Name": want}, enc.Fields["data"])
	}
}
