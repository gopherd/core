package ordered

import (
	"bytes"
	"fmt"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/internal/rbtree"
	"github.com/gopherd/core/container/tree"
	"github.com/gopherd/core/operator"
)

// LessFunc represents a comparation function which reports whether x "less" than y
type LessFunc[T any] func(x, y T) bool

// Map represents an ordered map
type Map[K comparable, V any] rbtree.RBTree[K, V]

// New creates an ordered map for ordered K
func NewMap[K constraints.Ordered, V any]() *Map[K, V] {
	return NewMapFunc[K, V](operator.Less[K])
}

// NewMapFunc creates an ordered map with custom less function
func NewMapFunc[K comparable, V any](less LessFunc[K]) *Map[K, V] {
	return (*Map[K, V])(rbtree.NewFunc[K, V](rbtree.LessFunc[K](less)))
}

// Len returns the number of elements
func (m Map[K, V]) Len() int {
	return (rbtree.RBTree[K, V])(m).Len()
}

// Clear clears the set
func (m *Map[K, V]) Clear() {
	(*rbtree.RBTree[K, V])(m).Clear()
}

// Keys collects all keys of the map as a slice
func (m *Map[K, V]) Keys() []K {
	return (*rbtree.RBTree[K, V])(m).Keys()
}

// Values collects all values of the map as a slice
func (m *Map[K, V]) Values() []V {
	return (*rbtree.RBTree[K, V])(m).Values()
}

// Find finds node by key, nil returned if the key not found.
func (m Map[K, V]) Find(key K) *MapIterator[K, V] {
	return (*MapIterator[K, V])((rbtree.RBTree[K, V])(m).Find(key))
}

// Get retrives value by key
func (m Map[K, V]) Get(key K) V {
	var iter = m.Find(key)
	if iter != nil {
		return iter.Value()
	}
	var zero V
	return zero
}

// Contains reports whether the set contains the key
func (m Map[K, V]) Contains(key K) bool {
	return (rbtree.RBTree[K, V])(m).Contains(key)
}

// Insert inserts a key-value pair, inserted node and true returned
// if the key not found, otherwise, existed node and false returned.
func (m *Map[K, V]) Insert(key K, value V) (*MapIterator[K, V], bool) {
	iter, ok := (*rbtree.RBTree[K, V])(m).Insert(key, value)
	return (*MapIterator[K, V])(iter), ok
}

// Remove removes an element by key, false returned if the key not found.
func (m *Map[K, V]) Remove(key K) bool {
	return (*rbtree.RBTree[K, V])(m).Remove(key)
}

// Erase deletes the node, false returned if the node not found.
func (m *Map[K, V]) Erase(iter *MapIterator[K, V]) bool {
	return (*rbtree.RBTree[K, V])(m).Erase((*rbtree.Node[K, V])(iter))
}

// First returns the first node.
//
// usage:
//
//	iter := m.First()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Next()
//	}
func (m Map[K, V]) First() *MapIterator[K, V] {
	return (*MapIterator[K, V])((rbtree.RBTree[K, V])(m).First())
}

// Last returns the last node.
//
// usage:
//
//	iter := m.Last()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Prev()
//	}
func (m Map[K, V]) Last() *MapIterator[K, V] {
	return (*MapIterator[K, V])((rbtree.RBTree[K, V])(m).Last())
}

// Stringify pretty stringifies the map with specified options or nil
func (m Map[K, V]) Stringify(options *tree.Options) string {
	return tree.Stringify[*mapNode[K, V]]((*mapNode[K, V])((rbtree.RBTree[K, V])(m).Root()), options)
}

// String returns content of the set as a plain string
func (m Map[K, V]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	iter := m.First()
	for iter != nil {
		fmt.Fprintf(&buf, "%v:%v", iter.Key(), iter.Value())
		iter = iter.Next()
		if iter != nil {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

// MapIterator represents an iterator of map to iterate nodes
type MapIterator[K comparable, V any] rbtree.Node[K, V]

// Prev returns previous node
func (iter *MapIterator[K, V]) Prev() *MapIterator[K, V] {
	return (*MapIterator[K, V])((*rbtree.Node[K, V])(iter).Prev())
}

// Next returns next node
func (iter *MapIterator[K, V]) Next() *MapIterator[K, V] {
	return (*MapIterator[K, V])((*rbtree.Node[K, V])(iter).Next())
}

// Key returns node's key
func (iter *MapIterator[K, V]) Key() K {
	return (*rbtree.Node[K, V])(iter).Key()
}

// Value returns node's value
func (iter *MapIterator[K, V]) Value() V {
	return (*rbtree.Node[K, V])(iter).Value()
}

// SetValue sets node's value
func (iter *MapIterator[K, V]) SetValue(value V) {
	(*rbtree.Node[K, V])(iter).SetValue(value)
}

// mapNode implements container.Node
type mapNode[K comparable, V any] rbtree.Node[K, V]

// String implements container.Node String method
func (node *mapNode[K, V]) String() string {
	if node == nil {
		return "<nil>"
	}
	var rn = (*rbtree.Node[K, V])(node)
	return fmt.Sprintf("%v:%v", rn.Key(), rn.Value())
}

// Parent implements container.Node Parent method
func (node *mapNode[K, V]) Parent() *mapNode[K, V] {
	return (*mapNode[K, V])((*rbtree.Node[K, V])(node).Parent())
}

// NumChild implements container.Node NumChild method
func (node *mapNode[K, V]) NumChild() int {
	return (*rbtree.Node[K, V])(node).NumChild()
}

// GetChildByIndex implements container.Node GetChildByIndex method
func (node *mapNode[K, V]) GetChildByIndex(i int) *mapNode[K, V] {
	return (*mapNode[K, V])((*rbtree.Node[K, V])(node).GetChildByIndex(i))
}
