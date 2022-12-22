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
