package service_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/service"
)

func TestNewBaseConfig(t *testing.T) {
	var ctx = struct {
		Project string
		Name    string
		ID      int
	}{
		Project: "test",
		Name:    "testapp",
		ID:      1,
	}
	cfg := service.NewBaseConfig(ctx, nil)

	if cfg == nil {
		t.Fatal("NewBaseConfig returned nil")
	}

	if !reflect.DeepEqual(cfg.GetContext(), ctx) {
		t.Errorf("CoreConfig does not match: got %v, want %v", cfg.GetContext(), ctx)
	}
}

func TestSetFlags(t *testing.T) {
	cfg := service.NewBaseConfig(struct{}{}, nil)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg.SetupFlags(fs)

	flags := []string{"-c", "config.json", "-o", "export.json", "-t", "-T"}
	err := fs.Parse(flags)
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	c := fs.Lookup("c")
	if c == nil || c.Value.String() != "config.json" {
		t.Errorf("Expected -c flag to be 'config.json', got %v", c)
	}

	e := fs.Lookup("o")
	if e == nil || e.Value.String() != "export.json" {
		t.Errorf("Expected -o flag to be 'export.json', got %v", e)
	}

	test := fs.Lookup("t")
	if test == nil || test.Value.String() != "true" {
		t.Errorf("Expected -i flag to be 'true', got %v", test)
	}
	disableTemplate := fs.Lookup("T")
	if disableTemplate == nil || disableTemplate.Value.String() != "true" {
		t.Errorf("Expected -T flag to be 'true', got %v", disableTemplate)
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() func()
		expectedExit   bool
		expectedErrStr string
	}{
		{
			name: "Load from file",
			setupMock: func() func() {
				tmpDir, err := os.MkdirTemp("", "config_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				configPath := filepath.Join(tmpDir, "config.json")
				configContent := `{"project":"test","name":"testapp","id":1}`
				err = os.WriteFile(configPath, []byte(configContent), 0644)
				if err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}
				oldArgs := os.Args
				os.Args = []string{"cmd", "-c", configPath}
				return func() {
					os.Args = oldArgs
					os.RemoveAll(tmpDir)
				}
			},
			expectedExit: false,
		},
		{
			name: "Load from HTTP",
			setupMock: func() func() {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"project":"test","name":"testapp","id":1}`))
				}))
				oldArgs := os.Args
				os.Args = []string{"cmd", "-c", server.URL}
				return func() {
					os.Args = oldArgs
					server.Close()
				}
			},
			expectedExit: false,
		},
		{
			name: "Export config",
			setupMock: func() func() {
				tmpDir, err := os.MkdirTemp("", "config_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				exportPath := filepath.Join(tmpDir, "export.json")
				oldArgs := os.Args
				os.Args = []string{"cmd", "-o", exportPath}
				return func() {
					os.Args = oldArgs
					os.RemoveAll(tmpDir)
				}
			},
			expectedExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			cfg := service.NewBaseConfig(struct{}{}, nil)
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			cfg.SetupFlags(fs)
			fs.Parse(os.Args[1:])

			err := cfg.Load()
			_, exit := errkit.ExitCode(err)
			if exit {
				err = nil
			}

			if exit != tt.expectedExit {
				t.Errorf("Expected exit to be %v, got %v", tt.expectedExit, exit)
			}

			if tt.expectedErrStr == "" && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.expectedErrStr != "" && (err == nil || err.Error() != tt.expectedErrStr) {
				t.Errorf("Expected error '%s', got '%v'", tt.expectedErrStr, err)
			}
		})
	}
}

func TestLoadFromHTTPWithRedirects(t *testing.T) {
	redirectCount := 0
	maxRedirects := 5

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectCount < maxRedirects {
			redirectCount++
			w.Header().Set("Location", "/redirect")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"project":"test","name":"testapp","id":1}`))
	}))
	defer server.Close()

	cfg := service.NewBaseConfig(struct{}{}, nil)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg.SetupFlags(fs)
	fs.Parse([]string{"-c", server.URL})

	err := cfg.Load()
	_, exit := errkit.ExitCode(err)

	if exit {
		t.Errorf("Expected exit to be false, got true")
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if redirectCount != maxRedirects {
		t.Errorf("Expected %d redirects, got %d", maxRedirects, redirectCount)
	}
}

type contextConfig struct {
	Project string
	Name    string
	ID      int
}

func TestExportConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	exportPath := filepath.Join(tmpDir, "export.json")
	ctx := contextConfig{
		Project: "test",
		Name:    "testapp",
		ID:      1,
	}
	components := []component.Config{
		{Name: "comp1"},
	}

	cfg := service.NewBaseConfig(ctx, components)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg.SetupFlags(fs)
	fs.Parse([]string{"-o", exportPath})

	err = cfg.Load()
	_, exit := errkit.ExitCode(err)
	if exit {
		err = nil
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exit {
		t.Errorf("Expected exit to be true, got false")
	}

	exportedData, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportedConfig struct {
		Context    contextConfig      `json:",omitempty"`
		Components []component.Config `json:",omitempty"`
	}
	err = json.Unmarshal(exportedData, &exportedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal exported config: %v", err)
	}

	// Use a custom comparison function to ignore the Enabled field
	if !compareContextConfig(exportedConfig.Context, ctx) {
		t.Errorf("Exported context config does not match original: got %v, want %v", exportedConfig.Context, ctx)
	}
	if !compareComponentsConfig(exportedConfig.Components, components) {
		t.Errorf("Exported components config does not match original: got %v, want %v", exportedConfig.Components, components)
	}
}

func compareContextConfig(a, b contextConfig) bool {
	return a.Project == b.Project && a.Name == b.Name && a.ID == b.ID
}

func compareComponentsConfig(a, b []component.Config) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].UUID != b[i].UUID {
			return false
		}
		if a[i].Name != b[i].Name {
			return false
		}
		if !bytes.Equal(a[i].Refs.Bytes(), b[i].Refs.Bytes()) {
			return false
		}
		if !bytes.Equal(a[i].Options.Bytes(), b[i].Options.Bytes()) {
			return false
		}
	}
	return true
}

func TestLoadConfigWithoutSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := service.NewBaseConfig(struct{}{}, nil)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg.SetupFlags(fs)

	// Don't set the -c flag, let it use the default empty string
	// fs.Parse([]string{"-c", nonExistentPath})

	err = cfg.Load()
	_, exit := errkit.ExitCode(err)

	if exit {
		t.Errorf("Expected exit to be false, got true")
	}
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
