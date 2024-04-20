package convert

import (
	newsletter "newsletter/src/pkg/utils/error"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ContainsInt(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsStr(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func ValueToStringPointer(s string) *string {
	p := s
	return &p
}

func StringPointerToValue(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func ValueToBoolPointer(s bool) *bool {
	p := s
	return &p
}

func ValueToIntPointer(s int) *int {
	p := s
	return &p
}

func ValueToInt32Pointer(s int32) *int32 {
	p := s
	return &p
}

func ValueToInt64Pointer(s int64) *int64 {
	p := s
	return &p
}

func ValueToFloat64Pointer(s float64) *float64 {
	p := s
	return &p
}

func ValueToTimePointer(s time.Time) *time.Time {
	p := s
	return &p
}

func ValueToErrorCodePointer(e newsletter.ErrorCode) *newsletter.ErrorCode {
	p := e
	return &p
}

func ValueArrayInt32ToString(i []int32) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(i)), ","), "[]")
}

func ValueArrayInt64ToString(i []int64) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(i)), ","), "[]")
}

func ValueArrayStringToString(s []string) string {
	strValues := []string{}
	for i := 0; i < len(s); i++ {
		strValues = append(strValues, fmt.Sprintf("'%v'", s[i]))
	}
	return strings.Trim(strings.Join(strValues, ","), "[]")
}

func ValueArrayStringToStringN(s []string) string {
	strValues := []string{}
	for i := 0; i < len(s); i++ {
		strValues = append(strValues, fmt.Sprintf("N'%v'", s[i]))
	}
	return strings.Trim(strings.Join(strValues, ","), "[]")
}

func ValueStringToFloat64(i string) float64 {
	numc := strings.Trim(i, " ")
	strReplace := fmt.Sprintf(`%[1]s`, ",")
	numc = strings.Replace(numc, strReplace, "", -1)

	res, err := strconv.ParseFloat(numc, 64)
	if err != nil {
		return float64(0)
	}

	return float64(res)
}

func ValueStringToInt32(i string) int32 {
	numc := strings.Trim(i, " ")
	strReplace := fmt.Sprintf(`%[1]s`, ",")
	numc = strings.Replace(numc, strReplace, "", -1)

	res, err := strconv.ParseFloat(numc, 64)
	if err != nil {
		return int32(0)
	}

	return int32(res)
}

func ValueStringToInt64(i string) int64 {
	numc := strings.Trim(i, " ")
	strReplace := fmt.Sprintf(`%[1]s`, ",")
	numc = strings.Replace(numc, strReplace, "", -1)

	res, err := strconv.ParseFloat(numc, 64)
	if err != nil {
		return int64(0)
	}

	return int64(res)
}

func AnyToStrings(values ...interface{}) string {
	results := []string{}

	for _, value := range values {
		if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
			continue
		}

		if reflect.ValueOf(value).Kind() == reflect.Ptr {
			v := reflect.ValueOf(value).Elem()
			results = append(results, fmt.Sprintf("%v", v))
		} else {
			results = append(results, fmt.Sprintf("%v", value))
		}

	}
	return strings.Join(results, " - ")
}

func TimeToFormat(t time.Time, f string) string {
	f = strings.Replace(f, "yyyy", strconv.Itoa((t.Year())), -1)
	f = strings.Replace(f, "MM", fmt.Sprintf("%02d", int(t.Month())), -1)
	f = strings.Replace(f, "dd", fmt.Sprintf("%02d", t.Day()), -1)
	f = strings.Replace(f, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	f = strings.Replace(f, "mm", fmt.Sprintf("%02d", t.Minute()), -1)
	f = strings.Replace(f, "ss", fmt.Sprintf("%02d", t.Second()), -1)
	return f
}

func RoundUp(x float64, decimalPlaces int) float64 {
	incremental := math.Pow10(decimalPlaces)
	xf := big.NewFloat(x)
	xfIncremental, _ := new(big.Float).Mul(xf, big.NewFloat(incremental)).Float64()
	return math.Ceil(xfIncremental) / float64(incremental)
}

func RenderFloat(format string, n float64) string {
	// Examples of format strings, given n = 12345.6789:
	// "#,###.##" => "12,345.67"
	// "#,###." => "12,345"
	// "#,###" => "12345,678"
	// "#\u202F###,##" => "12 345,67"
	// "#.###,###### => 12.345,678900
	// "" (aka default format) => 12,345.67

	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	if math.IsNaN(n) {
		return "NaN"
	}
	if n > math.MaxFloat64 {
		return "Infinity"
	}
	if n < -math.MaxFloat64 {
		return "-Infinity"
	}

	var renderFloatPrecisionMultipliers = [10]float64{
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
	}

	var renderFloatPrecisionRounders = [10]float64{
		0.5,
		0.05,
		0.005,
		0.0005,
		0.00005,
		0.000005,
		0.0000005,
		0.00000005,
		0.000000005,
		0.0000000005,
	}

	// default format
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = ""

		// collect indices of meaningful formatting directives
		formatDirectiveChars := []rune(format)
		formatDirectiveIndices := make([]int, 0)
		for i, char := range formatDirectiveChars {
			if char != '#' && char != '0' {
				formatDirectiveIndices = append(formatDirectiveIndices, i)
			}
		}

		if len(formatDirectiveIndices) > 0 {
			// Directive at index 0:
			//   Must be a '+'
			//   Raise an error if not the case
			// index: 0123456789
			//        +0.000,000
			//        +000,000.0
			//        +0000.00
			//        +0000
			if formatDirectiveIndices[0] == 0 {
				if formatDirectiveChars[formatDirectiveIndices[0]] != '+' {
					panic("RenderFloat(): invalid positive sign directive")
				}
				positiveStr = "+"
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// Two directives:
			//   First is thousands separator
			//   Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatDirectiveIndices) == 2 {
				if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
					panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				}
				thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// One directive:
			//   Directive is decimal separator
			//   The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatDirectiveIndices) == 1 {
				decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
			}
		}
	}

	// generate sign part
	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.Itoa(int(intf))

	// add thousand separator if required
	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if precision == 0 {
		return signStr + intStr
	}

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr
}
