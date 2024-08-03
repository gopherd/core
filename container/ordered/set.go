package ordered

import (
	"bytes"
	"fmt"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/internal/rbtree"
	"github.com/gopherd/core/container/tree"
	"github.com/gopherd/core/operator"
)

type empty struct{}

var emptyValue = empty{}

// Set represents an ordered set
type Set[K comparable] rbtree.RBTree[K, empty]

// New creates an ordered set for ordered K
func NewSet[K constraints.Ordered]() *Set[K] {
	return NewSetFunc[K](operator.Less[K])
}

// NewSetFunc creates an ordered set with custom less function
func NewSetFunc[K comparable](less LessFunc[K]) *Set[K] {
	return (*Set[K])(rbtree.NewFunc[K, empty](rbtree.LessFunc[K](less)))
}

// Len returns the number of elements
func (s Set[K]) Len() int {
	return (rbtree.RBTree[K, empty])(s).Len()
}

// Clear clears the set
func (s *Set[K]) Clear() {
	(*rbtree.RBTree[K, empty])(s).Clear()
}

// Keys collects all keys of the set as a slice
func (s *Set[K]) Keys() []K {
	return (*rbtree.RBTree[K, empty])(s).Keys()
}

// Find finds node by key, nil returned if the key not found.
func (s Set[K]) Find(key K) *SetIterator[K] {
	return (*SetIterator[K])((rbtree.RBTree[K, empty])(s).Find(key))
}

// Contains reports whether the set contains the key
func (s Set[K]) Contains(key K) bool {
	return (rbtree.RBTree[K, empty])(s).Contains(key)
}

// Insert inserts a key-value pair, inserted node and true returned
// if the key not found, otherwise, existed node and false returned.
func (s *Set[K]) Insert(key K) (*SetIterator[K], bool) {
	iter, ok := (*rbtree.RBTree[K, empty])(s).Insert(key, emptyValue)
	return (*SetIterator[K])(iter), ok
}

// Remove removes an element by key, false returned if the key not found.
func (s *Set[K]) Remove(key K) bool {
	return (*rbtree.RBTree[K, empty])(s).Remove(key)
}

// Erase deletes the node, false returned if the node not found.
func (s *Set[K]) Erase(iter *SetIterator[K]) bool {
	return (*rbtree.RBTree[K, empty])(s).Erase((*rbtree.Node[K, empty])(iter))
}

// First returns the first node.
//
// usage:
//
//	iter := s.First()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key()
//		iter = iter.Next()
//	}
func (s Set[K]) First() *SetIterator[K] {
	return (*SetIterator[K])((rbtree.RBTree[K, empty])(s).First())
}

// Last returns the last node.
//
// usage:
//
//	iter := s.Last()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key()
//		iter = iter.Prev()
//	}
func (s Set[K]) Last() *SetIterator[K] {
	return (*SetIterator[K])((rbtree.RBTree[K, empty])(s).Last())
}

// Stringify pretty stringifies the set with specified options or nil
func (s Set[K]) Stringify(options *tree.Options) string {
	return tree.Stringify[*setNode[K]]((*setNode[K])((rbtree.RBTree[K, empty])(s).Root()), options)
}

// String returns content of the set as a plain string
func (s Set[K]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	iter := s.First()
	for iter != nil {
		fmt.Fprintf(&buf, "%v", iter.Key())
		iter = iter.Next()
		if iter != nil {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

// SetIterator represents an iterator of OrderedSet to iterate nodes
type SetIterator[K comparable] rbtree.Node[K, empty]

// Prev returns previous node
func (iter *SetIterator[K]) Prev() *SetIterator[K] {
	return (*SetIterator[K])((*rbtree.Node[K, empty])(iter).Prev())
}

// Next returns next node
func (iter *SetIterator[K]) Next() *SetIterator[K] {
	return (*SetIterator[K])((*rbtree.Node[K, empty])(iter).Next())
}

// Key returns node's key
func (iter *SetIterator[K]) Key() K {
	return (*rbtree.Node[K, empty])(iter).Key()
}

// setNode implements container.Node
type setNode[K comparable] rbtree.Node[K, empty]

// String implements container.Node String method
func (node *setNode[K]) String() string {
	if node == nil {
		return "<nil>"
	}
	return fmt.Sprint((*rbtree.Node[K, empty])(node).Key())
}

// Parent implements container.Node Parent method
func (node *setNode[K]) Parent() *setNode[K] {
	return (*setNode[K])((*rbtree.Node[K, empty])(node).Parent())
}

// NumChild implements container.Node NumChild method
func (node *setNode[K]) NumChild() int {
	return (*rbtree.Node[K, empty])(node).NumChild()
}

// GetChildByIndex implements container.Node GetChildByIndex method
func (node *setNode[K]) GetChildByIndex(i int) *setNode[K] {
	return (*setNode[K])((*rbtree.Node[K, empty])(node).GetChildByIndex(i))
}
