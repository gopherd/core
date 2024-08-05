package history_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gopherd/core/container/history"
)

func TestMap_NewMap(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 10)

	if m == nil {
		t.Fatal("NewMap returned nil")
	}

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}
}

func TestMap_String(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	s := m.String()
	if len(s) == 0 || s[0] != '{' || s[len(s)-1] != '}' {
		t.Errorf("Invalid string representation: %s", s)
	}
	if !containsPair(s, "a", 1) || !containsPair(s, "b", 2) || !containsPair(s, "c", 3) {
		t.Errorf("String representation missing expected pairs: %s", s)
	}
}

func containsPair(s string, key string, value int) bool {
	return strings.Contains(s, fmt.Sprintf("%s:%d", key, value))
}

func TestMap_Clone(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)
	m.Set("b", 2)

	newRecorder := &history.BaseRecorder{}
	clone := m.Clone(newRecorder)

	if clone.Len() != m.Len() {
		t.Errorf("Expected clone length %d, got %d", m.Len(), clone.Len())
	}

	m.Set("c", 3)
	if clone.Len() == m.Len() {
		t.Error("Clone should not be affected by changes to original map")
	}

	if v, ok := clone.Get("a"); !ok || v != 1 {
		t.Errorf("Expected clone to have key 'a' with value 1, got %v, %v", v, ok)
	}
}

func TestMap_Len(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}

	m.Set("a", 1)
	m.Set("b", 2)

	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	m.Remove("a")

	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

func TestMap_Contains(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)

	if !m.Contains("a") {
		t.Error("Expected map to contain 'a'")
	}

	if m.Contains("b") {
		t.Error("Expected map to not contain 'b'")
	}
}

func TestMap_Get(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("Expected (1, true), got (%v, %v)", v, ok)
	}

	if v, ok := m.Get("b"); ok || v != 0 {
		t.Errorf("Expected (0, false), got (%v, %v)", v, ok)
	}
}

func TestMap_Set(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	replaced := m.Set("a", 1)
	if replaced {
		t.Error("Expected Set to return false for new key")
	}

	replaced = m.Set("a", 2)
	if !replaced {
		t.Error("Expected Set to return true for existing key")
	}

	if v, ok := m.Get("a"); !ok || v != 2 {
		t.Errorf("Expected (2, true), got (%v, %v)", v, ok)
	}
}

func TestMap_Remove(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)

	v, removed := m.Remove("a")
	if !removed || v != 1 {
		t.Errorf("Expected (1, true), got (%v, %v)", v, removed)
	}

	v, removed = m.Remove("b")
	if removed || v != 0 {
		t.Errorf("Expected (0, false), got (%v, %v)", v, removed)
	}
}

func TestMap_Range(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	expectedPairs := map[string]int{"a": 1, "b": 2, "c": 3}
	visitedPairs := make(map[string]int)

	m.Range(func(k string, v int) bool {
		visitedPairs[k] = v
		return true
	})

	if !reflect.DeepEqual(expectedPairs, visitedPairs) {
		t.Errorf("Range did not visit all expected pairs. Expected %v, got %v", expectedPairs, visitedPairs)
	}

	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return count < 2
	})

	if count != 2 {
		t.Errorf("Range did not stop after returning false. Expected 2 iterations, got %d", count)
	}
}

func TestMap_Clear(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 3)

	m.Set("a", 1)
	m.Set("b", 2)

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", m.Len())
	}

	if m.Contains("a") || m.Contains("b") {
		t.Error("Map should not contain any elements after Clear")
	}
}

func TestMap_Undo(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 3)
	m.Remove("b")
	m.Clear()

	recorder.Undo() // Undo Clear
	if m.Len() != 1 || !m.Contains("a") {
		t.Errorf("After Undo Clear, expected map with 'a', got %v", m)
	}

	recorder.Undo() // Undo Remove("b")
	if m.Len() != 2 || !m.Contains("a") || !m.Contains("b") {
		t.Errorf("After Undo Remove, expected map with 'a' and 'b', got %v", m)
	}

	recorder.Undo() // Undo Set("a", 3)
	if v, _ := m.Get("a"); v != 1 {
		t.Errorf("After Undo Set, expected 'a' to be 1, got %v", v)
	}

	recorder.Undo() // Undo Set("b", 2)
	if m.Contains("b") {
		t.Errorf("After Undo Set, 'b' should not be in the map")
	}

	recorder.Undo() // Undo Set("a", 1)
	if m.Len() != 0 {
		t.Errorf("After final Undo, expected empty map, got %v", m)
	}
}

func TestMap_UndoEmptyMap(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	// Undo on an empty map should not panic
	recorder.Undo()

	if m.Len() != 0 {
		t.Errorf("Expected empty map after Undo on empty map, got %v", m)
	}
}

func TestMap_SetThenUndo(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	m.Set("a", 1)
	recorder.Undo()

	if m.Contains("a") {
		t.Errorf("After Undo Set, 'a' should not be in the map")
	}
}

func TestMap_RemoveThenUndo(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	m.Set("a", 1)
	m.Remove("a")
	recorder.Undo()

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("After Undo Remove, expected 'a' to be 1, got %v, %v", v, ok)
	}
}

func TestMap_ClearThenUndo(t *testing.T) {
	recorder := &history.BaseRecorder{}
	m := history.NewMap[string, int](recorder, 5)

	m.Set("a", 1)
	m.Set("b", 2)
	m.Clear()
	recorder.Undo()

	if m.Len() != 2 || !m.Contains("a") || !m.Contains("b") {
		t.Errorf("After Undo Clear, expected map with 'a' and 'b', got %v", m)
	}
}
