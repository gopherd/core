package component

import (
	"context"
	"encoding/json"
	"sync"
)

// Options represents component options
type Options = json.RawMessage

// CreateOptions creates options from any value, panic if failed
func CreateOptions(v any) Options {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return Options(data)
}

// Entity represents a generic entity with components
type Entity interface {
	GetComponent(uuid string) Component
}

// ComponentCreator creates a component
type ComponentCreator func() Component

// Metadata represents metadata of a component
type Metadata interface {
	// UUID returns the unique id of the component
	UUID() string
	// Entity returns the entity of the component
	Entity() Entity
}

// Component represents a generic logic component
type Component interface {
	Metadata
	// OnCreated is called when the component is created
	OnCreated(entity Entity, uuid string, options Options) error
	// Init initializes the component
	Init(context.Context, Entity) error
	// Start starts the component
	Start(context.Context) error
	// Shutdown gracefully shuts down the component
	Shutdown(context.Context) error
}

var _ Component = (*BaseComponent[any])(nil)

// BaseComponent implements the Component interface
type BaseComponent[T any] struct {
	uuid    string
	entity  Entity
	options T
}

// Options returns the options of the component
func (com *BaseComponent[T]) Options() *T {
	return &com.options
}

// UUID implements Metadata UUID method
func (com *BaseComponent[T]) UUID() string {
	return com.uuid
}

// Entity implements Metadata Entity method
func (com *BaseComponent[T]) Entity() Entity {
	return com.entity
}

// OnCreated implements Component OnCreated method
func (com *BaseComponent[T]) OnCreated(entity Entity, uuid string, options Options) error {
	com.uuid = uuid
	com.entity = entity
	return json.Unmarshal(options, &com.options)
}

// Init implements Component Init method
func (com *BaseComponent[T]) Init(_ context.Context, _ Entity) error {
	return nil
}

// Start implements Component Start method
func (com *BaseComponent[T]) Start(_ context.Context) error {
	return nil
}

// Shutdown implements Component Shutdown method
func (com *BaseComponent[T]) Shutdown(_ context.Context) error {
	return nil
}

// Manager used to manages a group of components
type Manager struct {
	components        map[string]Component
	orderedComponents []Component
}

// NewManager creates a Manager
func NewManager() *Manager {
	return &Manager{
		components: make(map[string]Component),
	}
}

// Add adds a component to the manager, returns nil if the uuid is duplicated
func (m *Manager) AddComponent(com Component) Component {
	uuid := com.UUID()
	if _, dup := m.components[uuid]; dup {
		return nil
	}
	m.components[uuid] = com
	m.orderedComponents = append(m.orderedComponents, com)
	return com
}

func (m *Manager) GetComponent(uuid string) Component {
	return m.components[uuid]
}

// Init initializes all components
func (m *Manager) Init(ctx context.Context, entity Entity) error {
	for i := range m.orderedComponents {
		com := m.orderedComponents[i]
		if err := com.Init(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// Start starts all components
func (m *Manager) Start(ctx context.Context) error {
	for i := range m.orderedComponents {
		com := m.orderedComponents[i]
		if err := com.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shutdowns all components in reverse order
func (m *Manager) Shutdown(ctx context.Context) error {
	for i := len(m.orderedComponents) - 1; i >= 0; i-- {
		com := m.orderedComponents[i]
		com.Shutdown(ctx)
	}
	return nil
}

var registry struct {
	creatorsMu sync.RWMutex
	creators   map[string]ComponentCreator
}

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, creator ComponentCreator) {
	registry.creatorsMu.Lock()
	defer registry.creatorsMu.Unlock()
	if registry.creators == nil {
		registry.creators = make(map[string]ComponentCreator)
	}
	if creator == nil {
		panic("component: Register component " + name + " creator is nil")
	}
	if _, dup := registry.creators[name]; dup {
		panic("component: Register called twice for component " + name)
	}
	registry.creators[name] = creator
}

// Lookup returns the component creator by name
func Lookup(name string) ComponentCreator {
	registry.creatorsMu.RLock()
	defer registry.creatorsMu.RUnlock()
	if registry.creators == nil {
		return nil
	}
	return registry.creators[name]
}
