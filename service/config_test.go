package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/encoding"
	"github.com/gopherd/core/op"
	"github.com/gopherd/core/types"
)

type TestContext struct {
	Name string
}

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		decoder encoding.Decoder
		wantErr bool
		setup   func() error
		cleanup func()
		check   func(*testing.T, *Config[TestContext])
	}{
		{
			name:    "Empty source",
			source:  "",
			wantErr: false,
			check: func(t *testing.T, c *Config[TestContext]) {
				if c.Context.Name != "" || len(c.Components) != 0 {
					t.Errorf("Expected empty config, got: %+v", c)
				}
			},
		},
		{
			name:    "Stdin source",
			source:  "-",
			wantErr: false,
			setup: func() error {
				content := `{"Context":{"Name":"Test"},"Components":[{"Name":"TestComponent"}]}`
				r, w, _ := os.Pipe()
				os.Stdin = r
				go func() {
					defer w.Close()
					w.Write([]byte(content))
				}()
				return nil
			},
			cleanup: func() {
				os.Stdin = os.NewFile(0, "/dev/stdin")
			},
			check: func(t *testing.T, c *Config[TestContext]) {
				if c.Context.Name != "Test" || len(c.Components) != 1 || c.Components[0].Name != "TestComponent" {
					t.Errorf("Unexpected config: %+v", c)
				}
			},
		},
		{
			name:    "File source",
			source:  "testconfig.json",
			wantErr: false,
			setup: func() error {
				content := `{"Context":{"Name":"Test"},"Components":[{"Name":"TestComponent"}]}`
				return os.WriteFile("testconfig.json", []byte(content), 0644)
			},
			cleanup: func() {
				os.Remove("testconfig.json")
			},
			check: func(t *testing.T, c *Config[TestContext]) {
				if c.Context.Name != "Test" || len(c.Components) != 1 || c.Components[0].Name != "TestComponent" {
					t.Errorf("Unexpected config: %+v", c)
				}
			},
		},
		{
			name:    "Non-existent file",
			source:  "nonexistent.json",
			wantErr: true,
		},
		{
			name:   "Custom decoder",
			source: "testconfig.toml",
			decoder: func(data []byte, v interface{}) error {
				// Mock TOML decoder
				return json.Unmarshal(data, v)
			},
			setup: func() error {
				content := `{"Context":{"Name":"TestTOML"},"Components":[{"Name":"TOMLComponent"}]}`
				return os.WriteFile("testconfig.toml", []byte(content), 0644)
			},
			cleanup: func() {
				os.Remove("testconfig.toml")
			},
			check: func(t *testing.T, c *Config[TestContext]) {
				if c.Context.Name != "TestTOML" || len(c.Components) != 1 || c.Components[0].Name != "TOMLComponent" {
					t.Errorf("Unexpected config: %+v", c)
				}
			},
		},
		{
			name:    "Invalid JSON",
			source:  "invalid.json",
			wantErr: true,
			setup: func() error {
				content := `{"Context":{"Name":"Test"},"Components":[{"Name":"TestComponent"}`
				return os.WriteFile("invalid.json", []byte(content), 0644)
			},
			cleanup: func() {
				os.Remove("invalid.json")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			c := &Config[TestContext]{}
			err := c.load(os.Stdin, tt.decoder, tt.source)

			if (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

func TestConfig_LoadFromHTTP(t *testing.T) {
	timeout := 100 * time.Millisecond
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		redirects      int
		wantErr        bool
		errorCheck     func(error) bool
	}{
		{
			name:           "Successful request",
			responseStatus: http.StatusOK,
			responseBody:   `{"Context":{"Name":"Test"},"Components":[{"Name":"TestComponent"}]}`,
			wantErr:        false,
		},
		{
			name:           "Not Found",
			responseStatus: http.StatusNotFound,
			wantErr:        true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "HTTP request failed with status code: 404")
			},
		},
		{
			name:           "Single Redirect",
			responseStatus: http.StatusOK,
			redirects:      1,
			responseBody:   `{"Context":{"Name":"Test"},"Components":[{"Name":"TestComponent"}]}`,
			wantErr:        false,
		},
		{
			name:           "Too Many Redirects",
			responseStatus: http.StatusOK,
			redirects:      33,
			wantErr:        true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "too many redirects")
			},
		},
		{
			name:           "Timeout",
			responseStatus: http.StatusOK,
			wantErr:        true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "context deadline exceeded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			redirectCount := 0
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "Timeout" {
					time.Sleep(timeout + 50*time.Millisecond)
				}
				if redirectCount < tt.redirects {
					redirectCount++
					w.Header().Set("Location", fmt.Sprintf("%s/redirect", server.URL))
					w.WriteHeader(http.StatusFound)
					return
				}
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			c := &Config[TestContext]{}
			reader, err := c.loadFromHTTP(server.URL, timeout)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errorCheck != nil && err != nil {
				if !tt.errorCheck(err) {
					t.Errorf("loadFromHTTP() error = %v, does not match expected error condition", err)
				}
			}

			if err == nil {
				defer reader.Close()
				body, _ := io.ReadAll(reader)
				if string(body) != tt.responseBody {
					t.Errorf("loadFromHTTP() got body = %v, want %v", string(body), tt.responseBody)
				}
			}
		})
	}
}

func TestConfig_ProcessTemplate(t *testing.T) {
	type refs struct {
		A string
		B string
	}
	type options struct {
		C string
		D string
	}
	tests := []struct {
		name           string
		config         Config[TestContext]
		enableTemplate bool
		source         string
		wantErr        bool
		want           Config[TestContext]
	}{
		{
			name: "Process all templates",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name:            "Component1",
						UUID:            "{{.Name}}-UUID",
						TemplateUUID:    op.Addr(types.Bool(true)),
						TemplateRefs:    op.Addr(types.Bool(true)),
						TemplateOptions: op.Addr(types.Bool(true)),
						Refs: types.NewRawObject(op.MustValue(json.Marshal(refs{
							A: "{{.Name}}-A",
							B: "B",
						}))),
						Options: types.NewRawObject(op.MustValue(json.Marshal(options{
							C: "C",
							D: "{{.Name}}-D",
						}))),
					},
				},
			},
			enableTemplate: true,
			source:         "config.json",
			wantErr:        false,
			want: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name:            "Component1",
						UUID:            "TestName-UUID",
						TemplateUUID:    op.Addr(types.Bool(true)),
						TemplateRefs:    op.Addr(types.Bool(true)),
						TemplateOptions: op.Addr(types.Bool(true)),
						Refs: types.NewRawObject(op.MustValue(json.Marshal(refs{
							A: "TestName-A",
							B: "B",
						}))),
						Options: types.NewRawObject(op.MustValue(json.Marshal(options{
							C: "C",
							D: "TestName-D",
						}))),
					},
				},
			},
		},
		{
			name: "Template disabled",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name: "Component1",
						UUID: "{{.Name}}-UUID",
					},
				},
			},
			enableTemplate: false,
			source:         "config.json",
			wantErr:        false,
			want: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name: "Component1",
						UUID: "{{.Name}}-UUID",
					},
				},
			},
		},
		{
			name: "Invalid UUID template",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name:         "Component1",
						UUID:         "{{.NameXXX}}-UUID",
						TemplateUUID: op.Addr(types.Bool(true)),
					},
				},
			},
			enableTemplate: true,
			source:         "config.json",
			wantErr:        true,
		},
		{
			name: "Invalid Refs template",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name:         "Component1",
						TemplateRefs: op.Addr(types.Bool(true)),
						Refs: types.NewRawObject(op.MustValue(json.Marshal(map[string]string{
							"A": "{{.NameXXX}}-A",
							"B": "B",
						}))),
					},
				},
			},
			enableTemplate: true,
			source:         "config.json",
			wantErr:        true,
		},
		{
			name: "Invalid Options template",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name:            "Component1",
						TemplateOptions: op.Addr(types.Bool(true)),
						Options: types.NewRawObject(op.MustValue(json.Marshal(map[string]string{
							"A": "{{.NameXXX}}-A",
							"B": "B",
						}))),
					},
				},
			},
			enableTemplate: true,
			source:         "config.json",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.processTemplate(tt.enableTemplate, tt.source)

			if (err != nil) != tt.wantErr {
				t.Errorf("processTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.config, tt.want) {
				t.Errorf("processTemplate() got = %v, want %v", tt.config, tt.want)
			}
		})
	}
}

func TestConfig_Output(t *testing.T) {
	tests := []struct {
		name    string
		config  Config[TestContext]
		encoder encoding.Encoder
		want    string
	}{
		{
			name: "Default JSON output",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name: "Component1",
						UUID: "TestUUID",
					},
				},
			},
			encoder: nil,
			want: `{
    "Context": {
        "Name": "TestName"
    },
    "Components": [
        {
            "Name": "Component1",
            "UUID": "TestUUID"
        }
    ]
}`,
		},
		{
			name: "Custom encoder output",
			config: Config[TestContext]{
				Context: TestContext{Name: "TestName"},
				Components: []component.Config{
					{
						Name: "Component1",
						UUID: "TestUUID",
					},
				},
			},
			encoder: func(v interface{}) ([]byte, error) {
				return []byte("Custom encoded output"), nil
			},
			want: "Custom encoded output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			tt.config.output(os.Stdout, io.Discard, tt.encoder)
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := strings.TrimSpace(buf.String())

			if output != tt.want {
				t.Errorf("output() got = %v, want %v", output, tt.want)
			}
		})
	}
}

func TestJsonIdentEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name: "Simple struct",
			input: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			want: `{
    "Name": "John",
    "Age": 30
}`,
			wantErr: false,
		},
		{
			name:    "String with HTML",
			input:   "<h1>Hello, World!</h1>",
			want:    `"<h1>Hello, World!</h1>"`,
			wantErr: false,
		},
		{
			name: "Nested struct",
			input: struct {
				Person struct {
					Name string
					Age  int
				}
				Address string
			}{
				Person: struct {
					Name string
					Age  int
				}{
					Name: "Alice",
					Age:  25,
				},
				Address: "123 Main St",
			},
			want: `{
    "Person": {
        "Name": "Alice",
        "Age": 25
    },
    "Address": "123 Main St"
}`,
			wantErr: false,
		},
		{
			name:    "Unsupported type",
			input:   complex(1, 2),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonIndentEncoder(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("jsonIdentEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotStr := strings.TrimSpace(string(got))
				if gotStr != tt.want {
					t.Errorf("jsonIdentEncoder() got = %v, want %v", gotStr, tt.want)
				}
			}
		})
	}
}

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "Basic comment removal",
			input: `{
    // This is a comment
    "name": "John",
    "age": 30 // This is an inline comment
}`,
			expected: `{
    "name": "John",
    "age": 30 // This is an inline comment
}`,
			wantErr: false,
		},
		{
			name: "Multiple comments",
			input: `{
    // Comment 1
    "a": 1,
    // Comment 2
    "b": 2,
    // Comment 3
    "c": 3
}`,
			expected: `{
    "a": 1,
    "b": 2,
    "c": 3
}`,
			wantErr: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name: "Only comments",
			input: `// Comment 1
// Comment 2
// Comment 3`,
			expected: "",
			wantErr:  false,
		},
		{
			name: "Comments with varying indentation",
			input: `{
    "a": 1,
  // Indented comment
        // More indented comment
    "b": 2
}`,
			expected: `{
    "a": 1,
    "b": 2
}`,
			wantErr: false,
		},
		{
			name: "Preserve strings with //",
			input: `{
    "url": "https://example.com",
    "comment": "This string contains // which is not a comment"
}`,
			expected: `{
    "url": "https://example.com",
    "comment": "This string contains // which is not a comment"
}`,
			wantErr: false,
		},
		{
			name: "Comment at the end of file",
			input: `{
    "name": "John"
}
// Comment at the end`,
			expected: `{
    "name": "John"
}`,
			wantErr: false,
		},
		{
			name: "Comment with special characters",
			input: `{
    // Comment with special chars: !@#$%^&*()_+
    "data": "value"
}`,
			expected: `{
    "data": "value"
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := stripJSONComments(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("stripJSONComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("stripJSONComments() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

// TestStripJSONCommentsError tests the error handling of stripJSONComments
func TestStripJSONCommentsError(t *testing.T) {
	// Create a reader that always returns an error
	errReader := &errorReader{Err: errors.New("read error")}

	_, err := stripJSONComments(errReader)
	if err == nil {
		t.Errorf("stripJSONComments() error = nil, wantErr = true")
	}
	if err.Error() != "read error" {
		t.Errorf("stripJSONComments() error = %v, want 'read error'", err)
	}
}

// errorReader is a custom io.Reader that always returns an error
type errorReader struct {
	Err error
}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, er.Err
}
