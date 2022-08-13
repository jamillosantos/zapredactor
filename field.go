package zapredactor

import (
	"go.uber.org/zap/zapcore"
)

// Redactable is an interface for types that can be redacted.
type Redactable interface {
	Redact(encoder zapcore.ObjectEncoder) error
}

// redactableField is a proxy that transforms a Redactable into a ObjectMarshaler.
type redactableField struct {
	Redactable
}

// MarshalLogObject implements ObjectMarshaler proxying the call to the Redactable.
func (p *redactableField) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	return p.Redact(encoder)
}
