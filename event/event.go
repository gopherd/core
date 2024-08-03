package event

import (
	"fmt"

	"github.com/gopherd/core/container/pair"
)

// Event is the interface that wraps the basic Type method.
type Event[T comparable] interface {
	Typeof() T // Type gets type of event
}

// A Listener handles fired event
type Listener[T comparable] interface {
	EventType() T            // EventType gets type of listening event
	Handle(Event[T], ...any) // Handle handles fired event
}

// Listen creates a Listener by eventType and handler function
func Listen[T comparable, E Event[T], H ~func(E, ...any)](eventType T, handler H) Listener[T] {
	return listenerFunc[T, E, H]{eventType, handler}
}

type listenerFunc[T comparable, E Event[T], H ~func(E, ...any)] struct {
	eventType T
	handler   H
}

// EventType implements Listener EventType method
func (h listenerFunc[T, E, H]) EventType() T {
	return h.eventType
}

// Handle implements Listener Handle method
func (h listenerFunc[T, E, H]) Handle(event Event[T], arguments ...any) {
	if e, ok := event.(E); ok {
		h.handler(e, arguments...)
	} else {
		panic(fmt.Sprintf("unexpected event %T for type %v", event, event.Typeof()))
	}
}

// Dispatcher manages event listeners
type Dispatcher[T comparable] struct {
	nextid    int
	ordered   bool
	listeners map[T][]pair.Pair[int, Listener[T]]
	mapping   map[int]pair.Pair[T, int]
}

// Ordered reports whether the listeners fired by added order
func (dispatcher *Dispatcher[T]) Ordered() bool {
	return dispatcher.ordered
}

// SetOrdered sets whether the listeners fired by added order
func (dispatcher *Dispatcher[T]) SetOrdered(ordered bool) {
	dispatcher.ordered = ordered
}

// AddListener registers a Listener
func (dispatcher *Dispatcher[T]) AddListener(listener Listener[T]) int {
	if dispatcher.listeners == nil {
		dispatcher.listeners = make(map[T][]pair.Pair[int, Listener[T]])
		dispatcher.mapping = make(map[int]pair.Pair[T, int])
	}
	dispatcher.nextid++
	var id = dispatcher.nextid
	var eventType = listener.EventType()
	var listeners = dispatcher.listeners[eventType]
	var index = len(listeners)
	dispatcher.listeners[eventType] = append(listeners, pair.Make(id, listener))
	dispatcher.mapping[id] = pair.Make(eventType, index)
	return id
}

// RemoveListener removes specified listener
func (dispatcher *Dispatcher[T]) RemoveListener(id int) bool {
	if dispatcher.listeners == nil {
		return false
	}
	index, ok := dispatcher.mapping[id]
	if !ok {
		return false
	}
	var eventType = index.First
	var listeners = dispatcher.listeners[eventType]
	var last = len(listeners) - 1
	if index.Second != last {
		if dispatcher.ordered {
			copy(listeners[index.Second:last], listeners[index.Second+1:])
			for i := index.Second; i < last; i++ {
				dispatcher.mapping[listeners[i].First] = pair.Make(eventType, i)
			}
		} else {
			listeners[index.Second] = listeners[last]
			dispatcher.mapping[listeners[index.Second].First] = pair.Make(eventType, index.Second)
		}
	}
	listeners[last].Second = nil
	dispatcher.listeners[eventType] = listeners[:last]
	delete(dispatcher.mapping, id)
	return true
}

// HasListener reports whether dispatcher has specified listener
func (dispatcher *Dispatcher[T]) HasListener(id int) bool {
	if dispatcher.mapping == nil {
		return false
	}
	_, ok := dispatcher.mapping[id]
	return ok
}

// Fire fires event
func (dispatcher *Dispatcher[T]) Fire(event Event[T], arguments ...any) bool {
	if dispatcher.listeners == nil {
		return false
	}
	listeners, ok := dispatcher.listeners[event.Typeof()]
	if !ok || len(listeners) == 0 {
		return false
	}
	for i := range listeners {
		listeners[i].Second.Handle(event, arguments...)
	}
	return true
}
