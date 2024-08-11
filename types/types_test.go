package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gopherd/core/types"
)

func TestObjectLen(t *testing.T) {
	tests := []struct {
		name string
		o    types.RawObject
		want int
	}{
		{"Empty", types.RawObject{}, 0},
		{"NonEmpty", types.RawObject(`{"key":"value"}`), 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Len(); got != tt.want {
				t.Errorf("Object.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	s := `{"test":"value"}`
	o := types.String(s)
	if string(o) != s {
		t.Errorf("String() = %v, want %v", string(o), s)
	}
}

func TestBytes(t *testing.T) {
	b := []byte(`{"test":"value"}`)
	o := types.RawObject(b)
	if !reflect.DeepEqual([]byte(o), b) {
		t.Errorf("Bytes() = %v, want %v", []byte(o), b)
	}
}

func TestObjectString(t *testing.T) {
	o := types.RawObject(`{"test":"value"}`)
	if o.String() != `{"test":"value"}` {
		t.Errorf("Object.String() = %v, want %v", o.String(), `{"test":"value"}`)
	}
}

func TestObjectSetString(t *testing.T) {
	var o types.RawObject
	s := `{"test":"value"}`
	o.SetString(s)
	if string(o) != s {
		t.Errorf("After SetString(), Object = %v, want %v", string(o), s)
	}
}

func TestObjectBytes(t *testing.T) {
	o := types.RawObject(`{"test":"value"}`)
	expected := []byte(`{"test":"value"}`)
	if !reflect.DeepEqual(o.Bytes(), expected) {
		t.Errorf("Object.Bytes() = %v, want %v", o.Bytes(), expected)
	}
}

func TestObjectSetBytes(t *testing.T) {
	var o types.RawObject
	b := []byte(`{"test":"value"}`)
	o.SetBytes(b)
	if !reflect.DeepEqual([]byte(o), b) {
		t.Errorf("After SetBytes(), Object = %v, want %v", []byte(o), b)
	}
}

func TestObjectMarshalJSON(t *testing.T) {
	o := types.RawObject(`{"test":"value"}`)
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if string(b) != `{"test":"value"}` {
		t.Errorf("MarshalJSON() = %v, want %v", string(b), `{"test":"value"}`)
	}
}

func TestObjectUnmarshalJSON(t *testing.T) {
	var o types.RawObject
	err := o.UnmarshalJSON([]byte(`{"test":"value"}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if string(o) != `{"test":"value"}` {
		t.Errorf("After UnmarshalJSON(), Object = %v, want %v", string(o), `{"test":"value"}`)
	}
}

func TestObjectDecodeJSON(t *testing.T) {
	o := types.RawObject(`{"test":"value"}`)
	var v map[string]string
	err := o.Decode(json.Unmarshal, &v)
	if err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}
	if v["test"] != "value" {
		t.Errorf("DecodeJSON() result = %v, want map with 'test':'value'", v)
	}

	// Test with nil Object
	var nilO types.RawObject
	err = nilO.Decode(json.Unmarshal, &v)
	if err != nil {
		t.Errorf("DecodeJSON() with nil Object should not return error, got %v", err)
	}
}
