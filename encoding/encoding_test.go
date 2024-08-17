package encoding_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gopherd/core/encoding"
)

func TestTransform(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		encoder     encoding.Encoder
		decoder     encoding.Decoder
		expected    []byte
		expectError bool
	}{
		{
			name:  "JSON to JSON (both encoder and decoder are JSON)",
			input: []byte(`{"name":"John","age":30,"hobbies":["reading","cycling"]}`),
			encoder: func(v any) ([]byte, error) {
				return json.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    []byte(`{"age":30,"hobbies":["reading","cycling"],"name":"John"}`),
			expectError: false,
		},
		{
			name:  "Simple JSON to JSON",
			input: []byte(`{"name":"John","age":30}`),
			encoder: func(v any) ([]byte, error) {
				return json.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    []byte(`{"age":30,"name":"John"}`),
			expectError: false,
		},
		{
			name:  "Complex JSON to JSON",
			input: []byte(`{"person":{"name":"Alice","age":25},"hobbies":["reading","swimming"],"address":{"city":"New York","zip":"10001"}}`),
			encoder: func(v any) ([]byte, error) {
				return json.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    []byte(`{"address":{"city":"New York","zip":"10001"},"hobbies":["reading","swimming"],"person":{"age":25,"name":"Alice"}}`),
			expectError: false,
		},
		{
			name:  "JSON Array to JSON Array",
			input: []byte(`[1,2,3,4,5]`),
			encoder: func(v any) ([]byte, error) {
				return json.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    []byte(`[1,2,3,4,5]`),
			expectError: false,
		},
		{
			name:  "Simple XML to XML",
			input: []byte(`<person><name>Alice</name><age>25</age></person>`),
			encoder: func(v any) ([]byte, error) {
				return xml.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				var person struct {
					XMLName xml.Name `xml:"person"`
					Name    string   `xml:"name"`
					Age     int      `xml:"age"`
				}
				if err := xml.Unmarshal(data, &person); err != nil {
					return err
				}
				*(v.(*any)) = person
				return nil
			},
			expected:    []byte(`<person><name>Alice</name><age>25</age></person>`),
			expectError: false,
		},
		{
			name:  "Complex XML to XML",
			input: []byte(`<root><person><name>Bob</name><age>40</age></person><hobbies><hobby>gardening</hobby><hobby>cooking</hobby></hobbies><address><city>London</city><zip>E1 6AN</zip></address></root>`),
			encoder: func(v any) ([]byte, error) {
				return xml.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				var root struct {
					XMLName xml.Name `xml:"root"`
					Person  struct {
						Name string `xml:"name"`
						Age  int    `xml:"age"`
					} `xml:"person"`
					Hobbies struct {
						Hobby []string `xml:"hobby"`
					} `xml:"hobbies"`
					Address struct {
						City string `xml:"city"`
						Zip  string `xml:"zip"`
					} `xml:"address"`
				}
				if err := xml.Unmarshal(data, &root); err != nil {
					return err
				}
				*(v.(*any)) = root
				return nil
			},
			expected:    []byte(`<root><person><name>Bob</name><age>40</age></person><hobbies><hobby>gardening</hobby><hobby>cooking</hobby></hobbies><address><city>London</city><zip>E1 6AN</zip></address></root>`),
			expectError: false,
		},
		{
			name:  "Custom format to Custom format",
			input: []byte(`key1=value1;key2=value2;key3=value3`),
			encoder: func(v any) ([]byte, error) {
				m := v.(map[string]string)
				var result string
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				slices.Sort(keys)
				for _, k := range keys {
					result += fmt.Sprintf("%s=%s;", k, m[k])
				}
				return []byte(result), nil
			},
			decoder: func(data []byte, v any) error {
				pairs := strings.Split(string(data), ";")
				m := make(map[string]string)
				for _, pair := range pairs {
					if pair == "" {
						continue
					}
					kv := strings.Split(pair, "=")
					if len(kv) != 2 {
						return fmt.Errorf("invalid format")
					}
					m[kv[0]] = kv[1]
				}
				*(v.(*any)) = m
				return nil
			},
			expected:    []byte(`key1=value1;key2=value2;key3=value3;`),
			expectError: false,
		},
		{
			name:  "Decoding Error",
			input: []byte(`invalid json`),
			encoder: func(v any) ([]byte, error) {
				return json.Marshal(v)
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:  "Encoding Error",
			input: []byte(`{"name":"John","age":30}`),
			encoder: func(v any) ([]byte, error) {
				return nil, fmt.Errorf("encoding error")
			},
			decoder: func(data []byte, v any) error {
				return json.Unmarshal(data, v)
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.Transform(tt.input, tt.decoder, tt.encoder)

			if tt.expectError {
				if err == nil {
					t.Errorf("Transform() error = nil, expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("Transform() error = %v, expected no error", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Transform() = %v, expected %v", string(result), string(tt.expected))
				}
			}
		})
	}
}

func ExampleTransform() {
	// Input data in a simple key-value format
	input := []byte("name:John Doe;age:30;city:New York")

	// Decoder: converts the input format to a map
	decoder := func(data []byte, v any) error {
		pairs := strings.Split(string(data), ";")
		m := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) == 2 {
				m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		*(v.(*any)) = m
		return nil
	}

	// Encoder: converts the map to a different key-value format
	encoder := func(v any) ([]byte, error) {
		m, ok := v.(map[string]string)
		if !ok {
			return nil, fmt.Errorf("expected map[string]string, got %T", v)
		}
		var result strings.Builder
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			result.WriteString(fmt.Sprintf("%s = %s\n", k, m[k]))
		}
		return []byte(result.String()), nil
	}

	result, err := encoding.Transform(input, decoder, encoder)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(string(result))
	// Output:
	// age = 30
	// city = New York
	// name = John Doe
}

func TestGetLineAndColumn(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		offset   int
		wantLine int
		wantCol  int
	}{
		{"Empty string", "", 0, 1, 1},
		{"Single char", "a", 0, 1, 1},
		{"Single char end", "a", 1, 1, 2},
		{"Multiple chars", "abc", 1, 1, 2},
		{"Multiple chars end", "abc", 3, 1, 4},
		{"Single line with LF", "abc\n", 3, 1, 4},
		{"Single line with LF end", "abc\n", 4, 2, 1},
		{"Two lines with LF", "abc\ndef", 4, 2, 1},
		{"Two lines with LF middle", "abc\ndef", 5, 2, 2},
		{"Single line with CRLF", "abc\r\n", 3, 1, 4},
		{"Single line with CRLF end", "abc\r\n", 5, 2, 1},
		{"Two lines with CRLF", "abc\r\ndef", 5, 2, 1},
		{"Two lines with CRLF middle", "abc\r\ndef", 6, 2, 2},
		{"Single line with CR", "abc\r", 3, 1, 4},
		{"Single line with CR end", "abc\r", 4, 2, 1},
		{"Two lines with CR", "abc\rdef", 4, 2, 1},
		{"Two lines with CR middle", "abc\rdef", 5, 2, 2},
		{"Mixed newlines", "abc\ndef\r\nghi\rjkl", 10, 3, 2},
		{"Offset beyond data", "abc", 5, 1, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotCol := encoding.GetPosition([]byte(tt.data), tt.offset)
			if gotLine != tt.wantLine || gotCol != tt.wantCol {
				t.Errorf("getLineAndColumn() = (%v, %v), want (%v, %v)",
					gotLine, gotCol, tt.wantLine, tt.wantCol)
			}
		})
	}
}

func TestGetJSONSourceError(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		err      error
		wantLine int
		wantCol  int
	}{
		{"Syntax error", `{"name":"John"`, &json.SyntaxError{Offset: 12}, 1, 13},
		{"Syntax error at end", `{"name":"John"`, &json.SyntaxError{Offset: 13}, 1, 14},
		{"Unmarshal type error", `{"name":30}`, &json.UnmarshalTypeError{Offset: 8}, 1, 9},
		{"Unmarshal type error at end", `{"name":30}`, &json.UnmarshalTypeError{Offset: 9}, 1, 10},
		{"Unknown error", `{"name":"John"}`, fmt.Errorf("unknown error"), 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := encoding.GetJSONSourceError("test.json", []byte(tt.data), tt.err)
			if err == nil {
				if tt.wantLine != 0 || tt.wantCol != 0 {
					t.Errorf("GetJSONSourceError() = nil, want line %v, col %v",
						tt.wantLine, tt.wantCol)
				}
			} else {
				if e, ok := err.(*encoding.SourceError); ok {
					if e.Line != tt.wantLine || e.Column != tt.wantCol {
						t.Errorf("GetJSONSourceError() = (%v, %v), want (%v, %v)",
							e.Line, e.Column, tt.wantLine, tt.wantCol)
					}
				} else {
					if err.Error() != "unknown error" {
						t.Errorf("GetJSONSourceError() = %v, want *JSONErrorContext", err)
					} else if tt.wantLine != 0 || tt.wantCol != 0 {
						t.Errorf("GetJSONSourceError() = %v, want line %v, col %v",
							err, tt.wantLine, tt.wantCol)
					}
				}
			}
		})
	}
}
