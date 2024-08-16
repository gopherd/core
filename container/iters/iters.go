//go:build go1.23

// Package iters provides utility functions for working with iterators and sequences.
package iters

import (
	"cmp"
	"iter"
	"slices"

	"github.com/gopherd/core/constraints"
)

// Repeat returns an iterator that generates a sequence of n elements, each with the value v.
// It panics if n is negative.
func Repeat[T any](v T, n int) iter.Seq[T] {
	if n < 0 {
		panic("n must be non-negative")
	}
	return func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(v) {
				return
			}
		}
	}
}

// Enumerate returns an iterator that generates a sequence of index-value pairs
// for each element in the slice.
//
// Example:
//
//	for i, v := range Enumerate([]string{"a", "b", "c"}) {
//		fmt.Println(i, v) // Output: 0 a \n 1 b \n 2 c
//	}
func Enumerate[S ~[]E, E any](s S) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}

// EnumerateMap returns an iterator that generates a sequence of key-value pairs
// for each entry in the map.
func EnumerateMap[M ~map[K]V, K comparable, V any](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Of returns an iterator that generates a sequence of the provided values.
func Of[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// Sort returns an iterator that generates a sorted sequence of the elements in s.
func Sort[T cmp.Ordered](s iter.Seq[T]) iter.Seq[T] {
	return Of(slices.Sorted(s)...)
}

// Zip returns an iterator that generates pairs of elements from s1 and s2.
// If one sequence is longer, remaining elements are paired with zero values.
func Zip[T any, U any](s1 iter.Seq[T], s2 iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		next, stop := iter.Pull(s2)
		defer stop()
		var zero1 T
		var zero2 U
		for v1 := range s1 {
			v2, ok := next()
			if !ok {
				v2 = zero2
			}
			if !yield(v1, v2) {
				return
			}
		}
		for {
			if v2, ok := next(); !ok {
				return
			} else if !yield(zero1, v2) {
				return
			}
		}
	}
}

// Loop returns an iterator that generates a sequence of numbers [0, end) with a step size of 1.
// It panics if end is negative.
func Loop[T constraints.Real](end T) iter.Seq[T] {
	if end < 0 {
		panic("end must be non-negative")
	}
	return func(yield func(T) bool) {
		for i := T(0); i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range returns an iterator that generates a sequence of numbers from start to end
// with a step size of step. The sequence includes start but excludes end.
//
// It panics if step is zero, or if start < end and step is negative,
// or if start > end and step is positive.
//
// Example:
//
//	for v := range Range(1, 10, 2) {
//		fmt.Println(v) // Output: 1 3 5 7 9
//	}
func Range[T cmp.Ordered](start, end, step T) iter.Seq[T] {
	var zero T
	if step == zero {
		panic("step cannot be zero")
	}
	if start < end && step < zero {
		panic("step must be positive when start < end")
	} else if start > end && step > zero {
		panic("step must be negative when start > end")
	}
	return func(yield func(T) bool) {
		if start < end {
			for i := start; i < end; i += step {
				if !yield(i) {
					return
				}
			}
		} else {
			for i := start; i > end; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// Steps returns an iterator that generates a sequence of n numbers.
// The behavior of the sequence depends on the provided start value and optional steps.
//
// Parameters:
//   - n: The number of elements to generate in the sequence.
//   - start: The starting value of the sequence.
//   - steps: Optional variadic parameter defining the increments for the sequence.
//
// If no steps are provided, it generates a sequence starting from 'start' and incrementing by 1 each time.
// If steps are provided, it generates a sequence starting from 'start', then repeatedly applying
// the steps in a cyclic manner to generate subsequent values.
//
// It panics if n is negative.
//
// Examples:
//
//	// No steps provided (increment by 1):
//	for v := range Steps(5, 10) {
//	    fmt.Print(v, " ") // Output: 10 11 12 13 14
//	}
//
//	// Single step provided:
//	for v := range Steps(5, 1, 2) {
//	    fmt.Print(v, " ") // Output: 1 3 5 7 9
//	}
//
//	// Multiple steps provided:
//	for v := range Steps(6, 1, 2, 3, 4) {
//	    fmt.Print(v, " ") // Output: 1 3 6 10 12 15
//	}
//
//	// Using negative steps:
//	for v := range Steps(5, 20, -1, -2, -3) {
//	    fmt.Print(v, " ") // Output: 20 19 17 14 13
//	}
func Steps[T constraints.Number](n int, start T, steps ...T) iter.Seq[T] {
	if n < 0 {
		panic("n must be non-negative")
	}
	return func(yield func(T) bool) {
		if len(steps) == 0 {
			for i := 0; i < n; i++ {
				if !yield(start) {
					return
				}
				start++
			}
			return
		}
		for i := 0; i < n; i++ {
			if !yield(start) {
				return
			}
			start += steps[i%len(steps)]
		}
	}
}

// Keys returns an iterator that generates a sequence of keys from the key-value sequence m.
func Keys[K any, V any](m iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator that generates a sequence of values from the key-value sequence m.
func Values[K any, V any](m iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				return
			}
		}
	}
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

// Unique returns an iterator that generates a sequence of unique elements from s.
// Adjacent duplicate elements are removed.
func Unique[T comparable](s iter.Seq[T]) iter.Seq[T] {
	var last T
	var first = true
	return func(yield func(T) bool) {
		for v := range s {
			if first || v != last {
				first = false
				last = v
				if !yield(v) {
					return
				}
			}
		}
	}
}

// UniqueFunc returns an iterator that generates a sequence of unique elements from s.
// Adjacent elements are considered duplicates if the function f returns the same value for them.
func UniqueFunc[F ~func(T, T) bool, T comparable](s iter.Seq[T], eq F) iter.Seq[T] {
	var last T
	var first = true
	return func(yield func(T) bool) {
		for v := range s {
			if first || !eq(v, last) {
				first = false
				last = v
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Map returns an iterator that applies the function f to each element in s.
func Map[T, U any](s iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Map2 returns an iterator that applies the function f to each key-value pair in m.
func Map2[K, V, U any](m iter.Seq2[K, V], f func(K, V) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for k, v := range m {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Filter returns an iterator that generates a sequence of elements from s
// for which the function f returns true.
func Filter[T any](s iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Filter2 returns an iterator that generates a sequence of key-value pairs from m
// for which the function f returns true.
func Filter2[K, V any](m iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if f(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
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

// GroupBy returns an iterator that generates a sequence of key-value pairs,
// where the key is the result of applying the function f to each element in s,
// and the value is a slice of all elements in s that produced that key.
func GroupBy[K comparable, V any](s iter.Seq[V], f func(V) K) iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		groups := make(map[K][]V)
		for v := range s {
			k := f(v)
			groups[k] = append(groups[k], v)
		}
		for k, vs := range groups {
			if !yield(k, vs) {
				return
			}
		}
	}
}

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

// Split returns an iterator that generates a sequence of chunks from s.
// The sequence is split into n chunks of approximately equal size.
// It panics if n is less than 1.
func Split[S ~[]T, T any](s S, n int) iter.Seq[[]T] {
	if n < 1 {
		panic("n must be positive")
	}
	return func(yield func([]T) bool) {
		total := len(s)
		size := total / n
		remainder := total % n
		i := 0
		for i < total {
			var chunk []T
			if remainder > 0 {
				chunk = s[i : i+size+1]
				remainder--
				i += size + 1
			} else {
				chunk = s[i : i+size]
				i += size
			}
			if !yield(chunk) {
				return
			}
		}
	}
}
