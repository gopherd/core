// Package event provides a generic event handling system.
package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/gopherd/core/container/pair"
)

// ErrUnexpectedEventType is the error returned when an unexpected event type is received.
var ErrUnexpectedEventType = errors.New("unexpected event type")

// ListenerID represents a unique identifier for listeners.
type ListenerID int

// Event is the interface that wraps the basic Typeof method.
type Event[T comparable] interface {
	// Typeof returns the type of the event.
	Typeof() T
}

// Listener handles fired events.
type Listener[T comparable] interface {
	// EventType returns the type of event this listener handles.
	EventType() T
	// HandleEvent processes the fired event.
	HandleEvent(context.Context, Event[T]) error
}

// Listen creates a Listener for the given event type and handler function.
func Listen[T comparable, E Event[T], H ~func(context.Context, E) error](eventType T, handler H) Listener[T] {
	return listenerFunc[T, E, H]{eventType, handler}
}

type listenerFunc[T comparable, E Event[T], H ~func(context.Context, E) error] struct {
	eventType T
	handler   H
}

// EventType implements the Listener interface.
func (h listenerFunc[T, E, H]) EventType() T {
	return h.eventType
}

// HandleEvent implements the Listener interface.
func (h listenerFunc[T, E, H]) HandleEvent(ctx context.Context, event Event[T]) error {
	if e, ok := event.(E); ok {
		return h.handler(ctx, e)
	}
	return fmt.Errorf("%w: got %T for type %v", ErrUnexpectedEventType, event, event.Typeof())
}

// Dispatcher manages event listeners and dispatches events.
type Dispatcher[T comparable] interface {
	// AddListener registers a new listener and returns its ID.
	AddListener(Listener[T]) ListenerID
	// RemoveListener removes the listener with the given ID.
	RemoveListener(ListenerID) bool
	// HasListener checks if a listener with the given ID exists.
	HasListener(ListenerID) bool
	// DispatchEvent sends an event to all registered listeners of its type.
	DispatchEvent(context.Context, Event[T]) error
}

type dispatcher[T comparable] struct {
	nextID    ListenerID
	ordered   bool
	listeners map[T][]pair.Pair[ListenerID, Listener[T]]
	mapping   map[ListenerID]pair.Pair[T, int]
}

func newDispatcher[T comparable](ordered bool) *dispatcher[T] {
	return &dispatcher[T]{
		ordered:   ordered,
		listeners: make(map[T][]pair.Pair[ListenerID, Listener[T]]),
		mapping:   make(map[ListenerID]pair.Pair[T, int]),
	}
}

// NewDispatcher creates a new Dispatcher instance.
func NewDispatcher[T comparable](ordered bool) Dispatcher[T] {
	return newDispatcher[T](ordered)
}

// AddListener implements the Dispatcher interface.
func (d *dispatcher[T]) AddListener(listener Listener[T]) ListenerID {
	d.nextID++
	id := d.nextID
	eventType := listener.EventType()
	listeners := d.listeners[eventType]
	index := len(listeners)
	d.listeners[eventType] = append(listeners, pair.New(id, listener))
	d.mapping[id] = pair.New(eventType, index)
	return id
}

// RemoveListener implements the Dispatcher interface.
func (d *dispatcher[T]) RemoveListener(id ListenerID) bool {
	index, ok := d.mapping[id]
	if !ok {
		return false
	}
	eventType := index.First
	listeners := d.listeners[eventType]
	last := len(listeners) - 1
	if index.Second != last {
		if d.ordered {
			copy(listeners[index.Second:last], listeners[index.Second+1:])
			for i := index.Second; i < last; i++ {
				d.mapping[listeners[i].First] = pair.New(eventType, i)
			}
		} else {
			listeners[index.Second] = listeners[last]
			d.mapping[listeners[index.Second].First] = pair.New(eventType, index.Second)
		}
	}
	listeners[last].Second = nil
	d.listeners[eventType] = listeners[:last]
	delete(d.mapping, id)
	return true
}

// HasListener implements the Dispatcher interface.
func (d *dispatcher[T]) HasListener(id ListenerID) bool {
	_, ok := d.mapping[id]
	return ok
}

// DispatchEvent implements the Dispatcher interface.
func (d *dispatcher[T]) DispatchEvent(ctx context.Context, event Event[T]) error {
	listeners, ok := d.listeners[event.Typeof()]
	if !ok || len(listeners) == 0 {
		return nil
	}
	var errs []error
	for i := range listeners {
		errs = append(errs, listeners[i].Second.HandleEvent(ctx, event))
	}
	return errors.Join(errs...)
}
