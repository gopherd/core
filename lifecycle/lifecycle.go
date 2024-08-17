// Package lifecycle provides interfaces and types for managing the lifecycle of components.
package lifecycle

import (
	"context"
	"fmt"
)

// Status represents the lifecycle state of a component.
type Status int

const (
	// Created indicates that the component has been instantiated but not yet initialized.
	Created Status = iota
	// Starting indicates that the component is in the process of starting.
	Starting
	// Running indicates that the component is fully operational.
	Running
	// Stopping indicates that the component is in the process of shutting down.
	Stopping
	// Closed indicates that the component has been fully shut down.
	Closed
)

// String returns a string representation of the Status.
func (s Status) String() string {
	switch s {
	case Created:
		return "Created"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	case Closed:
		return "Closed"
	default:
		return fmt.Sprintf("Unknown(%d)", int(s))
	}
}

// Lifecycle defines the interface for components with lifecycle management.
type Lifecycle interface {
	// Init initializes the component.
	Init(context.Context) error
	// Uninit performs cleanup after the component is no longer needed.
	Uninit(context.Context) error
	// Start begins the component's main operations.
	Start(context.Context) error
	// Shutdown gracefully stops the component's operations.
	Shutdown(context.Context) error
}

// Funcs represents a set of lifecycle functions for a simple component.
type Funcs struct {
	Init     func(context.Context) error
	Start    func(context.Context) error
	Shutdown func(context.Context) error
	Uninit   func(context.Context) error
}

// BaseLifecycle provides a default implementation of the Lifecycle interface.
type BaseLifecycle struct{}

// Init implements the Init method of the Lifecycle interface.
func (*BaseLifecycle) Init(context.Context) error {
	return nil
}

// Uninit implements the Uninit method of the Lifecycle interface.
func (*BaseLifecycle) Uninit(context.Context) error {
	return nil
}

// Start implements the Start method of the Lifecycle interface.
func (*BaseLifecycle) Start(context.Context) error {
	return nil
}

// Shutdown implements the Shutdown method of the Lifecycle interface.
func (*BaseLifecycle) Shutdown(context.Context) error {
	return nil
}
