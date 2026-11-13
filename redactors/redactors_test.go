package redactors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cardNumber = "123456789012345"
)

func TestPAN64(t *testing.T) {
	tests := []struct {
		data interface{}
		want string
	}{
		{
			cardNumber,
			"123456*****2345",
		},
		{
			&cardNumber,
			"123456*****2345",
		},
		{
			"too short",
			DefaultRedactedString,
		},
		{
			12345,
			notCompatible,
		},
		{
			nil,
			"",
		},
		{
			(*string)(nil),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("should format %v to %s", tt.data, tt.want), func(t *testing.T) {
			got := PAN64Redactor(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBIN(t *testing.T) {
	tests := []struct {
		data interface{}
		want string
	}{
		{
			cardNumber,
			"123456",
		},
		{
			&cardNumber,
			"123456",
		},
		{
			"short",
			DefaultRedactedString,
		},
		{
			12345,
			notCompatible,
		},
		{
			nil,
			"",
		},
		{
			(*string)(nil),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("should format %v to %s", tt.data, tt.want), func(t *testing.T) {
			got := BINRedactor(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStar(t *testing.T) {
	tests := []struct {
		data interface{}
		want string
	}{
		{nil, ""},
		{(*string)(nil), ""},
		{12345, notCompatible},
		{"", ""},
		{"12345", "*****"},
		{&cardNumber, "***************"},
		{strings.Repeat("a", 100), strings.Repeat("*", 100)}, // Longer than the constant the stars are cut from.
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v to %v", tt.data, tt.want), func(t *testing.T) {
			got := StarRedactor(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		data interface{}
		want int
	}{
		{nil, 0},
		{make([]string, 3), 3},
		{make([]int, 3), 3},
		{make([]int8, 3), 3},
		{make([]int16, 3), 3},
		{make([]int32, 3), 3},
		{make([]int64, 3), 3},
		{make([]uint, 3), 3},
		{make([]uint8, 3), 3},
		{make([]uint16, 3), 3},
		{make([]uint32, 3), 3},
		{make([]uint64, 3), 3},
		{make([]float32, 3), 3},
		{make([]float64, 3), 3},
		{make([]bool, 3), 3},
		{make([]interface{}, 3), 3},
		{"123", 3},
		{&cardNumber, 15},
		{(*string)(nil), 0},
		{make([]struct{}, 3), 3},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.data), func(t *testing.T) {
			got := LenRedactor(tt.data)
			assert.Equal(t, fmt.Sprintf("[len:%d]", tt.want), got)
		})
	}

	t.Run("not compatible", func(t *testing.T) {
		got := LenRedactor(12345)
		assert.Equal(t, notCompatible, got)
	})
}

// AppendString is the fast path of the redactor functions, so it must give the same text as they do.
func TestAppendString(t *testing.T) {
	fncs := map[Redactor]func(interface{}) string{
		Default:  DefaultRedactor,
		PAN64:    PAN64Redactor,
		BIN:      BINRedactor,
		Star:     StarRedactor,
		Asterisk: StarRedactor,
		Len:      LenRedactor,
	}
	for redactor, fnc := range fncs {
		for _, s := range []string{"", "short", "1234567", cardNumber, strings.Repeat("a", 100)} {
			got, ok := AppendString([]byte("kept:"), s, redactor)
			assert.True(t, ok)
			assert.Equal(t, "kept:"+fnc(s), string(got), "redactor %q with %q", redactor, s)
		}
	}

	t.Run("not a built-in redactor", func(t *testing.T) {
		got, ok := AppendString(nil, cardNumber, "custom")
		assert.False(t, ok)
		assert.Empty(t, got)
	})

	t.Run("replaced DefaultRedactor", func(t *testing.T) {
		original := DefaultRedactor
		t.Cleanup(func() { DefaultRedactor = original })
		DefaultRedactor = func(interface{}) string { return "[replaced]" }

		got, ok := AppendString(nil, "short", PAN64)
		assert.False(t, ok)
		assert.Empty(t, got)
	})
}
