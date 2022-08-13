// main is an auto generated file containing the Redact implementation for the annotated structs.
// DO NOT EDIT.
package main

import (
	"go.uber.org/zap/zapcore"

	"github.com/jamillosantos/zapredactor"
	"github.com/jamillosantos/zapredactor/zaparray"
)

func (s Demo) Redact(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("String", s.String)
	encoder.AddArray("Strings", zaparray.StringArray(s.Strings))
	if s.Stringp == nil {
		encoder.AddReflected("Stringp", nil)
	} else {
		encoder.AddString("Stringp", *s.Stringp)
	}
	encoder.AddInt("Int", s.Int)
	encoder.AddArray("Ints", zaparray.Ints(s.Ints))
	if s.Intp == nil {
		encoder.AddReflected("Intp", nil)
	} else {
		encoder.AddInt("Intp", *s.Intp)
	}
	encoder.AddInt8("Int8", s.Int8)
	encoder.AddArray("Int8s", zaparray.Int8s(s.Int8s))
	if s.Int8p == nil {
		encoder.AddReflected("Int8p", nil)
	} else {
		encoder.AddInt8("Int8p", *s.Int8p)
	}
	encoder.AddInt16("Int16", s.Int16)
	encoder.AddArray("Int16s", zaparray.Int16s(s.Int16s))
	if s.Int16p == nil {
		encoder.AddReflected("Int16p", nil)
	} else {
		encoder.AddInt16("Int16p", *s.Int16p)
	}
	encoder.AddInt32("Int32", s.Int32)
	encoder.AddArray("Int32s", zaparray.Int32s(s.Int32s))
	if s.Int32p == nil {
		encoder.AddReflected("Int32p", nil)
	} else {
		encoder.AddInt32("Int32p", *s.Int32p)
	}
	encoder.AddInt64("Int64", s.Int64)
	encoder.AddArray("Int64s", zaparray.Int64s(s.Int64s))
	if s.Int64p == nil {
		encoder.AddReflected("Int64p", nil)
	} else {
		encoder.AddInt64("Int64p", *s.Int64p)
	}
	encoder.AddUint("Uint", s.Uint)
	encoder.AddArray("Uints", zaparray.Uints(s.Uints))
	if s.Uintp == nil {
		encoder.AddReflected("Uintp", nil)
	} else {
		encoder.AddUint("Uintp", *s.Uintp)
	}
	encoder.AddUint8("Uint8", s.Uint8)
	encoder.AddArray("Uint8s", zaparray.Uint8s(s.Uint8s))
	if s.Uint8p == nil {
		encoder.AddReflected("Uint8p", nil)
	} else {
		encoder.AddUint8("Uint8p", *s.Uint8p)
	}
	encoder.AddUint16("Uint16", s.Uint16)
	encoder.AddArray("Uint16s", zaparray.Uint16s(s.Uint16s))
	if s.Uint16p == nil {
		encoder.AddReflected("Uint16p", nil)
	} else {
		encoder.AddUint16("Uint16p", *s.Uint16p)
	}
	encoder.AddUint32("Uint32", s.Uint32)
	encoder.AddArray("Uint32s", zaparray.Uint32s(s.Uint32s))
	if s.Uint32p == nil {
		encoder.AddReflected("Uint32p", nil)
	} else {
		encoder.AddUint32("Uint32p", *s.Uint32p)
	}
	encoder.AddUint64("Uint64", s.Uint64)
	encoder.AddArray("Uint64s", zaparray.Uint64s(s.Uint64s))
	if s.Uint64p == nil {
		encoder.AddReflected("Uint64p", nil)
	} else {
		encoder.AddUint64("Uint64p", *s.Uint64p)
	}
	encoder.AddFloat32("Float32", s.Float32)
	encoder.AddArray("Float32s", zaparray.Float32s(s.Float32s))
	if s.Float32p == nil {
		encoder.AddReflected("Float32p", nil)
	} else {
		encoder.AddFloat32("Float32p", *s.Float32p)
	}
	encoder.AddFloat64("Float64", s.Float64)
	encoder.AddArray("Float64s", zaparray.Float64s(s.Float64s))
	if s.Float64p == nil {
		encoder.AddReflected("Float64p", nil)
	} else {
		encoder.AddFloat64("Float64p", *s.Float64p)
	}
	encoder.AddComplex64("Complex64", s.Complex64)
	encoder.AddArray("Complex64s", zaparray.Complex64s(s.Complex64s))
	if s.Complex64p == nil {
		encoder.AddReflected("Complex64p", nil)
	} else {
		encoder.AddComplex64("Complex64p", *s.Complex64p)
	}
	encoder.AddComplex128("Complex128", s.Complex128)
	encoder.AddArray("Complex128s", zaparray.Complex128s(s.Complex128s))
	if s.Complex128p == nil {
		encoder.AddReflected("Complex128p", nil)
	} else {
		encoder.AddComplex128("Complex128p", *s.Complex128p)
	}
	encoder.AddBool("Bool", s.Bool)
	if s.Boolp == nil {
		encoder.AddReflected("Boolp", nil)
	} else {
		encoder.AddBool("Boolp", *s.Boolp)
	}
	encoder.AddTime("Time", s.Time)
	encoder.AddArray("Times", zaparray.Times(s.Times))
	if s.Timep == nil {
		encoder.AddReflected("Timep", nil)
	} else {
		encoder.AddTime("Timep", *s.Timep)
	}
	encoder.AddDuration("Duration", s.Duration)
	encoder.AddArray("Durations", zaparray.Durations(s.Durations))
	if s.Durationp == nil {
		encoder.AddReflected("Durationp", nil)
	} else {
		encoder.AddDuration("Durationp", *s.Durationp)
	}
	encoder.AddBinary("Byte", s.Byte)
	encoder.AddString("RedactedString", zapredactor.RedactValue(s.RedactedString, ""))
	encoder.AddString("RedactedStringWithCardData", zapredactor.RedactValue(s.RedactedStringWithCardData, ""))
	encoder.AddString("RedactedStrings", zapredactor.RedactValue(s.RedactedStrings, ""))
	encoder.AddString("RedactedStringp", zapredactor.RedactValue(s.RedactedStringp, ""))
	encoder.AddString("RedactedInt", zapredactor.RedactValue(s.RedactedInt, ""))
	encoder.AddString("RedactedInts", zapredactor.RedactValue(s.RedactedInts, ""))
	encoder.AddString("RedactedIntp", zapredactor.RedactValue(s.RedactedIntp, ""))
	encoder.AddString("RedactedInt8", zapredactor.RedactValue(s.RedactedInt8, ""))
	encoder.AddString("RedactedInt8s", zapredactor.RedactValue(s.RedactedInt8s, ""))
	encoder.AddString("RedactedInt8p", zapredactor.RedactValue(s.RedactedInt8p, ""))
	encoder.AddString("RedactedInt16", zapredactor.RedactValue(s.RedactedInt16, ""))
	encoder.AddString("RedactedInt16s", zapredactor.RedactValue(s.RedactedInt16s, ""))
	encoder.AddString("RedactedInt16p", zapredactor.RedactValue(s.RedactedInt16p, ""))
	encoder.AddString("RedactedInt32", zapredactor.RedactValue(s.RedactedInt32, ""))
	encoder.AddString("RedactedInt32s", zapredactor.RedactValue(s.RedactedInt32s, ""))
	encoder.AddString("RedactedInt32p", zapredactor.RedactValue(s.RedactedInt32p, ""))
	encoder.AddString("RedactedInt64", zapredactor.RedactValue(s.RedactedInt64, ""))
	encoder.AddString("RedactedInt64s", zapredactor.RedactValue(s.RedactedInt64s, ""))
	encoder.AddString("RedactedInt64p", zapredactor.RedactValue(s.RedactedInt64p, ""))
	encoder.AddString("RedactedUint", zapredactor.RedactValue(s.RedactedUint, ""))
	encoder.AddString("RedactedUints", zapredactor.RedactValue(s.RedactedUints, ""))
	encoder.AddString("RedactedUintp", zapredactor.RedactValue(s.RedactedUintp, ""))
	encoder.AddString("RedactedUint8", zapredactor.RedactValue(s.RedactedUint8, ""))
	encoder.AddString("RedactedUint8s", zapredactor.RedactValue(s.RedactedUint8s, ""))
	encoder.AddString("RedactedUint8p", zapredactor.RedactValue(s.RedactedUint8p, ""))
	encoder.AddString("RedactedUint16", zapredactor.RedactValue(s.RedactedUint16, ""))
	encoder.AddString("RedactedUint16s", zapredactor.RedactValue(s.RedactedUint16s, ""))
	encoder.AddString("RedactedUint16p", zapredactor.RedactValue(s.RedactedUint16p, ""))
	encoder.AddString("RedactedUint32", zapredactor.RedactValue(s.RedactedUint32, ""))
	encoder.AddString("RedactedUint32s", zapredactor.RedactValue(s.RedactedUint32s, ""))
	encoder.AddString("RedactedUint32p", zapredactor.RedactValue(s.RedactedUint32p, ""))
	encoder.AddString("RedactedUint64", zapredactor.RedactValue(s.RedactedUint64, ""))
	encoder.AddString("RedactedUint64s", zapredactor.RedactValue(s.RedactedUint64s, ""))
	encoder.AddString("RedactedUint64p", zapredactor.RedactValue(s.RedactedUint64p, ""))
	encoder.AddString("RedactedFloat32", zapredactor.RedactValue(s.RedactedFloat32, ""))
	encoder.AddString("RedactedFloat32s", zapredactor.RedactValue(s.RedactedFloat32s, ""))
	encoder.AddString("RedactedFloat32p", zapredactor.RedactValue(s.RedactedFloat32p, ""))
	encoder.AddString("RedactedFloat64", zapredactor.RedactValue(s.RedactedFloat64, ""))
	encoder.AddString("RedactedFloat64s", zapredactor.RedactValue(s.RedactedFloat64s, ""))
	encoder.AddString("RedactedFloat64p", zapredactor.RedactValue(s.RedactedFloat64p, ""))
	encoder.AddString("RedactedComplex64", zapredactor.RedactValue(s.RedactedComplex64, ""))
	encoder.AddString("RedactedComplex64s", zapredactor.RedactValue(s.RedactedComplex64s, ""))
	encoder.AddString("RedactedComplex64p", zapredactor.RedactValue(s.RedactedComplex64p, ""))
	encoder.AddString("RedactedComplex128", zapredactor.RedactValue(s.RedactedComplex128, ""))
	encoder.AddString("RedactedComplex128s", zapredactor.RedactValue(s.RedactedComplex128s, ""))
	encoder.AddString("RedactedComplex128p", zapredactor.RedactValue(s.RedactedComplex128p, ""))
	encoder.AddString("RedactedBool", zapredactor.RedactValue(s.RedactedBool, ""))
	encoder.AddString("RedactedBoolp", zapredactor.RedactValue(s.RedactedBoolp, ""))
	encoder.AddString("RedactedTime", zapredactor.RedactValue(s.RedactedTime, ""))
	encoder.AddString("RedactedTimes", zapredactor.RedactValue(s.RedactedTimes, ""))
	encoder.AddString("RedactedTimep", zapredactor.RedactValue(s.RedactedTimep, ""))
	encoder.AddString("RedactedDuration", zapredactor.RedactValue(s.RedactedDuration, ""))
	encoder.AddString("RedactedDurations", zapredactor.RedactValue(s.RedactedDurations, ""))
	encoder.AddString("RedactedDurationp", zapredactor.RedactValue(s.RedactedDurationp, ""))
	encoder.AddString("RedactedByte", zapredactor.RedactValue(s.RedactedByte, ""))
	return nil
}
