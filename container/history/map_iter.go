//go:build go1.23

package history

import (
	"iter"
	"maps"
)

// All returns an iterator over key-value pairs in the map
func (m Map[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.data)
}

// Keys returns an iterator over the map keys
func (m Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.data)
}

// Values returns an iterator over the map values
func (m Map[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.data)
}
