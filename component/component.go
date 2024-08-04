// Package component provides a generic component system for building modular applications.
package component

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/gopherd/core/event"
)

// Options represents component-specific configuration options.
type Options = json.RawMessage

// Config defines the configuration structure for creating a component.
type Config struct {
	UUID    string  `json:"uuid,omitempty"`
	Name    string  `json:"name"`
	Options Options `json:"options,omitempty"`
}

// CreateOptions marshals any value into Options. It panics if marshaling fails.
func CreateOptions(v any) Options {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
	return Options(out.Bytes())
}

// Entity represents a generic entity that can hold components.
type Entity interface {
	event.Dispatcher[reflect.Type]
	GetComponent(uuid string) Component
}

// ComponentCreator is a function type that creates a new Component instance.
type ComponentCreator func() Component

// Metadata defines the interface for accessing component metadata.
type Metadata interface {
	// UUID returns the unique identifier of the component.
	UUID() string
	// Name returns the name of the component.
	Name() string
	// Entity returns the entity to which the component belongs.
	Entity() Entity
}

// Component defines the interface for a generic logic component.
type Component interface {
	Metadata
	// OnCreated is called when the component is created.
	OnCreated(Entity, Config) error
	// Init initializes the component.
	Init(context.Context) error
	// Uninit uninitializes the component.
	Uninit(context.Context) error
	// Start starts the component.
	Start(context.Context) error
	// Shutdown gracefully shuts down the component.
	Shutdown(context.Context) error
}

// BaseComponent provides a basic implementation of the Component interface.
type BaseComponent[T any] struct {
	uuid, name string
	entity     Entity
	options    T
}

// Options returns a pointer to the component's options.
func (com *BaseComponent[T]) Options() *T {
	return &com.options
}

// UUID implements the Metadata UUID method.
func (com *BaseComponent[T]) UUID() string {
	return com.uuid
}

// Name implements the Metadata Name method.
func (com *BaseComponent[T]) Name() string {
	return com.name
}

// Entity implements the Metadata Entity method.
func (com *BaseComponent[T]) Entity() Entity {
	return com.entity
}

// OnCreated implements the Component OnCreated method.
func (com *BaseComponent[T]) OnCreated(entity Entity, config Config) error {
	com.uuid = config.UUID
	com.name = config.Name
	com.entity = entity
	if len(config.Options) > 0 {
		return json.Unmarshal(config.Options, &com.options)
	}
	return nil
}

// Init implements the Component Init method.
func (com *BaseComponent[T]) Init(_ context.Context) error {
	return nil
}

// Uninit implements the Component Uninit method.
func (com *BaseComponent[T]) Uninit(_ context.Context) error {
	return nil
}

// Start implements the Component Start method.
func (com *BaseComponent[T]) Start(_ context.Context) error {
	return nil
}

// Shutdown implements the Component Shutdown method.
func (com *BaseComponent[T]) Shutdown(_ context.Context) error {
	return nil
}

// Manager manages a group of components.
type Manager struct {
	components     []Component
	uuid2component map[string]Component
	numInitialized int
	numStarted     int
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		uuid2component: make(map[string]Component),
	}
}

// AddComponent adds a component to the manager.
// It returns nil if a component with the same UUID already exists.
func (m *Manager) AddComponent(com Component) Component {
	uuid := com.UUID()
	if uuid != "" {
		if _, exists := m.uuid2component[uuid]; exists {
			return nil
		}
		m.uuid2component[uuid] = com
	}
	m.components = append(m.components, com)
	return com
}

// GetComponent retrieves a component by its UUID.
func (m *Manager) GetComponent(uuid string) Component {
	if uuid == "" {
		return nil
	}
	return m.uuid2component[uuid]
}

// Init initializes all components in the manager.
func (m *Manager) Init(ctx context.Context) error {
	for i := range m.components {
		com := m.components[i]
		slog.Info(
			"initializing component",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		if err := com.Init(ctx); err != nil {
			slog.Error(
				"failed to initialize component",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
				slog.Any("error", err),
			)
			return err
		}
		slog.Info(
			"component initialized",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		m.numInitialized++
	}
	return nil
}

// Uninit uninitializes all components in reverse order.
func (m *Manager) Uninit(ctx context.Context) error {
	for i := m.numInitialized - 1; i >= 0; i-- {
		com := m.components[i]
		slog.Info(
			"uninitializing component",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		if err := com.Uninit(ctx); err != nil {
			slog.Error(
				"failed to uninitialize component",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
				slog.Any("error", err),
			)
			return err
		}
		slog.Info(
			"component uninitialized",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
	}
	return nil
}

// Start starts all components in the manager.
func (m *Manager) Start(ctx context.Context) error {
	for i := range m.components {
		com := m.components[i]
		slog.Info(
			"starting component",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		if err := com.Start(ctx); err != nil {
			slog.Error(
				"failed to start component",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
				slog.Any("error", err),
			)
			return err
		}
		slog.Info(
			"component started",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		m.numStarted++
	}
	return nil
}

// Shutdown shuts down all components in reverse order.
func (m *Manager) Shutdown(ctx context.Context) error {
	for i := m.numStarted - 1; i >= 0; i-- {
		com := m.components[i]
		slog.Info(
			"shutting down component",
			slog.String("uuid", com.UUID()),
			slog.String("name", com.Name()),
		)
		if err := com.Shutdown(ctx); err != nil {
			slog.Error(
				"failed to shutdown component",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
				slog.Any("error", err),
			)
		} else {
			slog.Info(
				"component shutdown",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
			)
		}
	}
	return nil
}

var (
	creatorsMu sync.RWMutex
	creators   map[string]ComponentCreator
)

// Register makes a component creator available by the provided name.
// It panics if Register is called twice with the same name or if creator is nil.
func Register(name string, creator ComponentCreator) {
	creatorsMu.Lock()
	defer creatorsMu.Unlock()
	if creators == nil {
		creators = make(map[string]ComponentCreator)
	}
	if creator == nil {
		panic("component: Register component " + name + " creator is nil")
	}
	if _, dup := creators[name]; dup {
		panic("component: Register called twice for component " + name)
	}
	creators[name] = creator
}

// Lookup returns the component creator by name.
func Lookup(name string) ComponentCreator {
	creatorsMu.RLock()
	defer creatorsMu.RUnlock()
	return creators[name]
}

// Resolve resolves the component in entity by uuid.
func Resolve[T any](entity Entity, uuid string) (T, error) {
	var zero T
	if entity == nil {
		return zero, errors.New("entity is nil")
	}
	com := entity.GetComponent(uuid)
	if com == nil {
		return zero, errors.New("component not found")
	}
	if c, ok := com.(T); ok {
		return c, nil
	}
	return zero, fmt.Errorf("component %T type mismatch", com)
}

// MustResolve resolves the component in entity by uuid.
// It panics if the component is not found or type mismatched.
func MustResolve[T any](entity Entity, uuid string) T {
	c, err := Resolve[T](entity, uuid)
	if err != nil {
		panic(fmt.Errorf("resolve component %q error: %w", uuid, err))
	}
	return c
}
