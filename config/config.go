// Package config provides functionality for managing application configuration.
// It includes interfaces and implementations for parsing command-line arguments,
// loading configuration from various sources (local files or HTTP), and exporting
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

	"github.com/gopherd/core/buildinfo"
	"github.com/gopherd/core/component"
)

// Config represents the configuration interface for an application.
type Config interface {
	// SetupFlags sets the command-line flags for the configuration.
	SetupFlags(flagSet *flag.FlagSet)

	// Load loads the configuration and returns whether to exit and any error encountered.
	Load() (exit bool, err error)

	// GetComponents returns the components in the configuration.
	GetComponents() []component.Config
}

// BaseConfig implements the Config interface and provides a generic
// configuration structure for most applications.
type BaseConfig[Context any] struct {
	flags struct {
		source string
		export string
		stdin  bool
	}
	core struct {
		Context    Context            `json:",omitempty"`
		Components []component.Config `json:",omitempty"`
	}
}

// NewBaseConfig creates a new BaseConfig with the given context and components.
func NewBaseConfig[Context any](context Context, components ...component.Config) *BaseConfig[Context] {
	c := &BaseConfig[Context]{}
	c.core.Context = context
	c.core.Components = components
	return c
}

// GetContext returns the context of the BaseConfig.
func (c *BaseConfig[Context]) GetContext() Context {
	return c.core.Context
}

// GetComponents returns the components of the BaseConfig.
func (c *BaseConfig[Context]) GetComponents() []component.Config {
	return c.core.Components
}

// SetupFlags sets command-line arguments for the BaseConfig.
func (c *BaseConfig[Context]) SetupFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.flags.source, "c", "", "Config source (file path or URL)")
	flagSet.StringVar(&c.flags.export, "e", "", "Path to export the config")
	flagSet.BoolVar(&c.flags.stdin, "i", false, "Read config from stdin")
}

// Load processes the configuration based on command-line flags.
// It returns true if the program should exit after this call, along with any error encountered.
func (c *BaseConfig[Context]) Load() (bool, error) {
	if c.flags.stdin {
		if err := c.decode(os.Stdin); err != nil {
			return false, fmt.Errorf("failed to read config from stdin: %w", err)
		}
		return false, nil
	}

	optional := c.flags.source == ""
	source := c.flags.source
	if source == "" {
		source = buildinfo.AppName() + ".json"
	}

	if err := c.load(source, optional); err != nil {
		return false, err
	}

	if c.flags.export != "" {
		if err := c.exportConfig(c.flags.export); err != nil {
			return false, fmt.Errorf("failed to export config: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// load loads the configuration from a file or HTTP service.
func (c *BaseConfig[Context]) load(source string, optional bool) error {
	var r io.ReadCloser
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		r, err = c.loadFromHTTP(source)
	} else {
		r, err = c.loadFromFile(source)
	}

	if err != nil {
		if optional && os.IsNotExist(err) {
			return nil
		}
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

// exportConfig exports the current configuration to a JSON file.
func (c *BaseConfig[Context]) exportConfig(path string) error {
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

	if err := encoder.Encode(c.core); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	return nil
}

// decode reads the configuration from JSON using the given reader.
func (c *BaseConfig[Context]) decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&c.core)
}
