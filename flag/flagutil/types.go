package flagutil

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

func parseValue(ptr interface{}, value string) error {
	var err error
	switch x := ptr.(type) {
	case *bool:
		*x, err = strconv.ParseBool(value)
	case *int:
		*x, err = strconv.Atoi(value)
	case *uint:
		y, err := strconv.ParseUint(value, 10, strconv.IntSize)
		if err != nil {
			return err
		}
		*x = uint(y)
	case *int8:
		y, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		*x = int8(y)
	case *uint8:
		y, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		*x = uint8(y)
	case *int16:
		y, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		*x = int16(y)
	case *uint16:
		y, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		*x = uint16(y)
	case *int32:
		y, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		*x = int32(y)
	case *uint32:
		y, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		*x = uint32(y)
	case *int64:
		y, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*x = int64(y)
	case *uint64:
		y, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*x = uint64(y)
	case *float32:
		y, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		*x = float32(y)
	case *float64:
		y, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*x = y
	case *string:
		*x = value
	case *[]byte:
		*x = []byte(value)
	default:
		if v, ok := ptr.(flag.Value); ok {
			return v.Set(value)
		}
		return errParse
	}
	return err
}

// -- int32 Value
type Int32 int32

func NewInt32(val int32, p *int32) *Int32 {
	*p = val
	return (*Int32)(p)
}

func (i *Int32) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		err = numError(err)
	}
	*i = Int32(v)
	return err
}

func (i *Int32) Get() interface{} { return int32(*i) }

func (i *Int32) String() string { return strconv.FormatInt(int64(*i), 10) }

// -- uint32 Value
type Uint32 uint32

func NewUint32(val uint32, p *uint32) *Uint32 {
	*p = val
	return (*Uint32)(p)
}

func (i *Uint32) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		err = numError(err)
	}
	*i = Uint32(v)
	return err
}

func (i *Uint32) Get() interface{} { return uint32(*i) }

func (i *Uint32) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- float32 Value
type Float32 float32

func NewFloat32(val float32, p *float32) *Float32 {
	*p = val
	return (*Float32)(p)
}

func (f *Float32) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		err = numError(err)
	}
	*f = Float32(v)
	return err
}

func (f *Float32) Get() interface{} { return float32(*f) }

func (f *Float32) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

// --- slice[K]
type Slice[T any] []T

func (slice Slice[T]) Set(s string) error {
	var v T
	if err := parseValue(&v, s); err != nil {
		return err
	}
	slice = append(slice, v)
	return nil
}

// --- map[K]V
type Map[K comparable, V any] map[K]V

var errNoKey = errors.New("no key for map")

func (m Map[K, V]) Set(s string) error {
	var index = strings.IndexByte(s, '=')
	if index <= 0 {
		return errNoKey
	}
	var k K
	if err := parseValue(&k, s[:index]); err != nil {
		return err
	}
	var v V
	if err := parseValue(&v, s[index+1:]); err != nil {
		return err
	}
	m[k] = v
	return nil
}
