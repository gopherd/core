// Package maps provides utility functions for working with maps.
package maps

import (
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
func Map[D []T, M ~map[K]V, F ~func(K, V) T, K comparable, V any, T any](m M, f F) D {
	result := make(D, 0, len(m))
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
		if first || k < min {
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
		if first || k > max {
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
		} else if k < min {
			min = k
		} else if k > max {
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
		if first || v < result.Second {
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
		if first || v > result.Second {
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
			if v < min.Second {
				min = pair.New(k, v)
			}
			if v > max.Second {
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
		if first || t < min {
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
		if first || t > max {
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
			if t < min {
				min = t
			}
			if t > max {
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
		if first || t < result.Second {
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
		if first || t > result.Second {
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
			if t < min.Second {
				min = pair.New(k, t)
			}
			if t > max.Second {
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

// Clone returns a new map with the same key-value pairs as the input map.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}
	result := make(M, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Copy copies all key-value pairs from src to dst.
func Copy[D, S ~map[K]V, K comparable, V any](dst D, src S) {
	for k, v := range src {
		dst[k] = v
	}
}

// Equal reports whether two maps contain the same key-value pairs.
func Equal[D, S ~map[K]V, K, V comparable](dst D, src S) bool {
	if len(dst) != len(src) {
		return false
	}
	for k, v := range src {
		if dv, ok := dst[k]; !ok || dv != v {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two maps contain the same key-value pairs,
// using the provided function f to compare values.
func EqualFunc[D ~map[K]V, S ~map[K]U, F ~func(V, U) bool, K comparable, V any, U any](dst D, src S, f F) bool {
	if len(dst) != len(src) {
		return false
	}
	for k, u := range src {
		v, ok := dst[k]
		if !ok || !f(v, u) {
			return false
		}
	}
	return true
}
