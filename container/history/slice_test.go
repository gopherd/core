package history_test

import (
	"testing"

	"github.com/gopherd/core/container/history"
)

func TestNewSlice(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 5, 10)

	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}

	// Test initial values
	for i := 0; i < s.Len(); i++ {
		if s.Get(i) != 0 {
			t.Errorf("Expected 0 at index %d, got %d", i, s.Get(i))
		}
	}
}

func TestSliceString(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3)

	expected := "[1,2,3]"
	if s.String() != expected {
		t.Errorf("Expected %s, got %s", expected, s.String())
	}
}

func TestSliceClone(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3)

	newRecorder := &history.BaseRecorder{}
	clone := s.Clone(newRecorder)

	if clone.String() != s.String() {
		t.Errorf("Clone doesn't match original. Expected %s, got %s", s.String(), clone.String())
	}

	// Modify original, clone should remain unchanged
	s.Set(0, 99)
	if clone.Get(0) == 99 {
		t.Error("Clone was affected by changes to original slice")
	}
}

func TestSliceGrow(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Grow(5)

	// Append more elements than the initial capacity
	for i := 0; i < 10; i++ {
		s.Append(i)
	}

	if s.Len() != 10 {
		t.Errorf("Expected length 10 after Grow and Append, got %d", s.Len())
	}
}

func TestSliceBinarySearch(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 3, 5, 7, 9)

	tests := []struct {
		target        int
		expectedIdx   int
		expectedFound bool
	}{
		{5, 2, true},
		{1, 0, true},
		{9, 4, true},
		{0, 0, false},
		{10, 5, false},
		{6, 3, false},
	}

	for _, tt := range tests {
		idx, found := s.BinarySearch(tt.target, func(a, b int) int { return a - b })
		if idx != tt.expectedIdx || found != tt.expectedFound {
			t.Errorf("BinarySearch(%d) = (%d, %v), expected (%d, %v)",
				tt.target, idx, found, tt.expectedIdx, tt.expectedFound)
		}
	}
}

func TestSliceCompare(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s1 := history.NewSlice[int](recorder, 0, 0)
	s2 := history.NewSlice[int](recorder, 0, 0)

	s1.Append(1, 2, 3)
	s2.Append(1, 2, 3)

	if cmp := s1.Compare(s2, func(a, b int) int { return a - b }); cmp != 0 {
		t.Errorf("Expected Compare to return 0, got %d", cmp)
	}

	s2.Set(2, 4)
	if cmp := s1.Compare(s2, func(a, b int) int { return a - b }); cmp >= 0 {
		t.Errorf("Expected Compare to return negative value, got %d", cmp)
	}

	s1.Set(0, 0)
	if cmp := s1.Compare(s2, func(a, b int) int { return a - b }); cmp >= 0 {
		t.Errorf("Expected Compare to return negative value, got %d", cmp)
	}
}

func TestSliceContains(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3, 4, 5)

	if !s.Contains(3, func(a, b int) bool { return a == b }) {
		t.Error("Expected Contains(3) to return true")
	}

	if s.Contains(6, func(a, b int) bool { return a == b }) {
		t.Error("Expected Contains(6) to return false")
	}
}

func TestSliceSet(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 3, 3)

	s.Set(1, 42)
	if s.Get(1) != 42 {
		t.Errorf("Expected 42 at index 1, got %d", s.Get(1))
	}

	// Test undo
	recorder.Undo()
	if s.Get(1) != 0 {
		t.Errorf("Expected 0 at index 1 after undo, got %d", s.Get(1))
	}
}

func TestSliceRemoveAt(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3, 4, 5)

	removed := s.RemoveAt(2)
	if removed != 3 {
		t.Errorf("Expected removed value to be 3, got %d", removed)
	}
	if s.Len() != 4 {
		t.Errorf("Expected length 4 after RemoveAt, got %d", s.Len())
	}
	if s.Get(2) != 4 {
		t.Errorf("Expected 4 at index 2 after RemoveAt, got %d", s.Get(2))
	}

	// Test undo
	recorder.Undo()
	if s.Len() != 5 {
		t.Errorf("Expected length 5 after undo, got %d", s.Len())
	}
	if s.Get(2) != 3 {
		t.Errorf("Expected 3 at index 2 after undo, got %d", s.Get(2))
	}
}

func TestSliceAppend(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)

	s.Append(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Expected length 3 after Append, got %d", s.Len())
	}

	// Test undo
	recorder.Undo()
	if s.Len() != 0 {
		t.Errorf("Expected length 0 after undo, got %d", s.Len())
	}
}

func TestSliceInsert(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 5)

	s.Insert(2, 3, 4)
	if s.Len() != 5 {
		t.Errorf("Expected length 5 after Insert, got %d", s.Len())
	}
	if s.Get(2) != 3 || s.Get(3) != 4 {
		t.Errorf("Insert didn't place elements correctly")
	}

	// Test undo
	recorder.Undo()
	if s.Len() != 3 {
		t.Errorf("Expected length 3 after undo, got %d", s.Len())
	}
	if s.Get(2) != 5 {
		t.Errorf("Expected 5 at index 2 after undo, got %d", s.Get(2))
	}
}

func TestSliceReverse(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3, 4, 5)

	s.Reverse()
	for i := 0; i < s.Len(); i++ {
		if s.Get(i) != 5-i {
			t.Errorf("Reverse failed: expected %d at index %d, got %d", 5-i, i, s.Get(i))
		}
	}

	// Test undo
	recorder.Undo()
	for i := 0; i < s.Len(); i++ {
		if s.Get(i) != i+1 {
			t.Errorf("Undo Reverse failed: expected %d at index %d, got %d", i+1, i, s.Get(i))
		}
	}
}

func TestSliceRemoveFirst(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3, 2, 4)

	removed := s.RemoveFirst(2, func(a, b int) bool { return a == b })
	if !removed {
		t.Error("RemoveFirst should return true for existing element")
	}
	if s.Len() != 4 {
		t.Errorf("Expected length 4 after RemoveFirst, got %d", s.Len())
	}
	if s.Get(1) != 3 {
		t.Errorf("Expected 3 at index 1 after RemoveFirst, got %d", s.Get(1))
	}

	removed = s.RemoveFirst(5, func(a, b int) bool { return a == b })
	if removed {
		t.Error("RemoveFirst should return false for non-existing element")
	}

	// Test undo
	recorder.Undo()
	if s.Len() != 5 {
		t.Errorf("Expected length 5 after undo, got %d", s.Len())
	}
	if s.Get(1) != 2 {
		t.Errorf("Expected 2 at index 1 after undo, got %d", s.Get(1))
	}
}

func TestSliceClear(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 0, 0)
	s.Append(1, 2, 3, 4, 5)

	s.Clear()
	if s.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", s.Len())
	}

	// Test undo
	recorder.Undo()
	if s.Len() != 5 {
		t.Errorf("Expected length 5 after undo, got %d", s.Len())
	}
	for i := 0; i < s.Len(); i++ {
		if s.Get(i) != i+1 {
			t.Errorf("Undo Clear failed: expected %d at index %d, got %d", i+1, i, s.Get(1))
		}
	}
}

func TestSliceClip(t *testing.T) {
	recorder := &history.BaseRecorder{}
	s := history.NewSlice[int](recorder, 3, 10)

	initialCap := s.Cap()
	if initialCap != 10 {
		t.Errorf("Expected initial capacity 10, got %d", initialCap)
	}

	s.Append(4, 5)
	s.Clip()

	finalCap := s.Cap()
	if finalCap != 5 {
		t.Errorf("Expected final capacity 5 after Clip, got %d", finalCap)
	}

	if s.Len() != 5 {
		t.Errorf("Expected length 5 after Clip, got %d", s.Len())
	}

	expectedValues := []int{0, 0, 0, 4, 5}
	for i, v := range expectedValues {
		if s.Get(i) != v {
			t.Errorf("Expected %d at index %d after Clip, got %d", v, i, s.Get(i))
		}
	}

	// Test that Clip is not undoable
	recorder.Undo()
	if s.Cap() != 5 {
		t.Errorf("Expected capacity to remain 5 after undo, got %d", s.Cap())
	}
}
