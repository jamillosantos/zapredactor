//go:generate go run github.com/golang/mock/mockgen -destination=redactor_mock_test.go -package=redactreflection go.uber.org/zap/zapcore ObjectEncoder

package redactreflection

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

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
	hiddenField string
}

var defaultTestingValue = testWalkFieldsStruct{
	Name: "test",
	Sub: testWalkFieldsStructSub{
		SubField1: "sub1",
		SubField2: 2,
		SubField3: "sub3",
	},
	Sub2: testWalkFieldsStructSub2{
		SubField1: "sub2.sub1",
	},
	hiddenField: "hidden",
}

func Test_walkFields(t *testing.T) {
	obj := defaultTestingValue
	fields := make([]string, 0)
	walkFields(obj, "", DefaultNameExtractor, func(field string, value interface{}) {
		fields = append(fields, field)
	})
	require.Equal(t, []string{"Name", "Sub.SubField1", "Sub.SubField2", "Sub.sub_field_3", "sub_2.subField1"}, fields)
}

func Test_redactByReflection_Redact(t *testing.T) {
	t.Run("should redact all fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		data := defaultTestingValue

		gomock.InOrder(
			encoder.EXPECT().AddString("Name", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.SubField1", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.SubField2", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.sub_field_3", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("sub_2.subField1", redactors.DefaultRedactedString),
		)

		redactor := newRedactor(data)

		err := redactor.Redact(encoder)
		require.NoError(t, err)
	})

	t.Run("should redact all fields that are not allowed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		data := defaultTestingValue

		gomock.InOrder(
			encoder.EXPECT().AddString("Name", data.Name),
			encoder.EXPECT().AddString("Sub.SubField1", redactors.DefaultRedactedString),
			encoder.EXPECT().AddInt("Sub.SubField2", data.Sub.SubField2),
			encoder.EXPECT().AddString("Sub.sub_field_3", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("sub_2.subField1", data.Sub2.SubField1),
		)

		redactor := newRedactor(defaultTestingValue,
			WithAllowFields("Name", "Sub.SubField2", "sub_2.subField1"),
		)

		err := redactor.Redact(encoder)
		require.NoError(t, err)
	})

	t.Run("should redact fields with a custom redactor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		data := defaultTestingValue
		data.Sub.SubField1 = "1234567890124321" // fake PAN number

		gomock.InOrder(
			encoder.EXPECT().AddString("Name", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.SubField1", "123456******4321"),
			encoder.EXPECT().AddString("Sub.SubField2", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.sub_field_3", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("sub_2.subField1", redactors.DefaultRedactedString),
		)

		redactor := newRedactor(data,
			WithRedactor("Sub.SubField1", redactors.PAN64),
		)

		err := redactor.Redact(encoder)
		require.NoError(t, err)
	})

	t.Run("should hide a field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		encoder := NewMockObjectEncoder(ctrl)

		data := defaultTestingValue

		gomock.InOrder(
			encoder.EXPECT().AddString("Name", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.SubField2", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("Sub.sub_field_3", redactors.DefaultRedactedString),
			encoder.EXPECT().AddString("sub_2.subField1", redactors.DefaultRedactedString),
		)

		redactor := newRedactor(data,
			WithHiddenField("Sub.SubField1"),
		)

		err := redactor.Redact(encoder)
		require.NoError(t, err)
	})
}
