package service

import (
	"context"
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"log/slog"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
)

func TestMain(m *testing.M) {
	// disable logging during tests
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	os.Exit(m.Run())
}

func newBaseServiceTest[T any](config Config[T]) *BaseService[T] {
	s := NewBaseService(config)
	s.stderr = io.Discard
	s.versionFunc = func() {}
	return s
}

// resetFlagsAndArgs 重置 flag.CommandLine 和 os.Args
func resetFlagsAndArgs() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"test"}
}

func TestNewBaseService(t *testing.T) {
	type testContext struct {
		Field string
	}
	config := Config[testContext]{Context: testContext{Field: "test"}}
	service := newBaseServiceTest(config)

	if service == nil {
		t.Fatal("NewBaseService returned nil")
	}

	if !reflect.DeepEqual(service.config, config) {
		t.Errorf("Config mismatch. Got %v, want %v", service.config, config)
	}

	if service.components == nil {
		t.Error("Components group is nil")
	}
}

func TestSetVersionFunc(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	called := false
	versionFunc := func() { called = true }

	service.SetVersionFunc(versionFunc)
	service.versionFunc()

	if !called {
		t.Error("Version function was not called")
	}

	service.SetVersionFunc(nil)
	if service.versionFunc != nil {
		t.Error("Version function was not set to nil")
	}
}

func TestGetComponent(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{uuid: "test-uuid"}
	service.components.AddComponent("test-uuid", mockComponent)

	result := service.GetComponent("test-uuid")
	if result != mockComponent {
		t.Errorf("GetComponent returned wrong component. Got %v, want %v", result, mockComponent)
	}

	result = service.GetComponent("non-existent")
	if result != nil {
		t.Errorf("GetComponent should return nil for non-existent component, got %v", result)
	}
}

func TestConfig(t *testing.T) {
	type testContext struct {
		Field string
	}
	config := Config[testContext]{Context: testContext{Field: "test"}}
	service := newBaseServiceTest(config)

	result := service.Config()
	if !reflect.DeepEqual(*result, config) {
		t.Errorf("Config mismatch. Got %v, want %v", *result, config)
	}
}

func TestSetupCommandLineFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		exitCode int
	}{
		{"Version flag", []string{"-v"}, true, 0},
		{"No config", []string{}, true, 2},
		{"Too many args", []string{"config1", "config2"}, true, 2},
		{"Valid config", []string{"config.json"}, false, 0},
		{"Print config", []string{"-p", "config.json"}, false, 0},
		{"Test config", []string{"-t", "config.json"}, false, 0},
		{"Enable template", []string{"-T", "config.json"}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlagsAndArgs()
			service := newBaseServiceTest(Config[struct{}]{})
			os.Args = append(os.Args, tt.args...)

			err := service.setupCommandLineFlags()

			if (err != nil) != tt.wantErr {
				t.Errorf("setupCommandLineFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				code, ok := errkit.ExitCode(err)
				if !ok {
					t.Errorf("Expected exitError, got %T", err)
				} else if code != tt.exitCode {
					t.Errorf("Expected exit code %d, got %d", tt.exitCode, code)
				}
			}

			if !tt.wantErr && len(tt.args) > 0 && !strings.HasPrefix(tt.args[0], "-") {
				if service.flags.source != tt.args[0] {
					t.Errorf("Expected source %s, got %s", tt.args[0], service.flags.source)
				}
			}
		})
	}
}
func TestSetupConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		enableTemplate bool
		wantErr        bool
	}{
		{
			name:           "Valid JSON config",
			configContent:  `{"Context": {"Field": "test"}, "Components": []}`,
			enableTemplate: false,
			wantErr:        false,
		},
		{
			name:           "Invalid JSON config",
			configContent:  `{"Context": {"Field": "test"`,
			enableTemplate: false,
			wantErr:        true,
		},
		{
			name:           "Valid config with template placeholder in Components",
			configContent:  `{"Context": {"Field": "test"}, "Components": [{"Name": "TestComponent", "UUID": "test-uuid"}]}`,
			enableTemplate: true,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlagsAndArgs()
			service := newBaseServiceTest(Config[struct{ Field string }]{})
			service.flags.source = "test-config.json"
			service.flags.enableTemplate = tt.enableTemplate

			// Create a temporary file with the test config content
			tmpfile, err := os.CreateTemp("", "test-config-*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.configContent)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			service.flags.source = tmpfile.Name()

			err = service.setupConfig()

			if (err != nil) != tt.wantErr {
				t.Errorf("setupConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if service.config.Context.Field != "test" {
					t.Errorf("Expected Context.Field value 'test', got %s", service.config.Context.Field)
				}

				if tt.enableTemplate && len(service.config.Components) > 0 {
					if service.config.Components[0].UUID != "test-uuid" {
						t.Errorf("Expected Component UUID 'test-uuid', got %s", service.config.Components[0].UUID)
					}
				}
			}
		})
	}
}

func TestUninit(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{}
	service.components.AddComponent("test", mockComponent)

	// 模拟初始化过程
	service.components.Init(context.Background())

	err := service.Uninit(context.Background())

	if err != nil {
		t.Errorf("Uninit() error = %v, want nil", err)
	}

	if !mockComponent.uninitCalled {
		t.Error("Component Uninit was not called")
	}
}

func TestStart(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{}
	service.components.AddComponent("test", mockComponent)

	// 模拟初始化过程
	service.components.Init(context.Background())

	err := service.Start(context.Background())

	if err != nil {
		t.Errorf("Start() error = %v, want nil", err)
	}

	if !mockComponent.startCalled {
		t.Error("Component Start was not called")
	}
}

// mockComponent 用于测试的模拟组件
type mockComponent struct {
	component.BaseComponent[struct{}]
	uuid           string
	initCalled     bool
	uninitCalled   bool
	startCalled    bool
	shutdownCalled bool
	setupCalled    bool
	setupContainer component.Container
	setupConfig    component.Config
	logger         *slog.Logger
}

func (m *mockComponent) Init(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *mockComponent) Uninit(ctx context.Context) error {
	m.uninitCalled = true
	return nil
}

func (m *mockComponent) Start(ctx context.Context) error {
	m.startCalled = true
	return nil
}

func (m *mockComponent) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	return nil
}

func (m *mockComponent) Setup(container component.Container, config component.Config) error {
	m.setupCalled = true
	m.setupContainer = container
	m.setupConfig = config
	m.logger = slog.Default().With("component", m.String())
	return nil
}

func (m *mockComponent) String() string {
	return "mockComponent"
}

func (m *mockComponent) Logger() *slog.Logger {
	if m.logger == nil {
		m.logger = slog.Default().With("component", m.String())
	}
	return m.logger
}
