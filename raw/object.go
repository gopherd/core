// Package raw provides utilities for handling raw JSON objects with delayed decoding.
package raw

import (
	"bytes"
	"encoding/json"
)

var null = []byte("null")

// Object represents a raw object for delayed JSON decoding.
type Object struct {
	data []byte
}

// Object returns a new Object with the provided string.
func String(s string) Object {
	return Object{
		data: []byte(s),
	}
}

// Object returns a new Object with the provided byte slice.
func Bytes(b []byte) Object {
	return Object{
		data: b,
	}
}

// String returns the string representation of the Object.
func (o Object) String() string {
	return string(o.data)
}

// Bytes returns the raw byte slice of the Object.
func (o Object) Bytes() []byte {
	return o.data
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the raw JSON encoding of the Object.
func (o Object) MarshalJSON() ([]byte, error) {
	if o.data == nil {
		return []byte("null"), nil
	}
	return o.data, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It sets the Object's data to a copy of the input JSON data.
func (o *Object) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		o.data = nil
		return nil
	}
	o.data = append(o.data[:0], data...)
	return nil
}

// DecodeJSON decodes the Object's JSON data into the provided value.
// It does nothing and returns nil if the Object is empty.
func (o *Object) DecodeJSON(v any) error {
	if len(o.data) == 0 {
		return nil
	}
	return json.Unmarshal(o.data, v)
}

// MustJSON returns a new Object containing the MustJSON encoding of v.
// It panics if the encoding fails.
func MustJSON(v any) Object {
	o, err := JSON(v)
	if err != nil {
		panic(err)
	}
	return o
}

// JSON returns a new Object containing the JSON encoding of v.
func JSON(v any) (Object, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(v); err != nil {
		return Object{}, err
	}
	return Object{
		data: buf.Bytes(),
	}, nil
}
