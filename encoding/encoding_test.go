package encoding_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gopherd/core/container/maputil"
	"github.com/gopherd/core/container/sliceutil"
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
				for k, v := range m {
					result += fmt.Sprintf("%s=%s;", k, v)
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
			result, err := encoding.Transform(tt.input, tt.encoder, tt.decoder)

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
		for _, k := range sliceutil.Sort(maputil.Keys(m)) {
			result.WriteString(fmt.Sprintf("%s = %s\n", k, m[k]))
		}
		return []byte(result.String()), nil
	}

	result, err := encoding.Transform(input, encoder, decoder)
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
