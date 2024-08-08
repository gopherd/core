package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/text/templateutil"
)

// Config implements the Config interface and provides a generic
// configuration structure for most applications.
type Config[T any] struct {
	Context    T                  `json:",omitempty"`
	Components []component.Config `json:",omitempty"`
}

// Load processes the configuration based on command-line flags.
// It returns true if the program should exit after this call, along with any error encountered.
func (c *Config[T]) load(source string) error {
	var r io.ReadCloser
	var err error

	switch {
	case source == "":
		return nil
	case source == "-":
		r = os.Stdin
	case strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://"):
		r, err = c.loadFromHTTP(source)
		if err == nil {
			defer r.Close()
		}
	default:
		r, err = os.Open(source)
		if err == nil {
			defer r.Close()
		}
	}

	if err != nil {
		return err
	}

	return json.NewDecoder(r).Decode(&c)
}

// loadFromHTTP loads the configuration from an HTTP source.
// It handles redirects up to a maximum of 32 times.
func (c *Config[T]) loadFromHTTP(source string) (io.ReadCloser, error) {
	url := source
	const maxRedirects = 32

	for redirects := 0; redirects < maxRedirects; redirects++ {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			url = resp.Header.Get("Location")
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

// processTemplate processes the Refs and Options fields of each component.Config
// as text/template templates, using c.Context as the template context.
func (c *Config[T]) processTemplate() error {
	for i := range c.Components {
		com := &c.Components[i]
		if com.UUID != "" {
			if new, err := templateutil.Execute(com.UUID, c.Context); err != nil {
				return fmt.Errorf("process UUID for component %s: %w", com.UUID, err)
			} else {
				com.UUID = new
			}
		}
		if com.Refs.Len() > 0 {
			if new, err := templateutil.Execute(com.Refs.String(), c.Context); err != nil {
				return fmt.Errorf("process Refs for component %s: %w", com.UUID, err)
			} else {
				com.Refs.SetString(new)
			}
		}
		if com.Options.Len() > 0 {
			if new, err := templateutil.Execute(com.Options.String(), c.Context); err != nil {
				return fmt.Errorf("process Options for component %s: %w", com.UUID, err)
			} else {
				com.Options.SetString(new)
			}
		}
	}
	return nil
}

func (c Config[T]) output() {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode config: %e\n", err)
		return
	}
	fmt.Fprint(os.Stdout, buf.String())
}
