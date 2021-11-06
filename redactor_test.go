//go:generate go run github.com/golang/mock/mockgen -package zapredactor -destination object_encoder_mock_test.go go.uber.org/zap/zapcore ObjectEncoder
package zapredactor

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type examplePerson struct {
	Name          string
	Age           int `json:"age"`
	MaritalStatus int `yaml:"marital_status"`
}

var (
	testStruct = examplePerson{}
)

func Test_fieldPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		field  reflect.StructField
		want   string
	}{
		{"should get Name", "", reflect.TypeOf(testStruct).Field(0), "Name"},
		{"should get age", "", reflect.TypeOf(testStruct).Field(1), "age"},
		{"should get marital_status", "", reflect.TypeOf(testStruct).Field(2), "marital_status"},
		{"should get 123.Name", "123", reflect.TypeOf(testStruct).Field(0), "123.Name"},
		{"should get 123.age", "123", reflect.TypeOf(testStruct).Field(1), "123.age"},
		{"should get 123.marital_status", "123", reflect.TypeOf(testStruct).Field(2), "123.marital_status"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldPrefix(tt.prefix, tt.field)
			assert.Equal(t, tt.want, got)
		})
	}
}

func buildOE(t *testing.T) *MockObjectEncoder {
	ctrl := gomock.NewController(t)
	return NewMockObjectEncoder(ctrl)
}

func mockAddStringRedacted(name string) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddString(name, "[REDACTED]")
	}
}

func mockAddString(name, value string) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddString(name, value)
	}
}

func mockAddBool(name string, value bool) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddBool(name, value)
	}
}

func mockAddInt(name string, value int64) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddInt64(name, value)
	}
}

func mockAddUint(name string, value uint64) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddUint64(name, value)
	}
}

func mockAddUintptr(name string, value uintptr) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddUintptr(name, value)
	}
}

func mockAddFloat(name string, value float64) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddFloat64(name, value)
	}
}

func mockAddComplex(name string, value complex128) func(encoder *MockObjectEncoder) {
	return func(encoder *MockObjectEncoder) {
		encoder.EXPECT().AddComplex128(name, value)
	}
}

func Test_logRedacted(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		prefix    string
		assertion func(encoder *MockObjectEncoder)
	}{
		{"given nil", nil, "", mockAddStringRedacted("")},
		{"given string nil", (*string)(nil), "", mockAddStringRedacted("")},
		{"given string", "", "", mockAddStringRedacted("")},
		{"given string with name", "", "test.test", mockAddStringRedacted("test.test")},
		// Summaries
		{"given *array", &[1]string{""}, "", mockAddString("", "[REDACTED!len=1]")},
		{"given array", [3]string{}, "", mockAddString("", "[REDACTED!len=3]")},
		{"given slice", []string{""}, "", mockAddString("", "[REDACTED!len=1]")},
		{"given map", map[string]string{}, "", mockAddString("", "[REDACTED!len=0]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := buildOE(t)
			tt.assertion(oe)
			logRedacted(reflect.ValueOf(tt.val), reflect.TypeOf(tt.val), tt.prefix, oe)
		})
	}
}

func strToPtr(s string) *string {
	return &s
}

func Test_logValue(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		prefix    string
		assertion func(encoder *MockObjectEncoder)
	}{
		{"given nil", nil, "", mockAddString("", "[NIL]")},
		{"given string nil", (*string)(nil), "", mockAddString("", "[NIL]")},
		{"given string pointer", strToPtr("string"), "", mockAddString("", "string")},
		{"given string", "string", "", mockAddString("", "string")},
		{"given string with name", "string", "test.test", mockAddString("test.test", "string")},
		{"given bool", true, "name", mockAddBool("name", true)},
		{"given int", 123, "name", mockAddInt("name", 123)},
		{"given uint", uint(123), "name", mockAddUint("name", 123)},
		{"given float", float64(123), "name", mockAddFloat("name", 123)},
		{"given complex", complex128(123), "name", mockAddComplex("name", 123)},
		{"given uintptr", uintptr(123), "name", mockAddString("name", "[unsupported:uintptr]")},
		{"given array", [3]string{}, "name", mockAddString("name", "[unsupported:array]")},
		{"given slice", []string{""}, "name", mockAddString("name", "[unsupported:slice]")},
		{"given map", map[string]string{}, "name", mockAddString("name", "[unsupported:map]")},
		{"given struct", logObjectValue, "name", mockLogObjectValue("name.")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := buildOE(t)
			tt.assertion(oe)
			logValue(reflect.ValueOf(tt.val), reflect.TypeOf(tt.val), tt.prefix, oe)
		})
	}
}

type logObjectStruct struct {
	Name       string
	Age        int    `redact:"allow" json:"age"`
	Country    string `redact:"allow"`
	Irrelevant string `redact:"-"`
}

var (
	logObjectValue = logObjectStruct{
		Name:       "n1",
		Age:        123,
		Country:    "c1",
		Irrelevant: "i1",
	}
)

func Test_logObject(t *testing.T) {
	oe := buildOE(t)
	mockLogObjectValue("")(oe)
	logObject(reflect.ValueOf(logObjectValue), reflect.TypeOf(logObjectValue), "", oe)
}

func mockLogObjectValue(prefix string) func(oe *MockObjectEncoder) {
	return func(oe *MockObjectEncoder) {
		oe.EXPECT().AddString(prefix+"Name", redactedString)
		oe.EXPECT().AddInt64(prefix+"age", int64(123))
		oe.EXPECT().AddString(prefix+"Country", "c1")
	}
}

func TestTagRedactor_Redact(t *testing.T) {
	t.Run("when the given information is not a struct", func(t *testing.T) {
		oe := buildOE(t)
		tr := buildTagRedactor("not a struct")

		oe.EXPECT().AddString("", redactedString)

		err := tr.Redact(oe)
		assert.NoError(t, err)
	})

	t.Run("when the given information is a pointer to a non struct", func(t *testing.T) {
		oe := buildOE(t)
		str := "not a struct"
		tr := buildTagRedactor(&str)

		oe.EXPECT().AddString("", redactedString)

		err := tr.Redact(oe)
		assert.NoError(t, err)
	})

	t.Run("when the given information is a struct", func(t *testing.T) {
		oe := buildOE(t)
		tr := buildTagRedactor(logObjectValue)

		mockLogObjectValue("")(oe)

		err := tr.Redact(oe)
		assert.NoError(t, err)
	})

	t.Run("when the given information is nil", func(t *testing.T) {
		oe := buildOE(t)
		tr := buildTagRedactor(nil)

		oe.EXPECT().AddString("", "[NIL]")

		err := tr.Redact(oe)
		assert.NoError(t, err)
	})
}

func buildTagRedactor(s interface{}) *TagRedactor {
	return &TagRedactor{s}
}
