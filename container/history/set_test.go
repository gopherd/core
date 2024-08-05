package history_test

import (
	"testing"

	"github.com/gopherd/core/container/history"
)

func TestNewSet(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)
	if set == nil {
		t.Fatal("NewSet returned nil")
	}
	if set.Len() != 0 {
		t.Errorf("Expected empty set, got length %d", set.Len())
	}
}

func TestSetAdd(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)

	// Test adding a new element
	added := set.Add(1)
	if !added {
		t.Error("Add should return true for a new element")
	}
	if set.Len() != 1 {
		t.Errorf("Expected length 1, got %d", set.Len())
	}
	if !set.Contains(1) {
		t.Error("Set should contain 1")
	}

	// Test adding an existing element
	added = set.Add(1)
	if added {
		t.Error("Add should return false for an existing element")
	}
	if set.Len() != 1 {
		t.Errorf("Expected length 1, got %d", set.Len())
	}

	// Check if the action was recorded
	if recorder.Len() != 1 {
		t.Errorf("Expected 1 recorded action, got %d", recorder.Len())
	}
}

func TestSetRemove(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)
	set.Add(1)
	set.Add(2)

	// Test removing an existing element
	removed := set.Remove(1)
	if !removed {
		t.Error("Remove should return true for an existing element")
	}
	if set.Len() != 1 {
		t.Errorf("Expected length 1, got %d", set.Len())
	}
	if set.Contains(1) {
		t.Error("Set should not contain 1")
	}

	// Test removing a non-existing element
	removed = set.Remove(3)
	if removed {
		t.Error("Remove should return false for a non-existing element")
	}
	if set.Len() != 1 {
		t.Errorf("Expected length 1, got %d", set.Len())
	}

	// Check if the action was recorded
	if recorder.Len() != 3 {
		t.Errorf("Expected 3 recorded actions, got %d", recorder.Len())
	}
}

func TestSetClear(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)
	set.Add(1)
	set.Add(2)
	set.Add(3)

	set.Clear()
	if set.Len() != 0 {
		t.Errorf("Expected empty set after Clear, got length %d", set.Len())
	}
	if set.Contains(1) || set.Contains(2) || set.Contains(3) {
		t.Error("Set should not contain any elements after Clear")
	}

	// Check if the action was recorded
	if recorder.Len() != 4 {
		t.Errorf("Expected 4 recorded actions, got %d", recorder.Len())
	}
}

func TestSetRange(t *testing.T) {
	set := history.NewSet[int](&history.BaseRecorder{}, 10)
	elements := []int{1, 2, 3, 4, 5}
	for _, e := range elements {
		set.Add(e)
	}

	visited := make(map[int]bool)
	set.Range(func(k int) bool {
		visited[k] = true
		return true
	})

	for _, e := range elements {
		if !visited[e] {
			t.Errorf("Element %d was not visited during Range", e)
		}
	}

	// Test early termination
	count := 0
	set.Range(func(k int) bool {
		count++
		return count < 3
	})
	if count != 3 {
		t.Errorf("Range should have terminated after 3 elements, but processed %d", count)
	}
}

func TestSetString(t *testing.T) {
	set := history.NewSet[int](&history.BaseRecorder{}, 10)
	set.Add(1)
	set.Add(2)
	set.Add(3)

	str := set.String()
	if str != "{1,2,3}" && str != "{1,3,2}" && str != "{2,1,3}" && str != "{2,3,1}" && str != "{3,1,2}" && str != "{3,2,1}" {
		t.Errorf("Unexpected string representation: %s", str)
	}
}

func TestSetClone(t *testing.T) {
	originalRecorder := &history.BaseRecorder{}
	originalSet := history.NewSet[int](originalRecorder, 10)
	originalSet.Add(1)
	originalSet.Add(2)

	newRecorder := &history.BaseRecorder{}
	clonedSet := originalSet.Clone(newRecorder)

	if clonedSet.Len() != originalSet.Len() {
		t.Errorf("Cloned set length %d doesn't match original set length %d", clonedSet.Len(), originalSet.Len())
	}

	if !clonedSet.Contains(1) || !clonedSet.Contains(2) {
		t.Error("Cloned set doesn't contain all elements from the original set")
	}

	// Modify the cloned set
	clonedSet.Add(3)

	if originalSet.Contains(3) {
		t.Error("Modifying cloned set should not affect the original set")
	}

	if newRecorder.Len() != 1 {
		t.Errorf("Expected 1 recorded action in the new recorder, got %d", newRecorder.Len())
	}

	if originalRecorder.Len() != 2 {
		t.Errorf("Expected 2 recorded actions in the original recorder, got %d", originalRecorder.Len())
	}
}

func TestSetUndo(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)

	set.Add(1)
	set.Add(2)
	set.Remove(1)
	set.Add(3)

	// Undo Add(3)
	if !recorder.Undo() {
		t.Error("Undo should return true")
	}
	if set.Contains(3) {
		t.Error("Set should not contain 3 after undoing Add(3)")
	}

	// Undo Remove(1)
	if !recorder.Undo() {
		t.Error("Undo should return true")
	}
	if !set.Contains(1) {
		t.Error("Set should contain 1 after undoing Remove(1)")
	}

	// Undo Add(2)
	if !recorder.Undo() {
		t.Error("Undo should return true")
	}
	if set.Contains(2) {
		t.Error("Set should not contain 2 after undoing Add(2)")
	}

	// Undo Add(1)
	if !recorder.Undo() {
		t.Error("Undo should return true")
	}
	if set.Contains(1) {
		t.Error("Set should not contain 1 after undoing Add(1)")
	}

	// Try to undo when there are no more actions
	if recorder.Undo() {
		t.Error("Undo should return false when there are no more actions")
	}

	if set.Len() != 0 {
		t.Errorf("Expected empty set after undoing all actions, got length %d", set.Len())
	}
}

func TestSetUndoAll(t *testing.T) {
	recorder := &history.BaseRecorder{}
	set := history.NewSet[int](recorder, 10)

	set.Add(1)
	set.Add(2)
	set.Remove(1)
	set.Add(3)
	set.Clear()
	set.Add(4)

	actionsUndone := recorder.UndoAll()
	if actionsUndone != 6 {
		t.Errorf("Expected 6 actions undone, got %d", actionsUndone)
	}

	if set.Len() != 0 {
		t.Errorf("Expected empty set after UndoAll, got length %d", set.Len())
	}

	// Try to undo when there are no more actions
	if recorder.Undo() {
		t.Error("Undo should return false when there are no more actions")
	}
}
