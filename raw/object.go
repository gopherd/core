// Package raw provides utilities for handling raw JSON objects with delayed decoding.
package raw

import (
	"bytes"
	"encoding/json"
)

// Object represents a raw object for delayed JSON decoding.
type Object []byte

// Len returns the length of the Object's data.
func (o Object) Len() int {
	return len(o)
}

// Object returns a new Object with the provided string.
func String(s string) Object {
	return Object(s)
}

// Object returns a new Object with the provided byte slice.
func Bytes(b []byte) Object {
	return Object(b)
}

// String returns the string representation of the Object.
func (o Object) String() string {
	return string(o)
}

// SetString sets the string representation of the Object.
func (o *Object) SetString(s string) {
	*o = Object(s)
}

// Bytes returns the raw byte slice of the Object.
func (o Object) Bytes() []byte {
	return o
}

// SetBytes sets the raw byte slice of the Object.
func (o *Object) SetBytes(b []byte) {
	*o = Object(b)
}

// MarshalJSON implements the json.Marshaler interface.
// It returns the raw JSON encoding of the Object.
func (o Object) MarshalJSON() ([]byte, error) {
	return json.RawMessage(o).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It sets the Object's data to a copy of the input JSON data.
func (o *Object) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(o).UnmarshalJSON(data)
}

// DecodeJSON decodes the Object's JSON data into the provided value.
// It does nothing and returns nil if the Object is empty.
func (o Object) DecodeJSON(v any) error {
	if o == nil {
		return nil
	}
	return json.Unmarshal(o, v)
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
	return Object(buf.Bytes()), nil
}
