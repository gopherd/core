// Package component provides a flexible and extensible component system for building modular applications.
// It defines interfaces and structures for creating, managing, and controlling the lifecycle of components
// within an application.
package component

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
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
	fmt.Stringer

	// Setup configures the component with the given container and configuration.
	Setup(Container, Config) error
}

// Container represents a generic container that can hold components.
type Container interface {
	// GetComponent returns a component by its UUID.
	GetComponent(uuid string) Component
}

// Resolver resolves a reference for a component.
type Resolver interface {
	// UUID returns the UUID of the referenced component.
	UUID() string
	// Resolve resolves the reference for the component.
	Resolve(Container) error
}

// BaseComponent provides a basic implementation of the Component interface.
type BaseComponent[T any] struct {
	lifecycle.BaseLifecycle
	options    T
	identifier string
	container  Container
}

// Options returns a pointer to the component's options.
func (c *BaseComponent[T]) Options() *T {
	return &c.options
}

// String implements the fmt.Stringer interface.
func (c *BaseComponent[T]) String() string {
	return c.identifier
}

// Setup implements the Component Setup method.
func (c *BaseComponent[T]) Setup(container Container, config Config) error {
	c.container = container

	if config.UUID != "" {
		if strings.Contains(config.UUID, config.Name) {
			c.identifier = "#" + config.UUID
		} else {
			c.identifier = fmt.Sprintf("%s#%s", config.Name, config.UUID)
		}
	} else {
		c.identifier = config.Name
	}

	if err := config.Options.DecodeJSON(&c.options); err != nil {
		return fmt.Errorf("failed to unmarshal options: %w", err)
	}
	return nil
}

// Reference represents a reference to another component.
type Reference[T any] struct {
	component T
	uuid      string
}

// Ref creates a reference to a component with the given UUID.
func Ref[T any](uuid string) Reference[T] {
	return Reference[T]{uuid: uuid}
}

// UUID returns the UUID of the referenced component.
func (r Reference[T]) UUID() string {
	return r.uuid
}

// Component returns the referenced component.
func (r Reference[T]) Component() T {
	return r.component
}

// MarshalJSON marshals the referenced component UUID to JSON.
func (r Reference[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.uuid)
}

// UnmarshalJSON unmarshals the referenced component UUID from JSON.
func (r *Reference[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.uuid)
}

// Resolve resolves the reference for the component.
func (r *Reference[T]) Resolve(container Container) error {
	return Resolve(&r.component, container, r.uuid)
}

// Resolve resolves the target component for the given container and UUID.
func Resolve[T any](target *T, container Container, uuid string) error {
	com := container.GetComponent(uuid)
	if com == nil {
		return fmt.Errorf("component %q not found", uuid)
	}
	if c, ok := com.(T); ok {
		*target = c
		return nil
	}
	return fmt.Errorf("component %q type mismatch", uuid)
}

// BaseComponentWithRefs provides a basic implementation of the Component interface with references.
type BaseComponentWithRefs[T, R any] struct {
	BaseComponent[T]
	refs R
}

// Refs returns a pointer to the component's references.
func (c *BaseComponentWithRefs[T, R]) Refs() *R {
	return &c.refs
}

// Setup implements the Component Setup method.
func (c *BaseComponentWithRefs[T, R]) Setup(container Container, config Config) error {
	if err := c.BaseComponent.Setup(container, config); err != nil {
		return err
	}
	if err := config.Refs.DecodeJSON(&c.refs); err != nil {
		return fmt.Errorf("failed to unmarshal refs: %w", err)
	}
	return c.resolveRefs(container)
}

// resolveRefs iterates over the refs field and calls the Resolve method on fields that implement Resolver
func (c *BaseComponentWithRefs[T, R]) resolveRefs(container Container) error {
	t := reflect.TypeOf(&c.refs).Elem()
	v := reflect.ValueOf(&c.refs).Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}
	return c.recursiveResolveRefs(container, t, v)
}

// recursiveResolveRefs recursively resolves references in nested structs
func (c *BaseComponentWithRefs[T, R]) recursiveResolveRefs(container Container, t reflect.Type, v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		resolver := getResolver(fv)
		if resolver == nil && fv.CanAddr() {
			resolver = getResolver(fv.Addr())
		}
		if resolver == nil {
			if fv.Kind() == reflect.Struct {
				if err := c.recursiveResolveRefs(container, fv.Type(), fv); err != nil {
					return err
				}
			}
			continue
		}

		if err := resolver.Resolve(container); err != nil {
			return fmt.Errorf("failed to resolve reference %s: %w", ft.Name, err)
		}
		slog.Debug("resolved reference", "reference", ft.Name)
	}

	return nil
}

// getResolver safely checks if the field implements Resolver interface
func getResolver(field reflect.Value) Resolver {
	if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
		return nil
	}
	if !field.CanInterface() {
		return nil
	}
	if resolver, ok := field.Interface().(Resolver); ok {
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
func (g *Group) AddComponent(uuid string, com Component) Component {
	if uuid != "" {
		if _, exists := g.uuidToComponent[uuid]; exists {
			return nil
		}
		g.uuidToComponent[uuid] = com
	}
	g.components = append(g.components, com)
	return com
}

// GetComponent retrieves a component by its UUID.
func (g *Group) GetComponent(uuid string) Component {
	return g.uuidToComponent[uuid]
}

// Init initializes all components in the group.
func (g *Group) Init(ctx context.Context) error {
	for i := range g.components {
		com := g.components[i]
		slog.Info("initializing component")
		if err := com.Init(ctx); err != nil {
			slog.Error("failed to initialize component", "error", err)
			return err
		}
		slog.Info("component initialized")
		g.numInitialized++
	}
	return nil
}

// Uninit uninitializes all components in reverse order.
func (g *Group) Uninit(ctx context.Context) error {
	for i := g.numInitialized - 1; i >= 0; i-- {
		com := g.components[i]
		slog.Info("uninitializing component")
		if err := com.Uninit(ctx); err != nil {
			slog.Error("failed to uninitialize component", "error", err)
			return err
		}
		slog.Info("component uninitialized")
	}
	return nil
}

// Start starts all components in the group.
func (g *Group) Start(ctx context.Context) error {
	for i := range g.components {
		com := g.components[i]
		slog.Info("starting component")
		if err := com.Start(ctx); err != nil {
			slog.Error("failed to start component", "error", err)
			return err
		}
		slog.Info("component started")
		g.numStarted++
	}
	return nil
}

// Shutdown shuts down all components in reverse order.
func (g *Group) Shutdown(ctx context.Context) error {
	var errs []error
	for i := g.numStarted - 1; i >= 0; i-- {
		com := g.components[i]
		slog.Info("shutting down component")
		if err := com.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown component", "error", err)
			errs = append(errs, err)
		} else {
			slog.Info("component shutdown")
		}
	}
	return errors.Join(errs...)
}

var (
	creatorsMu sync.RWMutex
	creators   = make(map[string]func() Component)
)

// Register makes a component creator available by the provided name.
// It panics if Register is called twice with the same name or if creator is nil.
func Register(name string, creator func() Component) {
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

// Create creates a new component by its name.
func Create(name string) (Component, error) {
	creatorsMu.RLock()
	defer creatorsMu.RUnlock()
	creator, ok := creators[name]
	if !ok {
		return nil, fmt.Errorf("unknown component %q (forgotten import?)", name)
	}
	return creator(), nil
}
