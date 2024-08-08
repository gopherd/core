// Package service provides a framework for creating and managing service processes.
package service

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/gopherd/core/builder"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/lifecycle"
)

// Service represents a process with lifecycle management and component handling.
type Service interface {
	lifecycle.Lifecycle

	// GetComponent returns a component by its UUID.
	GetComponent(uuid string) component.Component
}

// BaseService implements the Service interface.
type BaseService[T any] struct {
	flags struct {
		version  bool   // print version information
		source   string // config source (file path, HTTP URL, or '-' for stdin)
		output   string // config output (file path or '-' for stdout)
		test     bool   // test the config for validity
		template bool   // enable template parsing for components config
	}
	versionFunc func()

	config     Config[T]
	components *component.Group
}

// NewBaseService creates a new BaseService with the given configuration.
func NewBaseService[T any](context T) *BaseService[T] {
	return &BaseService[T]{
		versionFunc: builder.PrintInfo,
		config:      Config[T]{Context: context},
		components:  component.NewGroup(),
	}
}

// SetVersionFunc sets the version function.
// The version function is called when the service is started with the -v flag.
// If the version function is not set, the default version function is used.
// And if set to nil, the version function is disabled.
func (s *BaseService[T]) SetVersionFunc(f func()) {
	s.versionFunc = f
}

// GetComponent returns a component by its UUID.
func (s *BaseService[T]) GetComponent(uuid string) component.Component {
	return s.components.GetComponent(uuid)
}

// Config returns the current configuration.
func (s *BaseService[T]) Config() *Config[T] {
	return &s.config
}

// setupCommandLineFlags sets command-line arguments for the service.
func (s *BaseService[T]) setupCommandLineFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-v] [-c <path>] [-o <path>] [-t] [T]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.BoolVar(&s.flags.version, "v", false, "Print version information including build details")
	flag.StringVar(&s.flags.source, "c", "", "Specify the config source (file path, HTTP URL, or '-' for stdin)")
	flag.StringVar(&s.flags.output, "o", "", "Specify the config output (file path or '-' for stdout) and exit")
	flag.BoolVar(&s.flags.test, "t", false, "Test the config for validity and exit")
	flag.BoolVar(&s.flags.template, "T", false, "Enable template processing for components config")

	flag.Parse()
}

func (s *BaseService[T]) processConfig() (err error) {
	defer func() {
		if s.flags.test {
			if err != nil {
				fmt.Println("Config test failed: ", err)
			} else {
				fmt.Println("Config test successful")
				// Exit after test
				err = errkit.NewExitError(0)
			}
		}
		if err == nil && s.flags.source == "" {
			// Config source is required unless testing or outputting
			fmt.Fprintf(os.Stderr, "No config source specified!\n\n")
			flag.Usage()
			err = errkit.NewExitError(2)
		}
	}()

	if err = s.config.load(s.flags.source); err != nil {
		return
	}

	if s.flags.template {
		if err = s.config.processTemplate(); err != nil {
			err = fmt.Errorf("failed to process template: %w", err)
			return
		}
	}

	if s.flags.output != "" {
		if err = s.config.output(s.flags.output); err != nil {
			err = fmt.Errorf("failed to output config: %w", err)
			return
		}
		// Exit after output
		return errkit.NewExitError(0)
	}

	return nil
}

// Init implements the Service Init method.
func (s *BaseService[T]) Init(ctx context.Context) error {
	s.setupCommandLineFlags()

	if s.flags.version {
		if s.versionFunc != nil {
			s.versionFunc()
		}
		return errkit.NewExitError(0)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := s.processConfig(); err != nil {
		return err
	}

	for _, c := range s.config.Components {
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
func (s *BaseService[T]) Uninit(ctx context.Context) error {
	return s.components.Uninit(ctx)
}

// Start implements the Service Start method.
func (s *BaseService[T]) Start(ctx context.Context) error {
	return s.components.Start(ctx)
}

// Shutdown implements the Service Shutdown method.
func (s *BaseService[T]) Shutdown(ctx context.Context) error {
	return s.components.Shutdown(ctx)
}

// Run is a shortcut for running a service with a default configuration.
func Run() {
	type context map[string]any
	if err := RunService(NewBaseService(Config[context]{Context: context{}})); err != nil {
		if exitCode, ok := errkit.ExitCode(err); ok {
			os.Exit(exitCode)
		}
		os.Exit(1)
	}
}

// RunService starts and manages the lifecycle of the given service.
// If the service returns an error, the program exits with the error code or 1.
func RunService(s Service) error {
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
	return err
}
