package zaparray

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Bools []bool

func (bs Bools) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bs {
		arr.AppendBool(bs[i])
	}
	return nil
}

type ByteStringsArray [][]byte

func (bss ByteStringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bss {
		arr.AppendByteString(bss[i])
	}
	return nil
}

type Complex128s []complex128

func (nums Complex128s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendComplex128(nums[i])
	}
	return nil
}

type Complex64s []complex64

func (nums Complex64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendComplex64(nums[i])
	}
	return nil
}

type Durations []time.Duration

func (ds Durations) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ds {
		arr.AppendDuration(ds[i])
	}
	return nil
}

type Float64s []float64

func (nums Float64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat64(nums[i])
	}
	return nil
}

type Float32s []float32

func (nums Float32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat32(nums[i])
	}
	return nil
}

type Ints []int

func (nums Ints) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt(nums[i])
	}
	return nil
}

type Int64s []int64

func (nums Int64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt64(nums[i])
	}
	return nil
}

type Int32s []int32

func (nums Int32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt32(nums[i])
	}
	return nil
}

type Int16s []int16

func (nums Int16s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt16(nums[i])
	}
	return nil
}

type Int8s []int8

func (nums Int8s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt8(nums[i])
	}
	return nil
}

type StringArray []string

func (ss StringArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i])
	}
	return nil
}

type Times []time.Time

func (ts Times) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ts {
		arr.AppendTime(ts[i])
	}
	return nil
}

type Uints []uint

func (nums Uints) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint(nums[i])
	}
	return nil
}

type Uint64s []uint64

func (nums Uint64s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint64(nums[i])
	}
	return nil
}

type Uint32s []uint32

func (nums Uint32s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint32(nums[i])
	}
	return nil
}

type Uint16s []uint16

func (nums Uint16s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint16(nums[i])
	}
	return nil
}

type Uint8s []uint8

func (nums Uint8s) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint8(nums[i])
	}
	return nil
}

type Uintptrs []uintptr

func (nums Uintptrs) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUintptr(nums[i])
	}
	return nil
}
