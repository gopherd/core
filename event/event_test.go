package event_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/gopherd/core/event"
)

type testEvent struct {
	eventType int
}

func (e testEvent) Typeof() int {
	return e.eventType
}

type anotherTestEvent struct {
	eventType int
}

func (e anotherTestEvent) Typeof() int {
	return e.eventType
}

func TestDispatcher(t *testing.T) {
	t.Run("AddListener", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		id := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))
		if !d.HasListener(id) {
			t.Errorf("Expected listener to be added")
		}
	})

	t.Run("RemoveListener", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		id := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))
		if !d.RemoveListener(id) {
			t.Errorf("Expected listener to be removed")
		}
		if d.HasListener(id) {
			t.Errorf("Expected listener to not exist after removal")
		}
	})

	t.Run("DispatchEvent", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		called := false
		d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			called = true
			return nil
		}))
		err := d.DispatchEvent(context.Background(), testEvent{eventType: 1})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !called {
			t.Errorf("Expected listener to be called")
		}
	})

	t.Run("DispatchEventError", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		expectedErr := errors.New("test error")
		d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return expectedErr
		}))
		err := d.DispatchEvent(context.Background(), testEvent{eventType: 1})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("TypeMismatch", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))
		err := d.DispatchEvent(context.Background(), anotherTestEvent{eventType: 1})
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if !errors.Is(err, event.ErrUnexpectedEventType) {
			t.Errorf("Expected ErrUnexpectedEventType, got %v", err)
		}
	})

	t.Run("DispatchNonExistentEventType", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		err := d.DispatchEvent(context.Background(), testEvent{eventType: 999})
		if err != nil {
			t.Errorf("Expected no error for non-existent event type, got %v", err)
		}
	})

	t.Run("MultipleListeners", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		count := 0
		for i := 0; i < 3; i++ {
			d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
				count++
				return nil
			}))
		}
		d.DispatchEvent(context.Background(), testEvent{eventType: 1})
		if count != 3 {
			t.Errorf("Expected 3 listeners to be called, got %d", count)
		}
	})

	t.Run("RemoveMiddleListenerUnordered", func(t *testing.T) {
		var called atomic.Int32
		d := event.NewEventSystem[int](false) // unordered dispatcher
		id1 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			called.Add(1)
			return nil
		}))
		id2 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			called.Add(1)
			return nil
		}))
		id3 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			called.Add(1)
			return nil
		}))

		// Remove the middle listener
		if !d.RemoveListener(id2) {
			t.Errorf("Failed to remove middle listener")
		}

		// Verify other listeners still exist
		if !d.HasListener(id1) || !d.HasListener(id3) {
			t.Errorf("Other listeners should still exist")
		}

		// Verify the removed listener no longer exists
		if d.HasListener(id2) {
			t.Errorf("Removed listener should not exist")
		}

		// Verify the remaining listeners still work
		d.DispatchEvent(context.Background(), testEvent{eventType: 1})
		if n := called.Load(); n != 2 {
			t.Errorf("Expected 2 listeners to be called, got %d", n)
		}
	})

	t.Run("RemoveMiddleListenerOrdered", func(t *testing.T) {
		d := event.NewEventSystem[int](true) // ordered dispatcher
		callOrder := []int{}
		id1 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			callOrder = append(callOrder, 1)
			return nil
		}))
		id2 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			callOrder = append(callOrder, 2)
			return nil
		}))
		id3 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			callOrder = append(callOrder, 3)
			return nil
		}))

		// Remove the middle listener
		if !d.RemoveListener(id2) {
			t.Errorf("Failed to remove middle listener")
		}

		// Verify other listeners still exist
		if !d.HasListener(id1) || !d.HasListener(id3) {
			t.Errorf("Other listeners should still exist")
		}

		// Verify the removed listener no longer exists
		if d.HasListener(id2) {
			t.Errorf("Removed listener should not exist")
		}

		// Verify the remaining listeners still work and are in the correct order
		d.DispatchEvent(context.Background(), testEvent{eventType: 1})
		if len(callOrder) != 2 || callOrder[0] != 1 || callOrder[1] != 3 {
			t.Errorf("Expected call order [1, 3], got %v", callOrder)
		}
	})

	t.Run("RemoveLastListener", func(t *testing.T) {
		d := event.NewEventSystem[int](false)
		id1 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))
		id2 := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))

		// Remove the last listener
		if !d.RemoveListener(id2) {
			t.Errorf("Failed to remove last listener")
		}

		// Remove a non-existent listener
		if d.RemoveListener(id2) {
			t.Errorf("Should not be able to remove non-existent listener")
		}

		// Verify the first listener still exists
		if !d.HasListener(id1) {
			t.Errorf("First listener should still exist")
		}

		// Verify the removed listener no longer exists
		if d.HasListener(id2) {
			t.Errorf("Removed listener should not exist")
		}
	})
}

type testEncodableEvent struct {
	Type    int
	Message string
}

func (e testEncodableEvent) Typeof() int {
	return e.Type
}

func TestRegister(t *testing.T) {
	event.Register(testEncodableEvent{})

	// Verify that the type is registered with gob
	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	dec := gob.NewDecoder(buffer)

	original := testEncodableEvent{Type: 1, Message: "Test"}
	if err := enc.Encode(&original); err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	var decoded testEncodableEvent
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if original != decoded {
		t.Errorf("Decoded event doesn't match original. Got %v, want %v", decoded, original)
	}
}

func TestEncoderDecoder(t *testing.T) {
	event.Register(testEncodableEvent{})

	testCases := []struct {
		name  string
		event testEncodableEvent
	}{
		{"Simple event", testEncodableEvent{Type: 1, Message: "Hello"}},
		{"Empty message", testEncodableEvent{Type: 2, Message: ""}},
		{"Zero type", testEncodableEvent{Type: 0, Message: "Zero type"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			// Test EncodeTo
			if err := event.EncodeTo(buffer, tc.event); err != nil {
				t.Fatalf("EncodeTo failed: %v", err)
			}

			// Test DecodeFrom
			var decoded testEncodableEvent
			if err := event.DecodeFrom(buffer, &decoded); err != nil {
				t.Fatalf("DecodeFrom failed: %v", err)
			}

			if tc.event != decoded {
				t.Errorf("Decoded event doesn't match original. Got %v, want %v", decoded, tc.event)
			}
		})
	}
}

func TestEncoder(t *testing.T) {
	event.Register(testEncodableEvent{})

	buffer := &bytes.Buffer{}
	encoder := event.NewEncoder(buffer)

	testEvent := testEncodableEvent{Type: 3, Message: "Encoder Test"}

	if err := event.Encode(encoder, testEvent); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded testEncodableEvent
	decoder := event.NewDecoder(buffer)

	if err := event.Decode(decoder, &decoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if testEvent != decoded {
		t.Errorf("Decoded event doesn't match original. Got %v, want %v", decoded, testEvent)
	}
}

func TestDecoderWithInvalidData(t *testing.T) {
	event.Register(testEncodableEvent{})

	invalidData := []byte{0xFF, 0xFF, 0xFF} // Invalid gob data
	decoder := event.NewDecoder(bytes.NewReader(invalidData))

	var decoded testEncodableEvent
	err := event.Decode(decoder, &decoded)
	if err == nil {
		t.Error("Expected an error when decoding invalid data, but got nil")
	}
}

func ExampleEncodeTo() {
	event.Register(testEncodableEvent{})

	testEvent := testEncodableEvent{Type: 4, Message: "Example Event"}
	buffer := &bytes.Buffer{}

	err := event.EncodeTo(buffer, testEvent)
	if err != nil {
		panic(err)
	}

	var decoded testEncodableEvent
	err = event.DecodeFrom(buffer, &decoded)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded.Type, decoded.Message)
	// Output: 4 Example Event
}
