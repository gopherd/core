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
		// print version information and exit
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

func (s *BaseService[T]) processConfig() (err error) {
	defer func() {
		if s.flags.testConfig {
			if err != nil {
				fmt.Println("Config test failed: ", err)
			} else {
				fmt.Println("Config test successful")
				// Exit after test
				err = errkit.NewExitError(0)
			}
		}
	}()

	if err = s.config.load(s.flags.source); err != nil {
		return
	}

	if s.flags.enableTemplate {
		if err = s.config.processTemplate(); err != nil {
			err = fmt.Errorf("failed to process template: %w", err)
			return
		}
	}

	if s.flags.printConfig {
		s.config.output()
		// Exit after output
		return errkit.NewExitError(0)
	}

	return nil
}

// Init implements the Service Init method.
func (s *BaseService[T]) Init(ctx context.Context) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := s.setupCommandLineFlags(); err != nil {
		return err
	}

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
