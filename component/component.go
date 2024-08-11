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
	"sync/atomic"

	"github.com/gopherd/core/lifecycle"
	"github.com/gopherd/core/types"
)

// Config defines the configuration structure for creating a component.
type Config struct {
	// Name is the component name. It's required.
	Name string

	// UUID is the unique identifier for the component. It can be empty.
	UUID string `json:",omitempty"`

	// Refs is the references to other components.
	Refs types.RawObject `json:",omitempty"`

	// Options is the configuration options for the component.
	Options types.RawObject `json:",omitempty"`

	// TemplateUUID determines if the UUID should be templated.
	// If not set, the default value is determined by the service.
	TemplateUUID *types.Bool `json:",omitempty"`

	// TemplateRefs determines if the Refs should be templated.
	// If not set, the default value is determined by the service.
	TemplateRefs *types.Bool `json:",omitempty"`

	// TemplateOptions determines if the Options should be templated.
	// If not set, the default value is determined by the service.
	TemplateOptions *types.Bool `json:",omitempty"`
}

// Component defines the interface for a generic logic component.
type Component interface {
	lifecycle.Lifecycle
	fmt.Stringer

	// Setup sets up the component with the given container and configuration.
	Setup(Container, Config) error

	// Logger returns the logger instance for the component.
	// Logger must be guranteed to return a non-nil logger instance after Setup is called.
	Logger() *slog.Logger
}

// Container represents a generic container that can hold components.
type Container interface {
	// GetComponent returns a component by its UUID.
	GetComponent(uuid string) Component

	// Logger returns the logger instance for the container.
	Logger() *slog.Logger
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
	logger     atomic.Pointer[slog.Logger]
}

// Options returns a pointer to the component's options.
func (c *BaseComponent[T]) Options() *T {
	return &c.options
}

// String implements the fmt.Stringer interface.
func (c *BaseComponent[T]) String() string {
	return c.identifier
}

// Logger implements the Component Logger method.
func (c *BaseComponent[T]) Logger() *slog.Logger {
	currentLogger := c.logger.Load()
	latestLogger := c.container.Logger()
	if currentLogger != latestLogger {
		currentLogger = latestLogger.With("component", c.identifier)
		c.logger.Store(currentLogger)
	}
	return currentLogger
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
	if loaded, ok := any(&c.options).(interface {
		OnLoaded() error
	}); ok {
		if err := loaded.OnLoaded(); err != nil {
			return fmt.Errorf("failed to load options: %w", err)
		}
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
	com := container.GetComponent(r.uuid)
	if com == nil {
		return fmt.Errorf("component %q not found", r.uuid)
	}
	if c, ok := com.(T); ok {
		r.component = c
		return nil
	}
	return fmt.Errorf("unexpected component %q type: %T", r.uuid, com)
}

// OptionalReference represents an optional reference to another component.
// If the UUID is empty, the reference is ignored, and Component returns nil.
type OptionalReference[T any] struct {
	Reference[T]
}

// OptionalRef creates an optional reference to a component with the given UUID.
func OptionalRef[T any](uuid string) OptionalReference[T] {
	return OptionalReference[T]{
		Reference: Reference[T]{uuid: uuid},
	}
}

// Resolve resolves the reference for the component.
func (r *OptionalReference[T]) Resolve(container Container) error {
	if r.uuid == "" {
		return nil
	}
	return r.Reference.Resolve(container)
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
			return fmt.Errorf("failed to resolve reference %s to %s: %w", ft.Name, resolver.UUID(), err)
		}
		c.Logger().Info("resolve referenced component", "current", c.identifier, "ref", resolver.UUID())
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
		com.Logger().Info("initializing component")
		if err := com.Init(ctx); err != nil {
			com.Logger().Error("failed to initialize component", "error", err)
			return err
		}
		com.Logger().Info("component initialized")
		g.numInitialized++
	}
	return nil
}

// Uninit uninitializes all components in reverse order.
func (g *Group) Uninit(ctx context.Context) error {
	for i := g.numInitialized - 1; i >= 0; i-- {
		com := g.components[i]
		com.Logger().Info("uninitializing component")
		if err := com.Uninit(ctx); err != nil {
			com.Logger().Error("failed to uninitialize component", "error", err)
			return err
		}
		com.Logger().Info("component uninitialized")
	}
	return nil
}

// Start starts all components in the group.
func (g *Group) Start(ctx context.Context) error {
	for i := range g.components {
		com := g.components[i]
		com.Logger().Info("starting component")
		if err := com.Start(ctx); err != nil {
			com.Logger().Error("failed to start component", "error", err)
			return err
		}
		com.Logger().Info("component started")
		g.numStarted++
	}
	return nil
}

// Shutdown shuts down all components in reverse order.
func (g *Group) Shutdown(ctx context.Context) error {
	var errs []error
	for i := g.numStarted - 1; i >= 0; i-- {
		com := g.components[i]
		com.Logger().Info("shutting down component")
		if err := com.Shutdown(ctx); err != nil {
			com.Logger().Error("failed to shutdown component", "error", err)
			errs = append(errs, err)
		} else {
			com.Logger().Info("component shutdown")
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
