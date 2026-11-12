package domain

import (
	"testing"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToRedactedStruct(t *testing.T) {
	env, err := myasthurts.NewEnvironment()
	require.NoError(t, err)
	pkg, err := env.Parse("./testdata/sample")
	require.NoError(t, err)
	require.Len(t, pkg.Structs, 1)

	got := ToRedactedStruct(pkg.Structs[0])

	type field struct {
		exportedName string
		isRedacted   bool
		skip         bool
		redactor     string
	}
	fields := make(map[string]field, len(got.Fields))
	for _, f := range got.Fields {
		fields[f.Name] = field{f.ExportedName, f.IsRedacted, f.Skip, f.Redactor}
	}
	assert.Equal(t, "Sample", got.Name)
	assert.Equal(t, map[string]field{
		"Default":  {exportedName: "Default", isRedacted: true},
		"Named":    {exportedName: "named", isRedacted: true},
		"Allowed":  {exportedName: "Allowed"},
		"Redactor": {exportedName: "Redactor", isRedacted: true, redactor: "pan64"},
		"All":      {exportedName: "all", redactor: "star"},
		"Skipped":  {exportedName: "Skipped", isRedacted: true, skip: true},
		"JSON":     {exportedName: "json_name", isRedacted: true},
		"Slice":    {exportedName: "Slice", isRedacted: true, redactor: "len"},
	}, fields)
}

func TestToRedactedStruct_types(t *testing.T) {
	env, err := myasthurts.NewEnvironment()
	require.NoError(t, err)
	pkg, err := env.Parse("./testdata/types")
	require.NoError(t, err)
	require.Len(t, pkg.Structs, 1)

	types := make(map[string]string)
	for _, f := range ToRedactedStruct(pkg.Structs[0]).Fields {
		types[f.Name] = f.Type
	}
	assert.Equal(t, map[string]string{
		"Value":            "builtin.string",
		"Pointer":          "*builtin.string",
		"Slice":            "[]builtin.string",
		"Time":             "*time.Time",
		"SliceOfPointers":  "[]*builtin.string",
		"PointerToSlice":   "*[]builtin.string",
		"SliceOfSlices":    "[][]builtin.string",
		"PointerToPointer": "**builtin.string",
		"Channel":          "chan builtin.string",
	}, types)
}
