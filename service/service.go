// Package service provides a framework for creating and managing service processes
// with support for configuration loading, lifecycle management, and component handling.
// It offers flexible configuration options, including file-based, HTTP-based, and
// stdin-based configuration loading, as well as template processing capabilities.
package service

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/gopherd/core/builder"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/container/pair"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/lifecycle"
)

// Service represents a process with lifecycle management and component handling capabilities.
type Service interface {
	lifecycle.Lifecycle
	component.Container
}

// BaseService implements the Service interface with a generic context type T.
type BaseService[T any] struct {
	flags struct {
		source         string // config source path, URL or "-" for stdin
		version        bool   // print version information and exit
		printConfig    bool   // output the config and exit
		testConfig     bool   // test the config for validity and exit
		enableTemplate bool   // enable template parsing for components config
	}
	versionFunc func()

	config     Config[T]
	components *component.Group
}

// NewBaseService creates a new BaseService with the given configuration.
func NewBaseService[T any](config Config[T]) *BaseService[T] {
	return &BaseService[T]{
		versionFunc: builder.PrintInfo,
		config:      config,
		components:  component.NewGroup(),
	}
}

// SetVersionFunc sets the version function to be called when the service is started with the -v flag.
// If set to nil, the version function is disabled.
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

// setupCommandLineFlags sets up and processes command-line arguments for the service.
func (s *BaseService[T]) setupCommandLineFlags() error {
	flag.Usage = func() {
		name := os.Args[0]
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <config>\n", name)
		fmt.Fprintf(os.Stderr, "       %s <path/to/file>   (Read configuration from file)\n", name)
		fmt.Fprintf(os.Stderr, "       %s <url>            (Read configuration from http)\n", name)
		fmt.Fprintf(os.Stderr, "       %s -                (Read configuration from stdin)\n", name)
		fmt.Fprintf(os.Stderr, "       %s -v               (Print version information)\n", name)
		fmt.Fprintf(os.Stderr, "       %s -p               (Print the configuration)\n", name)
		fmt.Fprintf(os.Stderr, "       %s -t               (Test the configuration for validity)\n", name)
		fmt.Fprintf(os.Stderr, "       %s -T               (Enable template processing for component configurations)\n", name)
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "       %s app.json\n", name)
		fmt.Fprintf(os.Stderr, "       %s http://example.com/app.json\n", name)
		fmt.Fprintf(os.Stderr, `       echo '{"Components":[]}' | %s -`+"\n", name)
		fmt.Fprintf(os.Stderr, "       %s -p app.json\n", name)
		fmt.Fprintf(os.Stderr, "       %s -t app.json\n", name)
		fmt.Fprintf(os.Stderr, "       %s -T app.json\n", name)
		fmt.Fprintf(os.Stderr, "       %s -p -T app.json\n", name)
		fmt.Fprintf(os.Stderr, "       %s -t -T app.json\n", name)
	}

	flag.BoolVar(&s.flags.version, "v", false, "")
	flag.BoolVar(&s.flags.printConfig, "p", false, "")
	flag.BoolVar(&s.flags.testConfig, "t", false, "")
	flag.BoolVar(&s.flags.enableTemplate, "T", false, "")

	flag.Parse()

	if s.flags.version {
		if s.versionFunc != nil {
			s.versionFunc()
		}
		return errkit.NewExitError(0)
	}

	if flag.NArg() == 0 || flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "No config source specified!\n\n")
		flag.Usage()
		return errkit.NewExitError(2)
	}
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments!\n\n")
		flag.Usage()
		return errkit.NewExitError(2)
	}
	s.flags.source = flag.Arg(0)

	return nil
}

// setupConfig loads and sets up the service configuration based on command-line flags.
func (s *BaseService[T]) setupConfig() error {
	if err := s.config.load(s.flags.source); err != nil {
		return err
	}
	if s.flags.enableTemplate {
		if err := s.config.processTemplate(s.flags.source); err != nil {
			return err
		}
	}
	return nil
}

// Init implements the Service Init method, setting up logging and initializing components.
func (s *BaseService[T]) Init(ctx context.Context) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := s.setupCommandLineFlags(); err != nil {
		return err
	}

	err := s.setupConfig()

	if s.flags.printConfig {
		if err != nil {
			return err
		}
		s.config.output()
		return errkit.NewExitError(0)
	}

	if err == nil {
		err = s.setupComponents()
	}

	if s.flags.testConfig {
		if err != nil {
			fmt.Printf("Config test failed: %e", err)
			err = errkit.NewExitError(2, err.Error())
		} else {
			fmt.Println("Config test successful")
			err = errkit.NewExitError(0)
		}
	}
	if err != nil {
		return err
	}

	return s.components.Init(ctx)
}

func (s *BaseService[T]) setupComponents() error {
	var components = make([]pair.Pair[component.Component, component.Config], 0, len(s.config.Components))
	for _, c := range s.config.Components {
		com, err := component.Create(c.Name)
		if err != nil {
			return err
		}
		if s.components.AddComponent(c.UUID, com) == nil {
			return fmt.Errorf("duplicate component uuid: %q", c.UUID)
		}
		components = append(components, pair.New(com, c))
	}
	for i := range components {
		if err := components[i].First.Setup(s, components[i].Second); err != nil {
			return fmt.Errorf("component %q setup error: %w", components[i].First.String(), err)
		}
	}
	return nil
}

// Uninit implements the Service Uninit method, uninitializing all components.
func (s *BaseService[T]) Uninit(ctx context.Context) error {
	return s.components.Uninit(ctx)
}

// Start implements the Service Start method, starting all components.
func (s *BaseService[T]) Start(ctx context.Context) error {
	return s.components.Start(ctx)
}

// Shutdown implements the Service Shutdown method, shutting down all components.
func (s *BaseService[T]) Shutdown(ctx context.Context) error {
	return s.components.Shutdown(ctx)
}

// Run is a convenience function for running a service with a default configuration.
// It creates and runs a BaseService with an empty context.
// This function always exits the program:
// - It exits with the error code if an error occurs.
// - It exits with code 0 if the service runs successfully.
// It is recommended to use this function unless you need to customize the Service
// or want to prevent the program from exiting.
func Run() {
	type context map[string]any
	if err := RunService(NewBaseService(Config[context]{Context: context{}})); err != nil {
		if exitCode, ok := errkit.ExitCode(err); ok {
			os.Exit(exitCode)
		}
		os.Exit(1)
	}
	os.Exit(0)
}

// RunService starts and manages the lifecycle of the given service.
// It handles initialization, starting, shutdown, and uninitialization of the service.
// This function returns any error encountered during the service lifecycle.
// Use this function if you need to run a custom Service implementation or
// if you want to handle errors without exiting the program.
func RunService(s Service) error {
	defer func() {
		slog.Info("uninitializing service")
		if err := s.Uninit(context.Background()); err != nil {
			slog.Error("failed to uninitialize service", slog.Any("error", err))
		}
		slog.Info("service exited")
	}()
	if err := s.Init(context.Background()); err != nil {
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
