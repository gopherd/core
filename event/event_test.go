package event_test

import (
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
	var dispatcher event.Dispatcher[string]
	dispatcher.AddListener(event.Listen("test", func(e testStringEvent, params ...any) {
		fired = true
	}))
	dispatcher.Fire(testStringEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

func TestDispatchEventPointer(t *testing.T) {
	var fired bool
	var dispatcher event.Dispatcher[string]
	dispatcher.AddListener(event.Listen("test", func(e *testStringEvent, params ...any) {
		fired = true
	}))
	dispatcher.Fire(&testStringEvent{})
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
	var dispatcher event.Dispatcher[int]
	dispatcher.AddListener(event.Listen(1, func(e testIntEvent, params ...any) {
		fired = true
	}))
	dispatcher.Fire(testIntEvent{})
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
	var dispatcher event.Dispatcher[reflect.Type]
	dispatcher.AddListener(event.Listen(eventType, func(e *testTypeEvent, params ...any) {
		fired = true
	}))
	dispatcher.Fire(&testTypeEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}
