// Package component provides a generic component system for building modular applications.
package component

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/gopherd/core/event"
	"github.com/gopherd/core/lifecycle"
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

// BaseOptions represents the base options structure for creating a component.
type BaseOptions[T any] struct {
	Deps T `json:"deps"`
}

// DependencyResolver resolves a dependency for a component.
type DependencyResolver interface {
	Resolve(Entity) error
}

// Dependency represents a dependency on another component.
type Dependency[T any] struct {
	component T
	uuid      string
}

// MarshalJSON marshals the dependency component uuid to JSON.
func (d *Dependency[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.uuid)
}

// UnmarshalJSON unmarshals the dependency component uuid from JSON.
func (d *Dependency[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.uuid)
}

// Component returns the target component.
func (d *Dependency[T]) Component() T {
	return d.component
}

// Resolve resolves the target component.
func (d *Dependency[T]) Resolve(entity Entity) error {
	return Resolve(&d.component, entity, d.uuid)
}

func Resolve[T any](target *T, entity Entity, uuid string) error {
	com := entity.GetComponent(uuid)
	if com == nil {
		return fmt.Errorf("component %q not found", uuid)
	}
	if c, ok := com.(T); ok {
		*target = c
		return nil
	}
	return fmt.Errorf("component %q type mismatch", uuid)
}

// Component defines the interface for a generic logic component.
type Component interface {
	Metadata
	lifecycle.Lifecycle
	// OnCreated is called when the component is created.
	OnCreated(Entity, Config) error
}

// BaseComponent provides a basic implementation of the Component interface.
type BaseComponent[T any] struct {
	lifecycle.BaseLifecycle
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
	return com.resolveDependencies()
}

// ResolveDependencies iterates over the Deps field in options and calls the Resolve method on fields that implement DependencyResolver
func (com *BaseComponent[T]) resolveDependencies() error {
	optionsValue := reflect.ValueOf(com.options)
	optionsType := reflect.TypeOf(com.options)

	// Find the Deps field
	depsField, found := optionsType.FieldByName("Deps")
	if !found {
		// If there is no Deps field, return directly
		return nil
	}

	// Get the value of the Deps field
	depsValue := optionsValue.FieldByName(depsField.Name)
	if depsValue.Kind() != reflect.Struct {
		return fmt.Errorf("field Deps should be a struct, but got %s", depsValue.Kind())
	}

	// Iterate over the fields of the Deps struct
	for i := 0; i < depsValue.NumField(); i++ {
		field := depsValue.Field(i)
		fieldType := depsValue.Type().Field(i)

		// Check if the field implements DependencyResolver interface
		if resolver, ok := isDependencyResolver(field); ok {
			if err := resolver.Resolve(com.entity); err != nil {
				return fmt.Errorf("failed to resolve dependency %s: %w", fieldType.Name, err)
			}
		} else if field.Kind() == reflect.Ptr && !field.IsNil() {
			// Check if the pointer field implements DependencyResolver interface
			if resolver, ok := isDependencyResolver(field.Elem()); ok {
				if err := resolver.Resolve(com.entity); err != nil {
					return fmt.Errorf("failed to resolve dependency %s: %w", fieldType.Name, err)
				}
			} else {
				return fmt.Errorf("dependency %s does not implement DependencyResolver", fieldType.Name)
			}
		} else {
			return fmt.Errorf("dependency %s does not implement DependencyResolver", fieldType.Name)
		}
	}

	return nil
}

// isDependencyResolver safely checks if the field implements DependencyResolver interface
func isDependencyResolver(field reflect.Value) (DependencyResolver, bool) {
	if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
		return nil, false
	}
	if !field.CanInterface() {
		return nil, false
	}
	if resolver, ok := field.Interface().(DependencyResolver); ok {
		return resolver, true
	}
	return nil, false
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
