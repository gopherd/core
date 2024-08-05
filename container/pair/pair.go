// Package pair provides a generic Pair type for holding two values of any types.
package pair

// Pair represents a tuple of two values of potentially different types.
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// New creates a new Pair with the given values.
func New[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}
