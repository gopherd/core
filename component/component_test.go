package component_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/event"
	"github.com/gopherd/core/operator"
)

func TestMain(m *testing.M) {
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	code := m.Run()
	slog.SetDefault(originalLogger)
	os.Exit(code)
}

type mockEntity struct {
	*component.Manager
	event.Dispatcher[reflect.Type]
}

func newMockEntity() *mockEntity {
	return &mockEntity{
		Manager:    component.NewManager(),
		Dispatcher: event.NewDispatcher[reflect.Type](true),
	}
}

// mockComponent is a test implementation of the Component interface
type mockComponent struct {
	component.BaseComponent[struct{}]
	initCalled     bool
	startCalled    bool
	shutdownCalled bool
	uninitCalled   bool
	shouldFail     bool
}

func (m *mockComponent) Init(ctx context.Context) error {
	m.initCalled = true
	if m.shouldFail {
		return errors.New("init failed")
	}
	return nil
}

func (m *mockComponent) Start(ctx context.Context) error {
	m.startCalled = true
	if m.shouldFail {
		return errors.New("start failed")
	}
	return nil
}

func (m *mockComponent) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	if m.shouldFail {
		return errors.New("shutdown failed")
	}
	return nil
}

func (m *mockComponent) Uninit(ctx context.Context) error {
	m.uninitCalled = true
	if m.shouldFail {
		return errors.New("uninit failed")
	}
	return nil
}

func TestBaseComponent(t *testing.T) {
	t.Run("BasicFunctionality", func(t *testing.T) {
		entity := newMockEntity()
		bc := &mockComponent{}
		config := component.Config{
			UUID: "test-uuid",
			Name: "test-component",
		}
		cb, err := bc.OnCreated(entity, config)
		if err != nil {
			t.Fatalf("OnCreated failed: %v", err)
		}
		if err := cb(); err != nil {
			t.Fatalf("OnCreated failed: %v", err)
		}

		if bc.UUID() != "test-uuid" {
			t.Errorf("Expected UUID 'test-uuid', got '%s'", bc.UUID())
		}

		if bc.Name() != "test-component" {
			t.Errorf("Expected Name 'test-component', got '%s'", bc.Name())
		}

		if bc.Entity() != entity {
			t.Error("Entity not set correctly")
		}
	})

	t.Run("OnCreated", func(t *testing.T) {
		entity := newMockEntity()
		bc := &component.BaseComponent[struct{ TestField string }]{}
		config := component.Config{
			UUID:    "new-uuid",
			Name:    "new-name",
			Options: json.RawMessage(`{"TestField":"test-value"}`),
		}

		cb, err := bc.OnCreated(entity, config)
		if err != nil {
			t.Fatalf("OnCreated failed: %v", err)
		}
		if err := cb(); err != nil {
			t.Fatalf("OnCreated failed: %v", err)
		}

		if bc.UUID() != "new-uuid" {
			t.Errorf("Expected UUID 'new-uuid', got '%s'", bc.UUID())
		}

		if bc.Name() != "new-name" {
			t.Errorf("Expected Name 'new-name', got '%s'", bc.Name())
		}

		if bc.Entity() != entity {
			t.Error("Entity not set correctly")
		}

		if bc.Options().TestField != "test-value" {
			t.Errorf("Options not unmarshaled correctly, got '%s'", bc.Options().TestField)
		}
	})
}

func TestManager(t *testing.T) {
	t.Run("AddComponent", func(t *testing.T) {
		entity := newMockEntity()
		comp1 := &mockComponent{}
		comp2 := &mockComponent{}

		// Add first component
		added := entity.AddComponent(comp1)
		if added != comp1 {
			t.Error("AddComponent should return the added component")
		}

		// Add second component
		added = entity.AddComponent(comp2)
		if added != comp2 {
			t.Error("AddComponent should return the added component")
		}

		// Try to add component with duplicate UUID
		config := component.Config{UUID: "duplicate"}
		comp1.OnCreated(entity, config)
		if added := entity.AddComponent(comp1); added == nil {
			t.Error("AddComponent should return the added component")
		} else {
			comp3 := &mockComponent{}
			comp3.OnCreated(entity, config)
			added = entity.AddComponent(comp3)
			if added != nil {
				t.Error("AddComponent should return nil for duplicate UUID")
			}
		}
	})

	t.Run("GetComponent", func(t *testing.T) {
		entity := newMockEntity()
		comp := &mockComponent{}
		config := component.Config{UUID: "test-uuid"}
		comp.OnCreated(entity, config)
		entity.AddComponent(comp)

		retrieved := entity.GetComponent("test-uuid")
		if retrieved != comp {
			t.Error("GetComponent failed to retrieve the correct component")
		}

		notFound := entity.GetComponent("non-existent")
		if notFound != nil {
			t.Error("GetComponent should return nil for non-existent UUID")
		}
	})

	t.Run("LifecycleOrder", func(t *testing.T) {
		manager := component.NewManager()
		comp1 := &mockComponent{}
		comp2 := &mockComponent{}
		comp3 := &mockComponent{}

		manager.AddComponent(comp1)
		manager.AddComponent(comp2)
		manager.AddComponent(comp3)

		ctx := context.Background()

		// Test Init
		err := manager.Init(ctx)
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		if !comp1.initCalled || !comp2.initCalled || !comp3.initCalled {
			t.Error("Not all components were initialized")
		}

		// Test Start
		err = manager.Start(ctx)
		if err != nil {
			t.Fatalf("Start failed: %v", err)
		}

		if !comp1.startCalled || !comp2.startCalled || !comp3.startCalled {
			t.Error("Not all components were started")
		}

		// Test Shutdown
		err = manager.Shutdown(ctx)
		if err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}

		if !comp1.shutdownCalled || !comp2.shutdownCalled || !comp3.shutdownCalled {
			t.Error("Not all components were shut down")
		}

		// Test Uninit
		err = manager.Uninit(ctx)
		if err != nil {
			t.Fatalf("Uninit failed: %v", err)
		}

		if !comp1.uninitCalled || !comp2.uninitCalled || !comp3.uninitCalled {
			t.Error("Not all components were uninitialized")
		}
	})

	t.Run("FailureHandling", func(t *testing.T) {
		manager := component.NewManager()
		comp1 := &mockComponent{}
		comp2 := &mockComponent{shouldFail: true}
		comp3 := &mockComponent{}

		manager.AddComponent(comp1)
		manager.AddComponent(comp2)
		manager.AddComponent(comp3)

		ctx := context.Background()

		// Test Init failure
		err := manager.Init(ctx)
		if err == nil {
			t.Fatal("Init should have failed")
		}

		if !comp1.initCalled || !comp2.initCalled {
			t.Error("Components before the failing component should have been initialized")
		}

		if comp3.initCalled {
			t.Error("Components after the failing component should not have been initialized")
		}

		// Reset
		comp1.initCalled = false
		comp2.initCalled = false
		comp2.shouldFail = false

		// Test Start failure
		err = manager.Init(ctx)
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		comp2.shouldFail = true
		err = manager.Start(ctx)
		if err == nil {
			t.Fatal("Start should have failed")
		}

		if !comp1.startCalled || !comp2.startCalled {
			t.Error("Components before the failing component should have been started")
		}

		if comp3.startCalled {
			t.Error("Components after the failing component should not have been started")
		}
	})
}

func TestSequentialComponentOperations(t *testing.T) {
	entity := newMockEntity()
	componentCount := 100

	// Sequentially add components
	for i := 0; i < componentCount; i++ {
		comp := &mockComponent{}
		config := component.Config{UUID: fmt.Sprintf("comp-%d", i)}
		comp.OnCreated(entity, config)
		added := entity.AddComponent(comp)
		if added == nil {
			t.Errorf("Failed to add component: %s", comp.UUID())
		}
	}

	ctx := context.Background()
	err := entity.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err = entity.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	err = entity.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	err = entity.Uninit(ctx)
	if err != nil {
		t.Fatalf("Uninit failed: %v", err)
	}

	// Check if all components went through all lifecycle stages
	for i := 0; i < componentCount; i++ {
		uuid := fmt.Sprintf("comp-%d", i)
		comp := entity.GetComponent(uuid)
		if comp == nil {
			t.Errorf("Component %s not found", uuid)
			continue
		}
		mockComp := comp.(*mockComponent)
		if !mockComp.initCalled || !mockComp.startCalled || !mockComp.shutdownCalled || !mockComp.uninitCalled {
			t.Errorf("Component %s did not complete all lifecycle stages", uuid)
		}
	}
}

func TestRegistry(t *testing.T) {
	t.Run("RegisterAndLookup", func(t *testing.T) {
		creator := func() component.Component { return &mockComponent{} }
		component.Register("test-component", creator)

		retrieved := component.Lookup("test-component")
		if retrieved == nil {
			t.Fatal("Failed to lookup registered component creator")
		}

		comp := retrieved()
		if _, ok := comp.(*mockComponent); !ok {
			t.Error("Retrieved creator did not create expected component type")
		}
	})

	t.Run("RegisterDuplicate", func(t *testing.T) {
		creator := func() component.Component { return &mockComponent{} }
		component.Register("unique-component", creator)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic on duplicate registration, but it did not occur")
			}
		}()

		component.Register("unique-component", creator)
	})

	t.Run("RegisterNil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic on nil creator registration, but it did not occur")
			}
		}()

		component.Register("nil-component", nil)
	})

	t.Run("LookupNonExistent", func(t *testing.T) {
		retrieved := component.Lookup("non-existent")
		if retrieved != nil {
			t.Error("Lookup of non-existent component should return nil")
		}
	})
}

func TestCreateOptions(t *testing.T) {
	t.Run("ValidOptions", func(t *testing.T) {
		type TestOptions struct {
			Field1 string
			Field2 int
		}

		opts := TestOptions{
			Field1: "test",
			Field2: 42,
		}

		createdOpts := component.CreateOptions(opts)
		if string(createdOpts) != `{
    "Field1": "test",
    "Field2": 42
}
` {
			t.Errorf("Unexpected created options: %q", string(createdOpts))
		}
	})

	t.Run("InvalidOptions", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic on invalid options, but it did not occur")
			}
		}()

		component.CreateOptions(make(chan int)) // channels cannot be marshaled to JSON
	})
}

func TestResolve(t *testing.T) {
	t.Run("ExistingComponent", func(t *testing.T) {
		entity := newMockEntity()
		comp := &mockComponent{}
		config := component.Config{UUID: "test-uuid"}
		comp.OnCreated(entity, config)
		entity.AddComponent(comp)

		var resolved *mockComponent
		err := component.Resolve(&resolved, entity, "test-uuid")
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		if resolved != comp {
			t.Error("Resolve did not return the correct component")
		}
	})

	t.Run("NonExistentComponent", func(t *testing.T) {
		entity := newMockEntity()

		var com *mockComponent
		err := component.Resolve(&com, entity, "non-existent")
		if err == nil {
			t.Error("Resolve should return an error for a non-existent component")
		}
	})

	t.Run("WrongComponentType", func(t *testing.T) {
		entity := newMockEntity()
		comp := &mockComponent{}
		config := component.Config{UUID: "test-uuid"}
		comp.OnCreated(entity, config)
		entity.AddComponent(comp)

		type wrongComponent struct{}
		var com *wrongComponent
		err := component.Resolve(&com, entity, "test-uuid")
		if err == nil {
			t.Error("Resolve should return an error for a component of the wrong type")
		}
	})
}

type DBComponent interface {
	Query(query string) (string, error)
}

type dbComponent struct {
	component.BaseComponent[struct {
		Driver string
		DSN    string
	}]
}

func (com *dbComponent) Query(query string) (string, error) {
	return "result", nil
}

type RedisComponent interface {
	HGet(key, field string) (string, error)
}

type redisComponent struct {
	component.BaseComponent[struct {
		Source string
	}]
}

func (com *redisComponent) HGet(key, field string) (string, error) {
	return "value", nil
}

func TestOptionsDeps(t *testing.T) {
	var entity = newMockEntity()

	// create db component
	var db dbComponent
	db.Options().Driver = "test-driver"
	db.Options().DSN = "test-dsn"
	operator.First(db.OnCreated(entity, component.Config{
		UUID: "@db",
		Name: "db",
	}))()
	entity.AddComponent(&db)

	// create redis component
	var redis redisComponent
	redis.Options().Source = "test-source"
	operator.First(redis.OnCreated(entity, component.Config{
		UUID: "@redis",
		Name: "redis",
	}))()
	entity.AddComponent(&redis)

	t.Run("ValidOptions", func(t *testing.T) {
		type usersComponent struct {
			component.BaseComponent[struct {
				DB     component.Dependency[DBComponent]
				Nested struct {
					Redis component.Dependency[RedisComponent]
					Hello string
				}

				Oops int
			}]
		}

		// create users component
		var users usersComponent
		users.Options().DB.SetUUID("@db")
		users.Options().Nested.Redis.SetUUID("@redis")
		users.Options().Nested.Hello = "world"
		users.Options().Oops = 42
		if err := operator.First(users.OnCreated(entity, component.Config{
			UUID: "@users1",
			Name: "users",
		}))(); err != nil {
			t.Fatalf("OnCreated failed: %v", err)
		}
		entity.AddComponent(&users)

		if users.Options().DB.Component() == nil {
			t.Error("DB component not resolved")
		}
		if users.Options().Nested.Redis.Component() == nil {
			t.Error("Redis component not resolved")
		}
	})

	t.Run("InvalidRedis", func(t *testing.T) {
		type usersComponent struct {
			component.BaseComponent[struct {
				InvalidRedis component.Dependency[DBComponent]
			}]
		}

		// create users component
		var users usersComponent
		users.Options().InvalidRedis.SetUUID("@redis")
		if err := operator.First(users.OnCreated(entity, component.Config{
			UUID: "@users2",
			Name: "users",
		}))(); err == nil {
			t.Errorf("OnCreated should have failed")
		}
		entity.AddComponent(&users)
	})

	t.Run("UUIDNotFound", func(t *testing.T) {
		type usersComponent struct {
			component.BaseComponent[struct {
				UUIDNotFound component.Dependency[DBComponent]
			}]
		}

		// create users component
		var users usersComponent
		users.Options().UUIDNotFound.SetUUID("@not-found")
		if err := operator.First(users.OnCreated(entity, component.Config{
			UUID: "@users3",
			Name: "users",
		}))(); err == nil {
			t.Errorf("OnCreated should have failed")
		}
		entity.AddComponent(&users)
	})

	t.Run("TypeNotFound", func(t *testing.T) {
		type usersComponent struct {
			component.BaseComponent[struct {
				TypeNotFound component.Dependency[struct {
					component.BaseComponent[struct{}]
				}]
			}]
		}

		// create users component
		var users usersComponent
		users.Options().TypeNotFound.SetUUID("@db")
		if err := operator.First(users.OnCreated(entity, component.Config{
			UUID: "@users4",
			Name: "users",
		}))(); err == nil {
			t.Errorf("OnCreated should have failed")
		}
		entity.AddComponent(&users)
	})
}
