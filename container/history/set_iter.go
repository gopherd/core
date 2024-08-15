//go:build go1.23

package history

import (
	"iter"
	"maps"
)

// All returns an iterator over keys in the set
func (s Set[K]) All() iter.Seq[K] {
	return maps.Keys(s.data)
}
