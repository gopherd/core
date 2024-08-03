package pair

// Pair contains two values like c++ std::pair
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// Make creates a pair by values
func Make[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}
