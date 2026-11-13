package redactreflection

import (
	"go.uber.org/zap"
)

// Redact redacts a given value using reflection.
func Redact(field string, val interface{}, opts ...ReflectionOption) zap.Field {
	return zap.Object(field, newRedactor(val, opts...))
}

func newRedactor(val interface{}, opts ...ReflectionOption) (r *redactByReflection) {
	r = &redactByReflection{
		val:                val,
		info:               make(map[string]redactInfo),
		fieldNameExtractor: DefaultNameExtractor,
	}
	for _, opt := range opts {
		opt(r)
	}
	return
}

// Redactor returns a factory of redactors based on the given configuration. The options are applied only once.
func Redactor(opts ...ReflectionOption) func(field string, val interface{}) zap.Field {
	base := newRedactor(nil, opts...)
	return func(field string, val interface{}) zap.Field {
		r := *base // The copy shares the info map, which is only read from here on.
		r.val = val
		return zap.Object(field, &r)
	}
}
