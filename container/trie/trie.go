// Package trie implements a prefix tree (trie) data structure.
package trie

import (
	"github.com/gopherd/core/container/tree"
)

// Trie represents a prefix tree data structure.
type Trie struct {
	root           *node
	hasEmptyString bool
}

// New creates and returns a new Trie.
func New() *Trie {
	return &Trie{
		root: newNode(0, nil),
	}
}

// String returns a string representation of the Trie.
func (t *Trie) String() string {
	return t.Stringify(nil)
}

// Stringify formats the Trie as a string using the provided options.
func (t *Trie) Stringify(options *tree.Options) string {
	return tree.Stringify[*node](t.root, options)
}

// search traverses the Trie for the given prefix and returns the last matching node,
// the depth of the match, and whether the prefix was fully matched.
func (t *Trie) search(prefix string) (lastMatchNode *node, depth int, match bool) {
	if prefix == "" {
		return t.root, 0, false
	}

	lastMatchNode = t.root
	depth = 0
	match = true
	current := t.root

	for i, r := range prefix {
		current = current.search(r)
		if current == nil {
			match = false
			return
		}
		lastMatchNode = current
		depth = i + 1
	}
	return
}

// Add inserts a word into the Trie.
func (t *Trie) Add(word string) {
	if word == "" {
		t.hasEmptyString = true
		return
	}
	n, depth, _ := t.search(word)
	for i, r := range word {
		if i >= depth {
			n = n.add(r)
		}
	}
	n.tail = true
}

// Remove deletes a word from the Trie if it exists.
// It returns true if the word was found and removed, false otherwise.
func (t *Trie) Remove(word string) bool {
	if word == "" {
		if t.hasEmptyString {
			t.hasEmptyString = false
			return true
		}
		return false
	}
	n, _, match := t.search(word)
	if match && n.tail {
		n.tail = false
		for n.parent != nil && !n.tail && len(n.children) == 0 {
			n.parent.remove(n.value)
			n = n.parent
		}
		return true
	}
	return false
}

// Has checks if the Trie contains the exact word.
func (t *Trie) Has(word string) bool {
	if word == "" {
		return t.hasEmptyString
	}
	n, _, match := t.search(word)
	return match && n.tail
}

// HasPrefix checks if the Trie contains any word with the given prefix.
func (t *Trie) HasPrefix(prefix string) bool {
	if prefix == "" {
		return true
	}
	_, _, match := t.search(prefix)
	return match
}

// Search retrieves words in the Trie that have the specified prefix.
// It returns up to 'limit' number of words. If limit is 0, it returns all matching words.
func (t *Trie) Search(prefix string, limit int) []string {
	return t.SearchAppend(nil, prefix, limit)
}

// SearchAppend is similar to Search but appends the results to the provided slice.
// It returns the updated slice containing the matching words.
func (t *Trie) SearchAppend(dst []string, prefix string, limit int) []string {
	if limit == 0 {
		limit = -1
	} else if limit > 0 {
		limit = limit - len(dst)
		if limit <= 0 {
			return dst
		}
	}
	n, depth, _ := t.search(prefix)
	if depth != len(prefix) {
		return dst
	}
	var buf []rune
	return n.words(dst, limit, append(buf, []rune(prefix)[:depth]...))
}

// node represents a node in the Trie.
type node struct {
	value    rune
	parent   *node
	children []*node
	tail     bool
}

// newNode creates a new node with the given rune value and parent.
func newNode(r rune, parent *node) *node {
	return &node{
		value:    r,
		parent:   parent,
		children: make([]*node, 0, 2),
	}
}

// String implements the container.Node String method.
func (n *node) String() string {
	if n.parent == nil {
		return "."
	}
	return string(n.value)
}

// Parent implements the container.Node Parent method.
func (n *node) Parent() *node {
	return n.parent
}

// NumChild implements the container.Node NumChild method.
func (n *node) NumChild() int {
	return len(n.children)
}

// GetChildByIndex implements the container.Node GetChildByIndex method.
func (n *node) GetChildByIndex(i int) *node {
	return n.children[i]
}

// indexof returns the index where a child node with the given rune should be inserted.
func (n *node) indexof(r rune) int {
	left, right := 0, len(n.children)
	for left < right {
		mid := int(uint(left+right) >> 1)
		if n.children[mid].value < r {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

// add inserts a new child node with the given rune value and returns it.
func (n *node) add(r rune) *node {
	i := n.indexof(r)
	if i < len(n.children) && n.children[i].value == r {
		return n.children[i]
	}
	child := newNode(r, n)
	n.children = append(n.children, nil)
	copy(n.children[i+1:], n.children[i:])
	n.children[i] = child
	return child
}

// remove deletes the child node with the given rune value.
func (n *node) remove(r rune) {
	i := n.indexof(r)
	if i == len(n.children) || n.children[i].value != r {
		return
	}
	n.children = append(n.children[:i], n.children[i+1:]...)
}

// search finds and returns the child node with the given rune value.
func (n *node) search(r rune) *node {
	i := n.indexof(r)
	if i < len(n.children) && n.children[i].value == r {
		return n.children[i]
	}
	return nil
}

// words recursively collects words from the current node and its children.
func (n *node) words(dst []string, limit int, buf []rune) []string {
	if n.tail {
		dst = append(dst, string(buf))
		if limit > 0 && len(dst) == limit {
			return dst
		}
	}
	for _, child := range n.children {
		dst = child.words(dst, limit, append(buf, child.value))
		if limit > 0 && len(dst) >= limit {
			break
		}
	}
	return dst
}
