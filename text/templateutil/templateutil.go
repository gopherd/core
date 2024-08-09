// Package templateutil provides utility functions for working with Go templates.
package templateutil

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strings"
)

var (
	// containerFuncs contains functions for working with containers.
	containerFuncs = template.FuncMap{
		"len":      func(vals []any) int { return len(vals) },
		"list":     func(vals ...any) []any { return vals },
		"bools":    func(vals ...bool) []bool { return vals },
		"strings":  func(vals ...string) []string { return vals },
		"ints":     func(vals ...int) []int { return vals },
		"int8s":    func(vals ...int8) []int8 { return vals },
		"int16s":   func(vals ...int16) []int16 { return vals },
		"int32s":   func(vals ...int32) []int32 { return vals },
		"int64s":   func(vals ...int64) []int64 { return vals },
		"uints":    func(vals ...uint) []uint { return vals },
		"uint8s":   func(vals ...uint8) []uint8 { return vals },
		"uint16s":  func(vals ...uint16) []uint16 { return vals },
		"uint32s":  func(vals ...uint32) []uint32 { return vals },
		"uint64s":  func(vals ...uint64) []uint64 { return vals },
		"float32s": func(vals ...float32) []float32 { return vals },
		"float64s": func(vals ...float64) []float64 { return vals },
	}

	// stringsFuncs contains functions for working with strings.
	stringsFuncs = template.FuncMap{
		"contains":    strings.Contains,
		"count":       strings.Count,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"index":       strings.Index,
		"join":        func(sep string, vals ...string) string { return strings.Join(vals, sep) },
		"lastIndex":   strings.LastIndex,
		"repeat":      strings.Repeat,
		"replace":     strings.Replace,
		"replaceAll":  strings.ReplaceAll,
		"split":       strings.Split,
		"toLower":     strings.ToLower,
		"toUpper":     strings.ToUpper,
		"toValidUTF8": strings.ToValidUTF8,
		"trim":        strings.Trim,
		"trimLeft":    strings.TrimLeft,
		"trimRight":   strings.TrimRight,
		"trimPrefix":  strings.TrimPrefix,
		"trimSuffix":  strings.TrimSuffix,
		"trimSpace":   strings.TrimSpace,
	}

	// convFuncs contains functions for type conversion.
	convFuncs = template.FuncMap{
		"float32": convToFloat32,
		"float64": convToFloat64,
		"int":     convToInt,
		"int8":    convToInt8,
		"int16":   convToInt16,
		"int32":   convToInt32,
		"int64":   convToInt64,
		"uint":    convToUint,
		"uint8":   convToUint8,
		"uint16":  convToUint16,
		"uint32":  convToUint32,
		"uint64":  convToUint64,
		"bool":    convToBool,
		"string":  convToString,
		"bytes":   convToBytes,
		"runes":   convToRunes,
		"rune":    convToRune,
		"byte":    convToByte,
	}

	// mathFuncs contains mathematical functions.
	mathFuncs = template.FuncMap{
		"sum": sum,
		"add": add,
		"sub": subtract,
		"mul": multiply,
		"div": divide,
		"mod": mod,
		"pow": pow,
	}
)

// withFuncs merges multiple FuncMaps into a single FuncMap.
func withFuncs(dst template.FuncMap, funcs ...template.FuncMap) template.FuncMap {
	for _, f := range funcs {
		for k, v := range f {
			dst[k] = v
		}
	}
	return dst
}

// DefaultFuncs returns the default template functions.
func DefaultFuncs() template.FuncMap {
	return withFuncs(template.FuncMap{},
		containerFuncs,
		stringsFuncs,
		mathFuncs,
		convFuncs,
	)
}

// DefaultTemplate creates a new template with the default functions.
func DefaultTemplate(name string) *template.Template {
	return template.New(name).Funcs(DefaultFuncs())
}

// Execute executes the default template with the given text and data.
func Execute(name, text string, data any, options ...string) (string, error) {
	var buf bytes.Buffer
	t := DefaultTemplate(name)
	if len(options) > 0 {
		t = t.Option(options...)
	}
	if t, err := t.Parse(text); err != nil {
		return "", err
	} else if err := t.Execute(&buf, data); err != nil {
		return "", err
	} else {
		return buf.String(), nil
	}
}

// toFloat64 converts various numeric types to float64.
func toFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float32, float64:
		return reflect.ValueOf(v).Float(), nil
	default:
		if i, err := toInt64(val, true); err == nil {
			return float64(i), nil
		}
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// toInt64 converts various numeric types to int64.
func toInt64(val any, strict bool) (int64, error) {
	switch v := val.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		uv := reflect.ValueOf(v).Uint()
		if uv > math.MaxInt64 {
			return 0, fmt.Errorf("value %v overflows int64", uv)
		}
		return int64(uv), nil
	case float32, float64:
		if strict {
			return 0, fmt.Errorf("cannot convert float to int")
		}
		return int64(reflect.ValueOf(v).Float()), nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// sum calculates the sum of given values.
func sum(vals ...any) (any, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("no values provided")
	}
	if _, ok := vals[0].(string); ok {
		var sum string
		for _, val := range vals {
			sum += fmt.Sprint(val)
		}
		return sum, nil
	}

	var sumI int64
	var err error
	for _, val := range vals {
		var i int64
		i, err = toInt64(val, true)
		if err != nil {
			break
		}
		sumI += i
	}
	if err == nil {
		return sumI, nil
	}

	var sumF float64
	for _, val := range vals {
		f, err := toFloat64(val)
		if err != nil {
			return 0, err
		}
		sumF += f
	}
	return sumF, nil
}

// add adds two values.
func add(a, b any) (any, error) {
	if as, ok := a.(string); ok {
		return as + fmt.Sprint(b), nil
	}

	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			return ai + bi, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	return af + bf, nil
}

// subtract subtracts two values.
func subtract(a, b any) (any, error) {
	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			return ai - bi, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	return af - bf, nil
}

// multiply multiplies two values.
func multiply(a, b any) (any, error) {
	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			return ai * bi, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	return af * bf, nil
}

// divide divides two values.
func divide(a, b any) (any, error) {
	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			if bi == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return ai / bi, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	if bf == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return af / bf, nil
}

// mod calculates the modulus of two values.
func mod(a, b any) (any, error) {
	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			if bi == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			return ai % bi, nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	if bf == 0 {
		return 0, fmt.Errorf("modulo by zero")
	}
	return math.Mod(af, bf), nil
}

// pow calculates the power of two values.
func pow(a, b any) (any, error) {
	ai, err := toInt64(a, true)
	if err == nil {
		bi, err := toInt64(b, true)
		if err == nil {
			if bi < 0 {
				return math.Pow(float64(ai), float64(bi)), nil
			}
			return int64(math.Pow(float64(ai), float64(bi))), nil
		}
	}

	af, err := toFloat64(a)
	if err != nil {
		return 0, err
	}
	bf, err := toFloat64(b)
	if err != nil {
		return 0, err
	}
	return math.Pow(af, bf), nil
}

// convToBool converts a value to bool.
func convToBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return v != "", nil
	default:
		if i, err := toInt64(val, true); err == nil {
			return i != 0, nil
		}
		if f, err := toFloat64(val); err == nil {
			return f != 0, nil
		}
		return false, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// convToFloat32 converts a value to float32.
func convToFloat32(val any) (float32, error) {
	f, err := toFloat64(val)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// convToFloat64 converts a value to float64.
func convToFloat64(val any) (float64, error) {
	return toFloat64(val)
}

// convToInt converts a value to int.
func convToInt(val any) (int, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// convToInt8 converts a value to int8.
func convToInt8(val any) (int8, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

// convToInt16 converts a value to int16.
func convToInt16(val any) (int16, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

// convToInt32 converts a value to int32.
func convToInt32(val any) (int32, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// convToInt64 converts a value to int64.
func convToInt64(val any) (int64, error) {
	return toInt64(val, false)
}

// convToUint converts a value to uint.
func convToUint(val any) (uint, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint(i), nil
}

// convToUint8 converts a value to uint8.
func convToUint8(val any) (uint8, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	if i < 0 || i > math.MaxUint8 {
		return 0, fmt.Errorf("value out of range for uint8")
	}
	return uint8(i), nil
}

// convToUint16 converts a value to uint16.
func convToUint16(val any) (uint16, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	if i < 0 || i > math.MaxUint16 {
		return 0, fmt.Errorf("value out of range for uint16")
	}
	return uint16(i), nil
}

// convToUint32 converts a value to uint32.
func convToUint32(val any) (uint32, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	if i < 0 || i > math.MaxUint32 {
		return 0, fmt.Errorf("value out of range for uint32")
	}
	return uint32(i), nil
}

// convToUint64 converts a value to uint64.
func convToUint64(val any) (uint64, error) {
	i, err := toInt64(val, false)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint64(i), nil
}

// convToString converts a value to string.
func convToString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprint(val), nil
	}
}

// convToBytes converts a value to []byte.
func convToBytes(val any) ([]byte, error) {
	switch v := val.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// convToRunes converts a value to []rune.
func convToRunes(val any) ([]rune, error) {
	switch v := val.(type) {
	case string:
		return []rune(v), nil
	case []rune:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// convToRune converts a value to rune.
func convToRune(val any) (rune, error) {
	switch v := val.(type) {
	case rune:
		return v, nil
	case string:
		if len(v) == 0 {
			return 0, fmt.Errorf("empty string")
		}
		return []rune(v)[0], nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

// convToByte converts a value to byte.
func convToByte(val any) (byte, error) {
	switch v := val.(type) {
	case byte:
		return v, nil
	case string:
		if len(v) == 0 {
			return 0, fmt.Errorf("empty string")
		}
		return v[0], nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}
