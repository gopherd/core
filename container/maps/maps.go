package maps

import (
	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/pair"
)

// Keys retrieves keys of map
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	var s = make([]K, len(m))
	if len(s) == 0 {
		return s
	}
	for k := range m {
		s = append(s, k)
	}
	return s
}

// Values retrieves values of map
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	var s = make([]V, len(m))
	if len(s) == 0 {
		return s
	}
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// Map creates a slice and inserts values of m by function f
func Map[
	D ~[]T,
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T any,
](m M, f F) D {
	d := make(D, 0, len(m))
	for k, v := range m {
		d = append(d, f(k, v))
	}
	return d
}

// MinKey retrieves mininum key of map
func MinKey[M ~map[K]V, K constraints.Ordered, V any](m M) K {
	var min K
	if m == nil {
		return min
	}
	var n int
	for k := range m {
		if n == 0 || k < min {
			min = k
		}
		n++
	}
	return min
}

// MaxKey retrieves maxinum key of map
func MaxKey[M ~map[K]V, K constraints.Ordered, V any](m M) K {
	var max K
	if m == nil {
		return max
	}
	var n int
	for k := range m {
		if n == 0 || k > max {
			max = k
		}
		n++
	}
	return max
}

// MinmaxKey retrieves mininum and maxinum key of map
func MinmaxKey[M ~map[K]V, K constraints.Ordered, V any](m M) (min, max K) {
	if m == nil {
		return
	}
	var n int
	for k := range m {
		if n == 0 || k < min {
			min = k
		}
		if n == 0 || k > max {
			max = k
		}
		n++
	}
	return
}

// MinValue retrieves mininum value of map
func MinValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) pair.Pair[K, V] {
	var key K
	var min V
	if m == nil {
		return pair.Make(key, min)
	}
	var n int
	for k, v := range m {
		if n == 0 || v < min {
			key = k
			min = v
		}
		n++
	}
	return pair.Make(key, min)
}

// MaxValue retrieves mininum value of map
func MaxValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) pair.Pair[K, V] {
	var key K
	var max V
	if m == nil {
		return pair.Make(key, max)
	}
	var n int
	for k, v := range m {
		if n == 0 || v > max {
			key = k
			max = v
		}
		n++
	}
	return pair.Make(key, max)
}

// MinmaxValue retrieves mininum and maxinum value of map
func MinmaxValue[M ~map[K]V, K comparable, V constraints.Ordered](m M) (min, max pair.Pair[K, V]) {
	if m == nil {
		return
	}
	var n int
	for k, v := range m {
		if n == 0 || v < min.Second {
			min.First = k
			min.Second = v
		}
		if n == 0 || v > max.Second {
			max.First = k
			max.Second = v
		}
		n++
	}
	return
}

// MinKeyFunc retrieves mininum key of map
func MinKeyFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) T {
	var min T
	if m == nil {
		return min
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t < min {
			min = t
		}
		n++
	}
	return min
}

// MaxKeyFunc retrieves maxinum key of map
func MaxKeyFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) T {
	var max T
	if m == nil {
		return max
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t > max {
			max = t
		}
		n++
	}
	return max
}

// MinmaxKeyFunc retrieves mininum and maxinum key of map
func MinmaxKeyFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) (min, max T) {
	if m == nil {
		return
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t < min {
			min = t
		}
		if n == 0 || t > max {
			max = t
		}
		n++
	}
	return
}

// MinValueFunc retrieves mininum value of map
func MinValueFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) pair.Pair[K, T] {
	var key K
	var min T
	if m == nil {
		return pair.Make(key, min)
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t < min {
			key = k
			min = t
		}
		n++
	}
	return pair.Make(key, min)
}

// MaxValueFunc retrieves mininum value of map
func MaxValueFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) pair.Pair[K, T] {
	var key K
	var max T
	if m == nil {
		return pair.Make(key, max)
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t > max {
			key = k
			max = t
		}
		n++
	}
	return pair.Make(key, max)
}

// MinmaxValueFunc retrieves mininum and maxinum value of map
func MinmaxValueFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Ordered,
](m M, f F) (min, max pair.Pair[K, T]) {
	if m == nil {
		return
	}
	var n int
	for k, v := range m {
		t := f(k, v)
		if n == 0 || t < min.Second {
			min.First = k
			min.Second = t
		}
		if n == 0 || t > max.Second {
			max.First = k
			max.Second = t
		}
		n++
	}
	return
}

// CopyFunc inserts pairs mapping from key-value pair of m by function f
func CopyFunc[
	D ~map[T]U,
	M ~map[K]V,
	F ~func(K, V) (T, U),
	K comparable,
	V any,
	T comparable,
	U any,
](d D, m M, f F) D {
	for k, v := range m {
		t, u := f(k, v)
		d[t] = u
	}
	return d
}

// SumKey sums keys of map
func SumKey[
	M ~map[K]V,
	K constraints.Number | ~string,
	V any,
](m M) K {
	var sum K
	if m == nil {
		return sum
	}
	for k := range m {
		sum += k
	}
	return sum
}

// SumValue sums values of map
func SumValue[
	M ~map[K]V,
	K comparable,
	V constraints.Number | ~string,
](m M) V {
	var sum V
	if m == nil {
		return sum
	}
	for _, v := range m {
		sum += v
	}
	return sum
}

// SumFunc sums mapped values by function f
func SumFunc[
	M ~map[K]V,
	F ~func(K, V) T,
	K comparable,
	V any,
	T constraints.Number | ~string,
](m M, f F) T {
	var sum T
	if m == nil {
		return sum
	}
	for k, v := range m {
		sum += f(k, v)
	}
	return sum
}

// Clone returns a copy of m
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}
	var d = make(M, len(m))
	for k, v := range m {
		d[k] = v
	}
	return d
}

// Copy copies all key/value pairs in src adding them to dst.
func Copy[D, S ~map[K]V, K comparable, V any](dst D, src S) {
	for k, v := range src {
		dst[k] = v
	}
}

// Equal reports whether two maps contain the same key/value pairs.
func Equal[D, S ~map[K]V, K, V comparable](dst D, src S) bool {
	if len(dst) != len(src) {
		return false
	}
	if len(dst) == 0 {
		return true
	}
	for k, v := range src {
		if x, ok := dst[k]; !ok || x != v {
			return false
		}
	}
	return true
}

// EqualFunc is like Equal, but compares values using
func EqualFunc[
	D ~map[K]V,
	S ~map[K]U,
	F ~func(V, U) bool,
	K comparable,
	V any,
	U any,
](dst D, src S, f F) bool {
	if len(dst) != len(src) {
		return false
	}
	if len(dst) == 0 {
		return true
	}
	for k, u := range src {
		if v, ok := dst[k]; !ok || !f(v, u) {
			return false
		}
	}
	return true
}
