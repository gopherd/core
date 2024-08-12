// Package types provides utilities for handling raw JSON objects with delayed decoding.
package types

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gopherd/core/encoding"
)

// RawObject represents a raw object for delayed JSON decoding.
type RawObject []byte

// NewRawObject creates a new RawObject with the provided data.
func NewRawObject[T ~string | ~[]byte](v T) RawObject {
	return RawObject(v)
}

// Len returns the length of the Object's data.
func (o RawObject) Len() int {
	return len(o)
}

// String returns the string representation of the Object.
func (o RawObject) String() string {
	return string(o)
}

// SetString sets the string representation of the Object.
func (o *RawObject) SetString(s string) {
	*o = RawObject(s)
}

// Bytes returns the raw byte slice of the Object.
func (o RawObject) Bytes() []byte {
	return o
}

// SetBytes sets the raw byte slice of the Object.
func (o *RawObject) SetBytes(b []byte) {
	*o = RawObject(b)
}

func (o RawObject) marshal(nilValue []byte) ([]byte, error) {
	if o == nil {
		return nilValue, nil
	}
	return o, nil
}

func (o *RawObject) unmarshal(name string, data []byte) error {
	if o == nil {
		return fmt.Errorf("types.RawObject: %s on nil pointer", name)
	}
	*o = append((*o)[0:0], data...)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the raw JSON encoding of the Object.
func (o RawObject) MarshalJSON() ([]byte, error) {
	return o.marshal([]byte("null"))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It sets the Object's data to a copy of the input JSON data.
func (o *RawObject) UnmarshalJSON(data []byte) error {
	return o.unmarshal("UnmarshalJSON", data)
}

// MarshalText implements the encoding.TextMarshaler interface.
// It returns the base64 encoding of the Object's data.
func (o RawObject) MarshalText() ([]byte, error) {
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(o)))
	enc.Encode(buf, o)
	return buf, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It decodes the input text as base64 and sets the Object's data to the result.
func (o *RawObject) UnmarshalText(data []byte) error {
	if o == nil {
		return errors.New("types.RawObject: UnmarshalText on nil pointer")
	}
	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(dbuf, []byte(data))
	*o = dbuf[:n]
	return err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (o RawObject) MarshalBinary() ([]byte, error) {
	return o.marshal(nil)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (o *RawObject) UnmarshalBinary(data []byte) error {
	return o.unmarshal("UnmarshalBinary", data)
}

// Decode decodes the Object's data using the provided decoder.
// It does nothing and returns nil if the Object is empty.
func (o RawObject) Decode(decoder encoding.Decoder, v any) error {
	if o == nil {
		return nil
	}
	return decoder(o, v)
}

// Bool wraps a boolean value.
type Bool bool

func (b Bool) Value() bool {
	return bool(b)
}

func (b *Bool) SetValue(v bool) {
	*b = Bool(v)
}

func (b *Bool) Deref() bool {
	return bool(*b)
}

func (b *Bool) Set(v string) error {
	x, err := strconv.ParseBool(v)
	if err != nil {
		return fmt.Errorf("parse bool: %w", err)
	}
	*b = Bool(x)
	return nil
}

// Int wraps an integer value.
type Int int

func (i Int) Value() int {
	return int(i)
}

func (i *Int) SetValue(v int) {
	*i = Int(v)
}

func (i *Int) Deref() int {
	return int(*i)
}

func (i *Int) Set(v string) error {
	x, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("parse int: %w", err)
	}
	*i = Int(x)
	return nil
}

// Int8 wraps an int8 value.
type Int8 int8

func (i Int8) Value() int8 {
	return int8(i)
}

func (i *Int8) SetValue(v int8) {
	*i = Int8(v)
}

func (i *Int8) Deref() int8 {
	return int8(*i)
}

func (i *Int8) Set(v string) error {
	x, err := strconv.ParseInt(v, 10, 8)
	if err != nil {
		return fmt.Errorf("parse int8: %w", err)
	}
	*i = Int8(x)
	return nil
}

// Int16 wraps an int16 value.
type Int16 int16

func (i Int16) Value() int16 {
	return int16(i)
}

func (i *Int16) SetValue(v int16) {
	*i = Int16(v)
}

func (i *Int16) Deref() int16 {
	return int16(*i)
}

func (i *Int16) Set(v string) error {
	x, err := strconv.ParseInt(v, 10, 16)
	if err != nil {
		return fmt.Errorf("parse int16: %w", err)
	}
	*i = Int16(x)
	return nil
}

// Int32 wraps an int32 value.
type Int32 int32

func (i Int32) Value() int32 {
	return int32(i)
}

func (i *Int32) SetValue(v int32) {
	*i = Int32(v)
}

func (i *Int32) Deref() int32 {
	return int32(*i)
}

func (i *Int32) Set(v string) error {
	x, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return fmt.Errorf("parse int32: %w", err)
	}
	*i = Int32(x)
	return nil
}

// Int64 wraps an int64 value.
type Int64 int64

func (i Int64) Value() int64 {
	return int64(i)
}

func (i *Int64) SetValue(v int64) {
	*i = Int64(v)
}

func (i *Int64) Deref() int64 {
	return int64(*i)
}

func (i *Int64) Set(v string) error {
	x, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fmt.Errorf("parse int64: %w", err)
	}
	*i = Int64(x)
	return nil
}

// Uint wraps an unsigned integer value.
type Uint uint

func (u Uint) Value() uint {
	return uint(u)
}

func (u *Uint) SetValue(v uint) {
	*u = Uint(v)
}

func (u *Uint) Deref() uint {
	return uint(*u)
}

func (u *Uint) Set(v string) error {
	x, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return fmt.Errorf("parse uint: %w", err)
	}
	*u = Uint(x)
	return nil
}

// Uint8 wraps an uint8 value.
type Uint8 uint8

func (u Uint8) Value() uint8 {
	return uint8(u)
}

func (u *Uint8) SetValue(v uint8) {
	*u = Uint8(v)
}

func (u *Uint8) Deref() uint8 {
	return uint8(*u)
}

func (u *Uint8) Set(v string) error {
	x, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		return fmt.Errorf("parse uint8: %w", err)
	}
	*u = Uint8(x)
	return nil
}

// Uint16 wraps an uint16 value.
type Uint16 uint16

func (u Uint16) Value() uint16 {
	return uint16(u)
}

func (u *Uint16) SetValue(v uint16) {
	*u = Uint16(v)
}

func (u *Uint16) Deref() uint16 {
	return uint16(*u)
}

func (u *Uint16) Set(v string) error {
	x, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		return fmt.Errorf("parse uint16: %w", err)
	}
	*u = Uint16(x)
	return nil
}

// Uint32 wraps an uint32 value.
type Uint32 uint32

func (u Uint32) Value() uint32 {
	return uint32(u)
}

func (u *Uint32) SetValue(v uint32) {
	*u = Uint32(v)
}

func (u *Uint32) Deref() uint32 {
	return uint32(*u)
}

func (u *Uint32) Set(v string) error {
	x, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return fmt.Errorf("parse uint32: %w", err)
	}
	*u = Uint32(x)
	return nil
}

// Uint64 wraps an uint64 value.
type Uint64 uint64

func (u Uint64) Value() uint64 {
	return uint64(u)
}

func (u *Uint64) SetValue(v uint64) {
	*u = Uint64(v)
}

func (u *Uint64) Deref() uint64 {
	return uint64(*u)
}

func (u *Uint64) Set(v string) error {
	x, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return fmt.Errorf("parse uint64: %w", err)
	}
	*u = Uint64(x)
	return nil
}

// Float32 wraps a float32 value.
type Float32 float32

func (f Float32) Value() float32 {
	return float32(f)
}

func (f *Float32) SetValue(v float32) {
	*f = Float32(v)
}

func (f *Float32) Deref() float32 {
	return float32(*f)
}

func (f *Float32) Set(v string) error {
	x, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return fmt.Errorf("parse float32: %w", err)
	}
	*f = Float32(x)
	return nil
}

// Float64 wraps a float64 value.
type Float64 float64

func (f Float64) Value() float64 {
	return float64(f)
}

func (f *Float64) SetValue(v float64) {
	*f = Float64(v)
}

func (f *Float64) Deref() float64 {
	return float64(*f)
}

func (f *Float64) Set(v string) error {
	x, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("parse float64: %w", err)
	}
	*f = Float64(x)
	return nil
}

// String wraps a string value.
type String string

func (s String) Value() string {
	return string(s)
}

func (s *String) SetValue(v string) {
	*s = String(v)
}

func (s *String) Deref() string {
	return string(*s)
}

func (s *String) Set(v string) error {
	*s = String(v)
	return nil
}

// Complex64 wraps a complex64 value.
type Complex64 complex64

func (c Complex64) Value() complex64 {
	return complex64(c)
}

func (c *Complex64) SetValue(v complex64) {
	*c = Complex64(v)
}

func (c *Complex64) Deref() complex64 {
	return complex64(*c)
}

func (c *Complex64) Set(v string) error {
	x, err := strconv.ParseComplex(v, 64)
	if err != nil {
		return fmt.Errorf("parse complex64: %w", err)
	}
	*c = Complex64(x)
	return nil
}

// Complex128 wraps a complex128 value.
type Complex128 complex128

func (c Complex128) Value() complex128 {
	return complex128(c)
}

func (c *Complex128) SetValue(v complex128) {
	*c = Complex128(v)
}

func (c *Complex128) Deref() complex128 {
	return complex128(*c)
}

func (c *Complex128) Set(v string) error {
	x, err := strconv.ParseComplex(v, 128)
	if err != nil {
		return fmt.Errorf("parse complex128: %w", err)
	}
	*c = Complex128(x)
	return nil
}

// Duration wraps a time.Duration value.
type Duration time.Duration

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

func (d *Duration) SetValue(v time.Duration) {
	*d = Duration(v)
}

func (d *Duration) Deref() time.Duration {
	return time.Duration(*d)
}

func (d *Duration) Set(v string) error {
	x, err := time.ParseDuration(v)
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}
	*d = Duration(x)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Duration(d).String())), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] != '"' {
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*d = Duration(time.Duration(i))
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	x, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(x)
	return nil
}

// Time wraps a time.Time value.
type Time time.Time

func (t Time) Value() time.Time {
	return time.Time(t)
}

func (t *Time) SetValue(v time.Time) {
	*t = Time(v)
}

func (t *Time) Deref() time.Time {
	return time.Time(*t)
}

func (t *Time) Set(v string) error {
	x, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return fmt.Errorf("parse time: %w", err)
	}
	*t = Time(x)
	return nil
}
