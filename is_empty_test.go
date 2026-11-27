package zapredactor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	text := ""
	for name, tt := range map[string]struct {
		value any
		want  bool
	}{
		"nil":              {nil, true},
		"nil pointer":      {(*string)(nil), true},
		"pointer to empty": {&text, false},
		"empty string":     {"", true},
		"string":           {"a", false},
		"zero":             {0, true},
		"number":           {1, false},
		"false":            {false, true},
		"empty slice":      {[]string{}, true},
		"slice":            {[]string{""}, false},
		"empty map":        {map[string]int{}, true},
		"zero struct":      {time.Time{}, false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsEmpty(tt.value))
		})
	}
}

// The generated code calls IsEmpty on the hot path, where it must not add allocations.
func TestIsEmpty_doesNotAllocate(t *testing.T) {
	text, number := "value", 42
	assert.Zero(t, testing.AllocsPerRun(100, func() {
		IsEmpty(text)
		IsEmpty(number)
		IsEmpty(&text)
	}))
}
