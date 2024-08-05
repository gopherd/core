package history

import (
	"fmt"
	"maps"
	"strings"
)

// Map is a generic map with keys of type K and values of type V that supports undo operations.
type Map[K comparable, V any] struct {
	recorder Recorder
	data     map[K]V
}

// NewMap creates a new Map with the given recorder and initial size.
func NewMap[K comparable, V any](recorder Recorder, size int) *Map[K, V] {
	return &Map[K, V]{
		recorder: recorder,
		data:     make(map[K]V, size),
	}
}

// String returns a string representation of the Map.
func (m Map[K, V]) String() string {
	var sb strings.Builder
	sb.Grow(len(m.data)*18 + 1)
	sb.WriteByte('{')
	i := 0
	for k, v := range m.data {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%v:%v", k, v)
		i++
	}
	sb.WriteByte('}')
	return sb.String()
}

// Clone creates a deep copy of the Map with a new recorder.
func (m *Map[K, V]) Clone(recorder Recorder) *Map[K, V] {
	return &Map[K, V]{
		recorder: recorder,
		data:     maps.Clone(m.data),
	}
}

// Len returns the number of elements in the Map.
func (m *Map[K, V]) Len() int {
	return len(m.data)
}

// Contains checks if the Map contains the given key.
func (m *Map[K, V]) Contains(k K) bool {
	_, ok := m.data[k]
	return ok
}

// Get retrieves the value for a key in the Map.
func (m *Map[K, V]) Get(k K) (v V, ok bool) {
	v, ok = m.data[k]
	return
}

// Set adds or updates a key-value pair in the Map.
// It returns true if an existing entry was updated.
func (m *Map[K, V]) Set(k K, v V) bool {
	old, replaced := m.data[k]
	if replaced {
		m.recorder.PushAction(&mapUndoSetAction[K, V]{m: m, k: k, v: old, replaced: true})
	} else {
		m.recorder.PushAction(&mapUndoSetAction[K, V]{m: m, k: k})
	}
	m.data[k] = v
	return replaced
}

type mapUndoSetAction[K comparable, V any] struct {
	m        *Map[K, V]
	k        K
	v        V
	replaced bool
}

func (r *mapUndoSetAction[K, V]) Undo() {
	if r.replaced {
		r.m.data[r.k] = r.v
	} else {
		delete(r.m.data, r.k)
	}
}

// Remove removes a key-value pair from the Map.
// It returns the removed value and a boolean indicating if the key was present.
func (m *Map[K, V]) Remove(k K) (v V, removed bool) {
	v, removed = m.data[k]
	if removed {
		m.recorder.PushAction(&mapUndoRemoveAction[K, V]{m: m, k: k, v: v})
	}
	delete(m.data, k)
	return
}

type mapUndoRemoveAction[K comparable, V any] struct {
	m *Map[K, V]
	k K
	v V
}

func (r *mapUndoRemoveAction[K, V]) Undo() {
	r.m.data[r.k] = r.v
}

// Range calls f sequentially for each key and value in the Map.
// If f returns false, Range stops the iteration.
func (m *Map[K, V]) Range(f func(K, V) bool) bool {
	for k, v := range m.data {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// Clear removes all elements from the Map.
func (m *Map[K, V]) Clear() {
	m.recorder.PushAction(&mapUndoClearAction[K, V]{m: m, values: maps.Clone(m.data)})
	clear(m.data)
}

type mapUndoClearAction[K comparable, V any] struct {
	m      *Map[K, V]
	values map[K]V
}

func (r *mapUndoClearAction[K, V]) Undo() {
	r.m.data = r.values
}
