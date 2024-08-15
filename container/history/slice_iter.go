//go:build go1.23

package history

import (
	"iter"
	"slices"
)

// All returns an iterator over index-value pairs in the slice
func (s Slice[T]) All() iter.Seq2[int, T] {
	return slices.All(s.data)
}

// Backward returns an iterator over index-value pairs in the slice,
func (s Slice[T]) Backward() iter.Seq2[int, T] {
	return slices.Backward(s.data)
}

// Values returns an iterator that yields the slice elements in order.
func (s Slice[T]) Values() iter.Seq[T] {
	return slices.Values(s.data)
}
