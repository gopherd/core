// Package sliceutil provides utility functions for slice operations.
package sliceutil

import (
	"math/rand"

	"github.com/gopherd/core/constraints"
)

// Map returns a new slice containing the results of applying the given function to each element of the original slice.
func Map[S ~[]T, F ~func(T) U, T any, U any](s S, f F) []U {
	result := make([]U, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}
	return result
}

// Sum returns the sum of all elements in the slice.
func Sum[S ~[]T, T constraints.Number | ~string](s S) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// Accumulate applies the given function to each element of the slice, accumulating the result starting from the given initial value.
func Accumulate[S ~[]T, F ~func(U, T) U, T, U any](s S, f F, initial U) U {
	result := initial
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

// LastIndex returns the index of the last occurrence of v in s, or -1 if not present.
func LastIndex[S ~[]T, T comparable](s S, v T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the index of the last element satisfying f(s[i]),
// or -1 if none do.
func LastIndexFunc[S ~[]T, F ~func(T) bool, T any](s S, f F) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// Unique returns a new slice containing only the unique elements from the sorted input slice.
func Unique[S ~[]T, T comparable](s S) S {
	if s == nil {
		return nil
	}
	if len(s) == 0 {
		return make(S, 0)
	}
	result := make(S, 0, len(s))
	result = append(result, s[0])
	for i := 1; i < len(s); i++ {
		if s[i] != result[len(result)-1] {
			result = append(result, s[i])
		}
	}
	return result
}

// Shuffle randomly shuffles the elements in the slice.
func Shuffle[S ~[]T, T any](s S) S {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

// ShuffleN randomly selects and shuffles the first n elements from the entire slice.
// It ensures that the first n elements are randomly chosen from the whole slice,
// not just shuffled among themselves. Elements after the nth position may also be affected.
// This differs from a complete shuffle as it only guarantees randomness for the first n elements.
// It panics if n is negative or greater than the length of the slice.
func ShuffleN[S ~[]T, T any](s S, n int) S {
	if n < 0 || n > len(s) {
		panic("ShuffleN: invalid number of elements to shuffle")
	}
	for i := 0; i < n; i++ {
		j := rand.Intn(len(s)-i) + i
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Repeat returns a new slice consisting of n copies of x.
func Repeat[T any](n int, x T) []T {
	if n < 0 {
		panic("Repeat: negative count")
	}
	s := make([]T, n)
	for i := range s {
		s[i] = x
	}
	return s
}
