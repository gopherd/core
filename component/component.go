package component

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
)

// Options represents component options
type Options = json.RawMessage

// Config represents component configuration to create a component
type Config struct {
	UUID    string  `json:"uuid,omitempty"`
	Name    string  `json:"name"`
	Options Options `json:"options"`
}

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
	// UUID returns the unique id of the component, may be empty
	UUID() string
	// Name returns the name of the component
	Name() string
	// Entity returns the entity of the component
	Entity() Entity
}

// Component represents a generic logic component
type Component interface {
	Metadata
	// OnCreated is called when the component is created
	OnCreated(Entity, Config) error
	// Init initializes the component
	Init(context.Context) error
	// Uninit uninitializes the component
	Uninit(context.Context) error
	// Start starts the component
	Start(context.Context) error
	// Shutdown gracefully shuts down the component
	Shutdown(context.Context) error
}

var _ Component = (*BaseComponent[any])(nil)

// BaseComponent implements the Component interface
type BaseComponent[T any] struct {
	uuid, name string
	entity     Entity
	options    T
}

// Options returns the options of the component
func (com *BaseComponent[T]) Options() *T {
	return &com.options
}

// UUID implements Metadata UUID method
func (com *BaseComponent[T]) UUID() string {
	return com.uuid
}

// Name implements Metadata Name method
func (com *BaseComponent[T]) Name() string {
	return com.name
}

// Entity implements Metadata Entity method
func (com *BaseComponent[T]) Entity() Entity {
	return com.entity
}

// OnCreated implements Component OnCreated method
func (com *BaseComponent[T]) OnCreated(entity Entity, config Config) error {
	com.uuid = config.UUID
	com.name = config.Name
	com.entity = entity
	return json.Unmarshal(config.Options, &com.options)
}

// Init implements Component Init method
func (com *BaseComponent[T]) Init(_ context.Context) error {
	return nil
}

// Uninit implements Component Uninit method
func (com *BaseComponent[T]) Uninit(_ context.Context) error {
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
	components     []Component
	uuid2component map[string]Component
	initialized    int
	started        int
}

// NewManager creates a Manager
func NewManager() *Manager {
	return &Manager{
		uuid2component: make(map[string]Component),
	}
}

// Add adds a component to the manager, returns nil if the uuid is duplicated
func (m *Manager) AddComponent(com Component) Component {
	uuid := com.UUID()
	if uuid != "" {
		if _, dup := m.uuid2component[uuid]; dup {
			return nil
		}
		m.uuid2component[uuid] = com
	}
	m.components = append(m.components, com)
	return com
}

func (m *Manager) GetComponent(uuid string) Component {
	if uuid == "" {
		return nil
	}
	return m.uuid2component[uuid]
}

// Init initializes all components
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
		m.initialized++
	}
	return nil
}

func (m *Manager) Uninit(ctx context.Context) error {
	for i := m.initialized - 1; i >= 0; i-- {
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

// Start starts all components
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
		m.started++
	}
	return nil
}

// Shutdown shutdowns all components in reverse order
func (m *Manager) Shutdown(ctx context.Context) error {
	for i := m.started - 1; i >= 0; i-- {
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
