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
		o    types.Bytes
		want int
	}{
		{"Empty", types.Bytes{}, 0},
		{"NonEmpty", types.Bytes(`{"key":"value"}`), 15},
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
	o := types.Bytes(b)
	if !reflect.DeepEqual([]byte(o), b) {
		t.Errorf("Bytes() = %v, want %v", []byte(o), b)
	}
}

func TestObjectString(t *testing.T) {
	o := types.Bytes(`{"test":"value"}`)
	if o.String() != `{"test":"value"}` {
		t.Errorf("Object.String() = %v, want %v", o.String(), `{"test":"value"}`)
	}
}

func TestObjectSetString(t *testing.T) {
	var o types.Bytes
	s := `{"test":"value"}`
	o.SetString(s)
	if string(o) != s {
		t.Errorf("After SetString(), Object = %v, want %v", string(o), s)
	}
}

func TestObjectBytes(t *testing.T) {
	o := types.Bytes(`{"test":"value"}`)
	expected := []byte(`{"test":"value"}`)
	if !reflect.DeepEqual(o.Bytes(), expected) {
		t.Errorf("Object.Bytes() = %v, want %v", o.Bytes(), expected)
	}
}

func TestObjectSetBytes(t *testing.T) {
	var o types.Bytes
	b := []byte(`{"test":"value"}`)
	o.SetBytes(b)
	if !reflect.DeepEqual([]byte(o), b) {
		t.Errorf("After SetBytes(), Object = %v, want %v", []byte(o), b)
	}
}

func TestObjectMarshalJSON(t *testing.T) {
	o := types.Bytes(`{"test":"value"}`)
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if string(b) != `{"test":"value"}` {
		t.Errorf("MarshalJSON() = %v, want %v", string(b), `{"test":"value"}`)
	}
}

func TestObjectUnmarshalJSON(t *testing.T) {
	var o types.Bytes
	err := o.UnmarshalJSON([]byte(`{"test":"value"}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if string(o) != `{"test":"value"}` {
		t.Errorf("After UnmarshalJSON(), Object = %v, want %v", string(o), `{"test":"value"}`)
	}
}

func TestObjectDecodeJSON(t *testing.T) {
	o := types.Bytes(`{"test":"value"}`)
	var v map[string]string
	err := o.DecodeJSON(&v)
	if err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}
	if v["test"] != "value" {
		t.Errorf("DecodeJSON() result = %v, want map with 'test':'value'", v)
	}

	// Test with nil Object
	var nilO types.Bytes
	err = nilO.DecodeJSON(&v)
	if err != nil {
		t.Errorf("DecodeJSON() with nil Object should not return error, got %v", err)
	}
}

func TestMustJSON(t *testing.T) {
	v := map[string]string{"test": "value"}
	o := types.MustJSON(v)
	var decoded map[string]string
	err := json.Unmarshal(o, &decoded)
	if err != nil {
		t.Fatalf("Error decoding MustJSON result: %v", err)
	}
	if !reflect.DeepEqual(v, decoded) {
		t.Errorf("MustJSON() result = %v, want %v", decoded, v)
	}

	// Test panic scenario
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustJSON() should panic with un-marshallable value")
		}
	}()
	types.MustJSON(make(chan int)) // This should panic
}

func TestJSON(t *testing.T) {
	v := map[string]string{"test": "value"}
	o, err := types.JSON(v)
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	var decoded map[string]string
	err = json.Unmarshal(o, &decoded)
	if err != nil {
		t.Fatalf("Error decoding JSON result: %v", err)
	}
	if !reflect.DeepEqual(v, decoded) {
		t.Errorf("JSON() result = %v, want %v", decoded, v)
	}

	// Test error scenario
	_, err = types.JSON(make(chan int))
	if err == nil {
		t.Errorf("JSON() should return error with un-marshallable value")
	}
}
