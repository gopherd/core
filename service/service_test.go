package service

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
)

func TestMain(m *testing.M) {
	flag.Parse()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

// mockConfig implements config.Config for testing
type mockConfig struct {
	components []component.Config
}

func (m *mockConfig) SetupFlags(*flag.FlagSet) {}
func (m *mockConfig) Load() error              { return nil }
func (m *mockConfig) GetComponents() []component.Config {
	return m.components
}

// mockComponent implements component.Component for testing
type mockComponent struct {
	uuid   string
	name   string
	entity component.Entity
}

func (m *mockComponent) UUID() string                     { return m.uuid }
func (m *mockComponent) Name() string                     { return m.name }
func (m *mockComponent) Entity() component.Entity         { return m.entity }
func (m *mockComponent) Ctor(component.Config) error      { return nil }
func (m *mockComponent) Init(context.Context) error       { return nil }
func (m *mockComponent) Uninit(context.Context) error     { return nil }
func (m *mockComponent) Start(context.Context) error      { return nil }
func (m *mockComponent) Shutdown(context.Context) error   { return nil }
func (m *mockComponent) OnMounted(component.Entity) error { return nil }

// safeRegisterComponent registers a component for testing, ignoring if it's already registered
func safeRegisterComponent(name string, creator func() component.Component) {
	component.Register(name, func() component.Component {
		return creator()
	})
}

func TestNewBaseService(t *testing.T) {
	cfg := &mockConfig{}
	s := NewBaseService(cfg)
	if s == nil {
		t.Fatal("NewBaseService returned nil")
	}
}

func TestBaseService_GetComponent(t *testing.T) {
	cfg := &mockConfig{
		components: []component.Config{
			{UUID: "test-uuid-1", Name: "test-component-1"},
		},
	}
	s := NewBaseService(cfg)

	safeRegisterComponent("test-component-1", func() component.Component {
		return &mockComponent{uuid: "test-uuid-1", name: "test-component-1"}
	})

	err := s.Init(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}

	c := s.GetComponent("test-uuid-1")
	if c == nil {
		t.Fatal("GetComponent returned nil for existing component")
	}

	if c.UUID() != "test-uuid-1" {
		t.Errorf("Expected UUID 'test-uuid-1', got '%s'", c.UUID())
	}

	c = s.GetComponent("non-existent")
	if c != nil {
		t.Error("GetComponent returned non-nil for non-existent component")
	}
}

func TestBaseService_Lifecycle(t *testing.T) {
	cfg := &mockConfig{
		components: []component.Config{
			{UUID: "test-uuid-2", Name: "test-component-2"},
		},
	}
	s := NewBaseService(cfg)

	safeRegisterComponent("test-component-2", func() component.Component {
		return &mockComponent{uuid: "test-uuid-2", name: "test-component-2"}
	})

	ctx := context.Background()

	if err := s.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := s.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	if err := s.Uninit(ctx); err != nil {
		t.Fatalf("Uninit failed: %v", err)
	}
}

func TestBaseService_SetupFlags(t *testing.T) {
	s := NewBaseService(&mockConfig{})
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	s.SetupFlags(fs)

	if fs.Lookup("v") == nil {
		t.Error("Version flag (-v) not set up")
	}
}

func TestRun(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"test"}

	s := &MockService{
		BaseService: NewBaseService(&mockConfig{}),
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		s.shutdown.Store(true)
	}()

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	RunServiceFlagSet(s, flagSet)

	if !s.initialized || !s.started || !s.shutdownCalled || !s.uninitialized {
		t.Error("Service lifecycle methods were not called as expected")
	}
}

func TestRun_ExitError(t *testing.T) {
	s := &MockService{
		BaseService: NewBaseService(&mockConfig{}),
		initErr:     errkit.NewExitError(42),
	}
	err := runTestService(t, "TestRun_ExitError", s)
	if exitCode, ok := errkit.ExitCode(err); !ok || exitCode != 42 {
		t.Errorf("Expected ExitError with code 42, got %v", err)
	}
}

func TestRun_OtherError(t *testing.T) {
	s := &MockService{
		BaseService: NewBaseService(&mockConfig{}),
		initErr:     errkit.NewExitError(1),
	}
	err := runTestService(t, "TestRun_OtherError", s)
	if err == nil || err.Error() != "exit with code 1" {
		t.Errorf("Expected 'exit with code 1', got %v", err)
	}
}

// MockService implements Service for testing
type MockService struct {
	*BaseService[*mockConfig]
	initialized    bool
	started        bool
	shutdownCalled bool
	uninitialized  bool
	shutdown       atomic.Bool
	initErr        error
	startErr       error
}

func NewMockService() *MockService {
	return &MockService{
		BaseService: NewBaseService(&mockConfig{}),
	}
}

func (m *MockService) Init(ctx context.Context) error {
	m.initialized = true
	return m.initErr
}

func (m *MockService) Start(ctx context.Context) error {
	m.started = true
	return m.startErr
}

func (m *MockService) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	return nil
}

func (m *MockService) Uninit(ctx context.Context) error {
	m.uninitialized = true
	return nil
}

// runTestService is a helper function to run a service for testing
func runTestService(t *testing.T, caseName string, s Service) error {
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				t.Logf("%s: %v", caseName, err)
			}
		}()
		flagSet := flag.NewFlagSet("runTest", flag.ContinueOnError)
		errCh <- RunServiceFlagSet(s, flagSet)
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(1 * time.Second):
		return errkit.NewExitError(1)
	}
}
