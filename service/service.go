// Package service provides a framework for creating and managing service processes
// with support for configuration loading, lifecycle management, and component handling.
// It offers flexible configuration options, including file-based, HTTP-based, and
// stdin-based configuration loading, as well as template processing capabilities.
package service

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gopherd/core/builder"
	"github.com/gopherd/core/component"
	"github.com/gopherd/core/container/pair"
	"github.com/gopherd/core/encoding"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/lifecycle"
)

// Service represents a process with lifecycle management and component handling capabilities.
type Service interface {
	lifecycle.Lifecycle
	component.Container

	// Logger returns the logger instance for the service.
	Logger() *slog.Logger
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
	flagSet     *flag.FlagSet
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	encoder     encoding.Encoder
	decoder     encoding.Decoder

	config     Config[T]
	components *component.Group
}

// NewBaseService creates a new BaseService with the given configuration.
func NewBaseService[T any](config Config[T]) *BaseService[T] {
	return &BaseService[T]{
		versionFunc: builder.PrintInfo,
		flagSet:     flag.CommandLine,
		stdin:       os.Stdin,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
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

// Logger returns the logger instance for the service.
func (s *BaseService[T]) Logger() *slog.Logger {
	return slog.Default()
}

// Config returns the current configuration.
func (s *BaseService[T]) Config() *Config[T] {
	return &s.config
}

// setupCommandLineFlags sets up and processes command-line arguments for the service.
func (s *BaseService[T]) setupCommandLineFlags() error {
	s.flagSet.Usage = func() {
		name := os.Args[0]
		var sb strings.Builder
		fmt.Fprintf(&sb, "Usage: %s [Options] <Config>\n", name)
		fmt.Fprintf(&sb, "       %s -h\n", name)
		fmt.Fprintf(&sb, "       %s -v\n", name)
		fmt.Fprintf(&sb, "\nConfig:\n")
		fmt.Fprintf(&sb, "       <path/to/file>   (Read configuration from file)\n")
		fmt.Fprintf(&sb, "       <url>            (Read configuration from http)\n")
		fmt.Fprintf(&sb, "       -                (Read configuration from stdin)\n")
		fmt.Fprintf(&sb, "\nOptions:\n")
		fmt.Fprintf(&sb, "       -p               (Print the configuration)\n")
		fmt.Fprintf(&sb, "       -t               (Test the configuration for validity)\n")
		fmt.Fprintf(&sb, "       -T               (Enable template processing for component configurations)\n")
		fmt.Fprintf(&sb, "\nExamples:\n")
		fmt.Fprintf(&sb, "       %s app.json\n", name)
		fmt.Fprintf(&sb, "       %s http://example.com/app.json\n", name)
		fmt.Fprintf(&sb, `       echo '{"Components":[]}' | %s -`+"\n", name)
		fmt.Fprintf(&sb, "       %s -p app.json\n", name)
		fmt.Fprintf(&sb, "       %s -t app.json\n", name)
		fmt.Fprintf(&sb, "       %s -T app.json\n", name)
		fmt.Fprintf(&sb, "       %s -p -T app.json\n", name)
		fmt.Fprintf(&sb, "       %s -t -T app.json\n", name)
		fmt.Fprint(s.stderr, sb.String())
	}

	s.flagSet.BoolVar(&s.flags.version, "v", false, "")
	s.flagSet.BoolVar(&s.flags.printConfig, "p", false, "")
	s.flagSet.BoolVar(&s.flags.testConfig, "t", false, "")
	s.flagSet.BoolVar(&s.flags.enableTemplate, "T", false, "")

	if err := s.flagSet.Parse(os.Args[1:]); err != nil {
		return errkit.NewExitError(2, err.Error())
	}

	if s.flags.version {
		if s.versionFunc != nil {
			s.versionFunc()
		}
		return errkit.NewExitError(0)
	}

	if s.flagSet.NArg() == 0 || s.flagSet.Arg(0) == "" {
		fmt.Fprintf(s.stderr, "No config source specified!\n\n")
		s.flagSet.Usage()
		return errkit.NewExitError(2)
	}
	if s.flagSet.NArg() > 1 {
		fmt.Fprintf(s.stderr, "Too many arguments!\n\n")
		s.flagSet.Usage()
		return errkit.NewExitError(2)
	}
	s.flags.source = s.flagSet.Arg(0)

	return nil
}

// setupConfig loads and sets up the service configuration based on command-line flags.
func (s *BaseService[T]) setupConfig() error {
	if err := s.config.load(s.stdin, s.decoder, s.flags.source); err != nil {
		return err
	}
	if err := s.config.processTemplate(s.flags.enableTemplate, s.flags.source); err != nil {
		return err
	}
	return nil
}

// Init implements the Service Init method, setting up logging and initializing components.
func (s *BaseService[T]) Init(ctx context.Context) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(s.stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := s.setupCommandLineFlags(); err != nil {
		return err
	}

	err := s.setupConfig()

	if s.flags.printConfig {
		if err != nil {
			fmt.Fprintf(s.stderr, "Config test failed: %v\n", err)
			return errkit.NewExitError(2, err.Error())
		}
		s.config.output(s.stdout, s.stderr, s.encoder)
		return errkit.NewExitError(0)
	}

	if err == nil {
		err = s.setupComponents()
	}

	if s.flags.testConfig {
		if err != nil {
			fmt.Fprintf(s.stderr, "Config test failed: %v\n", err)
			err = errkit.NewExitError(2, err.Error())
		} else {
			fmt.Fprintln(s.stdout, "Config test successful")
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
			return errkit.NewExitError(2, fmt.Sprintf("failed to create component %q: %v", c.Name, err))
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

type runOptions struct {
	encoder encoding.Encoder
	decoder encoding.Decoder
}

// apply applies the options to the given options.
func (o *runOptions) apply(opts []RunOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// RunOption is a functional option for configuring the Run function.
type RunOption func(*runOptions)

// WithEncoder sets the encoder function for the Run function.
func WithEncoder(encoder encoding.Encoder) RunOption {
	return func(o *runOptions) {
		o.encoder = encoder
	}
}

// WithDecoder sets the decoder function for the Run function.
func WithDecoder(decoder encoding.Decoder) RunOption {
	return func(o *runOptions) {
		o.decoder = decoder
	}
}

// Run is a convenience function for running a service with a default configuration.
// It creates and runs a BaseService with an empty context.
// This function always exits the program:
// - It exits with the error code if an error occurs.
// - It exits with code 0 if the service runs successfully.
// It is recommended to use this function unless you need to customize the Service
// or want to prevent the program from exiting.
func Run(opts ...RunOption) {
	type context map[string]any
	var o runOptions
	o.apply(opts)
	s := NewBaseService(Config[context]{Context: context{}})
	s.encoder = o.encoder
	s.decoder = o.decoder
	if err := RunService(s); err != nil {
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
		s.Logger().Info("uninitializing service")
		if err := s.Uninit(context.Background()); err != nil {
			s.Logger().Error("failed to uninitialize service", slog.Any("error", err))
		}
		s.Logger().Info("service exited")
	}()
	if err := s.Init(context.Background()); err != nil {
		if _, ok := errkit.ExitCode(err); ok {
			return err
		}
		s.Logger().Error("failed to initialize service", slog.Any("error", err))
		return err
	}

	s.Logger().Info("starting service")
	defer func() {
		s.Logger().Info("shutting down service")
		if err := s.Shutdown(context.Background()); err != nil {
			s.Logger().Error("failed to shutdown service", slog.Any("error", err))
		}
	}()
	err := s.Start(context.Background())
	if err != nil {
		s.Logger().Error("failed to start service", slog.Any("error", err))
	}
	return err
}
