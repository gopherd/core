// Package service provides a framework for creating and managing service processes.
package service

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/buildinfo"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/config"
	"github.com/gopherd/core/event"
	"github.com/gopherd/core/lifecycle"
)

// Service represents a process with lifecycle management and component handling.
type Service interface {
	event.Dispatcher[reflect.Type]
	lifecycle.Lifecycle

	// IsBusy reports whether the service is busy.
	IsBusy() bool

	// SetupFlags sets command-line flags for the service.
	SetupFlags(flagSet *flag.FlagSet)

	// GetComponent returns a component by its UUID.
	GetComponent(uuid string) component.Component
}

// BaseService implements the Service interface.
type BaseService[Config config.Config] struct {
	event.Dispatcher[reflect.Type]

	flags struct {
		version bool
	}

	config     atomic.Value
	components *component.Manager
}

// NewBaseService creates a new BaseService with the given configuration.
func NewBaseService[Config config.Config](cfg Config) *BaseService[Config] {
	s := &BaseService[Config]{
		Dispatcher: event.NewDispatcher[reflect.Type](true),
		components: component.NewManager(),
	}
	s.config.Store(cfg)
	return s
}

// GetComponent returns a component by its UUID.
func (s *BaseService[Config]) GetComponent(uuid string) component.Component {
	return s.components.GetComponent(uuid)
}

// IsBusy implements the Service IsBusy method.
func (s *BaseService[Config]) IsBusy() bool {
	return false
}

// Config returns the current configuration.
func (s *BaseService[Config]) Config() Config {
	return s.config.Load().(Config)
}

// SetupFlags implements the Service SetupFlags method.
func (s *BaseService[Config]) SetupFlags(flagSet *flag.FlagSet) {
	s.Config().SetupFlags(flagSet)
	flagSet.BoolVar(&s.flags.version, "v", false, "Print version information")
}

// Init implements the Service Init method.
func (s *BaseService[Config]) Init(ctx context.Context) error {
	if s.flags.version {
		buildinfo.PrintVersion()
		return &ExitError{Code: 0}
	}

	cfg := s.Config()
	if exit, err := cfg.Load(); err != nil {
		return err
	} else if exit {
		return &ExitError{Code: 0}
	}

	for _, c := range cfg.GetComponents() {
		creator := component.Lookup(c.Name)
		if creator == nil {
			return fmt.Errorf("unknown component name: %q", c.Name)
		}
		com := creator()
		if com == nil {
			return fmt.Errorf("create component %q error", c.UUID)
		}
		if err := com.Ctor(c); err != nil {
			return fmt.Errorf("create component %q error: %w", c.UUID, err)
		}
		if s.components.AddComponent(com) == nil {
			return fmt.Errorf("duplicate component id: %q", c.UUID)
		}
	}
	if err := s.components.OnMounted(s); err != nil {
		return fmt.Errorf("mount components error: %w", err)
	}

	return s.components.Init(ctx)
}

// Uninit implements the Service Uninit method.
func (s *BaseService[Config]) Uninit(ctx context.Context) error {
	return s.components.Uninit(ctx)
}

// Start implements the Service Start method.
func (s *BaseService[Config]) Start(ctx context.Context) error {
	return s.components.Start(ctx)
}

// Shutdown implements the Service Shutdown method.
func (s *BaseService[Config]) Shutdown(ctx context.Context) error {
	return s.components.Shutdown(ctx)
}

// ExitError represents an error that causes the service to exit.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit with code %d", e.Code)
}

// Run starts and manages the lifecycle of the given service.
func Run(s Service) {
	if err := run(s); err != nil {
		if exit := (*ExitError)(nil); errors.As(err, &exit) {
			os.Exit(exit.Code)
		}
		os.Exit(1)
	}
}

func run(s Service) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	s.SetupFlags(flag.CommandLine)
	flag.CommandLine.Parse(os.Args[1:])

	defer func() {
		slog.Info("uninitializing service")
		if err := s.Uninit(context.Background()); err != nil {
			slog.Error("failed to uninitialize service", slog.Any("error", err))
		}
		slog.Info("service exited")
	}()
	if err := s.Init(context.Background()); err != nil {
		// If the error is an ExitError, return it directly without logging.
		if exit := (*ExitError)(nil); errors.As(err, &exit) {
			return exit
		}
		slog.Error("failed to initialize service", slog.Any("error", err))
		return err
	}

	slog.Info("starting service")
	defer func() {
		slog.Info("shutting down service")
		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown service", slog.Any("error", err))
		}
	}()
	err := s.Start(context.Background())
	if err != nil {
		slog.Error("failed to start service", slog.Any("error", err))
	}

	slog.Info("stopping service")

	for i := 0; s.IsBusy() && i < 4; i++ {
		time.Sleep(time.Millisecond * time.Duration(1<<(i*2)))
	}

	if s.IsBusy() {
		slog.Info("waiting for service to stop")
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for range ticker.C {
			if !s.IsBusy() {
				break
			}
		}
	}

	return err
}
