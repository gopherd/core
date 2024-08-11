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
	"time"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/op"
	"github.com/gopherd/core/text/templateutil"
)

// Config represents a generic configuration structure for services.
// It includes a context of type T and a list of component configurations.
type Config[T any] struct {
	Context    T                  `json:",omitempty"`
	Components []component.Config `json:",omitempty"`
}

// load processes the configuration based on the provided source.
// It returns an error if the configuration cannot be loaded or decoded.
func (c *Config[T]) load(source string) error {
	if source == "" {
		return nil
	}

	var r io.ReadCloser
	var err error

	switch {
	case source == "-":
		r = os.Stdin
	case strings.HasPrefix(source, "http://"), strings.HasPrefix(source, "https://"):
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

	return json.NewDecoder(r).Decode(c)
}

// loadFromHTTP loads the configuration from an HTTP source.
// It handles redirects up to a maximum of 32 times.
func (c *Config[T]) loadFromHTTP(source string) (io.ReadCloser, error) {
	const maxRedirects = 32

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(source)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// processTemplate processes the UUID, Refs, and Options fields of each component.Config
// as text/template templates, using c.Context as the template context.
func (c *Config[T]) processTemplate(enableTemplate bool, source string) error {
	const option = "missingkey=error"
	for i := range c.Components {
		com := &c.Components[i]

		identifier := com.Name
		if com.UUID != "" {
			identifier += "#" + com.UUID
		}
		sourcePrefix := fmt.Sprintf("%s[%s].", source, identifier)
		if op.IfFunc(com.TemplateUUID == nil, enableTemplate, com.TemplateUUID.Deref) && com.UUID != "" {
			new, err := templateutil.Execute(sourcePrefix+"UUID", com.UUID, c.Context, option)
			if err != nil {
				return err
			}
			com.UUID = new
		}

		if op.IfFunc(com.TemplateRefs == nil, enableTemplate, com.TemplateRefs.Deref) && com.Refs.Len() > 0 {
			new, err := templateutil.Execute(sourcePrefix+"Refs", com.Refs.String(), c.Context, option)
			if err != nil {
				return err
			}
			com.Refs.SetString(new)
		}

		if op.IfFunc(com.TemplateOptions == nil, enableTemplate, com.TemplateOptions.Deref) && com.Options.Len() > 0 {
			new, err := templateutil.Execute(sourcePrefix+"Options", com.Options.String(), c.Context, option)
			if err != nil {
				return err
			}
			com.Options.SetString(new)
		}
	}

	return nil
}

// output encodes the configuration as JSON and writes it to stdout.
// It uses indentation for better readability.
func (c Config[T]) output() {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(c); err != nil {
		fmt.Fprintf(os.Stderr, "Encode config failed: %v\n", err)
		return
	}

	fmt.Fprint(os.Stdout, buf.String())
}
