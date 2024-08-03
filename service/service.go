package service

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/config"
	"github.com/gopherd/core/erron"
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

	// SetState sets state of service
	SetState(state State) error
	// Init initializes the service
	Init(context.Context) error
	// Start starts the service
	Start(context.Context) error
	// Shutdown shutdowns the service
	Shutdown(context.Context) error

	// GetComponent returns a component by id
	GetComponent(id string) component.Component
}

// Run runs the service
func Run(s Service) {
	slog.Info("initializing service")
	if err := s.Init(context.Background()); err != nil {
		slog.Error("failed to initialize service", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("starting service")
	s.SetState(Running)
	if err := s.Start(context.Background()); err != nil {
		slog.Error("failed to start service", slog.Any("error", err))
		s.SetState(Closed)
		s.Shutdown(context.Background())
		os.Exit(1)
	}

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

	slog.Info("shutting down service")
	s.SetState(Closed)
	s.Shutdown(context.Background())
}

// BaseService implements Service
type BaseService[Self Service, Config config.Configurator] struct {
	self  Self
	name  string
	id    int
	state State

	config     atomic.Value
	components *component.Manager
}

// NewBaseService creates a BaseService
func NewBaseService[Self Service, Config config.Configurator](self Self, cfg Config) *BaseService[Self, Config] {
	s := &BaseService[Self, Config]{
		self:       self,
		components: component.NewManager(),
	}
	s.config.Store(cfg)
	return s
}

// New creates a service with configurator
func New[Config config.Configurator](cfg Config) Service {
	type server struct {
		*BaseService[*server, Config]
	}
	s := &server{}
	s.BaseService = NewBaseService(s, cfg)
	return s
}

// AddComponent returns a component
func (s *BaseService[Self, Config]) AddComponent(com component.Component) component.Component {
	return s.components.AddComponent(com)
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

// Init implements Service Init method
func (s *BaseService[Self, Config]) Init(ctx context.Context) error {
	// load config
	cfg := s.Config()
	if err := cfg.ParseArgs(os.Args); err != nil {
		return err
	}
	core := cfg.CoreConfig()
	s.id = core.ID
	s.name = core.Name

	// create components
	for _, c := range core.Components {
		creator := component.Lookup(c.Name)
		if creator == nil {
			return erron.Throwf("unknown component name: %q", c.Name)
		}
		com := creator()
		if com == nil {
			return erron.Throwf("create component %q error", c.UUID)
		}
		if err := com.OnCreated(s.self, c); err != nil {
			return erron.Throw(err)
		}
		if s.AddComponent(com) == nil {
			return erron.Throwf("duplicate component id: %q", c.UUID)
		}
	}

	return s.components.Init(ctx)
}

// Start implements Service Start method
func (s *BaseService[Self, Config]) Start(ctx context.Context) error {
	return s.components.Start(ctx)
}

// Shutdown implements Service Shutdown method
func (s *BaseService[Self, Config]) Shutdown(ctx context.Context) error {
	s.components.Shutdown(ctx)
	return nil
}
