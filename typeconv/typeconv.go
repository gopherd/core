package typeconv

import (
	"fmt"
	"strconv"

	"github.com/gopherd/core/operator"
)

const Unused = 0

var _true = []byte("true")
var _false = []byte("false")

func ToBytes(v any) ([]byte, error) {
	if byter, ok := v.(interface{ Bytes() []byte }); ok {
		return byter.Bytes(), nil
	}
	if stringer, ok := v.(fmt.Stringer); ok {
		return []byte(stringer.String()), nil
	}
	var b []byte
	switch value := v.(type) {
	case bool:
		b = operator.Ternary(value, _true, _false)
	case int:
		b = strconv.AppendInt(b, int64(value), 10)
	case int8:
		b = strconv.AppendInt(b, int64(value), 10)
	case int16:
		b = strconv.AppendInt(b, int64(value), 10)
	case int32:
		b = strconv.AppendInt(b, int64(value), 10)
	case int64:
		b = strconv.AppendInt(b, int64(value), 10)
	case uint:
		b = strconv.AppendUint(b, uint64(value), 10)
	case uint8:
		b = strconv.AppendUint(b, uint64(value), 10)
	case uint16:
		b = strconv.AppendUint(b, uint64(value), 10)
	case uint32:
		b = strconv.AppendUint(b, uint64(value), 10)
	case uint64:
		b = strconv.AppendUint(b, uint64(value), 10)
	case float32:
		b = strconv.AppendFloat(b, float64(value), 'g', -1, 32)
	case float64:
		b = strconv.AppendFloat(b, value, 'g', -1, 64)
	case string:
		b = []byte(value)
	case []byte:
		b = value
	default:
		b = []byte(fmt.Sprint(v))
	}
	return b, nil
}

func ToString(v any) string {
	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}
	if byter, ok := v.(interface{ Bytes() []byte }); ok {
		return string(byter.Bytes())
	}
	switch value := v.(type) {
	case bool:
		return operator.Ternary(value, "true", "false")
	case int:
		return strconv.FormatInt(int64(value), 10)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(int64(value), 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprint(v)
	}
}
