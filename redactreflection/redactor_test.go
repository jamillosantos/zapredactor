package redactreflection

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor/redactors"
)

type testWalkFieldsStructSub struct {
	SubField1 string
	SubField2 int
	SubField3 string `json:"sub_field_3"`
}

type testWalkFieldsStructSub2 struct {
	SubField1 string `json:"subField1"`
}

type testWalkFieldsStruct struct {
	Name        string
	Sub         testWalkFieldsStructSub
	Sub2        testWalkFieldsStructSub2 `json:"sub_2"`
	Map         map[string]any           `redact:"this_map"`
	hiddenField string
}

var defaultTestingValue = testWalkFieldsStruct{
	Name: "name_value",
	Sub: testWalkFieldsStructSub{
		SubField1: "sub_field1_value",
		SubField2: 2222,
		SubField3: "sub_field_3_value",
	},
	Sub2: testWalkFieldsStructSub2{
		SubField1: "sub2_field1_value",
	},
	Map: map[string]any{
		"map_key1": "map_value1",
		"map_key2": "map_value2",
	},

	hiddenField: "hidden_value",
}

func Test_redactByReflection_Redact(t *testing.T) {
	t.Run("should redact all fields", func(t *testing.T) {
		var logData bytes.Buffer
		c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
		logger := zap.New(c)

		data := defaultTestingValue

		logger.Info("message", Redact("data", data))
		_ = logger.Sync()
		logString := logData.String()
		assert.NotContains(t, logString, `name_value`)
		assert.NotContains(t, logString, `sub_field1_value`)
		assert.NotContains(t, logString, `2222`)
		assert.NotContains(t, logString, `sub_field_3_value`)
		assert.NotContains(t, logString, `sub2_field1_value`)
		assert.NotContains(t, logString, `map_value1`)
		assert.NotContains(t, logString, `map_value2`)
		assert.Contains(t, logString, `"Name":"[redacted]"`)
		assert.Contains(t, logString, `"SubField1":"[redacted]"`)
		assert.Contains(t, logString, `"SubField2":"[redacted]"`)
		assert.Contains(t, logString, `"sub_field_3":"[redacted]"`)
		assert.Contains(t, logString, `"subField1":"[redacted]"`)
		assert.Contains(t, logString, `"map_key1":"[redacted]"`)
		assert.Contains(t, logString, `"map_key2":"[redacted]"`)
	})

	t.Run("when the input is a map", func(t *testing.T) {
		t.Run("should redact all fields", func(t *testing.T) {
			var logData bytes.Buffer
			c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
			logger := zap.New(c)

			logger.Info("message", Redact("data", map[string]any{
				"name": "name_value",
				"sub": map[string]any{
					"sub_field1":  "sub_field1_value",
					"sub_field2":  2222,
					"sub_field_3": "sub_field_3_value",
				},
				"sub2": testWalkFieldsStructSub2{
					SubField1: "sub2_field1_value",
				},
			}))
			_ = logger.Sync()
			logString := logData.String()
			assert.NotContains(t, logString, `name_value`)
			assert.NotContains(t, logString, `sub_field1_value`)
			assert.NotContains(t, logString, `2222`)
			assert.NotContains(t, logString, `sub_field_3_value`)
			assert.NotContains(t, logString, `sub2_field1_value`)
			assert.Contains(t, logString, `"name":"[redacted]"`)
			assert.Contains(t, logString, `"sub_field1":"[redacted]"`)
			assert.Contains(t, logString, `"sub_field2":"[redacted]"`)
			assert.Contains(t, logString, `"sub_field_3":"[redacted]"`)
			assert.Contains(t, logString, `"subField1":"[redacted]"`)
		})
	})

	t.Run("should redact all fields that are not allowed", func(t *testing.T) {
		var logData bytes.Buffer
		c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
		logger := zap.New(c)

		data := defaultTestingValue

		logger.Info("message", Redact("data", data,
			WithAllowFields("Name", "Sub.SubField2", "sub_2.subField1", "this_map.map_key1"),
		))
		_ = logger.Sync()
		logString := logData.String()
		assert.NotContains(t, logString, `sub_field1_value`)
		assert.NotContains(t, logString, `sub_field_3_value`)
		assert.Contains(t, logString, `"Name":"name_value"`)
		assert.Contains(t, logString, `"SubField1":"[redacted]"`)
		assert.Contains(t, logString, `"SubField2":2222`)
		assert.Contains(t, logString, `"sub_field_3":"[redacted]"`)
		assert.Contains(t, logString, `"subField1":"sub2_field1_value"`)
		assert.Contains(t, logString, `"map_key1":"map_value1"`)
	})

	t.Run("should redact fields with a custom redactor", func(t *testing.T) {
		var logData bytes.Buffer
		c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
		logger := zap.New(c)

		data := defaultTestingValue
		data.Sub.SubField1 = "1234567890124321" // fake PAN number

		logger.Info("message", Redact("data", data,
			WithRedactor("Sub.SubField1", redactors.PAN64),
		))

		_ = logger.Sync()
		logString := logData.String()

		assert.Contains(t, logString, `"SubField1":"123456******4321"`)
	})

	t.Run("should hide a field", func(t *testing.T) {
		var logData bytes.Buffer
		c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
		logger := zap.New(c)

		data := defaultTestingValue

		logger.Info("message", Redact("data", data,
			WithHiddenField("Sub.SubField1", "sub_2.subField1"),
		))

		_ = logger.Sync()
		logString := logData.String()
		assert.NotContains(t, logString, `"SubField1"`)
		assert.NotContains(t, logString, `"subField1"`)
	})

	t.Run("should render all fields types", func(t *testing.T) {
		var logData bytes.Buffer
		c := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&logData), zapcore.DebugLevel)
		logger := zap.New(c)

		date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		logger.Info("message", Redact("data",
			map[string]any{
				"string":     "string_value",
				"bool":       true,
				"int":        123,
				"int8":       int8(123),
				"int16":      int16(123),
				"int32":      int32(123),
				"int64":      int64(123),
				"uint":       uint(123),
				"uint8":      uint8(123),
				"uint16":     uint16(123),
				"uint32":     uint32(123),
				"uint64":     uint64(123),
				"uintptr":    uintptr(123),
				"float32":    float32(123.123),
				"float64":    123.123,
				"complex64":  complex64(123.123),
				"complex128": complex128(123.123),
				"bytearray":  []byte("bytearray_value"),
				"time":       date,
				"*time":      &date,
				"duration":   time.Second * 10,
				"error":      errors.New("error"),
			},
			WithAllowFields(
				"string",
				"bool",
				"int",
				"int8",
				"int16",
				"int32",
				"int64",
				"uint",
				"uint8",
				"uint16",
				"uint32",
				"uint64",
				"uintptr",
				"float32",
				"float64",
				"complex64",
				"complex128",
				"bytearray",
				"time",
				"*time",
				"duration",
				"error",
			),
		))

		_ = logger.Sync()
		logString := logData.String()
		assert.Contains(t, logString, `"string":"string_value"`)
		assert.Contains(t, logString, `"bool":true`)
		assert.Contains(t, logString, `"int":123`)
		assert.Contains(t, logString, `"int8":123`)
		assert.Contains(t, logString, `"int16":123`)
		assert.Contains(t, logString, `"int32":123`)
		assert.Contains(t, logString, `"int64":123`)
		assert.Contains(t, logString, `"uint":123`)
		assert.Contains(t, logString, `"uint8":123`)
		assert.Contains(t, logString, `"uint16":123`)
		assert.Contains(t, logString, `"uint32":123`)
		assert.Contains(t, logString, `"uint64":123`)
		assert.Contains(t, logString, `"uintptr":123`)
		assert.Contains(t, logString, `"float32":123.123`)
		assert.Contains(t, logString, `"float64":123.123`)
		assert.Contains(t, logString, `"complex64":"123.123+0i"`)
		assert.Contains(t, logString, `"complex128":"123.123+0i"`)
		assert.Contains(t, logString, `"bytearray":"Ynl0ZWFycmF5X3ZhbHVl"`)
		assert.Contains(t, logString, `"time":1609459200`)
		assert.Contains(t, logString, `"*time":1609459200`)
		assert.Contains(t, logString, `"duration":10`)
		assert.Contains(t, logString, `"error":"error"`)
	})
}
