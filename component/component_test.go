package component_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/op"
	"github.com/gopherd/core/types"
)

func TestMain(m *testing.M) {
	// disable logging during tests
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	m.Run()
}

// mockComponent is a mock implementation of the Component interface for testing
type mockComponent struct {
	component.BaseComponent[mockOptions]
	initCalled     bool
	uninitCalled   bool
	startCalled    bool
	shutdownCalled bool

	initError     bool
	startError    bool
	shutdownError bool
	uninitError   bool
}

type mockOptions struct {
	Value string
}

func (m *mockComponent) Init(ctx context.Context) error {
	m.initCalled = true
	if m.initError {
		return errors.New("mock init error")
	}
	return nil
}

func (m *mockComponent) Uninit(ctx context.Context) error {
	m.uninitCalled = true
	if m.uninitError {
		return errors.New("mock uninit error")
	}
	return nil
}

func (m *mockComponent) Start(ctx context.Context) error {
	m.startCalled = true
	if m.startError {
		return errors.New("mock start error")
	}
	return nil
}

func (m *mockComponent) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	if m.shutdownError {
		return errors.New("mock shutdown error")
	}
	return nil
}

// mockContainer is a mock implementation of the Container interface for testing
type mockContainer struct {
	components map[string]component.Component
	logger     *slog.Logger
}

func newMockContainer() *mockContainer {
	return &mockContainer{
		components: make(map[string]component.Component),
		logger:     slog.Default(),
	}
}

func (c *mockContainer) GetComponent(uuid string) component.Component {
	return c.components[uuid]
}

func (c *mockContainer) Decoder() types.Decoder {
	return json.Unmarshal
}

func (c *mockContainer) Logger() *slog.Logger {
	return c.logger
}

type failingBaseComponent struct {
	component.BaseComponent[struct{}]
}

func (f *failingBaseComponent) Setup(container component.Container, config component.Config) error {
	return errors.New("setup failed")
}

type nonStructRefs string

type componentWithNonStructRefs struct {
	component.BaseComponentWithRefs[struct{}, nonStructRefs]
}

type deeplyNestedRefs struct {
	Ref1   component.Reference[*mockComponent]
	Nested struct {
		Ref2       component.Reference[*mockComponent]
		DeepNested struct {
			Ref3 component.Reference[*mockComponent]
		}
	}
}

type componentWithDeepRefs struct {
	component.BaseComponentWithRefs[struct{}, deeplyNestedRefs]
}

type complexRefs struct {
	DirectRef    component.Reference[*mockComponent]
	PointerRef   *component.Reference[*mockComponent]
	NonRefField  string
	NestedStruct struct {
		NestedRef component.Reference[*mockComponent]
	}
	NestedPointer *struct {
		PointerNestedRef component.Reference[*mockComponent]
	}
}

type componentWithComplexRefs struct {
	component.BaseComponentWithRefs[struct{}, complexRefs]
}

type failingResolver struct {
	component.Reference[*mockComponent]
}

func (f *failingResolver) Resolve(container component.Container) error {
	return errors.New("resolution failed")
}

func TestConfigMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		config component.Config
		want   string
	}{
		{
			name: "Full config",
			config: component.Config{
				Name:            "TestComponent",
				UUID:            "test-uuid",
				Refs:            types.RawObject(`{"ref1":"uuid1"}`),
				Options:         types.RawObject(`{"option1":"value1"}`),
				TemplateUUID:    op.Addr(types.Bool(true)),
				TemplateRefs:    op.Addr(types.Bool(false)),
				TemplateOptions: op.Addr(types.Bool(true)),
			},
			want: `{"Name":"TestComponent","UUID":"test-uuid","Refs":{"ref1":"uuid1"},"Options":{"option1":"value1"},"TemplateUUID":true,"TemplateRefs":false,"TemplateOptions":true}`,
		},
		{
			name: "Minimal config",
			config: component.Config{
				Name: "MinimalComponent",
			},
			want: `{"Name":"MinimalComponent"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.config)
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}

			if string(data) != tt.want {
				t.Errorf("Marshaled JSON does not match expected.\nGot:  %s\nWant: %s", string(data), tt.want)
			}

			var unmarshaled component.Config
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal config: %v", err)
			}

			if !reflect.DeepEqual(unmarshaled, tt.config) {
				t.Errorf("Unmarshaled config does not match original.\nGot:  %+v\nWant: %+v", unmarshaled, tt.config)
			}
		})
	}
}

func TestBaseComponentSetup(t *testing.T) {
	tests := []struct {
		name               string
		config             component.Config
		wantErr            bool
		expectedIdentifier string
	}{
		{
			name: "Valid setup",
			config: component.Config{
				Name:    "TestComponent",
				UUID:    "test-uuid",
				Options: types.RawObject(`{"Value":"test"}`),
			},
			wantErr: false,
		},
		{
			name: "Invalid options",
			config: component.Config{
				Name:    "TestComponent",
				UUID:    "test-uuid",
				Options: types.RawObject(`{"InvalidKey":"test"}`),
			},
			wantErr: false,
		},
		{
			name: "Invalid options JSON",
			config: component.Config{
				Name:    "TestComponent",
				UUID:    "test-uuid",
				Options: types.RawObject(`{invalid json`),
			},
			wantErr: true,
		},
		{
			name: "UUID contains Name",
			config: component.Config{
				Name: "TestComponent",
				UUID: "TestComponent-123",
			},
			expectedIdentifier: "#TestComponent-123",
		},
		{
			name: "UUID doesn't contain Name",
			config: component.Config{
				Name: "TestComponent",
				UUID: "123",
			},
			expectedIdentifier: "TestComponent#123",
		},
		{
			name: "Empty UUID",
			config: component.Config{
				Name: "TestComponent",
				UUID: "",
			},
			expectedIdentifier: "TestComponent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &mockComponent{}
			container := newMockContainer()

			err := mc.Setup(container, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.name == "Valid setup" && mc.Options().Value != "test" {
					t.Errorf("Unexpected option value: got %s, want test", mc.Options().Value)
				}

				if mc.Logger() == nil {
					t.Error("Logger is nil after setup")
				}
			}

			if tt.expectedIdentifier == "" {
				if mc.String() != "TestComponent#test-uuid" {
					t.Errorf("Unexpected identifier: got %s, want TestComponent#test-uuid", mc.String())
				}
			} else if mc.String() != tt.expectedIdentifier {
				t.Errorf("Unexpected identifier: got %s, want %s", mc.String(), tt.expectedIdentifier)
			}
		})
	}
}

func TestReference(t *testing.T) {
	container := newMockContainer()
	mockComp := &mockComponent{}
	container.components["test-uuid"] = mockComp

	ref := component.Ref[*mockComponent]("test-uuid")

	if ref.UUID() != "test-uuid" {
		t.Errorf("Unexpected UUID: got %s, want test-uuid", ref.UUID())
	}

	err := ref.Resolve(container)
	if err != nil {
		t.Errorf("Unexpected error resolving reference: %v", err)
	}

	if ref.Component() != mockComp {
		t.Error("Reference did not resolve to the correct component")
	}

	// Test JSON marshaling and unmarshaling
	data, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("Failed to marshal reference: %v", err)
	}

	var unmarshaledRef component.Reference[*mockComponent]
	err = json.Unmarshal(data, &unmarshaledRef)
	if err != nil {
		t.Fatalf("Failed to unmarshal reference: %v", err)
	}

	if unmarshaledRef.UUID() != "test-uuid" {
		t.Errorf("Unexpected UUID after unmarshaling: got %s, want test-uuid", unmarshaledRef.UUID())
	}

	t.Run("Incorrect type assertion", func(t *testing.T) {
		wrongTypeRef := component.Ref[*struct{}]("test-uuid")
		err := wrongTypeRef.Resolve(container)
		if err == nil || !strings.Contains(err.Error(), "unexpected component") {
			t.Errorf("Expected error for incorrect type assertion, got: %v", err)
		}
	})
}

func TestOptionalReference(t *testing.T) {
	container := newMockContainer()
	mockComp := &mockComponent{}
	container.components["test-uuid"] = mockComp

	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{"Valid UUID", "test-uuid", false},
		{"Empty UUID", "", false},
		{"Invalid UUID", "invalid-uuid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := component.OptionalRef[*mockComponent](tt.uuid)

			err := ref.Resolve(container)

			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.uuid == "" && ref.Component() != nil {
				t.Error("Empty UUID should resolve to nil component")
			}

			if tt.uuid == "test-uuid" && ref.Component() != mockComp {
				t.Error("Valid UUID did not resolve to the correct component")
			}
		})
	}
}

func TestGroup(t *testing.T) {
	ctx := context.Background()
	group := component.NewGroup()

	// Add components to the group
	for i := 0; i < 3; i++ {
		mc := &mockComponent{}
		uuid := fmt.Sprintf("test-uuid-%d", i)
		container := newMockContainer() // Create a new container for each component
		err := mc.Setup(container, component.Config{Name: fmt.Sprintf("TestComponent-%d", i), UUID: uuid})
		if err != nil {
			t.Fatalf("Failed to setup mockComponent: %v", err)
		}
		if group.AddComponent(uuid, mc) == nil {
			t.Errorf("Failed to add component with UUID %s", uuid)
		}
	}

	// Test duplicate UUID
	if group.AddComponent("test-uuid-0", &mockComponent{}) != nil {
		t.Error("Adding component with duplicate UUID should return nil")
	}

	// Test lifecycle methods
	testGroupLifecycle(t, group, ctx)

	// Test GetComponent
	comp := group.GetComponent("test-uuid-1")
	if comp == nil {
		t.Error("GetComponent failed to retrieve existing component")
	}

	comp = group.GetComponent("non-existent")
	if comp != nil {
		t.Error("GetComponent should return nil for non-existent component")
	}

	t.Run("Init failure", func(t *testing.T) {
		group := component.NewGroup()
		failingComponent := &mockComponent{initError: true}
		failingComponent.Setup(newMockContainer(), component.Config{Name: "FailingComponent", UUID: "failing-uuid"})
		group.AddComponent("failing-component", failingComponent)

		err := group.Init(context.Background())
		if err == nil {
			t.Error("Expected error during group initialization, got nil")
		}
	})

	t.Run("Start failure", func(t *testing.T) {
		group := component.NewGroup()
		failingComponent := &mockComponent{startError: true}
		failingComponent.Setup(newMockContainer(), component.Config{Name: "FailingComponent", UUID: "failing-uuid"})
		group.AddComponent("failing-component", failingComponent)

		if err := group.Init(context.Background()); err != nil {
			t.Error("Unexpected error during group initialization")
		}
		if err := group.Start(context.Background()); err == nil {
			t.Error("Expected error during group start, got nil")
		}
	})

	t.Run("Shutdown failure", func(t *testing.T) {
		group := component.NewGroup()
		failingComponent := &mockComponent{shutdownError: true}
		failingComponent.Setup(newMockContainer(), component.Config{Name: "FailingComponent", UUID: "failing-uuid"})
		group.AddComponent("failing-component", failingComponent)

		if err := group.Init(context.Background()); err != nil {
			t.Error("Unexpected error during group initialization")
		}
		if err := group.Start(context.Background()); err != nil {
			t.Error("Unexpected error during group start")
		}
		if err := group.Shutdown(context.Background()); err == nil {
			t.Error("Expected error during group shutdown, got nil")
		}
	})

	t.Run("Uninit failure", func(t *testing.T) {
		group := component.NewGroup()
		failingComponent := &mockComponent{uninitError: true}
		failingComponent.Setup(newMockContainer(), component.Config{Name: "FailingComponent", UUID: "failing-uuid"})
		group.AddComponent("failing-component", failingComponent)

		if err := group.Init(context.Background()); err != nil {
			t.Error("Unexpected error during group initialization")
		}
		if err := group.Start(context.Background()); err != nil {
			t.Error("Unexpected error during group start")
		}
		if err := group.Shutdown(context.Background()); err != nil {
			t.Error("Unexpected error during group shutdown")
		}
		if err := group.Uninit(context.Background()); err == nil {
			t.Error("Expected error during group uninit, got nil")
		}
	})
}

func testGroupLifecycle(t *testing.T, group *component.Group, ctx context.Context) {
	// Test Init
	if err := group.Init(ctx); err != nil {
		t.Errorf("Group.Init failed: %v", err)
	}

	// Test Start
	if err := group.Start(ctx); err != nil {
		t.Errorf("Group.Start failed: %v", err)
	}

	// Test Shutdown
	if err := group.Shutdown(ctx); err != nil {
		t.Errorf("Group.Shutdown failed: %v", err)
	}

	// Test Uninit
	if err := group.Uninit(ctx); err != nil {
		t.Errorf("Group.Uninit failed: %v", err)
	}

	// Verify all lifecycle methods were called on each component
	for i := 0; i < 3; i++ {
		uuid := fmt.Sprintf("test-uuid-%d", i)
		mc, ok := group.GetComponent(uuid).(*mockComponent)
		if !ok {
			t.Errorf("Failed to get mockComponent for UUID %s", uuid)
			continue
		}

		if !mc.initCalled {
			t.Errorf("Init not called for component %s", uuid)
		}
		if !mc.startCalled {
			t.Errorf("Start not called for component %s", uuid)
		}
		if !mc.shutdownCalled {
			t.Errorf("Shutdown not called for component %s", uuid)
		}
		if !mc.uninitCalled {
			t.Errorf("Uninit not called for component %s", uuid)
		}
	}
}

func TestRegisterAndCreate(t *testing.T) {
	// Register a mock component
	component.Register("MockComponent", func() component.Component {
		return &mockComponent{}
	})

	// Test creating a registered component
	comp, err := component.Create("MockComponent")
	if err != nil {
		t.Errorf("Failed to create registered component: %v", err)
	}
	if _, ok := comp.(*mockComponent); !ok {
		t.Error("Created component is not of expected type")
	}

	// Test creating an unregistered component
	_, err = component.Create("UnregisteredComponent")
	if err == nil {
		t.Error("Expected error when creating unregistered component, got nil")
	}

	// Test registering a nil creator
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when registering nil creator, but it didn't panic")
		}
	}()
	component.Register("NilCreator", nil)
}

func TestBoolPointerBehavior(t *testing.T) {
	tests := []struct {
		name     string
		boolPtr  *types.Bool
		wantJSON string
	}{
		{
			name:     "True value",
			boolPtr:  op.Addr(types.Bool(true)),
			wantJSON: "true",
		},
		{
			name:     "False value",
			boolPtr:  op.Addr(types.Bool(false)),
			wantJSON: "false",
		},
		{
			name:     "Nil pointer",
			boolPtr:  nil,
			wantJSON: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.boolPtr)
			if err != nil {
				t.Fatalf("Failed to marshal bool pointer: %v", err)
			}

			if string(data) != tt.wantJSON {
				t.Errorf("Marshaled JSON does not match expected.\nGot:  %s\nWant: %s", string(data), tt.wantJSON)
			}

			var unmarshaled *types.Bool
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal bool pointer: %v", err)
			}

			if tt.boolPtr == nil {
				if unmarshaled != nil {
					t.Errorf("Expected nil, got %v", *unmarshaled)
				}
			} else {
				if unmarshaled == nil {
					t.Error("Expected non-nil value, got nil")
				} else if *unmarshaled != *tt.boolPtr {
					t.Errorf("Unmarshaled value does not match original.\nGot:  %v\nWant: %v", *unmarshaled, *tt.boolPtr)
				}
			}
		})
	}
}

// Example test to demonstrate usage of the component package
func ExampleBaseComponent() {
	// Create a mock component
	mc := &mockComponent{}

	// Setup the component
	container := newMockContainer()
	config := component.Config{
		Name:    "ExampleComponent",
		UUID:    "example-uuid",
		Options: types.RawObject(`{"Value":"example"}`),
	}

	err := mc.Setup(container, config)
	if err != nil {
		fmt.Printf("Failed to setup component: %v\n", err)
		return
	}

	// Use the component
	fmt.Println(mc.String())
	fmt.Println(mc.Options().Value)

	// Output:
	// ExampleComponent#example-uuid
	// example
}

func ExampleGroup() {
	// Create a new group
	group := component.NewGroup()

	// Create and add components to the group
	for i := 0; i < 3; i++ {
		mc := &mockComponent{}
		uuid := fmt.Sprintf("example-uuid-%d", i)
		container := newMockContainer()
		_ = mc.Setup(container, component.Config{Name: fmt.Sprintf("ExampleComponent-%d", i), UUID: uuid})
		group.AddComponent(uuid, mc)
	}

	// Use the group to manage component lifecycle
	ctx := context.Background()
	_ = group.Init(ctx)
	_ = group.Start(ctx)

	// Simulate some work
	fmt.Println("Components are running...")

	// Shutdown and uninitialize
	_ = group.Shutdown(ctx)
	_ = group.Uninit(ctx)

	fmt.Println("All components have been shut down")

	// Output:
	// Components are running...
	// All components have been shut down
}

// TestBaseComponentWithRefs tests the BaseComponentWithRefs functionality
func TestBaseComponentWithRefs(t *testing.T) {
	type testRefs struct {
		Ref1 component.Reference[*mockComponent]
		Ref2 component.OptionalReference[*mockComponent]
	}

	type testComponent struct {
		component.BaseComponentWithRefs[mockOptions, testRefs]
	}

	container := newMockContainer()
	mc1 := &mockComponent{}
	mc2 := &mockComponent{}
	container.components["ref1"] = mc1
	container.components["ref2"] = mc2

	tc := &testComponent{}
	config := component.Config{
		Name: "TestComponentWithRefs",
		UUID: "test-uuid",
		Refs: types.RawObject(`{"Ref1":"ref1","Ref2":"ref2"}`),
	}

	err := tc.Setup(container, config)
	if err != nil {
		t.Fatalf("Failed to setup component: %v", err)
	}

	if tc.Refs().Ref1.Component() != mc1 {
		t.Error("Ref1 did not resolve to the correct component")
	}

	if tc.Refs().Ref2.Component() != mc2 {
		t.Error("Ref2 did not resolve to the correct component")
	}

	t.Run("Invalid refs JSON", func(t *testing.T) {
		tc := &testComponent{}
		config := component.Config{
			Name: "TestComponentWithRefs",
			UUID: "test-uuid",
			Refs: types.RawObject(`{invalid json`),
		}
		err := tc.Setup(container, config)
		if err == nil || !strings.Contains(err.Error(), "failed to unmarshal refs") {
			t.Errorf("Expected error for invalid refs JSON, got: %v", err)
		}
	})

	t.Run("Failed to resolve reference", func(t *testing.T) {
		mockContainer := &mockContainer{
			components: make(map[string]component.Component),
		}
		tc := &testComponent{}
		config := component.Config{
			Name: "TestComponentWithRefs",
			UUID: "test-uuid",
			Refs: types.RawObject(`{"Ref1":"non-existent-uuid"}`),
		}
		err := tc.Setup(mockContainer, config)
		if err == nil || !strings.Contains(err.Error(), "failed to resolve reference") {
			t.Errorf("Expected error for failed reference resolution, got: %v", err)
		}
	})

	t.Run("Nested refs", func(t *testing.T) {
		type nestedRefs struct {
			Ref1         component.Reference[*mockComponent]
			NestedStruct struct {
				Ref2 component.Reference[*mockComponent]
			}
		}

		type nestedComponent struct {
			component.BaseComponentWithRefs[mockOptions, nestedRefs]
		}

		container := newMockContainer()
		mc1 := &mockComponent{}
		mc2 := &mockComponent{}
		container.components["ref1"] = mc1
		container.components["ref2"] = mc2

		nc := &nestedComponent{}
		config := component.Config{
			Name: "NestedComponent",
			UUID: "nested-uuid",
			Refs: types.RawObject(`{"Ref1":"ref1","NestedStruct":{"Ref2":"ref2"}}`),
		}

		err := nc.Setup(container, config)
		if err != nil {
			t.Fatalf("Failed to setup nested component: %v", err)
		}

		if nc.Refs().Ref1.Component() != mc1 {
			t.Error("Ref1 did not resolve to the correct component")
		}

		if nc.Refs().NestedStruct.Ref2.Component() != mc2 {
			t.Error("Nested Ref2 did not resolve to the correct component")
		}
	})

	t.Run("Different field types", func(t *testing.T) {
		type mixedRefs struct {
			Ref1 component.Reference[*mockComponent]
			Ref2 *component.Reference[*mockComponent]
			Ref3 interface{}
			Ref4 string
		}

		type mixedComponent struct {
			component.BaseComponentWithRefs[mockOptions, mixedRefs]
		}

		container := newMockContainer()
		mc := &mockComponent{}
		container.components["ref1"] = mc

		mixedComp := &mixedComponent{}
		config := component.Config{
			Name: "MixedComponent",
			UUID: "mixed-uuid",
			Refs: types.RawObject(`{"Ref1":"ref1","Ref2":null,"Ref3":{},"Ref4":"not-a-ref"}`),
		}

		err := mixedComp.Setup(container, config)
		if err != nil {
			t.Fatalf("Failed to setup mixed component: %v", err)
		}

		if mixedComp.Refs().Ref1.Component() != mc {
			t.Error("Ref1 did not resolve to the correct component")
		}

		if mixedComp.Refs().Ref2 != nil {
			t.Error("Ref2 should be nil")
		}

		// Ref3 and Ref4 should not cause errors, but won't be resolved
	})

	t.Run("BaseComponent Setup failure", func(t *testing.T) {
		fc := &failingBaseComponent{}
		config := component.Config{
			Name: "FailingComponent",
			UUID: "failing-uuid",
		}
		err := fc.Setup(newMockContainer(), config)
		if err == nil || err.Error() != "setup failed" {
			t.Errorf("Expected setup failed error, got: %v", err)
		}
	})

	t.Run("Invalid refs JSON", func(t *testing.T) {
		tc := &testComponent{}
		config := component.Config{
			Name: "TestComponentWithRefs",
			UUID: "test-uuid",
			Refs: types.RawObject(`{invalid json`),
		}
		err := tc.Setup(newMockContainer(), config)
		if err == nil || !strings.Contains(err.Error(), "failed to unmarshal refs") {
			t.Errorf("Expected error for invalid refs JSON, got: %v", err)
		}
	})

	t.Run("Non-struct refs", func(t *testing.T) {
		cnsr := &componentWithNonStructRefs{}
		config := component.Config{
			Name: "NonStructRefs",
			UUID: "non-struct-uuid",
			Refs: types.RawObject(`"string ref"`),
		}
		err := cnsr.Setup(newMockContainer(), config)
		if err != nil {
			t.Errorf("Expected no error for non-struct refs, got: %v", err)
		}
	})

	t.Run("Deeply nested refs", func(t *testing.T) {
		container := newMockContainer()
		mc1 := &mockComponent{}
		mc2 := &mockComponent{}
		mc3 := &mockComponent{}
		container.components["ref1"] = mc1
		container.components["ref2"] = mc2
		container.components["ref3"] = mc3

		cdr := &componentWithDeepRefs{}
		config := component.Config{
			Name: "DeepRefs",
			UUID: "deep-refs-uuid",
			Refs: types.RawObject(`{
				"Ref1": "ref1",
				"Nested": {
					"Ref2": "ref2",
					"DeepNested": {
						"Ref3": "ref3"
					}
				}
			}`),
		}

		err := cdr.Setup(container, config)
		if err != nil {
			t.Fatalf("Failed to setup component with deep refs: %v", err)
		}

		if cdr.Refs().Ref1.Component() != mc1 {
			t.Error("Ref1 did not resolve to the correct component")
		}
		if cdr.Refs().Nested.Ref2.Component() != mc2 {
			t.Error("Nested Ref2 did not resolve to the correct component")
		}
		if cdr.Refs().Nested.DeepNested.Ref3.Component() != mc3 {
			t.Error("Deep nested Ref3 did not resolve to the correct component")
		}
	})

	t.Run("Failing resolver", func(t *testing.T) {
		type refsWithFailingResolver struct {
			FailingRef failingResolver
		}

		type componentWithFailingResolver struct {
			component.BaseComponentWithRefs[struct{}, refsWithFailingResolver]
		}

		cfr := &componentWithFailingResolver{}
		config := component.Config{
			Name: "FailingResolver",
			UUID: "failing-resolver-uuid",
			Refs: types.RawObject(`{"FailingRef": "some-uuid"}`),
		}

		err := cfr.Setup(newMockContainer(), config)
		if err == nil || !strings.Contains(err.Error(), "failed to resolve reference") {
			t.Errorf("Expected error for failing resolver, got: %v", err)
		}
	})
}

// TestComponentLifecycle tests the full lifecycle of a component
func TestComponentLifecycle(t *testing.T) {
	ctx := context.Background()
	mc := &mockComponent{}
	container := newMockContainer()
	config := component.Config{
		Name: "TestLifecycleComponent",
		UUID: "lifecycle-uuid",
	}

	err := mc.Setup(container, config)
	if err != nil {
		t.Fatalf("Failed to setup component: %v", err)
	}

	if err := mc.Init(ctx); err != nil {
		t.Errorf("Init failed: %v", err)
	}
	if !mc.initCalled {
		t.Error("Init was not called")
	}

	if err := mc.Start(ctx); err != nil {
		t.Errorf("Start failed: %v", err)
	}
	if !mc.startCalled {
		t.Error("Start was not called")
	}

	if err := mc.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
	if !mc.shutdownCalled {
		t.Error("Shutdown was not called")
	}

	if err := mc.Uninit(ctx); err != nil {
		t.Errorf("Uninit failed: %v", err)
	}
	if !mc.uninitCalled {
		t.Error("Uninit was not called")
	}
}

type mockErrorComponent struct {
	component.BaseComponent[struct{}]
}

func (m *mockErrorComponent) Init(ctx context.Context) error {
	return errors.New("init error")
}

func (m *mockErrorComponent) Start(ctx context.Context) error {
	return errors.New("start error")
}

func (m *mockErrorComponent) Shutdown(ctx context.Context) error {
	return errors.New("shutdown error")
}

func (m *mockErrorComponent) Uninit(ctx context.Context) error {
	return errors.New("uninit error")
}

// TestComponentErrorHandling tests how components handle errors during lifecycle methods
func TestComponentErrorHandling(t *testing.T) {
	ctx := context.Background()

	// Create a component that returns errors for lifecycle methods
	errorComponent := &mockErrorComponent{}

	container := newMockContainer()
	config := component.Config{
		Name: "ErrorComponent",
		UUID: "error-uuid",
	}

	err := errorComponent.Setup(container, config)
	if err != nil {
		t.Fatalf("Failed to setup component: %v", err)
	}

	if err := errorComponent.Init(ctx); err == nil || err.Error() != "init error" {
		t.Errorf("Expected init error, got: %v", err)
	}

	if err := errorComponent.Start(ctx); err == nil || err.Error() != "start error" {
		t.Errorf("Expected start error, got: %v", err)
	}

	if err := errorComponent.Shutdown(ctx); err == nil || err.Error() != "shutdown error" {
		t.Errorf("Expected shutdown error, got: %v", err)
	}

	if err := errorComponent.Uninit(ctx); err == nil || err.Error() != "uninit error" {
		t.Errorf("Expected uninit error, got: %v", err)
	}
}

// BenchmarkComponentLifecycle benchmarks the performance of component lifecycle methods
func BenchmarkComponentLifecycle(b *testing.B) {
	ctx := context.Background()
	container := newMockContainer()
	config := component.Config{
		Name: "BenchmarkComponent",
		UUID: "benchmark-uuid",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc := &mockComponent{}
		_ = mc.Setup(container, config)
		_ = mc.Init(ctx)
		_ = mc.Start(ctx)
		_ = mc.Shutdown(ctx)
		_ = mc.Uninit(ctx)
	}
}
