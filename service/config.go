package service

import (
	"bufio"
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
	"github.com/gopherd/core/encoding"
	"github.com/gopherd/core/op"
	"github.com/gopherd/core/text/templates"
)

// Config represents a generic configuration structure for services.
// It includes a context of type T and a list of component configurations.
type Config[T any] struct {
	Context    T                  `json:",omitempty"`
	Components []component.Config `json:",omitempty"`
}

// load processes the configuration based on the provided source.
// It returns an error if the configuration cannot be loaded or decoded.
func (c *Config[T]) load(stdin io.Reader, decoder encoding.Decoder, source string) error {
	if source == "" {
		return nil
	}

	var r io.Reader
	var err error

	switch {
	case source == "-":
		r = stdin
	case strings.HasPrefix(source, "http://"), strings.HasPrefix(source, "https://"):
		var b io.ReadCloser
		b, err = c.loadFromHTTP(source, time.Second*10)
		if err == nil {
			defer b.Close()
		}
		r = b
	default:
		var f io.ReadCloser
		f, err = os.Open(source)
		if err == nil {
			defer f.Close()
		}
		r = f
	}

	if err != nil {
		return fmt.Errorf("open config source failed: %w", err)
	}

	var data []byte
	if decoder == nil {
		data, err = stripJSONComments(r)
		if err != nil {
			return fmt.Errorf("strip JSON comments failed: %w", err)
		}
	} else {
		data, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read config data failed: %w", err)
		}
		data, err = encoding.Transform(data, decoder, json.Marshal)
		if err != nil {
			return fmt.Errorf("decode config failed: %w", err)
		}
	}

	if err := json.Unmarshal(data, c); err != nil {
		if decoder == nil {
			err = encoding.GetJSONSourceError(source, data, err)
		} else {
			switch e := err.(type) {
			case *json.UnmarshalTypeError:
				if e.Struct != "" || e.Field != "" {
					err = errors.New("cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String())
				} else {
					err = errors.New("cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String())
				}
			}
		}
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	return nil
}

// loadFromHTTP loads the configuration from an HTTP source.
// It handles redirects up to a maximum of 32 times.
func (c *Config[T]) loadFromHTTP(source string, timeout time.Duration) (io.ReadCloser, error) {
	const maxRedirects = 32

	client := &http.Client{
		Timeout: timeout,
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
			new, err := templates.Execute(sourcePrefix+"UUID", com.UUID, c.Context, option)
			if err != nil {
				return err
			}
			com.UUID = new
		}

		if op.IfFunc(com.TemplateRefs == nil, enableTemplate, com.TemplateRefs.Deref) && com.Refs.Len() > 0 {
			new, err := templates.Execute(sourcePrefix+"Refs", com.Refs.String(), c.Context, option)
			if err != nil {
				return err
			}
			com.Refs.SetString(new)
		}

		if op.IfFunc(com.TemplateOptions == nil, enableTemplate, com.TemplateOptions.Deref) && com.Options.Len() > 0 {
			new, err := templates.Execute(sourcePrefix+"Options", com.Options.String(), c.Context, option)
			if err != nil {
				return err
			}
			com.Options.SetString(new)
		}
	}

	return nil
}

// output encodes the configuration with the encoder and writes it to stdout.
// It uses indentation for better readability.
func (c *Config[T]) output(components []component.Config, stdout, stderr io.Writer, encoder encoding.Encoder) {
	if len(components) > 0 {
		c.Components = components
	}

	if encoder == nil {
		if data, err := jsonIndentEncoder(c); err != nil {
			fmt.Fprintf(stderr, "Encode config failed: %v\n", err)
		} else {
			fmt.Fprint(stdout, string(data))
		}
		return
	}

	if data, err := json.Marshal(c); err != nil {
		fmt.Fprintf(stderr, "Encode config failed: %v\n", err)
	} else if data, err = encoding.Transform(data, json.Unmarshal, encoder); err != nil {
		fmt.Fprintf(stderr, "Encode config failed: %v\n", err)
	} else {
		fmt.Fprint(stdout, string(data))
	}
}

func stripJSONComments(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("//")) {
			if _, err := buf.Write(line); err != nil {
				return nil, err
			}
		}
		if err := buf.WriteByte('\n'); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// Remove the last newline if it exists
	bytes := buf.Bytes()
	if len(bytes) > 0 && bytes[len(bytes)-1] == '\n' {
		bytes = bytes[:len(bytes)-1]
	}
	return bytes, nil
}

func jsonIndentEncoder(v any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
