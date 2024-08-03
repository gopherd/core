package trie

import (
	"github.com/gopherd/core/container/tree"
)

// Trie implements a prefix-tree
type Trie struct {
	root *node
}

// New creates a trie
func New() *Trie {
	return &Trie{
		root: newNode(0, nil),
	}
}

// String returns trie as a string
func (trie *Trie) String() string {
	return trie.Stringify(nil)
}

// String formats trie as a string
func (trie *Trie) Stringify(options *tree.Options) string {
	return tree.Stringify[*node](trie.root, options)
}

func (trie *Trie) search(prefix string) (lastMatchNode *node, deep int, match bool) {
	if prefix == "" {
		return trie.root, 0, false
	}

	lastMatchNode = trie.root
	deep = 0
	match = true
	next := trie.root

	for i, r := range prefix {
		next = next.search(r)
		if next == nil {
			match = false
			return
		}
		lastMatchNode = next
		deep = i + 1
	}
	return
}

// Add adds word to trie
func (trie *Trie) Add(word string) {
	n, deep, _ := trie.search(word)
	for i, r := range word {
		if i >= deep {
			n = n.add(r)
		}
	}
	n.tail = true
}

// Remove removes word from trie
func (trie *Trie) Remove(word string) bool {
	n, _, match := trie.search(word)
	if match && n.tail {
		n.tail = false
		for !n.tail && len(n.children) == 0 {
			n.parent.remove(n.value)
			n = n.parent
		}
		return true
	}
	return false
}

// Has reports whether the trie contains the word
func (trie *Trie) Has(word string) bool {
	n, _, match := trie.search(word)
	return match && n.tail
}

// HasPrefix reports the trie has prefix
func (trie *Trie) HasPrefix(prefix string) bool {
	_, _, match := trie.search(prefix)
	return match
}

// Search retrives words which has specified prefix
func (trie *Trie) Search(prefix string, limit int) []string {
	return trie.SearchAppend(nil, prefix, limit)
}

// SearchAppend likes Search but append words to dst
func (trie *Trie) SearchAppend(dst []string, prefix string, limit int) []string {
	if limit == 0 {
		return dst
	}
	n, deep, _ := trie.search(prefix)
	if deep != len(prefix) {
		return nil
	}
	var buf []rune
	if limit > 0 {
		limit += len(dst)
	}
	return n.words(dst, limit, buf)
}

// node of trie
type node struct {
	value    rune
	parent   *node
	children []*node
	tail     bool
}

func newNode(r rune, parent *node) *node {
	n := new(node)
	n.value = r
	n.parent = parent
	n.children = make([]*node, 0, 2)
	return n
}

// String implements container.Node String method
func (n *node) String() string {
	if n.parent == nil {
		return "."
	}
	return string(n.value)
}

// Parent implements container.Node Parent method
func (n *node) Parent() *node {
	return n.parent
}

// NumChild implements container.Node NumChild method
func (n *node) NumChild() int {
	return len(n.children)
}

// GetChildByIndex implements container.Node GetChildByIndex method
func (n *node) GetChildByIndex(i int) *node {
	return n.children[i]
}

func (n *node) indexof(r rune) int {
	i, j := 0, len(n.children)
	for i < j {
		h := int(uint(i+j) >> 1)
		if n.children[h].value < r {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (n *node) add(r rune) *node {
	i := n.indexof(r)
	if i < len(n.children) && n.children[i].value == r {
		return n.children[i]
	}
	child := newNode(r, n)
	child.parent = n
	size := len(n.children)
	n.children = append(n.children, child)
	if size > 0 {
		copy(n.children[i+1:], n.children[i:size])
		n.children[i] = child
	}
	return child
}

func (n *node) remove(r rune) {
	i := n.indexof(r)
	if i == len(n.children) || n.children[i].value != r {
		return
	}
	n.children = append(n.children[:i], n.children[i+1:]...)
}

func (n *node) search(r rune) *node {
	i := n.indexof(r)
	if i < len(n.children) && n.children[i].value == r {
		return n.children[i]
	}
	return nil
}

func (n *node) words(dst []string, limit int, buf []rune) []string {
	if n.tail {
		buf = buf[:0]
		next := n
		for next.parent != nil {
			buf = append(buf, next.value)
			next = next.parent
		}
		n := len(buf)
		for i, m := 0, n/2; i < m; i++ {
			buf[i], buf[n-i-1] = buf[n-i-1], buf[i]
		}
		dst = append(dst, string(buf))
	}
	for _, child := range n.children {
		if limit > 0 && len(dst) >= limit {
			break
		}
		dst = child.words(dst, limit, buf)
	}
	return dst
}
