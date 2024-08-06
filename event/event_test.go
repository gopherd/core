package event_test

import (
	"context"
	"errors"
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
		d := event.NewDispatcher[int](false)
		id := d.AddListener(event.Listen(1, func(ctx context.Context, e testEvent) error {
			return nil
		}))
		if !d.HasListener(id) {
			t.Errorf("Expected listener to be added")
		}
	})

	t.Run("RemoveListener", func(t *testing.T) {
		d := event.NewDispatcher[int](false)
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
		d := event.NewDispatcher[int](false)
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
		d := event.NewDispatcher[int](false)
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
		d := event.NewDispatcher[int](false)
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
		d := event.NewDispatcher[int](false)
		err := d.DispatchEvent(context.Background(), testEvent{eventType: 999})
		if err != nil {
			t.Errorf("Expected no error for non-existent event type, got %v", err)
		}
	})

	t.Run("MultipleListeners", func(t *testing.T) {
		d := event.NewDispatcher[int](false)
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
		d := event.NewDispatcher[int](false) // unordered dispatcher
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
		d := event.NewDispatcher[int](true) // ordered dispatcher
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
		d := event.NewDispatcher[int](false)
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
