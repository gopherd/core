package history

import (
	"fmt"
	"maps"
	"strings"
)

// Set is a generic set with elements of type K that supports undo operations.
type Set[K comparable] struct {
	recorder Recorder
	data     map[K]struct{}
}

// NewSet creates a new Set with the given recorder and initial size.
func NewSet[K comparable](recorder Recorder, size int) *Set[K] {
	return &Set[K]{
		recorder: recorder,
		data:     make(map[K]struct{}, size),
	}
}

// String returns a string representation of the Set.
func (s Set[K]) String() string {
	var sb strings.Builder
	sb.Grow(len(s.data)*9 + 1)
	sb.WriteByte('{')
	i := 0
	for k := range s.data {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprint(&sb, k)
		i++
	}
	sb.WriteByte('}')
	return sb.String()
}

// Clone creates a deep copy of the Set with a new recorder.
func (s *Set[K]) Clone(recorder Recorder) *Set[K] {
	return &Set[K]{
		recorder: recorder,
		data:     maps.Clone(s.data),
	}
}

// Len returns the number of elements in the Set.
func (s *Set[K]) Len() int {
	return len(s.data)
}

// Contains checks if the Set contains the given element.
func (s *Set[K]) Contains(k K) bool {
	_, ok := s.data[k]
	return ok
}

// Add adds an element to the Set.
// It returns true if the element was not already present.
func (s *Set[K]) Add(k K) (added bool) {
	_, found := s.data[k]
	if !found {
		s.data[k] = struct{}{}
		s.recorder.PushAction(&setUndoAddAction[K]{s: s, k: k})
		added = true
	}
	return
}

type setUndoAddAction[K comparable] struct {
	s *Set[K]
	k K
}

func (r *setUndoAddAction[K]) Undo() {
	delete(r.s.data, r.k)
}

// Remove removes an element from the Set.
// It returns true if the element was present.
func (s *Set[K]) Remove(k K) (removed bool) {
	_, removed = s.data[k]
	if removed {
		s.recorder.PushAction(&setUndoRemoveAction[K]{s: s, k: k})
		delete(s.data, k)
	}
	return
}

type setUndoRemoveAction[K comparable] struct {
	s *Set[K]
	k K
}

func (r *setUndoRemoveAction[K]) Undo() {
	r.s.data[r.k] = struct{}{}
}

// Range calls f sequentially for each element in the Set.
// If f returns false, Range stops the iteration.
func (s *Set[K]) Range(f func(K) bool) bool {
	for k := range s.data {
		if !f(k) {
			return false
		}
	}
	return true
}

// Clear removes all elements from the Set.
func (s *Set[K]) Clear() {
	s.recorder.PushAction(&setUndoClearAction[K]{s: s, values: maps.Clone(s.data)})
	clear(s.data)
}

type setUndoClearAction[K comparable] struct {
	s      *Set[K]
	values map[K]struct{}
}

func (r *setUndoClearAction[K]) Undo() {
	r.s.data = r.values
}
