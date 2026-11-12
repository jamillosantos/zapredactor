package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// The generated code must redact through zap's own field constructors, not only through zapredactor.Redact.
func TestGeneratedCode_redactsThroughZap(t *testing.T) {
	card := CreditCard{PAN: "1234123412344321", CVV: "123"}
	for name, field := range map[string]zap.Field{
		"zap.Object":  zap.Object("card", card),
		"zap.Any":     zap.Any("card", card),
		"zap.Objects": zap.Objects("card", []CreditCard{card}),
	} {
		t.Run(name, func(t *testing.T) {
			var logData bytes.Buffer
			logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel))

			logger.Info("message", field)

			assert.Contains(t, logData.String(), `"pan":"123412******4321"`)
			assert.NotContains(t, logData.String(), card.PAN)
		})
	}
}
