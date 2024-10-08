//go:build go1.23

// Package iters provides utility functions for working with iterators and sequences.
package iters

import (
	"cmp"
	"iter"

	"github.com/gopherd/core/constraints"
)

// Enumerate returns an iterator that generates a sequence of values for each element
// in the slice. The index of the element is not provided. If you need the index, use
// slices.All instead.
//
// Example:
//
//	for v := range Enumerate([]string{"a", "b", "c"}) {
//		fmt.Println(v) // Output: a \n b \n c
//	}
func Enumerate[S ~[]E, E any](s S) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// List returns an iterator that generates a sequence of the provided values.
func List[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

// Infinite returns an iterator that generates an infinite sequence of integers starting from 0.
func Infinite() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

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

// LessThan returns an iterator that generates a sequence of numbers [0, end) with a step size of 1.
// It panics if end is negative.
func LessThan[T constraints.Real](end T) iter.Seq[T] {
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
