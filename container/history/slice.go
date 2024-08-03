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

// NewSlice creates a new slice
func NewSlice[T any](recorder Recorder, len, cap int) *Slice[T] {
	return &Slice[T]{
		recorder: recorder,
		data:     make([]T, len, cap),
	}
}

func (s Slice[T]) String() string {
	var sb strings.Builder
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

func (s *Slice[T]) Clone(recorder Recorder) *Slice[T] {
	return &Slice[T]{
		recorder: recorder,
		data:     slices.Clone(s.data),
	}
}

func (s *Slice[T]) Clip() {
	s.data = slices.Clip(s.data)
}

func (s *Slice[T]) Grow(n int) {
	s.data = slices.Grow(s.data, n)
}

func (s *Slice[T]) BinarySearch(target T, cmp func(T, T) int) (int, bool) {
	return slices.BinarySearchFunc(s.data, target, cmp)
}

func (s *Slice[T]) Compare(target *Slice[T], cmp func(T, T) int) int {
	return slices.CompareFunc(s.data, target.data, cmp)
}

func (s *Slice[T]) Contains(target T, cmp func(T, T) int) bool {
	return slices.ContainsFunc(s.data, func(y T) bool {
		return cmp(target, y) == 0
	})
}

func (s *Slice[T]) ContainsFunc(f func(T) bool) bool {
	return slices.ContainsFunc(s.data, f)
}

// Len returns length of the slice
func (s *Slice[T]) Len() int {
	return len(s.data)
}

// Get returns ith element of the slice
func (s *Slice[T]) Get(i int) T {
	return s.data[i]
}

// Set sets ith element of the slice to x
func (s *Slice[T]) Set(i int, x T) {
	_ = s.data[i:i]
	s.recorder.AddRecord(&sliceUndoSet[T]{s: s, i: i, x: s.data[i]})
	s.data[i] = x
}

type sliceUndoSet[T any] struct {
	s *Slice[T]
	i int
	x T
}

func (r *sliceUndoSet[T]) Undo() {
	r.s.data[r.i] = r.x
}

// Remove removes ith element of the slice
func (s *Slice[T]) Delete(i int) T {
	_ = s.data[i:i]
	var (
		zero T
		x    = s.data[i]
		n    = len(s.data) - 1
	)
	s.recorder.AddRecord(&sliceUndoRemove[T]{s: s, i: i, x: x})
	copy(s.data[i:], s.data[i+1:])
	s.data = s.data[:n]
	s.data[n] = zero
	return x
}

type sliceUndoRemove[T any] struct {
	s *Slice[T]
	i int
	x T
}

func (r *sliceUndoRemove[T]) Undo() {
	var zero T
	var n = len(r.s.data)
	r.s.data = append(r.s.data, zero)
	copy(r.s.data[r.i+1:], r.s.data[r.i:n])
	r.s.data[r.i] = r.x
}

// Append appends elements to the slice
func (s *Slice[T]) Append(elements ...T) {
	s.recorder.AddRecord(&sliceUndoAppend[T]{s: s, n: len(elements)})
	s.data = append(s.data, elements...)
}

type sliceUndoAppend[T any] struct {
	s *Slice[T]
	n int
}

func (r *sliceUndoAppend[T]) Undo() {
	r.s.data = r.s.data[:len(r.s.data)-r.n]
}

// Insert inserts elements at ith of the slice
func (s *Slice[T]) Insert(i int, elements ...T) {
	_ = s.data[i:i]
	s.recorder.AddRecord(&sliceUndoInsert[T]{s: s, i: i, j: i + len(elements)})
	s.data = slices.Insert(s.data, i, elements...)
}

type sliceUndoInsert[T any] struct {
	s *Slice[T]
	i int
	j int
}

func (r *sliceUndoInsert[T]) Undo() {
	r.s.data = slices.Delete(r.s.data, r.i, r.j)
}

func (s *Slice[T]) Reverse() {
	s.recorder.AddRecord(&sliceUndoReverse[T]{s: s})
	slices.Reverse(s.data)
}

type sliceUndoReverse[T any] struct {
	s *Slice[T]
}

func (r *sliceUndoReverse[T]) Undo() {
	slices.Reverse(r.s.data)
}
