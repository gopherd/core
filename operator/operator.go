// Package operator provides a set of generic functions that extend
// and complement Go's built-in operators and basic operations.
// It includes utilities for conditional logic, comparisons, and type manipulations.
package operator

import "github.com/gopherd/core/constraints"

// Or returns b if a is the zero value for T, otherwise returns a.
// It provides a generic "or" operation for any comparable type.
func Or[T comparable](a, b T) T {
	var zero T
	if a == zero {
		return b
	}
	return a
}

// OrFunc returns the result of calling new() if a is the zero value for T,
// otherwise returns a. It allows for lazy evaluation of the alternative value.
func OrFunc[T comparable](a T, new func() T) T {
	var zero T
	if a == zero {
		return new()
	}
	return a
}

// SetDefault sets the value of a to b if a is the zero value for T.
func SetDefault[T comparable](a *T, b T) {
	var zero T
	if *a == zero {
		*a = b
	}
}

// SetDefaultFunc sets the value of a to the result of calling new() if a is
// the zero value for T.
func SetDefaultFunc[T comparable](a *T, new func() T) {
	var zero T
	if *a == zero {
		*a = new()
	}
}

// Ternary returns a if condition is true, otherwise returns b.
// It provides a generic ternary operation for any type.
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// TernaryFunc returns the result of calling a() if condition is true,
// otherwise returns the result of calling b().
// It allows for lazy evaluation of both alternatives.
func TernaryFunc[T any](condition bool, a, b func() T) T {
	if condition {
		return a()
	}
	return b()
}

// Bool converts a boolean to a number (1 for true, 0 for false).
// It provides a generic way to convert boolean values to numeric types.
func Bool[T constraints.Number](ok bool) T {
	return Ternary[T](ok, 1, 0)
}

// Equal reports whether x and y are equal.
// It provides a generic equality check for any comparable type.
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
// and -0.0 is not greater than (is equal to) 0.0.
func Greater[T constraints.Ordered](x, y T) bool {
	return (!isNaN(x) && isNaN(y)) || x > y
}

// Asc compares x and y, returning:
//
//	-1 if x < y
//	 0 if x == y
//	+1 if x > y
//
// For floating-point types, NaN values are handled specially:
// a NaN is considered less than any non-NaN, equal to another NaN,
// and -0.0 is considered equal to 0.0.
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

// Dec compares x and y in descending order, returning:
//
//	-1 if x > y
//	 0 if x == y
//	+1 if x < y
//
// For floating-point types, NaN values are handled specially:
// a NaN is considered less than any non-NaN, equal to another NaN,
// and -0.0 is considered equal to 0.0.
func Dec[T constraints.Ordered](x, y T) int {
	return Asc(y, x)
}

// isNaN reports whether x is a NaN value.
// This function works for any ordered type, but will always return false
// for non-floating-point types.
func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}

// First returns the first argument.
// It's useful in contexts where you need to extract the first value from a set.
func First[T1 any](x1 T1, _ ...any) T1 {
	return x1
}

// Second returns the second argument.
// It's useful in contexts where you need to extract the second value from a set.
func Second[T1, T2 any](_ T1, x2 T2, _ ...any) T2 {
	return x2
}

// Third returns the third argument.
// It's useful in contexts where you need to extract the third value from a set.
func Third[T1, T2, T3 any](_ T1, _ T2, x3 T3, _ ...any) T3 {
	return x3
}

// Deref returns the value of p if it is not nil, otherwise it returns the zero value of T.
func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

// DerefOr returns the value of p if it is not nil, otherwise it returns defaultValue.
func DerefOr[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}

// AddressOf returns the address of x.
func AddressOf[T any](x T) *T {
	return &x
}
