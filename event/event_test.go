package event_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/gopherd/core/event"
)

type testStringEvent struct {
}

func (e testStringEvent) Typeof() string {
	return "test"
}

func TestDispatchEvent(t *testing.T) {
	var fired bool
	var dispatcher = event.NewDispatcher[string](true)
	dispatcher.AddListener(event.Listen("test", func(_ context.Context, e testStringEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(context.Background(), testStringEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

func TestDispatchEventPointer(t *testing.T) {
	var fired bool
	var dispatcher = event.NewDispatcher[string](true)
	dispatcher.AddListener(event.Listen("test", func(_ context.Context, e *testStringEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(context.Background(), &testStringEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

type testIntEvent struct {
}

func (e testIntEvent) Typeof() int {
	return 1
}

func TestDispatchIntEvent(t *testing.T) {
	var fired bool
	var dispatcher = event.NewDispatcher[int](true)
	dispatcher.AddListener(event.Listen(1, func(_ context.Context, e testIntEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(context.Background(), testIntEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

type testTypeEvent struct {
}

var eventType = reflect.TypeOf((*testTypeEvent)(nil))

func (e *testTypeEvent) Typeof() reflect.Type {
	return eventType
}

func TestDispatchTypeEvent(t *testing.T) {
	var fired bool
	var dispatcher = event.NewDispatcher[reflect.Type](true)
	dispatcher.AddListener(event.Listen(eventType, func(_ context.Context, e *testTypeEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(context.Background(), &testTypeEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}
