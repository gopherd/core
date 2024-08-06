package logging

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
)

// Provider represents a logger provider
type Provider interface {
	Handler() slog.Handler
	Close() error
}

// ProviderFunc creates a logger provider
type ProviderFunc func([]byte) (Provider, error)

var (
	proviersMu sync.RWMutex
	providers  = make(map[string]ProviderFunc)
)

// Register registers a logger provider creator
func Register(name string, provider ProviderFunc) {
	proviersMu.Lock()
	defer proviersMu.Unlock()
	if provider == nil {
		panic("logger: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("logger: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Lookup returns the logger provider creator by name
func Lookup(name string) ProviderFunc {
	proviersMu.RLock()
	defer proviersMu.RUnlock()
	return providers[name]
}

// Register standard providers
func init() {
	Register("stderr", func(options []byte) (Provider, error) {
		return newStdProvider(os.Stderr, "text", options)
	})
	Register("stdout", func(options []byte) (Provider, error) {
		return newStdProvider(os.Stdout, "text", options)
	})
	Register("stderr/json", func(options []byte) (Provider, error) {
		return newStdProvider(os.Stderr, "json", options)
	})
	Register("stdout/json", func(options []byte) (Provider, error) {
		return newStdProvider(os.Stdout, "json", options)
	})
}

type stdProvider struct {
	handler slog.Handler
}

func (p *stdProvider) Handler() slog.Handler {
	return p.handler
}

func (p *stdProvider) Close() error {
	return nil
}

// StdOptions represents standard logger options
type StdOptions struct {
	Level     slog.Level
	AddSource bool
}

func newStdProvider(writer io.Writer, formatter string, opts []byte) (Provider, error) {
	var options StdOptions
	if err := json.Unmarshal(opts, &options); err != nil {
		return nil, err
	}
	var handler slog.Handler
	if formatter == "json" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     options.Level,
			AddSource: options.AddSource,
		})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:     options.Level,
			AddSource: options.AddSource,
		})
	}
	return &stdProvider{
		handler: handler,
	}, nil
}
