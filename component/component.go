/*
Package component provides a flexible and extensible component system for building modular applications.
It defines interfaces and structures for creating, managing, and controlling the lifecycle of components
within an application.

# Components

A Component in this system is an entity that has a defined lifecycle and can be managed as part of a larger application.
Each Component goes through the following lifecycle stages:

 1. Creation (Ctor): The component is instantiated and its initial configuration is set.
 2. Mount (OnMounted): The component is mounted to an entity and can reference other components.
 3. Initialization (Init): The component prepares its internal state and resources.
 4. Starting (Start): The component begins its main operations or services.
 5. Shutdown: The component gracefully stops its operations.
 6. Uninitialization (Uninit): The component releases any resources it has acquired.

The lifecycle methods are called in the following order:

		Ctor -> OnMounted -> Init -> Start -> Shutdown -> Uninit
	                           |       |         |           |
	                           |       +---------+           |
	                           +=============================+

It's important to note that:
  - If a component has been initialized (Init called), it must be uninitialized (Uninit called).
  - If a component has been started (Start called), it must be shut down (Shutdown called).

# Component Group

The Group struct is responsible for handling multiple components. It ensures that:

 1. Components are initialized, started, shut down, and uninitialized in the correct order.
 2. Components are initialized and started in the order they were added to the Group.
 3. Components are shut down and uninitialized in the reverse order they were added.

This ordering ensures that dependencies between components are respected during the application's lifecycle.

# Usage

Here's a basic example of how to create a custom component:

	type MyComponent struct {
		component.BaseComponent[MyOptions]
	}

	func (c *MyComponent) Init(ctx context.Context) error {
		// Initialize the component
		return nil
	}

	func (c *MyComponent) Start(ctx context.Context) error {
		// Start the component's operations
		return nil
	}

	func (c *MyComponent) Shutdown(ctx context.Context) error {
		// Gracefully stop the component's operations
		return nil
	}

	func (c *MyComponent) Uninit(ctx context.Context) error {
		// Clean up any resources
		return nil
	}

The package also provides a registry for component creators, allowing for dynamic component creation and management.

For more detailed information on specific types and methods, refer to the individual type and function documentation.
*/
package component

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/gopherd/core/lifecycle"
	"github.com/gopherd/core/raw"
)

// Config defines the configuration structure for creating a component.
type Config struct {
	Name    string
	UUID    string     `json:",omitempty"`
	Refs    raw.Object `json:",omitempty"`
	Options raw.Object `json:",omitempty"`
}

// Component defines the interface for a generic logic component.
type Component interface {
	lifecycle.Lifecycle

	// UUID returns the unique identifier of the component.
	UUID() string

	// Name returns the name of the component.
	Name() string

	// Entity returns the entity to which the component belongs.
	// It's available after the component is mounted.
	Entity() Entity

	// Ctor is called when the component is created.
	// The entity currently not available, so the component should not reference other components.
	Ctor(Config) error

	// OnMounted is called when the component is mounted to the entity.
	OnMounted(Entity) error
}

// Entity represents a generic entity that can hold components.
type Entity interface {
	GetComponent(uuid string) Component
}

// ComponentCreator is a function type that creates a new Component instance.
type ComponentCreator func() Component

// ReferenceResolver resolves a reference for a component.
type ReferenceResolver interface {
	ResolveReference(Entity) error
}

// Reference represents a reference on another component.
type Reference[T any] struct {
	component T
	uuid      string
}

// UUID returns the UUID of the referenced component.
func (r Reference[T]) UUID() string {
	return r.uuid
}

// MarshalJSON marshals the referenced component UUID to JSON.
func (r Reference[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.uuid)
}

// UnmarshalJSON unmarshals the referenced component UUID from JSON.
func (r *Reference[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.uuid)
}

// Component returns the target component.
func (r Reference[T]) Component() T {
	return r.component
}

// Ref creates a new reference to a component.
func Ref[T any](uuid string) Reference[T] {
	return Reference[T]{uuid: uuid}
}

// ResolveReference resolves the target component.
func (r *Reference[T]) ResolveReference(entity Entity) error {
	return Resolve(&r.component, entity, r.uuid)
}

// Resolve resolves the target component for the given entity and UUID.
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

// BaseComponent provides a basic implementation of the Component interface.
type BaseComponent[T any] struct {
	lifecycle.BaseLifecycle
	uuid, name string
	entity     Entity
	options    T
}

// Options returns a pointer to the component's options.
func (c *BaseComponent[T]) Options() *T {
	return &c.options
}

// UUID implements the Component UUID method.
func (c *BaseComponent[T]) UUID() string {
	return c.uuid
}

// Name implements the Component Name method.
func (c *BaseComponent[T]) Name() string {
	return c.name
}

// Entity implements the Component Entity method.
func (c *BaseComponent[T]) Entity() Entity {
	return c.entity
}

// Ctor implements the Component Ctor method.
func (c *BaseComponent[T]) Ctor(config Config) error {
	c.uuid = config.UUID
	c.name = config.Name
	if err := config.Options.DecodeJSON(&c.options); err != nil {
		return fmt.Errorf("failed to unmarshal options: %w", err)
	}
	return nil
}

// OnMounted implements the Component Mount method.
func (c *BaseComponent[T]) OnMounted(entity Entity) error {
	c.entity = entity
	return nil
}

// BaseComponentWithRefs provides a basic implementation of the Component interface with references.
type BaseComponentWithRefs[R, T any] struct {
	BaseComponent[T]
	refs R
}

// Refs returns a pointer to the component's references.
func (c *BaseComponentWithRefs[R, T]) Refs() *R {
	return &c.refs
}

// Ctor implements the Component Ctor method.
func (c *BaseComponentWithRefs[R, T]) Ctor(config Config) error {
	if err := c.BaseComponent.Ctor(config); err != nil {
		return err
	}
	if err := config.Refs.DecodeJSON(&c.refs); err != nil {
		return fmt.Errorf("failed to unmarshal refs: %w", err)
	}
	return nil
}

// OnMounted implements the Component OnMounted method.
func (c *BaseComponentWithRefs[R, T]) OnMounted(entity Entity) error {
	if err := c.BaseComponent.OnMounted(entity); err != nil {
		return err
	}
	return c.resolveReferences()
}

// resolveReferences iterates over the refs field and calls the Resolve method on fields that implement ReferenceResolver
func (c *BaseComponentWithRefs[R, T]) resolveReferences() error {
	t := reflect.TypeOf(&c.refs).Elem()
	v := reflect.ValueOf(&c.refs).Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}
	return c.recursiveResolveReferences(t, v)
}

// recursiveResolveReferences recursively resolves references in nested structs
func (c *BaseComponentWithRefs[R, T]) recursiveResolveReferences(t reflect.Type, v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		resolver := reflectReferenceResolver(fv)
		if resolver == nil && fv.CanAddr() {
			resolver = reflectReferenceResolver(fv.Addr())
		}
		if resolver == nil {
			if fv.Kind() == reflect.Struct {
				if err := c.recursiveResolveReferences(fv.Type(), fv); err != nil {
					return err
				}
			}
			continue
		}

		if err := resolver.ResolveReference(c.entity); err != nil {
			return fmt.Errorf("failed to resolve reference %s: %w", ft.Name, err)
		}
		slog.Debug(
			"resolved reference",
			slog.String("component", c.UUID()),
			slog.String("reference", ft.Name),
		)
	}

	return nil
}

// reflectReferenceResolver safely checks if the field implements ReferenceResolver interface
func reflectReferenceResolver(field reflect.Value) ReferenceResolver {
	if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
		return nil
	}
	if !field.CanInterface() {
		return nil
	}
	if resolver, ok := field.Interface().(ReferenceResolver); ok {
		return resolver
	}
	return nil
}

// Group manages a group of components.
type Group struct {
	components      []Component
	uuidToComponent map[string]Component
	numInitialized  int
	numStarted      int
}

// NewGroup creates a new Group instance.
func NewGroup() *Group {
	return &Group{
		uuidToComponent: make(map[string]Component),
	}
}

// AddComponent adds a component to the group.
// It returns nil if a component with the same UUID already exists.
func (m *Group) AddComponent(com Component) Component {
	uuid := com.UUID()
	if uuid != "" {
		if _, exists := m.uuidToComponent[uuid]; exists {
			return nil
		}
		m.uuidToComponent[uuid] = com
	}
	m.components = append(m.components, com)
	return com
}

// GetComponent retrieves a component by its UUID.
func (m *Group) GetComponent(uuid string) Component {
	return m.uuidToComponent[uuid]
}

// OnMounted calls the OnMounted method on all components.
func (m *Group) OnMounted(entity Entity) error {
	for i := range m.components {
		com := m.components[i]
		if err := com.OnMounted(entity); err != nil {
			return err
		}
	}
	return nil
}

// Init initializes all components in the group.
func (m *Group) Init(ctx context.Context) error {
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
func (m *Group) Uninit(ctx context.Context) error {
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

// Start starts all components in the group.
func (m *Group) Start(ctx context.Context) error {
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
func (m *Group) Shutdown(ctx context.Context) error {
	var errs []error
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
			errs = append(errs, err)
		} else {
			slog.Info(
				"component shutdown",
				slog.String("uuid", com.UUID()),
				slog.String("name", com.Name()),
			)
		}
	}
	return errors.Join(errs...)
}

var (
	creatorsMu sync.RWMutex
	creators   = make(map[string]ComponentCreator)
)

// Register makes a component creator available by the provided name.
// It panics if Register is called twice with the same name or if creator is nil.
func Register(name string, creator ComponentCreator) {
	creatorsMu.Lock()
	defer creatorsMu.Unlock()
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
