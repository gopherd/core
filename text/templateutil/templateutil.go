package templateutil

import (
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strings"
)

var (
	stringsFuncs = template.FuncMap{
		"contains":    strings.Contains,
		"count":       strings.Count,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"index":       strings.Index,
		"join":        strings.Join,
		"lastIndex":   strings.LastIndex,
		"repeat":      strings.Repeat,
		"replace":     strings.Replace,
		"replaceN":    strings.ReplaceAll,
		"split":       strings.Split,
		"toLower":     strings.ToLower,
		"toUpper":     strings.ToUpper,
		"toTitle":     strings.ToTitle,
		"toValidUTF8": strings.ToValidUTF8,
		"trim":        strings.Trim,
		"trimLeft":    strings.TrimLeft,
		"trimRight":   strings.TrimRight,
		"trimPrefix":  strings.TrimPrefix,
		"trimSuffix":  strings.TrimSuffix,
		"trimSpace":   strings.TrimSpace,
	}
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
		"rune":    convToRune,
		"byte":    convToByte,
	}
	mathFuncs = template.FuncMap{
		"sum": sum,      // sum
		"add": add,      // +
		"sub": subtract, // -
		"mul": multiply, // *
		"div": divide,   // /
		"mod": mod,      // %
		"pow": pow,      // math.Pow
	}
)

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
	return withFuncs(template.FuncMap{}, stringsFuncs, mathFuncs, convFuncs)
}

// DefaultTemplate creates a new template with the default functions.
func DefaultTemplate(name string) *template.Template {
	return template.New(name).Funcs(DefaultFuncs())
}

func toFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

func toInt64(val any) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		if v <= math.MaxInt64 {
			return int64(v), nil
		}
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v <= math.MaxInt64 {
			return int64(v), nil
		}
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
	return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
}

func sum(vals ...any) (any, error) {
	if len(vals) == 0 {
		return 0, fmt.Errorf("no values")
	}
	if _, ok := vals[0].(string); ok {
		var sum string
		for _, val := range vals {
			sum += fmt.Sprint(val)
		}
		return sum, nil
	}

	var sumi int64
	var err error
	for _, val := range vals {
		var i int64
		i, err = toInt64(val)
		if err != nil {
			break
		}
		sumi += i
	}
	if err == nil {
		return sumi, nil
	}
	var sumf float64
	for _, val := range vals {
		f, err := toFloat64(val)
		if err != nil {
			return 0, err
		}
		sumf += f
	}
	return sumf, nil
}

func add(a, b any) (any, error) {
	if as, ok := a.(string); ok {
		return as + fmt.Sprint(b), nil
	}

	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
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

func subtract(a, b any) (any, error) {
	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
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

func multiply(a, b any) (any, error) {
	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
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

func divide(a, b any) (any, error) {
	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
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

func mod(a, b any) (any, error) {
	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
		if err == nil {
			return 0, err
		}
		return ai % bi, nil
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

func pow(a, b any) (any, error) {
	ai, err := toInt64(a)
	if err == nil {
		bi, err := toInt64(b)
		if err == nil {
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

func convToBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return v != "", nil
	default:
		if i, err := toInt64(val); err == nil {
			return i != 0, nil
		}
		if f, err := toFloat64(val); err == nil {
			return f != 0, nil
		}
		return false, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

func convToFloat32(val any) (float32, error) {
	f, err := toFloat64(val)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func convToFloat64(val any) (float64, error) {
	return toFloat64(val)
}

func convToInt(val any) (int, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func convToInt8(val any) (int8, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func convToInt16(val any) (int16, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func convToInt32(val any) (int32, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func convToInt64(val any) (int64, error) {
	return toInt64(val)
}

func convToUint(val any) (uint, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint(i), nil
}

func convToUint8(val any) (uint8, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint8(i), nil
}

func convToUint16(val any) (uint16, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint16(i), nil
}

func convToUint32(val any) (uint32, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint32(i), nil
}

func convToUint64(val any) (uint64, error) {
	i, err := toInt64(val)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value")
	}
	return uint64(i), nil
}

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

func convToRune(val any) (rune, error) {
	switch v := val.(type) {
	case rune:
		return v, nil
	case string:
		if len(v) == 0 {
			return 0, fmt.Errorf("empty string")
		}
		return rune(v[0]), nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(val))
	}
}

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
