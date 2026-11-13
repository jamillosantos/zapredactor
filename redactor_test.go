package zapredactor

import (
	"fmt"
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

func TestRedactorManager_zeroValue(t *testing.T) {
	var rm RedactorManager
	assert.Equal(t, redactors.DefaultRedactedString, rm.RedactValue("secret", "custom"))

	rm.RegisterRedactor("custom", func(interface{}) string { return "[custom]" })

	assert.Equal(t, "[custom]", rm.RedactValue("secret", "custom"))
	assert.Equal(t, redactors.DefaultRedactedString, rm.RedactValue("secret", "unknown"))
	assert.Equal(t, redactors.DefaultRedactedString, RedactValue("secret", "custom"), "a manager must not register on the default one")
}

func TestAddRedacted(t *testing.T) {
	pan := "1234567890124321"
	redacted := func(add func(encoder zapcore.ObjectEncoder)) interface{} {
		encoder := zapcore.NewMapObjectEncoder()
		add(encoder)
		return encoder.Fields["key"]
	}

	assert.Equal(t, "123456******4321", redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", pan, redactors.PAN64) }))
	assert.Equal(t, "123456******4321", redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", &pan, redactors.PAN64) }))
	assert.Equal(t, redactors.DefaultRedactedString, redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", 42, redactors.Default) }))
	assert.Equal(t, redactors.DefaultRedactedString, redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", pan, "unknown") }))

	t.Run("a registered redactor replaces the built-in one", func(t *testing.T) {
		t.Cleanup(func() {
			defaultRedactorManager.redactors[redactors.Len] = redactors.LenRedactor
			delete(registeredRedactors, redactors.Len)
		})
		RegisterRedactor(redactors.Len, func(interface{}) string { return "[custom]" })

		assert.Equal(t, "[custom]", redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", pan, redactors.Len) }))
	})

	t.Run("a replaced DefaultRedactor still receives the value", func(t *testing.T) {
		original := redactors.DefaultRedactor
		t.Cleanup(func() { redactors.DefaultRedactor = original })
		redactors.DefaultRedactor = func(data interface{}) string { return fmt.Sprint("seen:", data) }

		assert.Equal(t, "seen:42", redacted(func(e zapcore.ObjectEncoder) { AddRedacted(e, "key", 42, redactors.Default) }))
	})
}
