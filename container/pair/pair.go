// Package pair provides a generic Pair type for holding two values of any types.
package pair

import (
	"cmp"
	"fmt"
)

// Pair represents a tuple of two values of potentially different types.
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// New creates a new Pair with the given values.
func New[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v,%v)", p.First, p.Second)
}

func Compare[T1, T2 cmp.Ordered](p1, p2 Pair[T1, T2]) int {
	if c := cmp.Compare(p1.First, p2.First); c != 0 {
		return c
	}
	return cmp.Compare(p1.Second, p2.Second)
}

func CompareFirst[T1 cmp.Ordered, T2 any](p1, p2 Pair[T1, T2]) int {
	return cmp.Compare(p1.First, p2.First)
}

func CompareSecond[T1 any, T2 cmp.Ordered](p1, p2 Pair[T1, T2]) int {
	return cmp.Compare(p1.Second, p2.Second)
}
