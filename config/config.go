// Package config provides functionality for managing application configuration.
// It includes interfaces and implementations for parsing command-line arguments,
// loading configuration from various sources (local files or HTTP), and outputting
// configuration data to JSON format.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/errkit"
	"github.com/gopherd/core/text/templateutil"
)

// Config represents the configuration interface for an application.
type Config interface {
	// SetupFlags sets the command-line flags for the configuration.
	SetupFlags(flagSet *flag.FlagSet)

	// Load loads the configuration
	Load() error

	// GetComponents returns the components in the configuration.
	GetComponents() []component.Config
}

// BaseConfig implements the Config interface and provides a generic
// configuration structure for most applications.
type BaseConfig[Context any] struct {
	flags struct {
		source          string
		output          string
		test            bool
		disableTemplate bool
	}
	data struct {
		Context    Context            `json:",omitempty"`
		Components []component.Config `json:",omitempty"`
	}
}

// NewBaseConfig creates a new BaseConfig with the given context and components.
func NewBaseConfig[Context any](context Context, components []component.Config) *BaseConfig[Context] {
	c := &BaseConfig[Context]{}
	c.data.Context = context
	c.data.Components = components
	return c
}

// GetContext returns the context of the BaseConfig.
func (c *BaseConfig[Context]) GetContext() Context {
	return c.data.Context
}

// GetComponents returns the components of the BaseConfig.
func (c *BaseConfig[Context]) GetComponents() []component.Config {
	return c.data.Components
}

// SetupFlags sets command-line arguments for the BaseConfig.
func (c *BaseConfig[Context]) SetupFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.flags.source, "c", "", "Specify the config source (file path, HTTP URL, or '-' for stdin)")
	flagSet.StringVar(&c.flags.output, "o", "", "Specify the config output (file path or '-' for stdout) and exit")
	flagSet.BoolVar(&c.flags.test, "t", false, "Test the config for validity and exit")
	flagSet.BoolVar(&c.flags.disableTemplate, "T", false, "Disable template parsing for components config")
}

// Load processes the configuration based on command-line flags.
// It returns true if the program should exit after this call, along with any error encountered.
func (c *BaseConfig[Context]) Load() (err error) {
	defer func() {
		if c.flags.test {
			if err != nil {
				fmt.Println("Config test failed: ", err)
			} else {
				fmt.Println("Config test successful")
				// Exit after test
				err = errkit.NewExitError(0)
			}
		}
		if err == nil && c.flags.source == "" {
			// Config source is required unless testing or outputting
			err = errors.New("no config source specified")
		}
	}()

	err = c.load()
	if err != nil {
		return
	}

	if !c.flags.disableTemplate {
		if err = c.parseComponentTemplates(); err != nil {
			err = fmt.Errorf("failed to parse components: %w", err)
			return
		}
	}

	if c.flags.output != "" {
		if err = c.outputConfig(c.flags.output); err != nil {
			err = fmt.Errorf("failed to output config: %w", err)
			return
		}
		// Exit after output
		return errkit.NewExitError(0)
	}

	return nil
}

func (c *BaseConfig[Context]) load() error {
	switch c.flags.source {
	case "":
	case "-":
		if err := c.decode(os.Stdin); err != nil {
			return fmt.Errorf("failed to read config from stdin: %w", err)
		}
	default:
		if err := c.loadFromSource(c.flags.source); err != nil {
			return fmt.Errorf("failed to load config from source: %w", err)
		}
	}

	return nil
}

// loadFromSource loads the configuration from a file or HTTP service.
func (c *BaseConfig[Context]) loadFromSource(source string) error {
	var r io.ReadCloser
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		r, err = c.loadFromHTTP(source)
	} else {
		r, err = c.loadFromFile(source)
	}

	if err != nil {
		return err
	}
	defer r.Close()

	return c.decode(r)
}

// loadFromFile loads the configuration from a local file.
func (c *BaseConfig[Context]) loadFromFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

// loadFromHTTP loads the configuration from an HTTP source.
// It handles redirects up to a maximum of 32 times.
func (c *BaseConfig[Context]) loadFromHTTP(source string) (io.ReadCloser, error) {
	url := source
	const maxRedirects = 32

	for redirects := 0; redirects < maxRedirects; redirects++ {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			url = resp.Header.Get("Location")
			if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusPermanentRedirect {
				source = url
			}
			resp.Body.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
		}

		return resp.Body, nil
	}

	return nil, errors.New("too many redirects")
}

// outputConfig outputs the current configuration to a JSON file.
func (c *BaseConfig[Context]) outputConfig(path string) error {
	if path == "-" {
		return c.encode(os.Stdout)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	return c.encode(f)
}

// encode writes the configuration as JSON to the given writer.
func (c *BaseConfig[Context]) encode(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(c.data); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return nil
}

// decode reads the configuration from JSON using the given reader.
func (c *BaseConfig[Context]) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&c.data)
}

// parseComponentTemplates processes the Refs and Options fields of each component.Config
// as text/template templates, using c.data.Context as the template context.
func (c *BaseConfig[Context]) parseComponentTemplates() error {
	for i := range c.data.Components {
		com := &c.data.Components[i]
		if com.UUID != "" {
			if new, err := templateutil.Execute(com.UUID, c.data.Context); err != nil {
				return fmt.Errorf("parse UUID for component %s: %w", com.UUID, err)
			} else {
				com.UUID = new
			}
		}
		if com.Refs.Len() > 0 {
			if new, err := templateutil.Execute(com.Refs.String(), c.data.Context); err != nil {
				return fmt.Errorf("parse Refs for component %s: %w", com.UUID, err)
			} else {
				com.Refs.SetString(new)
			}
		}
		if com.Options.Len() > 0 {
			if new, err := templateutil.Execute(com.Options.String(), c.data.Context); err != nil {
				return fmt.Errorf("parse Options for component %s: %w", com.UUID, err)
			} else {
				com.Options.SetString(new)
			}
		}
	}
	return nil
}
