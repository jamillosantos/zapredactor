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
