package event

import (
	"fmt"

	"github.com/gopherd/core/container/pair"
)

type ID = int

// Event is the interface that wraps the basic Type method.
type Event[T comparable] interface {
	Typeof() T // Type gets type of event
}

// A Listener handles fired event
type Listener[T comparable] interface {
	// Typeof gets type of listening event
	EventType() T
	// HandleEvent handles fired event
	HandleEvent(Event[T])
}

// Listen creates a Listener by eventType and handler function
func Listen[T comparable, E Event[T], H ~func(E)](eventType T, handler H) Listener[T] {
	return listenerFunc[T, E, H]{eventType, handler}
}

type listenerFunc[T comparable, E Event[T], H ~func(E)] struct {
	eventType T
	handler   H
}

// EventType implements Listener EventType method
func (h listenerFunc[T, E, H]) EventType() T {
	return h.eventType
}

// HandleEvent implements Listener HandleEvent method
func (h listenerFunc[T, E, H]) HandleEvent(event Event[T]) {
	if e, ok := event.(E); ok {
		h.handler(e)
	} else {
		panic(fmt.Sprintf("unexpected event %T for type %v", event, event.Typeof()))
	}
}

type Dispatcher[T comparable] interface {
	AddListener(listener Listener[T]) ID
	RemoveListener(id ID) bool
	HasListener(id ID) bool
	FireEvent(event Event[T]) bool
}

// dispatcher manages event listeners
type dispatcher[T comparable] struct {
	nextid    ID
	ordered   bool
	listeners map[T][]pair.Pair[ID, Listener[T]]
	mapping   map[ID]pair.Pair[T, int]
}

func NewDispatcher[T comparable](ordered bool) *dispatcher[T] {
	return &dispatcher[T]{
		ordered:   ordered,
		listeners: make(map[T][]pair.Pair[ID, Listener[T]]),
		mapping:   make(map[ID]pair.Pair[T, int]),
	}
}

// AddListener registers a Listener
func (dispatcher *dispatcher[T]) AddListener(listener Listener[T]) ID {
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
func (dispatcher *dispatcher[T]) RemoveListener(id ID) bool {
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
func (dispatcher *dispatcher[T]) HasListener(id ID) bool {
	if dispatcher.mapping == nil {
		return false
	}
	_, ok := dispatcher.mapping[id]
	return ok
}

// FireEvent fires event
func (dispatcher *dispatcher[T]) FireEvent(event Event[T]) bool {
	if dispatcher.listeners == nil {
		return false
	}
	listeners, ok := dispatcher.listeners[event.Typeof()]
	if !ok || len(listeners) == 0 {
		return false
	}
	for i := range listeners {
		listeners[i].Second.HandleEvent(event)
	}
	return true
}
