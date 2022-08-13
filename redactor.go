package zapredactor

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor/redactors"
)

type RedactorFnc func(data interface{}) string

// RedactorManager is a manager of redactors.
type RedactorManager struct {
	redactors map[string]RedactorFnc
}

var (
	defaultRedactorManager = &RedactorManager{}
)

func init() {
	defaultRedactorManager.redactors = map[string]RedactorFnc{
		"": func(data interface{}) string {
			return redactors.DefaultRedactor(data)
		},
		"pan64": redactors.PAN64,
		"bin":   redactors.BIN,
		"len":   redactors.Len,
		"star":  redactors.Star,
		"*":     redactors.Star,
	}
}

// RedactValue redacts a given data with the given redactor. If the redactor is not found, the default redactor is used.
func (rm *RedactorManager) RedactValue(data interface{}, redactor string) string {
	r, ok := defaultRedactorManager.redactors[redactor]
	if !ok {
		r = rm.redactors[""]
	}
	return r(data)
}

// RegisterRedactor registers a new redactor.
func (rm *RedactorManager) RegisterRedactor(name string, redactor RedactorFnc) {
	rm.redactors[name] = redactor
}

// RedactValue redacts a given data with the given redactor using the default RedactorManager.
func RedactValue(data interface{}, redactor string) string {
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
func RegisterRedactor(name string, redactor RedactorFnc) {
	defaultRedactorManager.RegisterRedactor(name, redactor)
}
