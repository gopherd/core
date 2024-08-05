package history

import (
	"fmt"
	"slices"
	"strings"
)

// Slice contains a slice data with recorder
type Slice[T any] struct {
	recorder Recorder
	data     []T
}

// NewSlice creates a new slice with the given recorder, length, and capacity.
func NewSlice[T any](recorder Recorder, len, cap int) *Slice[T] {
	return &Slice[T]{
		recorder: recorder,
		data:     make([]T, len, cap),
	}
}

// String returns a string representation of the Slice.
func (s Slice[T]) String() string {
	var sb strings.Builder
	sb.Grow(len(s.data)*9 + 1)
	sb.WriteByte('[')
	for i, x := range s.data {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprint(&sb, x)
	}
	sb.WriteByte(']')
	return sb.String()
}

// Clone creates a deep copy of the Slice with a new recorder.
func (s *Slice[T]) Clone(recorder Recorder) *Slice[T] {
	return &Slice[T]{
		recorder: recorder,
		data:     slices.Clone(s.data),
	}
}

// Clip reduces the slice's capacity to match its length.
func (s *Slice[T]) Clip() {
	s.data = slices.Clip(s.data)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for n more elements.
func (s *Slice[T]) Grow(n int) {
	s.data = slices.Grow(s.data, n)
}

// BinarySearch searches for target in the sorted slice and returns its index and a bool indicating if it was found.
func (s *Slice[T]) BinarySearch(target T, cmp func(T, T) int) (int, bool) {
	return slices.BinarySearchFunc(s.data, target, cmp)
}

// Compare compares the slice with another slice using the provided comparison function.
func (s *Slice[T]) Compare(target *Slice[T], cmp func(T, T) int) int {
	return slices.CompareFunc(s.data, target.data, cmp)
}

// Contains reports whether target is present in the slice.
func (s *Slice[T]) Contains(target T, eq func(T, T) bool) bool {
	return slices.ContainsFunc(s.data, func(x T) bool { return eq(target, x) })
}

// Len returns the length of the slice.
func (s *Slice[T]) Len() int {
	return len(s.data)
}

// Cap returns the capacity of the slice.
func (s *Slice[T]) Cap() int {
	return cap(s.data)
}

// Get returns the i-th element of the slice.
func (s *Slice[T]) Get(i int) T {
	return s.data[i]
}

// Set sets the i-th element of the slice to x.
func (s *Slice[T]) Set(i int, x T) {
	s.recorder.PushAction(&sliceUndoSetAction[T]{s: s, i: i, oldValue: s.data[i]})
	s.data[i] = x
}

type sliceUndoSetAction[T any] struct {
	s        *Slice[T]
	i        int
	oldValue T
}

func (r *sliceUndoSetAction[T]) Undo() {
	r.s.data[r.i] = r.oldValue
}

// RemoveAt removes the i-th element from the slice and returns it.
func (s *Slice[T]) RemoveAt(i int) T {
	removedValue := s.data[i]
	s.recorder.PushAction(&sliceUndoRemoveAtAction[T]{s: s, i: i, value: removedValue})
	s.data = slices.Delete(s.data, i, i+1)
	return removedValue
}

type sliceUndoRemoveAtAction[T any] struct {
	s     *Slice[T]
	i     int
	value T
}

func (r *sliceUndoRemoveAtAction[T]) Undo() {
	r.s.data = slices.Insert(r.s.data, r.i, r.value)
}

// Append appends elements to the slice.
func (s *Slice[T]) Append(elements ...T) {
	s.recorder.PushAction(&sliceUndoAppendAction[T]{s: s, n: len(elements)})
	s.data = append(s.data, elements...)
}

type sliceUndoAppendAction[T any] struct {
	s *Slice[T]
	n int
}

func (r *sliceUndoAppendAction[T]) Undo() {
	r.s.data = r.s.data[:len(r.s.data)-r.n]
}

// Insert inserts elements at i-th position of the slice.
func (s *Slice[T]) Insert(i int, elements ...T) {
	s.recorder.PushAction(&sliceUndoInsertAction[T]{s: s, i: i, n: len(elements)})
	s.data = slices.Insert(s.data, i, elements...)
}

type sliceUndoInsertAction[T any] struct {
	s *Slice[T]
	i int
	n int
}

func (r *sliceUndoInsertAction[T]) Undo() {
	r.s.data = slices.Delete(r.s.data, r.i, r.i+r.n)
}

// Reverse reverses the order of elements in the slice.
func (s *Slice[T]) Reverse() {
	s.recorder.PushAction(&sliceUndoReverseAction[T]{s: s})
	slices.Reverse(s.data)
}

type sliceUndoReverseAction[T any] struct {
	s *Slice[T]
}

func (r *sliceUndoReverseAction[T]) Undo() {
	slices.Reverse(r.s.data)
}

// RemoveFirst removes the first occurrence of the specified value from the slice.
// It returns true if the value was found and removed.
func (s *Slice[T]) RemoveFirst(value T, eq func(a, b T) bool) bool {
	for i, v := range s.data {
		if eq(v, value) {
			s.RemoveAt(i)
			return true
		}
	}
	return false
}

// Clear removes all elements from the slice.
func (s *Slice[T]) Clear() {
	oldData := slices.Clone(s.data)
	s.recorder.PushAction(&sliceUndoClearAction[T]{s: s, oldData: oldData})
	s.data = s.data[:0]
}

type sliceUndoClearAction[T any] struct {
	s       *Slice[T]
	oldData []T
}

func (r *sliceUndoClearAction[T]) Undo() {
	r.s.data = slices.Clone(r.oldData)
}
