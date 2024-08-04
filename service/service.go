package service

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/buildinfo"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/config"
	"github.com/gopherd/core/event"
	"github.com/gopherd/core/lifecycle"
)

// State represents service state
type State int

const (
	Closed   State = iota // Closed service
	Running               // Running service
	Stopping              // Stopping service
)

func (state State) String() string {
	switch state {
	case Closed:
		return "Closed"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	default:
		return "Unknown(" + strconv.Itoa(int(state)) + ")"
	}
}

// Metadata represents metadata of service
type Metadata interface {
	// ID returns id of service
	ID() int
	// Name of service
	Name() string
	// Busy reports whether the service is busy
	Busy() bool
	// State returns state of service
	State() State
}

// Service represents a process
type Service interface {
	Metadata
	event.Dispatcher[reflect.Type]
	lifecycle.Lifecycle

	// SetState sets state of service
	SetState(state State) error

	// SetFlags sets command-line flags
	SetFlags(flagSet *flag.FlagSet)

	// GetComponent returns a component by id
	GetComponent(id string) component.Component
}

// BaseService implements Service
type BaseService[Self Service, Config config.Config] struct {
	event.Dispatcher[reflect.Type]
	self  Self
	name  string
	id    int
	state State

	flags struct {
		version bool
	}

	config     atomic.Value
	components *component.Manager
}

// NewBaseService creates a BaseService
func NewBaseService[Self Service, Config config.Config](self Self, cfg Config) *BaseService[Self, Config] {
	s := &BaseService[Self, Config]{
		Dispatcher: event.NewDispatcher[reflect.Type](true),
		self:       self,
		components: component.NewManager(),
	}
	s.config.Store(cfg)
	return s
}

// New creates a service with configurator
func New[Config config.Config](cfg Config) Service {
	type server struct {
		*BaseService[*server, Config]
	}
	s := &server{}
	s.BaseService = NewBaseService(s, cfg)
	return s
}

// GetComponent returns a component by uuid
func (s *BaseService[Self, Config]) GetComponent(uuid string) component.Component {
	return s.components.GetComponent(uuid)
}

// Name implements Service Name method
func (s *BaseService[Self, Config]) Name() string {
	return s.name
}

// ID implements Service ID method
func (s *BaseService[Self, Config]) ID() int {
	return s.id
}

// Busy implements Service Busy method
func (s *BaseService[Self, Config]) Busy() bool {
	return false
}

// State returns state of service
func (s *BaseService[Self, Config]) State() State {
	return s.state
}

// SetState implements Service SetState method
func (s *BaseService[Self, Config]) SetState(state State) error {
	s.state = state
	return nil
}

// Config returns current config
func (s *BaseService[Self, Config]) Config() Config {
	return s.config.Load().(Config)
}

// SetFlags implements Service SetFlags method
func (s *BaseService[Self, Config]) SetFlags(flagSet *flag.FlagSet) {
	s.Config().SetFlags(flagSet)
	flagSet.BoolVar(&s.flags.version, "v", false, "Print version information")
}

// Init implements Service Init method
func (s *BaseService[Self, Config]) Init(ctx context.Context) error {
	if s.flags.version {
		buildinfo.PrintVersion()
		return &ExitError{Code: 0}
	}

	// load config
	cfg := s.Config()
	if exit, err := cfg.Load(); err != nil {
		return err
	} else if exit {
		return &ExitError{Code: 0}
	}
	core := cfg.CoreConfig()
	s.id = core.ID
	s.name = core.Name

	// create components
	for _, c := range core.Components {
		creator := component.Lookup(c.Name)
		if creator == nil {
			return fmt.Errorf("unknown component name: %q", c.Name)
		}
		com := creator()
		if com == nil {
			return fmt.Errorf("create component %q error", c.UUID)
		}
		if err := com.OnCreated(s.self, c); err != nil {
			return fmt.Errorf("create component %q error: %w", c.UUID, err)
		}
		if s.components.AddComponent(com) == nil {
			return fmt.Errorf("duplicate component id: %q", c.UUID)
		}
	}

	return s.components.Init(ctx)
}

// Uninit implements Service Uninit method
func (s *BaseService[Self, Config]) Uninit(ctx context.Context) error {
	return s.components.Uninit(ctx)
}

// Start implements Service Start method
func (s *BaseService[Self, Config]) Start(ctx context.Context) error {
	return s.components.Start(ctx)
}

// Shutdown implements Service Shutdown method
func (s *BaseService[Self, Config]) Shutdown(ctx context.Context) error {
	return s.components.Shutdown(ctx)
}

type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit with code %d", e.Code)
}

// Run runs the service
func Run(s Service) {
	if err := run(s); err != nil {
		os.Exit(1)
	}
}

func run(s Service) error {
	// discard log output below warn level before service initialized
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	// parse command-line flags
	s.SetFlags(flag.CommandLine)
	flag.CommandLine.Parse(os.Args[1:])

	// initialize service and defer uninitialize
	defer func() {
		slog.Info("uninitializing service")
		if err := s.Uninit(context.Background()); err != nil {
			slog.Error("failed to uninitialize service", slog.Any("error", err))
		}
		slog.Info("service exited")
	}()
	if err := s.Init(context.Background()); err != nil {
		if e, ok := err.(*ExitError); ok {
			os.Exit(e.Code)
		}
		slog.Error("failed to initialize service", slog.Any("error", err))
		return err
	}

	// start service and defer shutdown
	slog.Info("starting service")
	defer func() {
		slog.Info("shutting down service")
		s.SetState(Closed)
		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown service", slog.Any("error", err))
		}
	}()
	s.SetState(Running)
	err := s.Start(context.Background())
	if err != nil {
		slog.Error("failed to start service", slog.Any("error", err))
	}

	// wait for service to stop
	slog.Info("stopping service")
	s.SetState(Stopping)
	if s.Busy() {
		slog.Info("waiting for service to stop")
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for range ticker.C {
			if !s.Busy() {
				break
			}
		}
	}

	return err
}
