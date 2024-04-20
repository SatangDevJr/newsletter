package convert_test

import (
	"newsletter/src/pkg/utils/convert"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkContainsInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convert.ContainsInt([]int{1, 2, 3, 4, 5}, 3)
	}
}

func BenchmarkContainsString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convert.ContainsStr([]string{"1", "2", "3", "4", "5"}, "3")
	}
}

func TestConvert_StringToInt32(t *testing.T) {

	t.Run("should parse string int to int", func(t *testing.T) {
		got := convert.ValueStringToInt32("10")

		assert.Equal(t, int32(10), got)
	})

	t.Run("should parse string float to int", func(t *testing.T) {
		got := convert.ValueStringToInt32("10.000")

		assert.Equal(t, int32(10), got)
	})
}

func TestConvert_ValueArrayStringToString(t *testing.T) {
	t.Run("should parse slice string to string format", func(t *testing.T) {
		r := convert.ValueArrayStringToString([]string{"A", "B"})

		assert.Equal(t, "'A','B'", r)
	})
}

func TestConvert_ValueArrayStringToStringN(t *testing.T) {
	t.Run("should parse slice string to string format N", func(t *testing.T) {
		r := convert.ValueArrayStringToStringN([]string{"A", "B"})

		assert.Equal(t, "N'A',N'B'", r)
	})
}

func TestConvert_AnyToStrings(t *testing.T) {

	t.Run("should return empty string when dose not input any value", func(t *testing.T) {
		res := convert.AnyToStrings()

		assert.Empty(t, res)
	})

	t.Run("should return empty string when input nil", func(t *testing.T) {
		var input *string

		res := convert.AnyToStrings(input)

		assert.Empty(t, res)
	})

	t.Run("should return value of string pointer when input string pointer", func(t *testing.T) {
		input := convert.ValueToStringPointer("hellooooo")

		res := convert.AnyToStrings(input)

		assert.Equal(t, *input, res)
	})

	t.Run("should return value string  when input string", func(t *testing.T) {
		input := "hellooooo"

		res := convert.AnyToStrings(input)

		assert.Equal(t, input, res)
	})

	t.Run("should return value of string pointers when input string pointers", func(t *testing.T) {
		firstInput := convert.ValueToStringPointer("hellooooo")
		secondInput := convert.ValueToStringPointer("worlddddd")

		res := convert.AnyToStrings(firstInput, secondInput)

		assert.Equal(t, "hellooooo - worlddddd", res)
	})

	t.Run("should return value strings  when input strings", func(t *testing.T) {
		firstInput := "hellooooo"
		secondInput := "worlddddd"

		res := convert.AnyToStrings(firstInput, secondInput)

		assert.Equal(t, "hellooooo - worlddddd", res)
	})

	t.Run("should return value strings when input onces string pointer and onces string", func(t *testing.T) {
		firstInput := convert.ValueToStringPointer("hellooooo")
		secondInput := "worlddddd"

		res := convert.AnyToStrings(firstInput, secondInput)

		assert.Equal(t, "hellooooo - worlddddd", res)
	})

	t.Run("should return value boolean strings when input boolean", func(t *testing.T) {
		firstInput := true

		res := convert.AnyToStrings(firstInput)

		assert.Equal(t, "true", res)
	})

	t.Run("should return value boolean strings when input pointer boolean", func(t *testing.T) {
		firstInput := convert.ValueToBoolPointer(false)

		res := convert.AnyToStrings(firstInput)

		assert.Equal(t, "false", res)
	})

}

func TestConvert_ValueToStringPointer(t *testing.T) {
	t.Run("should return pointer of string when input string", func(t *testing.T) {
		param := "hello"

		result := convert.ValueToStringPointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_StringPointerToValue(t *testing.T) {
	t.Run("should return value of string pointer when input string pointer", func(t *testing.T) {
		param := "hello"

		result := convert.StringPointerToValue(&param)

		assert.Equal(t, param, result)
	})

	t.Run("should return empty string when input nill", func(t *testing.T) {
		result := convert.StringPointerToValue(nil)

		assert.Equal(t, "", result)
	})
}

func TestConvert_ValueToBoolPointer(t *testing.T) {
	t.Run("should return pointer of boolean when input boolean", func(t *testing.T) {
		param := true

		result := convert.ValueToBoolPointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueToIntPointer(t *testing.T) {
	t.Run("should return pointer of int when input int", func(t *testing.T) {
		param := int(16)

		result := convert.ValueToIntPointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueToInt32Pointer(t *testing.T) {
	t.Run("should return pointer of int32 when input int32", func(t *testing.T) {
		param := int32(16)

		result := convert.ValueToInt32Pointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueToInt64Pointer(t *testing.T) {
	t.Run("should return pointer of int64 when input int64", func(t *testing.T) {
		param := int64(16)

		result := convert.ValueToInt64Pointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueToFloat64Pointer(t *testing.T) {
	t.Run("should return pointer of float64 when input float64", func(t *testing.T) {
		param := float64(16)

		result := convert.ValueToFloat64Pointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueToTimePointer(t *testing.T) {
	t.Run("should return pointer of time when input time", func(t *testing.T) {
		param := time.Date(2021, 12, 22, 0, 0, 0, 0, time.UTC)

		result := convert.ValueToTimePointer(param)

		assert.Equal(t, param, *result)
	})
}

func TestConvert_ValueArrayInt32ToString(t *testing.T) {
	t.Run("should return string that display list of slice of int32", func(t *testing.T) {
		param := []int32{1, 2, 3}

		result := convert.ValueArrayInt32ToString(param)

		assert.Equal(t, "1,2,3", result)
	})
}

func TestConvert_ValueArrayInt64ToString(t *testing.T) {
	t.Run("should return string that display list of slice of int64", func(t *testing.T) {
		param := []int64{1, 2, 3}

		result := convert.ValueArrayInt64ToString(param)

		assert.Equal(t, "1,2,3", result)
	})
}

func TestConvert_ValueStringToInt32(t *testing.T) {
	t.Run("should return int32 value when input string int32", func(t *testing.T) {
		param := "16"

		result := convert.ValueStringToInt32(param)

		assert.Equal(t, int32(16), result)
	})

	t.Run("should return 0 when input string that is not number", func(t *testing.T) {
		param := "hello"

		result := convert.ValueStringToInt32(param)

		assert.Equal(t, int32(0), result)
	})

	t.Run("should delete separate ',' when input string int32 that include separate ',' ", func(t *testing.T) {
		param := "16,000"

		result := convert.ValueStringToInt32(param)

		assert.Equal(t, int32(16000), result)
	})
}

func TestConvert_TestConvert_ValueStringToFloat64(t *testing.T) {
	t.Run("should return float32 value when input string float32", func(t *testing.T) {
		param := "16.50"

		result := convert.ValueStringToFloat64(param)

		assert.Equal(t, float64(16.5), result)
	})

	t.Run("should return 0 when input string that is not number", func(t *testing.T) {
		param := "hello"

		result := convert.ValueStringToFloat64(param)

		assert.Equal(t, float64(0), result)
	})

	t.Run("should delete separate ',' when input string float32 that include separate ',' ", func(t *testing.T) {
		param := "16,000.50"

		result := convert.ValueStringToFloat64(param)

		assert.Equal(t, float64(16000.5), result)
	})
}

func TestConvert_RoundUp(t *testing.T) {
	t.Run("should alway round up when input 100.000001 and decimal is 6", func(t *testing.T) {
		res := convert.RoundUp(100.0000001, 6)

		expected := float64(100.000001)
		assert.Equal(t, expected, res)
	})

	t.Run("should alway round up when input 100.00000005 and decimal is 2", func(t *testing.T) {
		res := convert.RoundUp(100.00000005, 2)

		expected := float64(100.01)
		assert.Equal(t, expected, res)
	})
}

func Test_RenderFloat(t *testing.T) {
	t.Run("should return 1,000,000.00 where input 1000000 and format like #,###.##", func(t *testing.T) {
		res := convert.RenderFloat("#,###.##", 1000000)

		expected := "1,000,000.00"
		assert.Equal(t, expected, res)
	})

	t.Run("should return 123,456.78 where input 123456.781 and format like #,###.##", func(t *testing.T) {
		res := convert.RenderFloat("#,###.##", 123456.781)

		expected := "123,456.78"
		assert.Equal(t, expected, res)
	})

	t.Run("should return 123,456.79 where input 123456.789 and format like #,###.##", func(t *testing.T) {
		res := convert.RenderFloat("#,###.##", 123456.789)

		expected := "123,456.79"
		assert.Equal(t, expected, res)
	})

	t.Run("should return 301 where input 301 and format like #,###.", func(t *testing.T) {
		res := convert.RenderFloat("#,###.", 301)

		expected := "301"
		assert.Equal(t, expected, res)
	})
}
