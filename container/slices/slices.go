package slices

import (
	"math/rand"

	"github.com/gopherd/core/constraints"
)

// Min retrieves mininum value of slice
func Min[S ~[]T, T constraints.Ordered](s S) T {
	var min T
	for i, v := range s {
		if i == 0 || v < min {
			min = v
		}
	}
	return min
}

// Max retrieves maxinum value of slice
func Max[S ~[]T, T constraints.Ordered](s S) T {
	var max T
	for i, v := range s {
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

// Minmax retrieves mininum and maxinum value of slice
func Minmax[S ~[]T, T constraints.Ordered](s S) T {
	var min, max T
	for i, v := range s {
		if i == 0 || v < min {
			min = v
		}
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

// MinFunc retrieves mininum value of slice
func MinFunc[
	S ~[]T,
	F ~func(T) U,
	T any,
	U constraints.Ordered,
](s S, f F) U {
	var min U
	for i := range s {
		v := f(s[i])
		if i == 0 || v < min {
			min = v
		}
	}
	return min
}

// MaxFunc retrieves maxinum value of slice
func MaxFunc[
	S ~[]T,
	F ~func(T) U,
	T any,
	U constraints.Ordered,
](s S, f F) U {
	var max U
	for i := range s {
		v := f(s[i])
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

// MinmaxFunc retrieves mininum and maxinum value of slice
func MinmaxFunc[
	S ~[]T,
	F ~func(T) U,
	T any,
	U constraints.Ordered,
](s S, f F) U {
	var min, max U
	for i := range s {
		v := f(s[i])
		if i == 0 || v < min {
			min = v
		}
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

// Map fixes elements of slice by function f
func Map[
	S ~[]T,
	F ~func(T) U,
	T any,
	U any,
](s S, f F) []U {
	d := make([]U, len(s))
	for i, v := range s {
		d[i] = f(v)
	}
	return d
}

// CopyFunc copies mapped values by function f to slice d
func CopyFunc[
	D ~[]U,
	S ~[]T,
	F ~func(T) U,
	T any,
	U any,
](d D, s S, f F) {
	for i, v := range s {
		d[i] = f(v)
	}
}

// Sum sums slice
func Sum[S ~[]T, T constraints.Number | ~string](s S) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// Sum sums slice mapped values by function f
func SumFunc[
	S ~[]T,
	F ~func(T) U,
	T constraints.Number | ~string,
	U constraints.Number | ~string,
](s S, f F) U {
	var sum U
	for _, v := range s {
		sum += f(v)
	}
	return sum
}

// Accumulate accumulates slice from begin by function f
func Accumulate[
	S ~[]T,
	F ~func(U, T) U,
	T constraints.Number | ~string,
	U constraints.Number | ~string,
](s S, f F, begin U) U {
	for _, v := range s {
		begin = f(begin, v)
	}
	return begin
}

// Mean computes mean value of slice
func Mean[S ~[]T, T constraints.Real](s S) T {
	if len(s) == 0 {
		return 0
	}
	return Sum[S, T](s) / T(len(s))
}

// MeanFunc computes mean value mapped by function f of slice
func MeanFunc[
	S ~[]T,
	F ~func(T) U,
	T constraints.Number,
	U constraints.Real,
](s S, f F) U {
	if len(s) == 0 {
		return 0
	}
	return SumFunc[S, F, T, U](s, f) / U(len(s))
}

// Equal reports whether two slices are equal: the same length and all
// elements equal.
func Equal[X ~[]T, Y ~[]T, T comparable](x X, y Y) bool {
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

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements.
func EqualFunc[
	X ~[]U,
	Y ~[]V,
	F ~func(U, V) bool,
	U any,
	V any,
](x X, y Y, f F) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if f(x[i], y[i]) {
			return false
		}
	}
	return true
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]T, T comparable](s S, v T) int {
	for i, x := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[S ~[]T, F ~func(T) bool, T any](s S, f F) int {
	for i, v := range s {
		if f(v) {
			return i
		}
	}
	return -1
}

// LastIndex returns the index of the last occurrence of v in s,
// or -1 if not present.
func LastIndex[S ~[]T, T comparable](s S, v T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the first index i satisfying f(s[i]),
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

// Splice is like String.splice of javascript
func Splice[S ~[]T, T any](s S, i, n int, inserted ...T) S {
	if len(inserted) == 0 {
		return append(s[:i], s[i+n:]...)
	}
	if len(inserted) == n {
		copy(s[i:i+n], inserted)
		return s
	}
	return append(append(s[:i], inserted...), s[i+n:]...)
}

// Shrink removes unused capacity from the slice, returning s[:len(s):len(s)].
func Shrink[S ~[]T, T any](s S) S {
	var n = len(s)
	return s[:n:n]
}

// Unique retrieves unique set from sorted slice
func Unique[S ~[]T, T comparable](s S) S {
	var n = len(s)
	if n == 0 {
		return nil
	}
	var d = make(S, 0, n)
	d = append(d, s[0])
	for i := 1; i < n; i++ {
		if s[i] != d[len(d)-1] {
			d = append(d, s[i])
		}
	}
	return d
}

// Shuffle shuffles slice
func Shuffle[S ~[]T, T any](s S) S {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

func ShuffleN[S ~[]T, T any](s S, n int) S {
	var size = len(s)
	if n < 0 || n > size {
		panic("invalid argument to ShuffleN")
	}
	if n == size {
		return Shuffle(s)
	}
	for i := 0; i < n; i++ {
		var x = size - i
		var j int
		if x > 1<<31-1 {
			j = int(rand.Int63n(int64(x))) + i
		} else {
			j = int(rand.Int31n(int32(x))) + i
		}
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Clone returns a copy of s
func Clone[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}
	var d = make(S, len(s))
	copy(d, s)
	return d
}

func Repeat[T any](n int, x T) []T {
	s := make([]T, n)
	for i := 0; i < n; i++ {
		s[i] = x
	}
	return s
}
