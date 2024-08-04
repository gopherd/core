/*
Package component provides a flexible and extensible component system for building modular applications.
It defines interfaces and structures for creating, managing, and controlling the lifecycle of components
within an application.

# Components

A Component in this system is an entity that has a defined lifecycle and can be managed as part of a larger application.
Each Component goes through the following lifecycle stages:

 1. Creation (OnCreated): The component is instantiated and its initial configuration is set.
 2. Initialization (Init): The component prepares its internal state and resources.
 3. Starting (Start): The component begins its main operations or services.
 4. Shutdown: The component gracefully stops its operations.
 5. Uninitialization (Uninit): The component releases any resources it has acquired.

The lifecycle methods are called in the following order:

	OnCreated -> Init -> Start -> Shutdown -> Uninit

It's important to note that:
  - If a component has been initialized (Init called), it must be uninitialized (Uninit called).
  - If a component has been started (Start called), it must be shut down (Shutdown called).

# Component Manager

The Manager struct is responsible for handling multiple components. It ensures that:

 1. Components are initialized, started, shut down, and uninitialized in the correct order.
 2. Components are initialized and started in the order they were added to the Manager.
 3. Components are shut down and uninitialized in the reverse order they were added.

This ordering ensures that dependencies between components are respected during the application's lifecycle.

# Usage

Here's a basic example of how to use this package:

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

	// In your main application:
	manager := component.NewManager()
	myComponent := &MyComponent{}
	manager.AddComponent(myComponent)

	ctx := context.Background()
	if err := manager.Init(ctx); err != nil {
		// Handle initialization error
	}
	if err := manager.Start(ctx); err != nil {
		// Handle start error
	}

	// Application runs...

	if err := manager.Shutdown(ctx); err != nil {
		// Handle shutdown error
	}
	if err := manager.Uninit(ctx); err != nil {
		// Handle uninitialization error
	}

The package also provides a registry for component creators, allowing for dynamic component creation and management.

For more detailed information on specific types and methods, refer to the individual type and function documentation.
*/
package component
