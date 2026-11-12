package redactreflection

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNameExtractor_jsonDashFallsBackToFieldName(t *testing.T) {
	field := reflect.TypeOf(struct {
		Secret string `json:"-"`
	}{}).Field(0)
	assert.Equal(t, "Secret", DefaultNameExtractor(field))
}
