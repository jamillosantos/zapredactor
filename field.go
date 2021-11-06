package zapredactor

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Redactable interface {
	Redact(encoder zapcore.ObjectEncoder) error
}

type redactableProxy struct {
	Redactable
}

func (p *redactableProxy) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	return p.Redact(encoder)
}

// Redact adds a field to a zap log entry which its content is redacted.
// If the object implements the Redactable, its Redact method will be used to add the fields;
// If the object implements a zapcore.ObjectMarshaler it will be used to log the structure;
// Otherwise, the TagRedactor will wrap the given value redacting all fields by default, unless the ones explicitly
// marked as allowed.
func Redact(key string, val interface{}) zap.Field {
	switch v := val.(type) {
	case Redactable:
		return zap.Object(key, &redactableProxy{v})
	case zapcore.ObjectMarshaler:
		return zap.Object(key, v)
	default:
		return zap.Object(key, &redactableProxy{TagRedactor{val}})
	}
}
