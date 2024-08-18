package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// ValueType is an interface that describes the common methods of our value types
type ValueType[T any] interface {
	Value() T
	SetValue(T)
	Deref() T
	Set(string) error
}

// testValueType is a generic function to test Value, SetValue, Deref, and Set methods
func testValueType[T ValueType[V], V comparable](t *testing.T, name string, zero, nonZero V, str string, newFunc func() T) {
	t.Run(name, func(t *testing.T) {
		v := newFunc()

		t.Run("Value", func(t *testing.T) {
			if v.Value() != zero {
				t.Errorf("Value() = %v, want %v", v.Value(), zero)
			}
		})

		t.Run("SetValue", func(t *testing.T) {
			v.SetValue(nonZero)
			if v.Value() != nonZero {
				t.Errorf("After SetValue(%v), got %v, want %v", nonZero, v.Value(), nonZero)
			}
		})

		t.Run("Deref", func(t *testing.T) {
			if v.Deref() != nonZero {
				t.Errorf("Deref() = %v, want %v", v.Deref(), nonZero)
			}
		})

		t.Run("Set", func(t *testing.T) {
			v = newFunc() // Reset to zero value
			err := v.Set(str)
			if err != nil {
				t.Fatalf("Set() error = %v", err)
			}
			if v.Value() != nonZero {
				t.Errorf("After Set(%q), got %v, want %v", str, v.Value(), nonZero)
			}

			if name != "String" { // String type doesn't return an error for invalid input
				err = v.Set("invalid")
				if err == nil {
					t.Errorf("Set(\"invalid\") should return error")
				}
			}
		})
	})
}

func TestValueTypes(t *testing.T) {
	testValueType(t, "Bool", false, true, "true", func() *Bool { return new(Bool) })
	testValueType(t, "Int", 0, 42, "42", func() *Int { return new(Int) })
	testValueType(t, "Int8", int8(0), int8(42), "42", func() *Int8 { return new(Int8) })
	testValueType(t, "Int16", int16(0), int16(42), "42", func() *Int16 { return new(Int16) })
	testValueType(t, "Int32", int32(0), int32(42), "42", func() *Int32 { return new(Int32) })
	testValueType(t, "Int64", int64(0), int64(42), "42", func() *Int64 { return new(Int64) })
	testValueType(t, "Uint", uint(0), uint(42), "42", func() *Uint { return new(Uint) })
	testValueType(t, "Uint8", uint8(0), uint8(42), "42", func() *Uint8 { return new(Uint8) })
	testValueType(t, "Uint16", uint16(0), uint16(42), "42", func() *Uint16 { return new(Uint16) })
	testValueType(t, "Uint32", uint32(0), uint32(42), "42", func() *Uint32 { return new(Uint32) })
	testValueType(t, "Uint64", uint64(0), uint64(42), "42", func() *Uint64 { return new(Uint64) })
	testValueType(t, "Float32", float32(0), float32(3.14), "3.14", func() *Float32 { return new(Float32) })
	testValueType(t, "Float64", float64(0), float64(3.14), "3.14", func() *Float64 { return new(Float64) })
	testValueType(t, "String", "", "hello", "hello", func() *String { return new(String) })
	testValueType(t, "Complex64", complex64(0), complex64(1+2i), "(1+2i)", func() *Complex64 { return new(Complex64) })
	testValueType(t, "Complex128", complex128(0), complex128(1+2i), "(1+2i)", func() *Complex128 { return new(Complex128) })
	testValueType(t, "Duration", time.Duration(0), 5*time.Second, "5s", func() *Duration { return new(Duration) })
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		str string
		d   time.Duration
		err bool
	}{
		{"5s", 5 * time.Second, false},
		{"5m", 5 * time.Minute, false},
		{"5h", 5 * time.Hour, false},
		{"5d", 5 * 24 * time.Hour, false},
		{"5", 5, false},
		{"0", 0, false},
		{"5d6h7m8s9ms10µs11ns", 5*24*time.Hour + 6*time.Hour + 7*time.Minute + 8*time.Second + 9*time.Millisecond + 10*time.Microsecond + 11*time.Nanosecond, false},
		{"", 0, true},
		{"d", 0, true},
		{"1h2d", 0, true},
		{"w", 0, true},
		{"y", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			d, err := parseDuration(tt.str)
			if err != nil && !tt.err {
				t.Errorf("ParseDuration() error = %v", err)
			}
			if d != tt.d {
				t.Errorf("ParseDuration() = %v, want %v", d, tt.d)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d   time.Duration
		str string
	}{
		{5 * time.Second, "5s"},
		{5 * time.Minute, "5m"},
		{5 * time.Hour, "5h"},
		{5*time.Hour + 7*time.Second, "5h7s"},
		{-5*time.Hour - 7*time.Second, "-5h7s"},
		{5 * 24 * time.Hour, "5d"},
		{5, "5ns"},
		{5*24*time.Hour + 6*time.Hour + 7*time.Minute + 8*time.Second + 9*time.Millisecond + 10*time.Microsecond + 11*time.Nanosecond, "5d6h7m8s9ms10µs11ns"},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if s := formatDuration(tt.d); s != tt.str {
				t.Errorf("FormatDuration() = %v, want %v", s, tt.str)
			}
		})
	}
}

func TestRawObject(t *testing.T) {
	t.Run("RawObject", func(t *testing.T) {
		data := []byte(`{"key":"value"}`)
		ro := NewRawObject(data)

		t.Run("NewRawObject", func(t *testing.T) {
			if string(ro) != string(data) {
				t.Errorf("NewRawObject() = %v, want %v", ro, data)
			}
		})

		t.Run("Len", func(t *testing.T) {
			if ro.Len() != len(data) {
				t.Errorf("Len() = %d, want %d", ro.Len(), len(data))
			}
		})

		t.Run("String", func(t *testing.T) {
			if ro.String() != string(data) {
				t.Errorf("String() = %s, want %s", ro.String(), string(data))
			}
		})

		t.Run("SetString", func(t *testing.T) {
			newData := `{"new":"data"}`
			ro.SetString(newData)
			if string(ro) != newData {
				t.Errorf("After SetString(), got %s, want %s", string(ro), newData)
			}
		})

		t.Run("Bytes", func(t *testing.T) {
			if string(ro.Bytes()) != string(ro) {
				t.Errorf("Bytes() = %v, want %v", ro.Bytes(), []byte(ro))
			}
		})

		t.Run("SetBytes", func(t *testing.T) {
			newData := []byte(`{"new":"bytes"}`)
			ro.SetBytes(newData)
			if string(ro) != string(newData) {
				t.Errorf("After SetBytes(), got %v, want %v", ro, newData)
			}
		})

		testMarshalUnmarshal := func(t *testing.T, marshal func(RawObject) ([]byte, error), unmarshal func(*RawObject, []byte) error, isText bool) {
			t.Run("Marshal", func(t *testing.T) {
				encoded, err := marshal(ro)
				if err != nil {
					t.Fatalf("Marshal error: %v", err)
				}
				if isText {
					decoded, err := base64.StdEncoding.DecodeString(string(encoded))
					if err != nil {
						t.Fatalf("base64 decode error: %v", err)
					}
					if string(decoded) != string(ro) {
						t.Errorf("Marshal() (decoded) = %v, want %v", string(decoded), string(ro))
					}
				} else {
					if string(encoded) != string(ro) {
						t.Errorf("Marshal() = %v, want %v", string(encoded), string(ro))
					}
				}
			})

			t.Run("Unmarshal", func(t *testing.T) {
				var decoded RawObject
				var input []byte
				if isText {
					input = []byte(base64.StdEncoding.EncodeToString(ro))
				} else {
					input = ro
				}
				err := unmarshal(&decoded, input)
				if err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if string(decoded) != string(ro) {
					t.Errorf("Unmarshal() = %v, want %v", string(decoded), string(ro))
				}
			})

			t.Run("UnmarshalNilPointer", func(t *testing.T) {
				var nilRO *RawObject
				err := unmarshal(nilRO, ro)
				if err == nil {
					t.Error("Unmarshal to nil pointer should return error")
				}
			})
		}

		t.Run("JSON", func(t *testing.T) {
			testMarshalUnmarshal(t,
				func(ro RawObject) ([]byte, error) { return ro.MarshalJSON() },
				func(ro *RawObject, data []byte) error { return ro.UnmarshalJSON(data) },
				false)
		})

		t.Run("Text", func(t *testing.T) {
			testMarshalUnmarshal(t,
				func(ro RawObject) ([]byte, error) { return ro.MarshalText() },
				func(ro *RawObject, data []byte) error { return ro.UnmarshalText(data) },
				true)
		})

		t.Run("Binary", func(t *testing.T) {
			testMarshalUnmarshal(t,
				func(ro RawObject) ([]byte, error) { return ro.MarshalBinary() },
				func(ro *RawObject, data []byte) error { return ro.UnmarshalBinary(data) },
				false)
		})

		t.Run("Decode", func(t *testing.T) {
			var result map[string]string
			err := ro.Decode(json.Unmarshal, &result)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			expected := map[string]string{"new": "bytes"}
			if result["key"] != expected["key"] {
				t.Errorf("Decode() = %v, want %v", result, expected)
			}

			var nilRO RawObject
			err = nilRO.Decode(json.Unmarshal, &result)
			if err != nil {
				t.Errorf("Decode() on nil RawObject should not return error, got %v", err)
			}
		})
	})
}

func TestDuration(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		d := Duration(5 * time.Second)
		data, err := d.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(data) != `"5s"` {
			t.Errorf("MarshalJSON() = %s, want \"5s\"", data)
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var d Duration
		err := d.UnmarshalJSON([]byte(`"5s"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if d != Duration(5*time.Second) {
			t.Errorf("After UnmarshalJSON(), got %v, want 5s", d)
		}

		err = d.UnmarshalJSON([]byte(`5000000000`))
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if d != Duration(5*time.Second) {
			t.Errorf("After UnmarshalJSON(), got %v, want 5s", d)
		}

		err = d.UnmarshalJSON([]byte(`"invalid"`))
		if err == nil {
			t.Errorf("UnmarshalJSON(\"invalid\") should return error")
		}
	})
}

func ExampleRawObject() {
	ro := NewRawObject(`{"name":"John","age":30}`)
	fmt.Println("Raw JSON:", ro)

	var person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := ro.Decode(json.Unmarshal, &person)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Decoded: Name=%s, Age=%d\n", person.Name, person.Age)

	// Output:
	// Raw JSON: {"name":"John","age":30}
	// Decoded: Name=John, Age=30
}

func ExampleDuration() {
	d := Duration(5 * time.Second)
	fmt.Println("Duration:", d.Value())

	jsonData, _ := json.Marshal(d)
	fmt.Println("JSON:", string(jsonData))

	var parsed Duration
	_ = parsed.Set("10m")
	fmt.Println("Parsed:", parsed.Deref())

	// Output:
	// Duration: 5s
	// JSON: "5s"
	// Parsed: 10m0s
}

func BenchmarkRawObjectDecode(b *testing.B) {
	ro := NewRawObject(`{"name":"John","age":30,"city":"New York"}`)
	var result map[string]interface{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ro.Decode(json.Unmarshal, &result)
	}
}

func BenchmarkDurationMarshalJSON(b *testing.B) {
	d := Duration(5 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.MarshalJSON()
	}
}
