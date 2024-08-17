// Package event provides a generic event handling system.
package event

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"

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
func Listen[H ~func(context.Context, E) error, E Event[T], T comparable](eventType T, handler H) Listener[T] {
	return listenerFunc[H, E, T]{eventType, handler}
}

type listenerFunc[H ~func(context.Context, E) error, E Event[T], T comparable] struct {
	eventType T
	handler   H
}

// EventType implements the Listener interface.
func (h listenerFunc[H, E, T]) EventType() T {
	return h.eventType
}

// HandleEvent implements the Listener interface.
func (h listenerFunc[H, E, T]) HandleEvent(ctx context.Context, event Event[T]) error {
	if e, ok := event.(E); ok {
		return h.handler(ctx, e)
	}
	return fmt.Errorf("%w: got %T for type %v", ErrUnexpectedEventType, event, event.Typeof())
}

// ListenerAdder adds a new listener and returns its ID.
type ListenerAdder[T comparable] interface {
	AddListener(Listener[T]) ListenerID
}

// ListenerRemover removes a listener by its ID.
type ListenerRemover interface {
	RemoveListener(ListenerID) bool
}

// ListenerChecker checks if a listener exists by its ID.
type ListenerChecker interface {
	HasListener(ListenerID) bool
}

// Dispatcher dispatches events.
type Dispatcher[T comparable] interface {
	DispatchEvent(context.Context, Event[T]) error
}

// EventSystem is the interface that manages listeners and dispatches events.
type EventSystem[T comparable] interface {
	ListenerAdder[T]
	ListenerRemover
	ListenerChecker
	Dispatcher[T]
}

type eventSystem[T comparable] struct {
	nextID    ListenerID
	ordered   bool
	listeners map[T][]pair.Pair[ListenerID, Listener[T]]
	mapping   map[ListenerID]pair.Pair[T, int]
}

func newDispatcher[T comparable](ordered bool) *eventSystem[T] {
	return &eventSystem[T]{
		ordered:   ordered,
		listeners: make(map[T][]pair.Pair[ListenerID, Listener[T]]),
		mapping:   make(map[ListenerID]pair.Pair[T, int]),
	}
}

// NewEventSystem creates a new EventSystem instance.
func NewEventSystem[T comparable](ordered bool) EventSystem[T] {
	return newDispatcher[T](ordered)
}

// AddListener implements the ListenerAdder interface.
func (es *eventSystem[T]) AddListener(listener Listener[T]) ListenerID {
	es.nextID++
	id := es.nextID
	eventType := listener.EventType()
	listeners := es.listeners[eventType]
	index := len(listeners)
	es.listeners[eventType] = append(listeners, pair.New(id, listener))
	es.mapping[id] = pair.New(eventType, index)
	return id
}

// RemoveListener implements the ListenerRemover interface.
func (es *eventSystem[T]) RemoveListener(id ListenerID) bool {
	index, ok := es.mapping[id]
	if !ok {
		return false
	}
	eventType := index.First
	listeners := es.listeners[eventType]
	last := len(listeners) - 1
	if index.Second != last {
		if es.ordered {
			copy(listeners[index.Second:last], listeners[index.Second+1:])
			for i := index.Second; i < last; i++ {
				es.mapping[listeners[i].First] = pair.New(eventType, i)
			}
		} else {
			listeners[index.Second] = listeners[last]
			es.mapping[listeners[index.Second].First] = pair.New(eventType, index.Second)
		}
	}
	listeners[last].Second = nil
	es.listeners[eventType] = listeners[:last]
	delete(es.mapping, id)
	return true
}

// HasListener implements the Dispatcher interface.
func (es *eventSystem[T]) HasListener(id ListenerID) bool {
	_, ok := es.mapping[id]
	return ok
}

// DispatchEvent implements the Dispatcher interface.
func (es *eventSystem[T]) DispatchEvent(ctx context.Context, event Event[T]) error {
	listeners, ok := es.listeners[event.Typeof()]
	if !ok || len(listeners) == 0 {
		return nil
	}
	var errs []error
	for i := range listeners {
		errs = append(errs, listeners[i].Second.HandleEvent(ctx, event))
	}
	return errors.Join(errs...)
}

// Register registers an event type for encoding and decoding.
func Register[T comparable](event Event[T]) {
	gob.Register(event)
}

// Encoder encodes events to a writer.
type Encoder struct {
	encoder *gob.Encoder
}

// NewEncoder creates a new Encoder instance.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{gob.NewEncoder(w)}
}

// Encode encodes an event to the given encoder.
func Encode[T comparable](enc *Encoder, event Event[T]) error {
	return enc.encoder.Encode(event)
}

// EncodeTo encodes an event to the given writer.
func EncodeTo[T comparable](w io.Writer, event Event[T]) error {
	return Encode(NewEncoder(w), event)
}

// Decoder decodes events from a reader.
type Decoder struct {
	decoder *gob.Decoder
}

// NewDecoder creates a new Decoder instance.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{gob.NewDecoder(r)}
}

// Decode decodes an event from the given reader.
func Decode[T comparable](dec *Decoder, event Event[T]) error {
	return dec.decoder.Decode(event)
}

// DecodeFrom decodes an event from the given reader.
func DecodeFrom[T comparable](r io.Reader, event Event[T]) error {
	return Decode(NewDecoder(r), event)
}
