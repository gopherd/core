package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gopherd/core/build"
	"github.com/gopherd/core/component"
)

// Configurator represents config of application
type Configurator interface {
	// ParseArgs parses command line arguments
	ParseArgs(args []string) error
	// CoreConfig returns core configuration
	CoreConfig() *CoreConfig
}

// Core configuration
type CoreConfig struct {
	Project    string             `json:"project"`
	Name       string             `json:"name"`
	ID         int                `json:"id"`
	Components []component.Config `json:"components"`
}

// GenericConfig implments Configurator, which is a generic config for almost all applications
type GenericConfig struct {
	source string
	core   CoreConfig
}

func NewGenericConfig(core CoreConfig) *GenericConfig {
	return &GenericConfig{
		core: core,
	}
}

// CoreConfig implements Configurator CoreConfig method
func (c *GenericConfig) CoreConfig() *CoreConfig {
	return &c.core
}

func (c *GenericConfig) ParseArgs(args []string) error {
	var flags = flag.NewFlagSet(args[0], flag.ExitOnError)
	var input = flags.String("c", "", "Config source")
	var output = flags.String("e", "", "Exported config filename")
	var version = flags.Bool("v", false, "Print version information")
	flags.Parse(args[1:])

	if *version {
		build.Print()
		return nil
	}

	var optional = *input == ""
	var source = *input
	if source == "" {
		source = build.Name() + ".json"
	}
	if err := c.load(source, optional); err != nil {
		return nil
	}

	if *output != "" {
		if f, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
			return err
		} else {
			encoder := json.NewEncoder(f)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(c.core)
			f.Close()
			if err != nil {
				return err
			}
		}
		os.Exit(0)
	}

	return nil
}

// load loads config from source: file or http service
func (c *GenericConfig) load(source string, optional bool) error {
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
	c.source = source
	return json.NewDecoder(r).Decode(&c.core)
}

func (c *GenericConfig) loadFromFile(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

func (c *GenericConfig) loadFromHTTP(source string) (io.ReadCloser, error) {
	var url = source
	for maxRedirects := 32; maxRedirects >= 0; maxRedirects-- {
		// read config from http service
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			url = resp.Header.Get("Location")
			if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusPermanentRedirect {
				source = url
			}
			continue
		}
		return resp.Body, nil
	}
	return nil, errors.New("too many redirects")
}
