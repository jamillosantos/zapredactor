package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor"
	"github.com/jamillosantos/zapredactor/redactors"
	"github.com/jamillosantos/zapredactor/redactreflection"
)

var benchCard = CreditCard{
	PAN:            "1234123412344321",
	CardHolder:     "Snake Eyes Joe",
	CVV:            "123",
	Expiry:         Expiry{Month: 2, Year: 29},
	First4:         "1234",
	Last4:          "4321",
	InternalStatus: "pending",
}

// benchReflectionOptions makes the reflection redactor log the same object as the redact tags of CreditCard.
var benchReflectionOptions = []redactreflection.ReflectionOption{
	redactreflection.WithRedactor("pan", redactors.PAN64),
	redactreflection.WithRedactor("CVV", redactors.Star),
	redactreflection.WithAllowFields("Expiry.Year", "First4", "Last4"),
}

// handwrittenCard writes the same object as benchCard by hand, with pointer receivers and no redactor calls. What it
// costs is what zap itself costs, so it is the floor the other variants are compared to.
type handwrittenCard CreditCard

func (s *handwrittenCard) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("pan", "123412******4321")
	encoder.AddString("card_holder", "[redacted]")
	encoder.AddString("CVV", "***")
	_ = encoder.AddObject("Expiry", (*handwrittenExpiry)(&s.Expiry))
	encoder.AddString("First4", s.First4)
	encoder.AddString("Last4", s.Last4)
	return nil
}

type handwrittenExpiry Expiry

func (s *handwrittenExpiry) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("Month", "[redacted]")
	encoder.AddInt("Year", s.Year)
	return nil
}

type benchVariant struct {
	name string
	// maxAllocs is the limit of allocations for one log call. It includes the one zap makes for its fields.
	maxAllocs float64
	field     func() zap.Field
}

// benchVariants lists the ways of logging benchCard. All of them log the same object.
func benchVariants() []benchVariant {
	card := benchCard
	handwritten := handwrittenCard(benchCard)
	reflection := redactreflection.Redactor(benchReflectionOptions...)
	return []benchVariant{
		{"handwritten", 1, func() zap.Field { return zap.Object("card", &handwritten) }},
		{"codegen_zap.Object", 7, func() zap.Field { return zap.Object("card", &card) }},
		{"codegen_Redact", 8, func() zap.Field { return zapredactor.Redact("card", &card) }},
		{"codegen_Redact_by_value", 9, func() zap.Field { return zapredactor.Redact("card", card) }},
		{"reflection_Redactor", 23, func() zap.Field { return reflection("card", &card) }},
		{"reflection_Redact", 25, func() zap.Field {
			return redactreflection.Redact("card", &card, benchReflectionOptions...)
		}},
	}
}

func benchLogger() *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(io.Discard), zapcore.DebugLevel))
}

func BenchmarkCard(b *testing.B) {
	logger := benchLogger()
	for _, v := range benchVariants() {
		b.Run(v.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				logger.Info("message", v.field())
			}
		})
	}
}

// TestAllocations fails when a way of logging the card starts to allocate more. Lower a limit when it allocates less.
func TestAllocations(t *testing.T) {
	logger := benchLogger()
	for _, v := range benchVariants() {
		t.Run(v.name, func(t *testing.T) {
			allocs := testing.AllocsPerRun(100, func() { logger.Info("message", v.field()) })
			assert.LessOrEqual(t, allocs, v.maxAllocs)
		})
	}
}

// The benchmark only compares like with like while every variant logs the same object.
func TestBenchVariants_logTheSameObject(t *testing.T) {
	var want string
	for i, v := range benchVariants() {
		buf, err := zapcore.NewJSONEncoder(zapcore.EncoderConfig{}).EncodeEntry(zapcore.Entry{}, []zap.Field{v.field()})
		require.NoError(t, err)
		if i == 0 {
			want = buf.String()
			continue
		}
		assert.Equal(t, want, buf.String(), v.name)
	}
}
