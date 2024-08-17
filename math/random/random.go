package random

// Uint64NGenerator is an interface for generating random uint64 numbers in the range [0, n).
type Uint64NGenerator interface {
	Uint64N(n uint64) uint64
}

// Shuffle randomly shuffles the elements in the slice.
func Shuffle[R Uint64NGenerator, S ~[]T, T any](r R, s S) S {
	return ShuffleN(r, s, len(s))
}

// ShuffleN randomly selects and shuffles the first n elements from the entire slice.
// It ensures that the first n elements are randomly chosen from the whole slice,
// not just shuffled among themselves. Elements after the nth position may also be affected.
// This differs from a complete shuffle as it only guarantees randomness for the first n elements.
// It panics if n is negative or greater than the length of the slice.
func ShuffleN[R Uint64NGenerator, S ~[]T, T any](r R, s S, n int) S {
	if n < 0 || n > len(s) {
		panic("random.ShuffleN: invalid number of elements to shuffle")
	}
	for i := 0; i < n; i++ {
		j := int(r.Uint64N(uint64(i + 1)))
		s[i], s[j] = s[j], s[i]
	}
	return s
}
