package redactors

import (
	"fmt"
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
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("should format %v to %s", tt.data, tt.want), func(t *testing.T) {
			got := PAN64(tt.data)
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
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("should format %v to %s", tt.data, tt.want), func(t *testing.T) {
			got := BIN(tt.data)
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
		{12345, notCompatible},
		{"", ""},
		{"12345", "*****"},
		{&cardNumber, "***************"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v to %v", tt.data, tt.want), func(t *testing.T) {
			got := Star(tt.data)
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
		{make([]struct{}, 3), 3},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.data), func(t *testing.T) {
			got := Len(tt.data)
			assert.Equal(t, fmt.Sprintf("[len:%d]", tt.want), got)
		})
	}

	t.Run("not compatible", func(t *testing.T) {
		got := Len(12345)
		assert.Equal(t, notCompatible, got)
	})
}
