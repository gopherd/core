//go:build go1.23

// Package iters provides utility functions for working with iterators and sequences.
package iters

import (
	"cmp"
	"iter"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/pair"
)

// Min returns the minimum element in the sequence s.
// It panics if s is empty.
func Min[T cmp.Ordered](s iter.Seq[T]) T {
	var min T
	first := true
	for v := range s {
		if first || cmp.Less(v, min) {
			min = v
			first = false
		}
	}
	if first {
		panic("empty sequence")
	}
	return min
}

// MinFunc returns the element in the sequence s for which the function f returns the minimum value.
func MinFunc[T any](s iter.Seq[T], f func(T, T) int) T {
	var min T
	first := true
	for v := range s {
		if first || f(v, min) < 0 {
			min = v
			first = false
		}
	}
	if first {
		panic("empty sequence")
	}
	return min
}

// MinKey returns the minimum key in the key-value sequence m.
func MinKey[K cmp.Ordered, V any](m iter.Seq2[K, V]) K {
	var min K
	first := true
	for k := range m {
		if first || cmp.Less(k, min) {
			min = k
			first = false
		}
	}
	return min
}

// MinValue returns the minimum value in the key-value sequence m.
func MinValue[K any, V cmp.Ordered](m iter.Seq2[K, V]) V {
	var min V
	first := true
	for _, v := range m {
		if first || cmp.Less(v, min) {
			min = v
			first = false
		}
	}
	return min
}

// Min2 returns the minimum key-value pair in the key-value sequence m.
func Min2[K, V cmp.Ordered](s iter.Seq2[K, V]) (K, V) {
	var min pair.Pair[K, V]
	first := true
	for k, v := range s {
		if first || pair.Compare(pair.New(k, v), min) < 0 {
			min = pair.New(k, v)
			first = false
		}
	}
	return min.First, min.Second
}

// MinFunc2 returns the key-value pair in the key-value sequence m for which the function f returns the minimum value.
func MinFunc2[K, V any](s iter.Seq2[K, V], f func(pair.Pair[K, V], pair.Pair[K, V]) int) (K, V) {
	var min pair.Pair[K, V]
	first := true
	for k, v := range s {
		if first || f(pair.New(k, v), min) < 0 {
			min = pair.New(k, v)
			first = false
		}
	}
	return min.First, min.Second
}

// Max returns the maximum element in the sequence s.
// It panics if s is empty.
func Max[T cmp.Ordered](s iter.Seq[T]) T {
	var max T
	first := true
	for v := range s {
		if first || cmp.Less(max, v) {
			max = v
			first = false
		}
	}
	if first {
		panic("empty sequence")
	}
	return max
}

// MaxFunc returns the element in the sequence s for which the function f returns the maximum value.
func MaxFunc[T any](s iter.Seq[T], f func(T, T) int) T {
	var max T
	first := true
	for v := range s {
		if first || f(max, v) < 0 {
			max = v
			first = false
		}
	}
	if first {
		panic("empty sequence")
	}
	return max
}

// MaxKey returns the maximum key in the key-value sequence m.
func MaxKey[K cmp.Ordered, V any](m iter.Seq2[K, V]) K {
	var max K
	first := true
	for k := range m {
		if first || cmp.Less(max, k) {
			max = k
			first = false
		}
	}
	return max
}

// MaxValue returns the maximum value in the key-value sequence m.
func MaxValue[K any, V cmp.Ordered](m iter.Seq2[K, V]) V {
	var max V
	first := true
	for _, v := range m {
		if first || cmp.Less(max, v) {
			max = v
			first = false
		}
	}
	return max
}

// Max2 returns the maximum key-value pair in the key-value sequence m.
func Max2[K, V cmp.Ordered](s iter.Seq2[K, V]) (K, V) {
	var max pair.Pair[K, V]
	first := true
	for k, v := range s {
		if first || pair.Compare(max, pair.New(k, v)) < 0 {
			max = pair.New(k, v)
			first = false
		}
	}
	return max.First, max.Second
}

// MaxFunc2 returns the key-value pair in the key-value sequence m for which the function f returns the maximum value.
func MaxFunc2[K, V any](s iter.Seq2[K, V], f func(pair.Pair[K, V], pair.Pair[K, V]) int) (K, V) {
	var max pair.Pair[K, V]
	first := true
	for k, v := range s {
		if first || f(max, pair.New(k, v)) < 0 {
			max = pair.New(k, v)
			first = false
		}
	}
	return max.First, max.Second
}

// MinMax returns the minimum and maximum elements in the sequence s.
// It panics if s is empty.
func MinMax[T cmp.Ordered](s iter.Seq[T]) (min, max T) {
	first := true
	for v := range s {
		if first {
			min, max = v, v
			first = false
		} else if cmp.Less(v, min) {
			min = v
		} else if cmp.Less(max, v) {
			max = v
		}
	}
	if first {
		panic("empty sequence")
	}
	return
}

// MinMaxFunc returns the elements in the sequence s for which the function f returns the minimum and maximum values.
func MinMaxFunc[T any](s iter.Seq[T], f func(T, T) int) (min, max T) {
	first := true
	for v := range s {
		if first {
			min, max = v, v
			first = false
		} else if f(v, min) < 0 {
			min = v
		} else if f(max, v) < 0 {
			max = v
		}
	}
	if first {
		panic("empty sequence")
	}
	return
}

// MinMaxKey returns the minimum and maximum keys in the key-value sequence m.
func MinMaxKey[K cmp.Ordered, V any](m iter.Seq2[K, V]) (min, max K) {
	first := true
	for k := range m {
		if first {
			min, max = k, k
			first = false
		} else if cmp.Less(k, min) {
			min = k
		} else if cmp.Less(max, k) {
			max = k
		}
	}
	return
}

// MinMaxValue returns the key-value pairs with the minimum and maximum values in the key-value sequence m.
func MinMaxValue[K any, V cmp.Ordered](m iter.Seq2[K, V]) (min, max V) {
	first := true
	for _, v := range m {
		if first {
			min, max = v, v
			first = false
		} else if cmp.Less(v, min) {
			min = v
		} else if cmp.Less(max, v) {
			max = v
		}
	}
	return
}

// MinMax2 returns the key-value pairs with the minimum and maximum values in the key-value sequence m.
func MinMax2[K, V cmp.Ordered](s iter.Seq2[K, V]) (min, max pair.Pair[K, V]) {
	first := true
	for k, v := range s {
		if first {
			min, max = pair.New(k, v), pair.New(k, v)
			first = false
		} else if pair.Compare(pair.New(k, v), min) < 0 {
			min = pair.New(k, v)
		} else if pair.Compare(max, pair.New(k, v)) < 0 {
			max = pair.New(k, v)
		}
	}
	return
}

// MinMaxFunc2 returns the key-value pairs with the minimum and maximum values in the key-value sequence m.
func MinMaxFunc2[K, V any](s iter.Seq2[K, V], f func(pair.Pair[K, V], pair.Pair[K, V]) int) (min, max pair.Pair[K, V]) {
	first := true
	for k, v := range s {
		if first {
			min, max = pair.New(k, v), pair.New(k, v)
			first = false
		} else if f(pair.New(k, v), min) < 0 {
			min = pair.New(k, v)
		} else if f(max, pair.New(k, v)) < 0 {
			max = pair.New(k, v)
		}
	}
	return
}

// Sum returns the sum of all elements in the sequence s.
func Sum[T constraints.Number | string](s iter.Seq[T]) T {
	var sum T
	for v := range s {
		sum += v
	}
	return sum
}

// SumKeys returns the sum of all keys in the key-value sequence m.
func SumKeys[K constraints.Number | string, V any](m iter.Seq2[K, V]) K {
	var sum K
	for k := range m {
		sum += k
	}
	return sum
}

// SumValues returns the sum of all values in the key-value sequence m.
func SumValues[K any, V constraints.Number | string](m iter.Seq2[K, V]) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

// Accumulate returns the result of adding all elements in the sequence s to the initial value.
func Accumulate[T constraints.Number | string](s iter.Seq[T], initial T) T {
	for v := range s {
		initial += v
	}
	return initial
}

// AccumulateFunc applies the function f to each element in the sequence s,
// accumulating the result starting from the initial value.
func AccumulateFunc[T, U any](s iter.Seq[T], f func(U, T) U, initial U) U {
	for v := range s {
		initial = f(initial, v)
	}
	return initial
}

// AccumulateKeys returns the result of adding all keys in the key-value sequence m to the initial value.
func AccumulateKeys[K constraints.Number | string, V any](m iter.Seq2[K, V], initial K) K {
	for k := range m {
		initial += k
	}
	return initial
}

// AccumulateValues returns the result of adding all values in the key-value sequence m to the initial value.
func AccumulateValues[K any, V constraints.Number | string](m iter.Seq2[K, V], initial V) V {
	for _, v := range m {
		initial += v
	}
	return initial
}

// AccumulateKeysFunc applies the function f to each key in the key-value sequence m,
// accumulating the result starting from the initial value.
func AccumulateKeysFunc[K, V, U any](m iter.Seq2[K, V], f func(U, K) U, initial U) U {
	for k := range m {
		initial = f(initial, k)
	}
	return initial
}

// AccumulateValuesFunc applies the function f to each value in the key-value sequence m,
// accumulating the result starting from the initial value.
func AccumulateValuesFunc[K, V, U any](m iter.Seq2[K, V], f func(U, V) U, initial U) U {
	for _, v := range m {
		initial = f(initial, v)
	}
	return initial
}

// Contains returns true if the sequence s contains the target element.
func Contains[T comparable](s iter.Seq[T], target T) bool {
	for v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// ContainsFunc returns true if the sequence s contains an element for which the function f returns true.
func ContainsFunc[T any](s iter.Seq[T], f func(T) bool) bool {
	for v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// Count returns the number of elements in the sequence s.
func Count[T any](s iter.Seq[T]) int {
	var count int
	for range s {
		count++
	}
	return count
}

// CountFunc returns the number of elements in the sequence s for which the function f returns true.
func CountFunc[T any](s iter.Seq[T], f func(T) bool) int {
	var count int
	for v := range s {
		if f(v) {
			count++
		}
	}
	return count
}

// AppendSeq2 appends the key-value pairs from seq to the slice and
// returns the extended slice.
func AppendSeq2[S ~[]pair.Pair[K, V], K, V any](s S, m iter.Seq2[K, V]) S {
	for k, v := range m {
		s = append(s, pair.New(k, v))
	}
	return s
}

// Collect2 collects key-value pairs from seq into a new slice and returns it.
func Collect2[K, V any](m iter.Seq2[K, V]) []pair.Pair[K, V] {
	return AppendSeq2([]pair.Pair[K, V]{}, m)
}
