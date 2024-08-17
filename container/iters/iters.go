//go:build go1.23

// Package iters provides utility functions for working with iterators and sequences.
package iters

import (
	"cmp"
	"iter"
	"slices"

	"github.com/gopherd/core/container/pair"
)

// Sort returns an iterator that generates a sorted sequence of the elements in s.
func Sort[T cmp.Ordered](s iter.Seq[T]) iter.Seq[T] {
	return SortFunc(s, cmp.Compare)
}

// SortFunc returns an iterator that generates a sorted sequence of the elements in s using the comparison function.
func SortFunc[T any](s iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(yield func(T) bool) {
		s := slices.Collect(s)
		slices.SortFunc(s, cmp)
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// Sort2 returns an iterator that generates a sorted sequence of the key-value pairs in m.
func Sort2[K, V cmp.Ordered](m iter.Seq2[K, V]) iter.Seq2[K, V] {
	return SortFunc2(m, pair.Compare)
}

// SortFunc2 returns an iterator that generates a sorted sequence of the key-value pairs in m using the comparison function.
func SortFunc2[K, V any](m iter.Seq2[K, V], cmp func(pair.Pair[K, V], pair.Pair[K, V]) int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s := Collect2(m)
		slices.SortFunc(s, cmp)
		for _, p := range s {
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}

// SortKeys returns an iterator that generates a sorted sequence of the key-value pairs in m by key.
func SortKeys[K cmp.Ordered, V any](m iter.Seq2[K, V]) iter.Seq2[K, V] {
	return SortFunc2(m, pair.CompareFirst)
}

// SortValues returns an iterator that generates a sorted sequence of the key-value pairs in m by value.
func SortValues[K any, V cmp.Ordered](m iter.Seq2[K, V]) iter.Seq2[K, V] {
	return SortFunc2(m, pair.CompareSecond)
}

// Zip returns an iterator that generates pairs of elements from s1 and s2.
// If one sequence is longer, remaining elements are paired with zero values.
func Zip[T any, U any](s1 iter.Seq[T], s2 iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		next, stop := iter.Pull(s2)
		defer stop()
		for v1 := range s1 {
			v2, _ := next()
			if !yield(v1, v2) {
				return
			}
		}
		var zero1 T
		for {
			if v2, ok := next(); !ok {
				return
			} else if !yield(zero1, v2) {
				return
			}
		}
	}
}

// Concat returns an iterator that generates a sequence of elements from all input sequences.
func Concat[T any](ss ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, s := range ss {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Concat2 returns an iterator that generates a sequence of key-value pairs from all input sequences.
func Concat2[K, V any](ms ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, m := range ms {
			for k, v := range m {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// WithIndex returns an iterator that generates a sequence of index-value pairs from s.
func WithIndex[T any](s iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range s {
			if !yield(i, v) {
				return
			}
			i++
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

// Unique returns an iterator that generates a sequence of unique elements from s.
// Adjacent duplicate elements are removed. The sequence must be sorted. If not, use
// the Sort function first or the Distinct function directly.
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
// The sequence must be sorted. If not, use the SortFunc first or the DistinctFunc directly.
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

// Distinct returns an iterator that generates a sequence of distinct elements from s.
func Distinct[T comparable](s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// DistinctFunc returns an iterator that generates a sequence of distinct elements from s.
func DistinctFunc[F ~func(T) K, T any, K comparable](s iter.Seq[T], f F) iter.Seq[K] {
	return func(yield func(K) bool) {
		seen := make(map[K]struct{})
		for v := range s {
			k := f(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(k) {
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
