package zapredactor

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor/redactors"
)

type RedactorFnc func(data interface{}) string

// RedactorManager is a manager of redactors.
type RedactorManager struct {
	redactors map[redactors.Redactor]RedactorFnc
}

var (
	defaultRedactorManager = &RedactorManager{}
)

func init() {
	defaultRedactorManager.redactors = map[redactors.Redactor]RedactorFnc{
		redactors.Default: func(data interface{}) string {
			return redactors.DefaultRedactor(data)
		},
		redactors.PAN64:    redactors.PAN64Redactor,
		redactors.BIN:      redactors.BINRedactor,
		redactors.Len:      redactors.LenRedactor,
		redactors.Star:     redactors.StarRedactor,
		redactors.Asterisk: redactors.StarRedactor,
	}
}

// RedactValue redacts a given data with the given redactor. If the redactor is not found, the default redactor is used.
func (rm *RedactorManager) RedactValue(data interface{}, redactor redactors.Redactor) string {
	r, ok := defaultRedactorManager.redactors[redactor]
	if !ok {
		r = rm.redactors[""]
	}
	return r(data)
}

// RegisterRedactor registers a new redactor.
func (rm *RedactorManager) RegisterRedactor(name redactors.Redactor, redactor RedactorFnc) {
	rm.redactors[name] = redactor
}

// RedactValue redacts a given data with the given redactor using the default RedactorManager.
func RedactValue(data interface{}, redactor redactors.Redactor) string {
	return defaultRedactorManager.RedactValue(data, redactor)
}

// Redact constructs a field with the given key and value.
func Redact(key string, val Redactable) zap.Field {
	return zap.Object(key, &redactableField{val})
}

// RedactObject constructs zapcore.ObjectMashaler from a Redactable.
func RedactObject(val Redactable) zapcore.ObjectMarshaler {
	return &redactableField{val}
}

// RegisterRedactor registers a new redactor on the default RedactorManager.
func RegisterRedactor(name redactors.Redactor, redactor RedactorFnc) {
	defaultRedactorManager.RegisterRedactor(name, redactor)
}
