// Package maputil provides utility functions for working with maps.
package maputil

import (
	"cmp"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/pair"
)

// Keys retrieves all keys from the given map.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values retrieves all values from the given map.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Map creates a slice by applying function f to each key-value pair in the map.
func Map[S []T, M ~map[K]V, F ~func(K, V) T, K comparable, V any, T any](m M, f F) S {
	result := make(S, 0, len(m))
	for k, v := range m {
		result = append(result, f(k, v))
	}
	return result
}

// MinKey returns the minimum key in the map. If the map is empty, it returns the zero value of K.
func MinKey[M ~map[K]V, K constraints.Ordered, V any](m M) K {
	var min K
	first := true
	for k := range m {
		if first || cmp.Less(k, min) {
			min = k
			first = false
		}
	}
	return min
}

// MaxKey returns the maximum key in the map. If the map is empty, it returns the zero value of K.
func MaxKey[M ~map[K]V, K constraints.Ordered, V any](m M) K {
	var max K
	first := true
	for k := range m {
		if first || cmp.Less(max, k) {
			max = k
			first = false
		}
	}
	return max
}

// MinMaxKey returns the minimum and maximum keys in the map.
// If the map is empty, it returns zero values for both min and max.
func MinMaxKey[M ~map[K]V, K constraints.Ordered, V any](m M) (min, max K) {
	first := true
	for k := range m {
		if first {
			min, max = k, k
			first = false
		} else if cmp.Less(k, min) {
			min = k
		} else if cmp.Less(max, k) {
			max = k
		}
	}
	return
}

// MinValue returns the key-value pair with the minimum value in the map.
// If the map is empty, it returns zero values for both key and value.
func MinValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) pair.Pair[K, V] {
	var result pair.Pair[K, V]
	first := true
	for k, v := range m {
		if first || cmp.Less(v, result.Second) {
			result = pair.New(k, v)
			first = false
		}
	}
	return result
}

// MaxValue returns the key-value pair with the maximum value in the map.
// If the map is empty, it returns zero values for both key and value.
func MaxValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) pair.Pair[K, V] {
	var result pair.Pair[K, V]
	first := true
	for k, v := range m {
		if first || cmp.Less(result.Second, v) {
			result = pair.New(k, v)
			first = false
		}
	}
	return result
}

// MinMaxValue returns the key-value pairs with the minimum and maximum values in the map.
// If the map is empty, it returns zero values for both pairs.
func MinMaxValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) (min, max pair.Pair[K, V]) {
	first := true
	for k, v := range m {
		if first {
			min = pair.New(k, v)
			max = pair.New(k, v)
			first = false
		} else {
			if cmp.Less(v, min.Second) {
				min = pair.New(k, v)
			}
			if cmp.Less(max.Second, v) {
				max = pair.New(k, v)
			}
		}
	}
	return
}

// MinKeyFunc returns the minimum value obtained by applying function f to each key-value pair.
// If the map is empty, it returns the zero value of T.
func MinKeyFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) T {
	var min T
	first := true
	for k, v := range m {
		t := f(k, v)
		if first || cmp.Less(t, min) {
			min = t
			first = false
		}
	}
	return min
}

// MaxKeyFunc returns the maximum value obtained by applying function f to each key-value pair.
// If the map is empty, it returns the zero value of T.
func MaxKeyFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) T {
	var max T
	first := true
	for k, v := range m {
		t := f(k, v)
		if first || cmp.Less(max, t) {
			max = t
			first = false
		}
	}
	return max
}

// MinMaxKeyFunc returns the minimum and maximum values obtained by applying function f to each key-value pair.
// If the map is empty, it returns zero values for both min and max.
func MinMaxKeyFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) (min, max T) {
	first := true
	for k, v := range m {
		t := f(k, v)
		if first {
			min, max = t, t
			first = false
		} else {
			if cmp.Less(t, min) {
				min = t
			}
			if cmp.Less(max, t) {
				max = t
			}
		}
	}
	return
}

// MinValueFunc returns the key and minimum value obtained by applying function f to each key-value pair.
// If the map is empty, it returns zero values for both key and result.
func MinValueFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) pair.Pair[K, T] {
	var result pair.Pair[K, T]
	first := true
	for k, v := range m {
		t := f(k, v)
		if first || cmp.Less(t, result.Second) {
			result = pair.New(k, t)
			first = false
		}
	}
	return result
}

// MaxValueFunc returns the key and maximum value obtained by applying function f to each key-value pair.
// If the map is empty, it returns zero values for both key and result.
func MaxValueFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) pair.Pair[K, T] {
	var result pair.Pair[K, T]
	first := true
	for k, v := range m {
		t := f(k, v)
		if first || cmp.Less(result.Second, t) {
			result = pair.New(k, t)
			first = false
		}
	}
	return result
}

// MinMaxValueFunc returns the keys and values with the minimum and maximum results obtained by applying function f to each key-value pair.
// If the map is empty, it returns zero values for both pairs.
func MinMaxValueFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Ordered](m M, f F) (min, max pair.Pair[K, T]) {
	first := true
	for k, v := range m {
		t := f(k, v)
		if first {
			min = pair.New(k, t)
			max = pair.New(k, t)
			first = false
		} else {
			if cmp.Less(t, min.Second) {
				min = pair.New(k, t)
			}
			if cmp.Less(max.Second, t) {
				max = pair.New(k, t)
			}
		}
	}
	return
}

// CopyFunc creates a new map by applying function f to each key-value pair in the source map.
func CopyFunc[D ~map[T]U, M ~map[K]V, F ~func(K, V) (T, U), K comparable, V any, T comparable, U any](d D, m M, f F) D {
	for k, v := range m {
		t, u := f(k, v)
		d[t] = u
	}
	return d
}

// SumKey returns the sum of all keys in the map.
func SumKey[M ~map[K]V, K constraints.Number | ~string, V any](m M) K {
	var sum K
	for k := range m {
		sum += k
	}
	return sum
}

// SumValue returns the sum of all values in the map.
func SumValue[M ~map[K]V, K comparable, V constraints.Number | ~string](m M) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

// SumFunc returns the sum of values obtained by applying function f to each key-value pair in the map.
func SumFunc[M ~map[K]V, F ~func(K, V) T, K comparable, V any, T constraints.Number | ~string](m M, f F) T {
	var sum T
	for k, v := range m {
		sum += f(k, v)
	}
	return sum
}
