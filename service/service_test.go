package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
)

func TestMain(m *testing.M) {
	// Disable logging during tests
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	os.Exit(m.Run())
}

func newBaseServiceTest[T any](config Config[T]) *BaseService[T] {
	s := NewBaseService(config)
	s.stderr = io.Discard
	s.stdout = io.Discard
	s.flagSet = flag.NewFlagSet("test", flag.ContinueOnError)
	s.versionFunc = func() {}
	return s
}

// resetFlagsAndArgs resets os.Args
func resetFlagsAndArgs() {
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

	if service.stdout == nil || service.stderr == nil {
		t.Error("stdout or stderr is nil")
	}
}

func TestSetVersionFunc(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})

	t.Run("Set and call version function", func(t *testing.T) {
		called := false
		versionFunc := func() { called = true }

		service.SetVersionFunc(versionFunc)
		service.versionFunc()

		if !called {
			t.Error("Version function was not called")
		}
	})

	t.Run("Set version function to nil", func(t *testing.T) {
		service.SetVersionFunc(nil)
		if service.versionFunc != nil {
			t.Error("Version function was not set to nil")
		}
	})
}

func TestGetComponent(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{uuid: "test-uuid"}
	service.components.AddComponent("test-uuid", mockComponent)

	t.Run("Get existing component", func(t *testing.T) {
		result := service.GetComponent("test-uuid")
		if result != mockComponent {
			t.Errorf("GetComponent returned wrong component. Got %v, want %v", result, mockComponent)
		}
	})

	t.Run("Get non-existent component", func(t *testing.T) {
		result := service.GetComponent("non-existent")
		if result != nil {
			t.Errorf("GetComponent should return nil for non-existent component, got %v", result)
		}
	})
}

func TestLogger(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	logger := service.Logger()
	if logger == nil {
		t.Error("Logger returned nil")
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
		check    func(*testing.T, *BaseService[struct{}])
	}{
		{
			name:     "Version command",
			args:     []string{"version"},
			wantErr:  true,
			exitCode: 0,
		},
		{
			name:     "No config",
			args:     []string{},
			wantErr:  true,
			exitCode: 2,
		},
		{
			name:     "Too many args",
			args:     []string{"config1", "config2"},
			wantErr:  true,
			exitCode: 2,
		},
		{
			name:    "Valid config",
			args:    []string{"config.json"},
			wantErr: false,
			check: func(t *testing.T, s *BaseService[struct{}]) {
				if s.flags.source != "config.json" {
					t.Errorf("Expected source config.json, got %s", s.flags.source)
				}
			},
		},
		{
			name:    "Print config",
			args:    []string{"-p", "config.json"},
			wantErr: false,
			check: func(t *testing.T, s *BaseService[struct{}]) {
				if !s.flags.printConfig {
					t.Error("Print config flag not set")
				}
			},
		},
		{
			name:    "Test config",
			args:    []string{"-t", "config.json"},
			wantErr: false,
			check: func(t *testing.T, s *BaseService[struct{}]) {
				if !s.flags.testConfig {
					t.Error("Test config flag not set")
				}
			},
		},
		{
			name:    "Enable template",
			args:    []string{"-T", "config.json"},
			wantErr: false,
			check: func(t *testing.T, s *BaseService[struct{}]) {
				if !s.flags.enableTemplate {
					t.Error("Enable template flag not set")
				}
			},
		},
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

			if tt.check != nil {
				tt.check(t, service)
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
		check          func(*testing.T, *BaseService[struct{ Field string }])
	}{
		{
			name:           "Valid JSON config",
			configContent:  `{"Context": {"Field": "test"}, "Components": []}`,
			enableTemplate: false,
			wantErr:        false,
			check: func(t *testing.T, s *BaseService[struct{ Field string }]) {
				if s.config.Context.Field != "test" {
					t.Errorf("Expected Context.Field value 'test', got %s", s.config.Context.Field)
				}
				if len(s.config.Components) != 0 {
					t.Errorf("Expected 0 components, got %d", len(s.config.Components))
				}
			},
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
			check: func(t *testing.T, s *BaseService[struct{ Field string }]) {
				if s.config.Context.Field != "test" {
					t.Errorf("Expected Context.Field value 'test', got %s", s.config.Context.Field)
				}
				if len(s.config.Components) != 1 {
					t.Errorf("Expected 1 component, got %d", len(s.config.Components))
				}
				if s.config.Components[0].UUID != "test-uuid" {
					t.Errorf("Expected Component UUID 'test-uuid', got %s", s.config.Components[0].UUID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlagsAndArgs()
			service := newBaseServiceTest(Config[struct{ Field string }]{})
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

			if !tt.wantErr && tt.check != nil {
				tt.check(t, service)
			}
		})
	}
}

func TestSetupComponents(t *testing.T) {
	tests := []struct {
		name          string
		components    []component.Config
		expectedError bool
	}{
		{
			name: "Valid components",
			components: []component.Config{
				{Name: "TestComponent1", UUID: "uuid1"},
				{Name: "TestComponent2", UUID: "uuid2"},
			},
			expectedError: false,
		},
		{
			name: "Duplicate UUID",
			components: []component.Config{
				{Name: "TestComponent1", UUID: "uuid1"},
				{Name: "TestComponent2", UUID: "uuid1"},
			},
			expectedError: true,
		},
		{
			name: "Invalid component name",
			components: []component.Config{
				{Name: "NonExistentComponent", UUID: "uuid1"},
			},
			expectedError: true,
		},
	}

	// Register mock components
	component.Register("TestComponent1", func() component.Component { return &mockComponent{} })
	component.Register("TestComponent2", func() component.Component { return &mockComponent{} })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newBaseServiceTest(Config[struct{}]{
				Components: tt.components,
			})

			_, err := service.setupComponents()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		args           []string
		expectedError  bool
		expectedExit   int
		checkComponent bool
	}{
		{
			name:           "Valid init",
			configContent:  `{"Context": {}, "Components": [{"Name": "TestComponent", "UUID": "test-uuid"}]}`,
			args:           []string{"-"},
			checkComponent: true,
		},
		{
			name:          "Print config",
			configContent: `{"Context": {}, "Components": []}`,
			args:          []string{"-p", "-"},
			expectedError: true,
			expectedExit:  0,
		},
		{
			name:          "Test config",
			configContent: `{"Context": {}, "Components": []}`,
			args:          []string{"-t", "-"},
			expectedError: true,
			expectedExit:  0,
		},
		{
			name:          "Invalid config",
			configContent: `{"Context": {}, "Components": [{"Name": "InvalidComponent", "UUID": "test-uuid"}]}`,
			args:          []string{"-"},
			expectedError: true,
			expectedExit:  2,
		},
	}

	// Register mock component
	component.Register("TestComponent", func() component.Component { return &mockComponent{} })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlagsAndArgs()
			os.Args = append(os.Args, tt.args...)

			// Create a new BaseService with a buffer for stdin
			var stdin bytes.Buffer
			stdin.WriteString(tt.configContent)
			service := newBaseServiceTest(Config[struct{}]{})
			service.stdin = &stdin

			err := service.Init(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if code, ok := errkit.ExitCode(err); !ok || code != tt.expectedExit {
					t.Errorf("Expected exit code %d, got %v", tt.expectedExit, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.checkComponent {
				component := service.GetComponent("test-uuid")
				if component == nil {
					t.Error("Expected component not found")
				}
				mockComp, ok := component.(*mockComponent)
				if !ok {
					t.Error("Component is not a mockComponent")
				} else if !mockComp.initCalled {
					t.Error("Component Init was not called")
				}
			}
		})
	}
}

func TestUninit(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{}
	service.components.AddComponent("test", mockComponent)

	// Simulate initialization process
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

	// Simulate initialization process
	service.components.Init(context.Background())

	err := service.Start(context.Background())

	if err != nil {
		t.Errorf("Start() error = %v, want nil", err)
	}

	if !mockComponent.startCalled {
		t.Error("Component Start was not called")
	}
}

func TestShutdown(t *testing.T) {
	service := newBaseServiceTest(Config[struct{}]{})
	mockComponent := &mockComponent{}
	service.components.AddComponent("test", mockComponent)

	// Simulate initialization and start processes
	service.components.Init(context.Background())
	service.components.Start(context.Background())

	err := service.Shutdown(context.Background())

	if err != nil {
		t.Errorf("Shutdown() error = %v, want nil", err)
	}

	if !mockComponent.shutdownCalled {
		t.Error("Component Shutdown was not called")
	}
}

func TestRunOptions(t *testing.T) {
	t.Run("WithEncoder", func(t *testing.T) {
		encoder := func(v any) ([]byte, error) {
			return json.Marshal(v)
		}
		opt := WithEncoder(encoder)
		options := &runOptions{}
		opt(options)
		if options.encoder == nil {
			t.Error("Encoder was not set")
		}
	})

	t.Run("WithDecoder", func(t *testing.T) {
		decoder := func(data []byte, v any) error {
			return json.Unmarshal(data, v)
		}
		opt := WithDecoder(decoder)
		options := &runOptions{}
		opt(options)
		if options.decoder == nil {
			t.Error("Decoder was not set")
		}
	})
}

func TestRunService(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(*mockService)
		expectedError string
	}{
		{
			name: "Successful run",
			mockBehavior: func(m *mockService) {
				m.initFunc = func(context.Context) error { return nil }
				m.startFunc = func(context.Context) error { return nil }
				m.shutdownFunc = func(context.Context) error { return nil }
				m.uninitFunc = func(context.Context) error { return nil }
			},
			expectedError: "",
		},
		{
			name: "Init error",
			mockBehavior: func(m *mockService) {
				m.initFunc = func(context.Context) error { return errors.New("init error") }
			},
			expectedError: "init error",
		},
		{
			name: "Start error",
			mockBehavior: func(m *mockService) {
				m.initFunc = func(context.Context) error { return nil }
				m.startFunc = func(context.Context) error { return errors.New("start error") }
			},
			expectedError: "start error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockService{
				logger: slog.Default(),
			}
			tt.mockBehavior(mockService)

			err := RunService(mockService)

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
				}
			}
		})
	}
}

// mockComponent is a mock implementation of component.Component for testing
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
	mu             sync.Mutex
}

func (m *mockComponent) Init(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.initCalled = true
	return nil
}

func (m *mockComponent) Uninit(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.uninitCalled = true
	return nil
}

func (m *mockComponent) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startCalled = true
	return nil
}

func (m *mockComponent) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shutdownCalled = true
	return nil
}

func (m *mockComponent) Setup(container component.Container, config component.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.logger == nil {
		m.logger = slog.Default().With("component", m.String())
	}
	return m.logger
}

// mockService is a mock implementation of Service for testing
type mockService struct {
	initFunc     func(context.Context) error
	uninitFunc   func(context.Context) error
	startFunc    func(context.Context) error
	shutdownFunc func(context.Context) error
	logger       *slog.Logger
}

func (m *mockService) Init(ctx context.Context) error {
	if m.initFunc != nil {
		return m.initFunc(ctx)
	}
	return nil
}

func (m *mockService) Uninit(ctx context.Context) error {
	if m.uninitFunc != nil {
		return m.uninitFunc(ctx)
	}
	return nil
}

func (m *mockService) Start(ctx context.Context) error {
	if m.startFunc != nil {
		return m.startFunc(ctx)
	}
	return nil
}

func (m *mockService) Shutdown(ctx context.Context) error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc(ctx)
	}
	return nil
}

func (m *mockService) GetComponent(uuid string) component.Component {
	return nil
}

func (m *mockService) Logger() *slog.Logger {
	return m.logger
}
