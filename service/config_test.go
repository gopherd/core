package service

import (
	"bytes"
	"encoding/json"
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
		wantErr bool
		setup   func() error
		cleanup func()
	}{
		{
			name:    "Empty source",
			source:  "",
			wantErr: false,
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
		},
		{
			name:    "Non-existent file",
			source:  "nonexistent.json",
			wantErr: true,
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
			err := c.load(json.Unmarshal, tt.source, true)

			if (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.source != "" && tt.source != "nonexistent.json" {
				if c.Context.Name != "Test" {
					t.Errorf("load() did not properly populate Context, got: %v", c.Context)
				}
				if len(c.Components) != 1 || c.Components[0].Name != "TestComponent" {
					t.Errorf("load() did not properly populate Components, got: %v", c.Components)
				}
			}
		})
	}
}

func TestConfig_LoadFromHTTP(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		redirects      int
		wantErr        bool
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
		},
	}

	for i, tt := range tests {
		index := i
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			redirectCount := 0
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if redirectCount < tt.redirects {
					redirectCount++
					w.Header().Set("Location", fmt.Sprintf("%s/redirect/%d/%d", server.URL, index, redirectCount))
					w.WriteHeader(http.StatusFound)
					return
				}
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer func() {
				time.Sleep(100 * time.Millisecond)
				server.Close()
			}()

			c := &Config[TestContext]{}
			reader, err := c.loadFromHTTP(server.URL)

			if (err != nil) != tt.wantErr {
				t.Fatalf("test %v: loadFromHTTP() error = %v, wantErr %v, URL=%v", tt, err, tt.wantErr, server.URL)
				return
			}

			if err == nil {
				defer reader.Close()
				body, _ := io.ReadAll(reader)
				if string(body) != tt.responseBody {
					t.Fatalf("test %v: loadFromHTTP() got body = %v, want %v", tt, string(body), tt.responseBody)
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
						Name: "Component1",
						UUID: "{{.Name}}-UUID",
						Refs: types.RawObject(op.MustValue(json.Marshal(refs{
							A: "{{.Name}}-A",
							B: "B",
						}))),
						Options: types.RawObject(op.MustValue(json.Marshal(options{
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
						Name: "Component1",
						UUID: "TestName-UUID",
						Refs: types.RawObject(op.MustValue(json.Marshal(refs{
							A: "TestName-A",
							B: "B",
						}))),
						Options: types.RawObject(op.MustValue(json.Marshal(options{
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
						Name: "Component1",
						UUID: "{{.NameXXX}}-UUID",
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
						Name: "Component1",
						Refs: types.RawObject(op.MustValue(json.Marshal(map[string]string{
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
						Name: "Component1",
						Options: types.RawObject(op.MustValue(json.Marshal(map[string]string{
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
				t.Errorf("%s: processTemplate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.config, tt.want) {
				t.Errorf("%s: processTemplate() got = %v, want %v", tt.name, tt.config, tt.want)
			}
		})
	}
}

func TestConfig_Output(t *testing.T) {
	config := Config[TestContext]{
		Context: TestContext{Name: "TestName"},
		Components: []component.Config{
			{
				Name: "Component1",
				UUID: "TestUUID",
			},
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	config.output(jsonIdentEncoder)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	expected := `{
    "Context": {
        "Name": "TestName"
    },
    "Components": [
        {
            "Name": "Component1",
            "UUID": "TestUUID"
        }
    ]
}`

	if output != expected {
		t.Errorf("output() got = %v, want %v", output, expected)
	}
}

func ExampleConfig_output() {
	config := Config[TestContext]{
		Context: TestContext{Name: "ExampleName"},
		Components: []component.Config{
			{
				Name: "ExampleComponent",
				UUID: "ExampleUUID",
			},
		},
	}

	config.output(jsonIdentEncoder)
	// Output:
	// {
	//     "Context": {
	//         "Name": "ExampleName"
	//     },
	//     "Components": [
	//         {
	//             "Name": "ExampleComponent",
	//             "UUID": "ExampleUUID"
	//         }
	//     ]
	// }
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
	errReader := &rrrorReader{Err: io.ErrUnexpectedEOF}

	_, err := stripJSONComments(errReader)
	if err == nil {
		t.Errorf("stripJSONComments() error = nil, wantErr = true")
	}
}

// rrrorReader is a custom io.Reader that always returns an error
type rrrorReader struct {
	Err error
}

func (er *rrrorReader) Read(p []byte) (n int, err error) {
	return 0, er.Err
}
