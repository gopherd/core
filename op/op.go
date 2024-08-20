// Package op provides a set of generic functions that extend
// and complement Go's built-in operators and basic operations.
// It includes utilities for conditional logic, comparisons, and type manipulations.
package op

// Or returns b if a is the zero value for T, otherwise returns a.
func Or[T comparable](a, b T) T {
	var zero T
	if a == zero {
		return b
	}
	return a
}

// OrFunc returns the result of calling b() if a is the zero value for T,
// otherwise returns a. It allows for lazy evaluation of the alternative value.
func OrFunc[T comparable](a T, b func() T) T {
	var zero T
	if a == zero {
		return b()
	}
	return a
}

// SetOr sets the value of a to b if a is the zero value for T.
// It returns the final value of a.
func SetOr[T comparable](a *T, b T) T {
	var zero T
	if *a == zero {
		*a = b
	}
	return *a
}

// SetOrFunc sets the value of a to the result of calling b() if a is
// the zero value for T. It returns the final value of a.
func SetOrFunc[T comparable](a *T, b func() T) T {
	var zero T
	if *a == zero {
		*a = b()
	}
	return *a
}

// If returns a if condition is true, otherwise returns b.
// It provides a generic ternary operation for any type.
func If[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// IfFunc returns a if condition is true, otherwise returns the result of calling b().
func IfFunc[T any](condition bool, a T, b func() T) T {
	if condition {
		return a
	}
	return b()
}

// IfFunc2 returns the result of calling a() if condition is true,
// otherwise returns the result of calling b().
// It allows for lazy evaluation of both alternatives.
func IfFunc2[T any](condition bool, a, b func() T) T {
	if condition {
		return a()
	}
	return b()
}

// Bin converts a comparable value to a binary number (0 or 1).
// It returns 0 if the input is equal to its zero value, and 1 otherwise.
func Bin[T comparable](x T) int {
	var zero T
	if x == zero {
		return 0
	}
	return 1
}

// First returns the first argument.
// It extracts the first value from a set of arguments.
func First[T any](first T, _ ...any) T {
	return first
}

// Second returns the second argument.
// It extracts the second value from a set of arguments.
func Second[T1, T2 any](first T1, second T2, _ ...any) T2 {
	return second
}

// Third returns the third argument.
// It extracts the third value from a set of arguments.
func Third[T1, T2, T3 any](first T1, second T2, third T3, _ ...any) T3 {
	return third
}

// Deref returns the value of p if it is not nil, otherwise it returns the zero value of T.
func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

// DerefOr returns the value of p if it is not nil, otherwise it returns x.
func DerefOr[T any](p *T, x T) T {
	if p == nil {
		return x
	}
	return *p
}

// DerefOr returns the value of p if it is not nil, otherwise it returns result of calling x().
func DerefOrFunc[T any](p *T, x func() T) T {
	if p == nil {
		return x()
	}
	return *p
}

// Addr returns the address of x.
func Addr[T any](x T) *T {
	return &x
}

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Result returns err if it is not nil, otherwise it returns value.
func Result(value any, err error) any {
	if err != nil {
		return err
	}
	return value
}

// MustResult panics if err is not nil, otherwise it returns value.
// It is a convenient way to handle errors in a single line.
func MustResult[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// MustResult2 panics if err is not nil, otherwise it returns value1 and value2.
func MustResult2[T1, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return value1, value2
}

// ReverseCompare returns a comparison function that reverses the order of the original comparison function.
func ReverseCompare[T any](cmp func(T, T) int) func(T, T) int {
	return func(x, y T) int {
		return cmp(y, x)
	}
}

// Zero returns the zero value of type T.
func Zero[T any]() T {
	var zero T
	return zero
}

// Identity returns a function that returns the input value.
func Identity[T any](v T) func() T {
	return func() T {
		return v
	}
}
