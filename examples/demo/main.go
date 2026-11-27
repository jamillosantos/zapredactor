//go:generate go run ../../cli/zapredactor/main.go --destination redactor_gen.go
package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/jamillosantos/zapredactor"
)

type Demo struct {
	SkipThisField              string          `redact:"-"`
	String                     string          `redact:",allow"`
	Strings                    []string        `redact:",allow"`
	Stringp                    *string         `redact:",allow"`
	Int                        int             `redact:",allow"`
	Ints                       []int           `redact:",allow"`
	Intp                       *int            `redact:",allow"`
	Int8                       int8            `redact:",allow"`
	Int8s                      []int8          `redact:",allow"`
	Int8p                      *int8           `redact:",allow"`
	Int16                      int16           `redact:",allow"`
	Int16s                     []int16         `redact:",allow"`
	Int16p                     *int16          `redact:",allow"`
	Int32                      int32           `redact:",allow"`
	Int32s                     []int32         `redact:",allow"`
	Int32p                     *int32          `redact:",allow"`
	Int64                      int64           `redact:",allow"`
	Int64s                     []int64         `redact:",allow"`
	Int64p                     *int64          `redact:",allow"`
	Uint                       uint            `redact:",allow"`
	Uints                      []uint          `redact:",allow"`
	Uintp                      *uint           `redact:",allow"`
	Uint8                      uint8           `redact:",allow"`
	Uint8s                     []uint8         `redact:",allow"`
	Uint8p                     *uint8          `redact:",allow"`
	Uint16                     uint16          `redact:",allow"`
	Uint16s                    []uint16        `redact:",allow"`
	Uint16p                    *uint16         `redact:",allow"`
	Uint32                     uint32          `redact:",allow"`
	Uint32s                    []uint32        `redact:",allow"`
	Uint32p                    *uint32         `redact:",allow"`
	Uint64                     uint64          `redact:",allow"`
	Uint64s                    []uint64        `redact:",allow"`
	Uint64p                    *uint64         `redact:",allow"`
	Float32                    float32         `redact:",allow"`
	Float32s                   []float32       `redact:",allow"`
	Float32p                   *float32        `redact:",allow"`
	Float64                    float64         `redact:",allow"`
	Float64s                   []float64       `redact:",allow"`
	Float64p                   *float64        `redact:",allow"`
	Complex64                  complex64       `redact:",allow"`
	Complex64s                 []complex64     `redact:",allow"`
	Complex64p                 *complex64      `redact:",allow"`
	Complex128                 complex128      `redact:",allow"`
	Complex128s                []complex128    `redact:",allow"`
	Complex128p                *complex128     `redact:",allow"`
	Bool                       bool            `redact:",allow"`
	Boolp                      *bool           `redact:",allow"`
	Time                       time.Time       `redact:",allow"`
	Times                      []time.Time     `redact:",allow"`
	Timep                      *time.Time      `redact:",allow"`
	Duration                   time.Duration   `redact:",allow"`
	Durations                  []time.Duration `redact:",allow"`
	Durationp                  *time.Duration  `redact:",allow"`
	Byte                       []byte          `redact:",allow"`
	NamedStringp               *string         `json:"named_stringp" redact:",allow"`
	RedactedString             string          `redact:""`
	RedactedStringWithCardData string          `redact:",,pan64"`
	RedactedStrings            []string        `redact:""`
	RedactedStringp            *string         `redact:""`
	RedactedInt                int             `redact:""`
	RedactedInts               []int           `redact:""`
	RedactedIntp               *int            `redact:""`
	RedactedInt8               int8            `redact:""`
	RedactedInt8s              []int8          `redact:""`
	RedactedInt8p              *int8           `redact:""`
	RedactedInt16              int16           `redact:""`
	RedactedInt16s             []int16         `redact:""`
	RedactedInt16p             *int16          `redact:""`
	RedactedInt32              int32           `redact:""`
	RedactedInt32s             []int32         `redact:""`
	RedactedInt32p             *int32          `redact:""`
	RedactedInt64              int64           `redact:""`
	RedactedInt64s             []int64         `redact:""`
	RedactedInt64p             *int64          `redact:""`
	RedactedUint               uint            `redact:""`
	RedactedUints              []uint          `redact:""`
	RedactedUintp              *uint           `redact:""`
	RedactedUint8              uint8           `redact:""`
	RedactedUint8s             []uint8         `redact:""`
	RedactedUint8p             *uint8          `redact:""`
	RedactedUint16             uint16          `redact:""`
	RedactedUint16s            []uint16        `redact:""`
	RedactedUint16p            *uint16         `redact:""`
	RedactedUint32             uint32          `redact:""`
	RedactedUint32s            []uint32        `redact:""`
	RedactedUint32p            *uint32         `redact:""`
	RedactedUint64             uint64          `redact:""`
	RedactedUint64s            []uint64        `redact:""`
	RedactedUint64p            *uint64         `redact:""`
	RedactedFloat32            float32         `redact:""`
	RedactedFloat32s           []float32       `redact:""`
	RedactedFloat32p           *float32        `redact:""`
	RedactedFloat64            float64         `redact:""`
	RedactedFloat64s           []float64       `redact:""`
	RedactedFloat64p           *float64        `redact:""`
	RedactedComplex64          complex64       `redact:""`
	RedactedComplex64s         []complex64     `redact:""`
	RedactedComplex64p         *complex64      `redact:""`
	RedactedComplex128         complex128      `redact:""`
	RedactedComplex128s        []complex128    `redact:""`
	RedactedComplex128p        *complex128     `redact:""`
	RedactedBool               bool            `redact:""`
	RedactedBoolp              *bool           `redact:""`
	RedactedTime               time.Time       `redact:""`
	RedactedTimes              []time.Time     `redact:""`
	RedactedTimep              *time.Time      `redact:""`
	RedactedDuration           time.Duration   `redact:""`
	RedactedDurations          []time.Duration `redact:""`
	RedactedDurationp          *time.Duration  `redact:""`
	RedactedByte               []byte          `redact:""`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	demo := Demo{
		String:   "string",
		Int:      1,
		Int8:     2,
		Int16:    3,
		Int32:    4,
		Int64:    5,
		Uint:     6,
		Uint8:    7,
		Uint16:   8,
		Uint32:   9,
		Uint64:   10,
		Float32:  11.0,
		Float64:  12.0,
		Bool:     true,
		Time:     time.UnixMilli(13),
		Duration: 14,
	}

	logger.Info("demo entry", zapredactor.Redact("demo", &demo))
}
