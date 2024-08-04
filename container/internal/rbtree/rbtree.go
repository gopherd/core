// Package rbtree implements a Red-Black Tree data structure.
package rbtree

import (
	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/operator"
)

// LessFunc represents a comparison function which reports whether x is "less" than y.
type LessFunc[T any] func(x, y T) bool

// RBTree represents a Red-Black Tree.
type RBTree[K comparable, V any] struct {
	size int
	root *Node[K, V]
	less LessFunc[K]
}

// New creates an RBTree for ordered K.
func New[K constraints.Ordered, V any]() *RBTree[K, V] {
	return NewFunc[K, V](operator.Less[K])
}

// NewFunc creates an RBTree with a custom less function.
func NewFunc[K comparable, V any](less LessFunc[K]) *RBTree[K, V] {
	if less == nil {
		panic("rbtree: less function is nil")
	}
	return &RBTree[K, V]{
		less: less,
	}
}

// Root returns the root node of the tree.
func (tree RBTree[K, V]) Root() *Node[K, V] {
	return tree.root
}

// Len returns the number of elements in the tree.
func (tree RBTree[K, V]) Len() int {
	return tree.size
}

// Clear removes all elements from the tree.
func (tree *RBTree[K, V]) Clear() {
	tree.root = nil
	tree.size = 0
}

// Keys returns a slice containing all keys in the tree.
func (tree *RBTree[K, V]) Keys() []K {
	size := tree.Len()
	if size == 0 {
		return nil
	}
	keys := make([]K, 0, size)
	iter := tree.First()
	for iter != nil {
		keys = append(keys, iter.key)
		iter = iter.Next()
	}
	return keys
}

// Values returns a slice containing all values in the tree.
func (tree *RBTree[K, V]) Values() []V {
	size := tree.Len()
	if size == 0 {
		return nil
	}
	values := make([]V, 0, size)
	iter := tree.First()
	for iter != nil {
		values = append(values, iter.value)
		iter = iter.Next()
	}
	return values
}

// Find returns the node with the given key, or nil if not found.
func (tree RBTree[K, V]) Find(key K) *Node[K, V] {
	return tree.find(key)
}

// Contains reports whether the tree contains the given key.
func (tree RBTree[K, V]) Contains(key K) bool {
	return tree.find(key) != nil
}

// Get retrieves the value associated with the given key.
// If the key is not found, it returns the zero value for V.
func (tree RBTree[K, V]) Get(key K) V {
	node := tree.find(key)
	if node != nil {
		return node.value
	}
	var zero V
	return zero
}

// Insert inserts a key-value pair into the tree.
// It returns the inserted node and true if the key was not found,
// otherwise it returns the existing node and false.
func (tree *RBTree[K, V]) Insert(key K, value V) (*Node[K, V], bool) {
	node, ok := tree.insert(key, value)
	if ok {
		tree.size++
	}
	return node, ok
}

// Remove removes the node with the given key from the tree.
// It returns false if the key was not found.
func (tree *RBTree[K, V]) Remove(key K) bool {
	node := tree.find(key)
	if node == nil || node.null() {
		return false
	}
	tree.remove(node, true)
	tree.size--
	return true
}

// Erase removes the given node from the tree.
// It returns false if the node was not found in the tree.
func (tree *RBTree[K, V]) Erase(node *Node[K, V]) bool {
	if node == nil || node.null() {
		return false
	}
	ok := tree.remove(node, false)
	if ok {
		tree.size--
	}
	return ok
}

// First returns the node with the smallest key in the tree.
func (tree RBTree[K, V]) First() *Node[K, V] {
	if tree.root == nil {
		return nil
	}
	return tree.root.smallest()
}

// Last returns the node with the largest key in the tree.
func (tree RBTree[K, V]) Last() *Node[K, V] {
	if tree.root == nil {
		return nil
	}
	return tree.root.biggest()
}

func (tree *RBTree[K, V]) insert(key K, value V) (*Node[K, V], bool) {
	if tree.root == nil {
		tree.root = &Node[K, V]{
			color: black,
			key:   key,
			value: value,
		}
		tree.root.left = makeNull(tree.root)
		tree.root.right = makeNull(tree.root)
		return tree.root, true
	}

	next := tree.root
	var inserted *Node[K, V]
	for {
		if key == next.key {
			next.value = value
			return next, false
		}
		if tree.less(key, next.key) {
			if next.left.null() {
				inserted = &Node[K, V]{
					parent: next,
					color:  red,
					key:    key,
					value:  value,
				}
				inserted.left = makeNull(inserted)
				inserted.right = makeNull(inserted)
				next.left = inserted
				break
			} else {
				next = next.left
			}
		} else {
			if next.right.null() {
				inserted = &Node[K, V]{
					parent: next,
					color:  red,
					key:    key,
					value:  value,
				}
				inserted.left = makeNull(inserted)
				inserted.right = makeNull(inserted)
				next.right = inserted
				break
			} else {
				next = next.right
			}
		}
	}

	next = inserted
	for {
		next = tree.doInsert(next)
		if next == nil {
			break
		}
	}
	return inserted, true
}

func (tree RBTree[K, V]) find(key K) *Node[K, V] {
	next := tree.root
	for next != nil && !next.null() {
		if next.key == key {
			return next
		}
		if tree.less(key, next.key) {
			next = next.left
		} else {
			next = next.right
		}
	}
	return nil
}

func (tree *RBTree[K, V]) remove(n *Node[K, V], must bool) bool {
	if !must {
		if tree.root == nil || n == nil || n.ancestor() != tree.root {
			return false
		}
	}
	if !n.right.null() {
		smallest := n.right.smallest()
		n.value, smallest.value = smallest.value, n.value
		n.key, smallest.key = smallest.key, n.key
		n = smallest
	}
	child := n.left
	if child.null() {
		child = n.right
	}
	if n.parent == nil {
		if n.left.null() && n.right.null() {
			tree.root = nil
			return true
		}
		child.parent = nil
		tree.root = child
		tree.root.color = black
		return true
	}

	if n.parent.left == n {
		n.parent.left = child
	} else {
		n.parent.right = child
	}
	child.parent = n.parent
	if n.color == red {
		return true
	}
	if child.color == red {
		child.color = black
		return true
	}
	for child != nil {
		child = tree.doRemove(child)
	}
	return true
}

func (tree *RBTree[K, V]) doInsert(n *Node[K, V]) *Node[K, V] {
	if n.parent == nil {
		tree.root = n
		n.color = black
		return nil
	}
	if n.parent.color == black {
		return nil
	}
	uncle := n.uncle()
	if uncle.color == red {
		n.parent.color = black
		uncle.color = black
		gp := n.grandparent()
		gp.color = red
		return gp
	}
	if n.parent.right == n && n.grandparent().left == n.parent {
		tree.rotateLeft(n.parent)
		n.color = black
		n.parent.color = red
		tree.rotateRight(n.parent)
	} else if n.parent.left == n && n.grandparent().right == n.parent {
		tree.rotateRight(n.parent)
		n.color = black
		n.parent.color = red
		tree.rotateLeft(n.parent)
	} else if n.parent.left == n && n.grandparent().left == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		tree.rotateRight(n.grandparent())
	} else if n.parent.right == n && n.grandparent().right == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		tree.rotateLeft(n.grandparent())
	}
	return nil
}

func (tree *RBTree[K, V]) doRemove(n *Node[K, V]) *Node[K, V] {
	if n.parent == nil {
		n.color = black
		return nil
	}
	sibling := n.sibling()
	if sibling.color == red {
		n.parent.color = red
		sibling.color = black
		if n == n.parent.left {
			tree.rotateLeft(n.parent)
		} else {
			tree.rotateRight(n.parent)
		}
	}
	sibling = n.sibling()
	if n.parent.color == black &&
		sibling.color == black &&
		sibling.left.color == black &&
		sibling.right.color == black {
		sibling.color = red
		return n.parent
	}
	if n.parent.color == red &&
		sibling.color == black &&
		sibling.left.color == black &&
		sibling.right.color == black {
		sibling.color = red
		n.parent.color = black
		return nil
	}
	if sibling.color == black {
		if n == n.parent.left &&
			sibling.left.color == red &&
			sibling.right.color == black {
			sibling.color = red
			sibling.left.color = black
			tree.rotateRight(sibling.left.parent)
		} else if n == n.parent.right &&
			sibling.left.color == black &&
			sibling.right.color == red {
			sibling.color = red
			sibling.right.color = black
			tree.rotateLeft(sibling.right.parent)
		}
	}
	sibling = n.sibling()
	sibling.color = n.parent.color
	n.parent.color = black
	if n == n.parent.left {
		sibling.right.color = black
		tree.rotateLeft(sibling.parent)
	} else {
		sibling.left.color = black
		tree.rotateRight(sibling.parent)
	}
	return nil
}

const (
	left  = 0
	right = 1
)

func (tree *RBTree[K, V]) rotate(p *Node[K, V], dir int) *Node[K, V] {
	g := p.parent
	s := p.child(1 - dir)
	c := s.child(dir)
	p.setChild(1-dir, c)
	if !c.null() {
		c.parent = p
	}
	s.setChild(dir, p)
	p.parent = s
	s.parent = g
	if g != nil {
		if p == g.right {
			g.right = s
		} else {
			g.left = s
		}
	} else {
		tree.root = s
	}
	return s
}

func (tree *RBTree[K, V]) rotateLeft(p *Node[K, V]) {
	tree.rotate(p, left)
}

func (tree *RBTree[K, V]) rotateRight(p *Node[K, V]) {
	tree.rotate(p, right)
}

type color byte

const (
	red color = iota
	black
)

// Node represents a node in the Red-Black Tree.
type Node[K comparable, V any] struct {
	parent      *Node[K, V]
	left, right *Node[K, V]
	color       color
	key         K
	value       V
}

func makeNull[K comparable, V any](parent *Node[K, V]) *Node[K, V] {
	return &Node[K, V]{
		parent: parent,
		color:  black,
	}
}

// Key returns the key of the node.
func (node *Node[K, V]) Key() K { return node.key }

// Value returns the value of the node.
func (node *Node[K, V]) Value() V { return node.value }

// SetValue sets the value of the node.
func (node *Node[K, V]) SetValue(value V) { node.value = value }

// Parent returns the parent node.
func (node *Node[K, V]) Parent() *Node[K, V] {
	if node == nil {
		return nil
	}
	return node.parent
}

// NumChild returns the number of non-null children of the node.
func (node *Node[K, V]) NumChild() int {
	if node == nil {
		return 0
	}
	return operator.Bool[int](node.left != nil && !node.left.null()) +
		operator.Bool[int](node.right != nil && !node.right.null())
}

// GetChildByIndex returns the i-th child of the node.
// It panics if i is not 0 or 1.
func (node *Node[K, V]) GetChildByIndex(i int) *Node[K, V] {
	switch i {
	case 0:
		return operator.Ternary(node.left != nil && !node.left.null(), node.left, node.right)
	case 1:
		return node.right
	default:
		panic("rbtree: invalid child index")
	}
}

// Prev returns the predecessor of the node, or nil if there is none.
func (node *Node[K, V]) Prev() *Node[K, V] {
	if node == nil || node.null() {
		return nil
	}
	if !node.left.null() {
		return node.left.biggest()
	}
	parent := node.parent
	for parent != nil && node == parent.left {
		node = parent
		parent = node.parent
	}
	return parent
}

// Next returns the successor of the node, or nil if there is none.
func (node *Node[K, V]) Next() *Node[K, V] {
	if node == nil || node.null() {
		return nil
	}
	if !node.right.null() {
		return node.right.smallest()
	}
	parent := node.parent
	for parent != nil && node == parent.right {
		node = parent
		parent = node.parent
	}
	return parent
}

func (node *Node[K, V]) null() bool {
	return node.left == nil && node.right == nil
}

func (node *Node[K, V]) child(dir int) *Node[K, V] {
	if dir == left {
		return node.left
	}
	return node.right
}

func (node *Node[K, V]) setChild(dir int, child *Node[K, V]) {
	if dir == left {
		node.left = child
	} else {
		node.right = child
	}
}

func (node *Node[K, V]) ancestor() *Node[K, V] {
	ancestor := node
	for ancestor.parent != nil {
		ancestor = ancestor.parent
	}
	return ancestor
}

func (node *Node[K, V]) grandparent() *Node[K, V] {
	if node.parent == nil {
		return nil
	}
	return node.parent.parent
}

func (node *Node[K, V]) sibling() *Node[K, V] {
	if node.parent == nil {
		return nil
	}
	if node.parent.left == node {
		return node.parent.right
	}
	return node.parent.left
}

func (node *Node[K, V]) uncle() *Node[K, V] {
	if node.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *Node[K, V]) smallest() *Node[K, V] {
	next := node
	for !next.left.null() {
		next = next.left
	}
	return next
}

func (node *Node[K, V]) biggest() *Node[K, V] {
	next := node
	for !next.right.null() {
		next = next.right
	}
	return next
}
