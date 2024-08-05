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
	// SetFlags sets the command-line flags for the configuration.
	SetFlags(flagSet *flag.FlagSet)

	// Load loads the configuration and returns whether to exit and any error encountered.
	Load() (exit bool, err error)

	// CoreConfig returns the core configuration.
	CoreConfig() *CoreConfig
}

// CoreConfig represents the core configuration structure.
type CoreConfig struct {
	Project    string
	Name       string
	ID         int
	Components []component.Config
}

// BaseConfig implements the Config interface and provides a generic
// configuration structure for most applications.
type BaseConfig struct {
	flags struct {
		source string
		export string
	}
	core CoreConfig
}

// NewBaseConfig creates a new BaseConfig with the given CoreConfig.
func NewBaseConfig(core CoreConfig) *BaseConfig {
	return &BaseConfig{
		core: core,
	}
}

// CoreConfig returns a pointer to the core configuration.
func (c *BaseConfig) CoreConfig() *CoreConfig {
	return &c.core
}

// SetFlags sets command-line arguments for the BaseConfig.
func (c *BaseConfig) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.flags.source, "c", "", "Config source (file path or URL)")
	flagSet.StringVar(&c.flags.export, "e", "", "Path to export the config")
}

// Load processes the configuration based on command-line flags.
// It returns true if the program should exit after this call, along with any error encountered.
func (c *BaseConfig) Load() (bool, error) {
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
func (c *BaseConfig) load(source string, optional bool) error {
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

	c.flags.source = source
	return json.NewDecoder(r).Decode(&c.core)
}

// loadFromFile loads the configuration from a local file.
func (c *BaseConfig) loadFromFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

// loadFromHTTP loads the configuration from an HTTP source.
// It handles redirects up to a maximum of 32 times.
func (c *BaseConfig) loadFromHTTP(source string) (io.ReadCloser, error) {
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
func (c *BaseConfig) exportConfig(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(c.core); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
