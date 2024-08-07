// Package service provides a framework for creating and managing service processes.
package service

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/builder"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/config"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/lifecycle"
)

// Service represents a process with lifecycle management and component handling.
type Service interface {
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
	flags struct {
		version bool
	}
	versionFunc func()
	config      atomic.Value
	components  *component.Manager
}

// NewBaseService creates a new BaseService with the given configuration.
func NewBaseService[Config config.Config](cfg Config) *BaseService[Config] {
	s := &BaseService[Config]{
		components:  component.NewManager(),
		versionFunc: builder.PrintInfo,
	}
	s.config.Store(cfg)
	return s
}

// SetVersionFunc sets the version function.
// The version function is called when the service is started with the -v flag.
// If the version function is not set, the default version function is used.
// And if set to nil, the version function is disabled.
func (s *BaseService[Config]) SetVersionFunc(f func()) {
	s.versionFunc = f
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
	flagSet.BoolVar(&s.flags.version, "v", false, "Print version information including build details")
}

// Init implements the Service Init method.
func (s *BaseService[Config]) Init(ctx context.Context) error {
	if s.flags.version {
		if s.versionFunc != nil {
			s.versionFunc()
		}
		return errkit.NewExitError(0)
	}

	cfg := s.Config()
	if err := cfg.Load(); err != nil {
		return err
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

// Run starts and manages the lifecycle of the given service.
func Run(s Service) {
	if err := run(s, flag.CommandLine); err != nil {
		if exitCode, ok := errkit.ExitCode(err); ok {
			os.Exit(exitCode)
		}
		os.Exit(1)
	}
}

func run(s Service, flagSet *flag.FlagSet) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	s.SetupFlags(flagSet)
	flagSet.Parse(os.Args[1:])

	defer func() {
		slog.Info("uninitializing service")
		if err := s.Uninit(context.Background()); err != nil {
			slog.Error("failed to uninitialize service", slog.Any("error", err))
		}
		slog.Info("service exited")
	}()
	if err := s.Init(context.Background()); err != nil {
		// If the error is an ExitError, return it directly without logging.
		if _, ok := errkit.ExitCode(err); ok {
			return err
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
