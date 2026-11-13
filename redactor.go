package zapredactor

import (
	"reflect"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor/redactors"
)

type RedactorFnc func(data interface{}) string

// RedactorManager is a manager of redactors. The zero value is ready to use and has no redactors: every value is
// redacted by the default redactor until RegisterRedactor is called.
type RedactorManager struct {
	redactors map[redactors.Redactor]RedactorFnc
}

var (
	defaultRedactorManager = &RedactorManager{}

	// registeredRedactors has the names given to RegisterRedactor. AddRedacted leaves those to RedactValue, so that
	// replacing a built-in redactor keeps working.
	registeredRedactors = map[redactors.Redactor]struct{}{}

	// redactedBuffers are the buffers AddRedacted writes the redacted text to.
	redactedBuffers = sync.Pool{New: func() any {
		b := make([]byte, 0, 64)
		return &b
	}}
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
	r, ok := rm.redactors[redactor]
	if !ok {
		r, ok = rm.redactors[redactors.Default]
	}
	if !ok {
		return redactors.DefaultRedactor(data)
	}
	return r(data)
}

// RegisterRedactor registers a new redactor. It is not safe for concurrent use: call it at startup.
func (rm *RedactorManager) RegisterRedactor(name redactors.Redactor, redactor RedactorFnc) {
	if rm.redactors == nil {
		rm.redactors = make(map[redactors.Redactor]RedactorFnc)
	}
	rm.redactors[name] = redactor
}

// RedactValue redacts a given data with the given redactor using the default RedactorManager.
func RedactValue(data interface{}, redactor redactors.Redactor) string {
	return defaultRedactorManager.RedactValue(data, redactor)
}

// AddRedacted adds the value to the encoder, redacted by the given redactor. The result is the same as
// encoder.AddString(key, RedactValue(value, redactor)), but the built-in redactors do not allocate.
func AddRedacted[T any](encoder zapcore.ObjectEncoder, key string, value T, redactor redactors.Redactor) {
	s, isString := any(value).(string)
	// The default redactor ignores the value, so it takes any type. The others only have a fast path for strings.
	if _, registered := registeredRedactors[redactor]; !registered && (isString || redactor == redactors.Default) {
		buf := redactedBuffers.Get().(*[]byte)
		b, ok := redactors.AppendString((*buf)[:0], s, redactor)
		if ok {
			// zap copies the bytes, so the buffer can go back to the pool.
			encoder.AddByteString(key, b)
		}
		// A buffer that grew for a long value is dropped, not to keep its memory forever.
		if cap(b) <= 1024 {
			*buf = b
			redactedBuffers.Put(buf)
		}
		if ok {
			return
		}
	}
	encoder.AddString(key, RedactValue(value, redactor))
}

// Redact constructs a field with the given key and value. A nil value is logged as null.
func Redact(key string, val Redactable) zap.Field {
	// A nil pointer would panic when its value-receiver Redact method is called.
	if v := reflect.ValueOf(val); !v.IsValid() || (v.Kind() == reflect.Pointer && v.IsNil()) {
		return zap.Reflect(key, nil)
	}
	return zap.Object(key, RedactObject(val))
}

// RedactObject constructs zapcore.ObjectMashaler from a Redactable.
func RedactObject(val Redactable) zapcore.ObjectMarshaler {
	// A pointer to a generated type brings its own marshaler, which saves the proxy and the copy of the struct.
	if m, ok := val.(interface {
		RedactObject() zapcore.ObjectMarshaler
	}); ok {
		return m.RedactObject()
	}
	return &redactableField{val}
}

// RegisterRedactor registers a new redactor on the default RedactorManager.
func RegisterRedactor(name redactors.Redactor, redactor RedactorFnc) {
	registeredRedactors[name] = struct{}{}
	defaultRedactorManager.RegisterRedactor(name, redactor)
}
