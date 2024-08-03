package operator

import "github.com/gopherd/core/constraints"

// Or returns a || b
func Or[T comparable](a, b T) T {
	var zero T
	if a == zero {
		return b
	}
	return a
}

// OrFunc returns a || new()
func OrFunc[T comparable](a T, new func() T) T {
	var zero T
	if a == zero {
		return new()
	}
	return a
}

// Ternary returns yes ? a : b
func Ternary[T any](yes bool, a, b T) T {
	if yes {
		return a
	}
	return b
}

// TernaryFunc returns yes ? a() : b()
func TernaryFunc[T any](yes bool, a, b func() T) T {
	if yes {
		return a()
	}
	return b()
}

// Bool converts bool to number
func Bool[T constraints.Number](ok bool) T {
	return Ternary[T](ok, 1, 0)
}

// Equal reports whether x and y are equal.
func Equal[T comparable](x, y T) bool {
	return x == y
}

// Less reports whether x is less than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
func Less[T constraints.Ordered](x, y T) bool {
	return (isNaN(x) && !isNaN(y)) || x < y
}

// Greater reports whether x is greater than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
func Greater[T constraints.Ordered](x, y T) bool {
	return (isNaN(x) && !isNaN(y)) || x > y
}

// Asc returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Asc[T constraints.Ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

// Dec returns
//
//	-1 if x is greater than y,
//	 0 if x equals y,
//	+1 if x is less than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Dec[T constraints.Ordered](x, y T) int {
	return Asc(y, x)
}

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}

// First returns the first argument
func First[T1 any](x1 T1, others ...any) T1 {
	return x1
}

// Second returns the second argument
func Second[T1, T2 any](_ T1, x2 T2, others ...any) T2 {
	return x2
}

// Third returns the third argument
func Third[T1, T2, T3 any](_ T1, _ T2, x3 T3, others ...any) T3 {
	return x3
}
