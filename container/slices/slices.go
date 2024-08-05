// Package slices provides utility functions for slice operations.
package slices

import (
	"math/rand"

	"github.com/gopherd/core/constraints"
)

// Min returns the minimum value in the slice.
// It panics if the slice is empty.
func Min[S ~[]T, T constraints.Ordered](s S) T {
	if len(s) == 0 {
		panic("Min called on empty slice")
	}
	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum value in the slice.
// It panics if the slice is empty.
func Max[S ~[]T, T constraints.Ordered](s S) T {
	if len(s) == 0 {
		panic("Max called on empty slice")
	}
	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Minmax returns both the minimum and maximum values in the slice.
// It panics if the slice is empty.
func Minmax[S ~[]T, T constraints.Ordered](s S) (min, max T) {
	if len(s) == 0 {
		panic("Minmax called on empty slice")
	}
	min, max = s[0], s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

// MinFunc returns the minimum value in the slice after applying the given function to each element.
// It panics if the slice is empty.
func MinFunc[S ~[]T, F ~func(T) U, T any, U constraints.Ordered](s S, f F) U {
	if len(s) == 0 {
		panic("MinFunc called on empty slice")
	}
	min := f(s[0])
	for _, v := range s[1:] {
		if current := f(v); current < min {
			min = current
		}
	}
	return min
}

// MaxFunc returns the maximum value in the slice after applying the given function to each element.
// It panics if the slice is empty.
func MaxFunc[S ~[]T, F ~func(T) U, T any, U constraints.Ordered](s S, f F) U {
	if len(s) == 0 {
		panic("MaxFunc called on empty slice")
	}
	max := f(s[0])
	for _, v := range s[1:] {
		if current := f(v); current > max {
			max = current
		}
	}
	return max
}

// MinmaxFunc returns both the minimum and maximum values in the slice after applying the given function to each element.
// It panics if the slice is empty.
func MinmaxFunc[S ~[]T, F ~func(T) U, T any, U constraints.Ordered](s S, f F) (min, max U) {
	if len(s) == 0 {
		panic("MinmaxFunc called on empty slice")
	}
	min, max = f(s[0]), f(s[0])
	for _, v := range s[1:] {
		current := f(v)
		if current < min {
			min = current
		}
		if current > max {
			max = current
		}
	}
	return min, max
}

// Map returns a new slice containing the results of applying the given function to each element of the original slice.
func Map[S ~[]T, F ~func(T) U, T any, U any](s S, f F) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// CopyFunc copies the results of applying the given function to each element of the source slice into the destination slice.
// It panics if the destination slice is shorter than the source slice.
func CopyFunc[D ~[]U, S ~[]T, F ~func(T) U, T any, U any](dst D, src S, f F) {
	if len(dst) < len(src) {
		panic("destination slice is shorter than source slice")
	}
	for i, v := range src {
		dst[i] = f(v)
	}
}

// Sum returns the sum of all elements in the slice.
func Sum[S ~[]T, T constraints.Number | ~string](s S) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// SumFunc returns the sum of the results of applying the given function to each element of the slice.
func SumFunc[S ~[]T, F ~func(T) U, T any, U constraints.Number | ~string](s S, f F) U {
	var sum U
	for _, v := range s {
		sum += f(v)
	}
	return sum
}

// Accumulate applies the given function to each element of the slice, accumulating the result starting from the given initial value.
func Accumulate[S ~[]T, F ~func(U, T) U, T any, U any](s S, f F, initial U) U {
	result := initial
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

// Mean returns the arithmetic mean of the slice elements.
// It returns 0 if the slice is empty.
func Mean[S ~[]T, T constraints.Real](s S) T {
	if len(s) == 0 {
		return 0
	}
	return Sum(s) / T(len(s))
}

// MeanFunc returns the arithmetic mean of the results of applying the given function to each element of the slice.
// It returns 0 if the slice is empty.
func MeanFunc[S ~[]T, F ~func(T) U, T any, U constraints.Real](s S, f F) U {
	if len(s) == 0 {
		return 0
	}
	return SumFunc(s, f) / U(len(s))
}

// Equal reports whether two slices are equal: the same length and all elements equal.
func Equal[X, Y ~[]T, T comparable](x X, y Y) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two slices are equal using a comparison function on each pair of elements.
func EqualFunc[X ~[]T, Y ~[]U, F ~func(T, U) bool, T, U any](x X, y Y, f F) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if !f(x[i], y[i]) {
			return false
		}
	}
	return true
}

// Index returns the index of the first occurrence of v in s, or -1 if not present.
func Index[S ~[]T, T comparable](s S, v T) int {
	for i, x := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// IndexFunc returns the index of the first element satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[S ~[]T, F ~func(T) bool, T any](s S, f F) int {
	for i, v := range s {
		if f(v) {
			return i
		}
	}
	return -1
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

// Contains reports whether v is present in s.
func Contains[S ~[]T, T comparable](s S, v T) bool {
	return Index(s, v) >= 0
}

// Shrink removes unused capacity from the slice, returning s[:len(s):len(s)].
func Shrink[S ~[]T, T any](s S) S {
	return s[:len(s):len(s)]
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

// Clone returns a copy of the slice.
// If s is nil, it returns nil.
func Clone[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}
	if len(s) == 0 {
		return make(S, 0)
	}
	return append(S(nil), s...)
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
