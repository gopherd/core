// Package tree provides functionality for working with tree-like data structures.
package tree

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gopherd/core/container/pair"
	"github.com/gopherd/core/operator"
)

// Node represents a generic printable node in a tree structure.
type Node[T comparable] interface {
	// String returns the node's information as a string.
	String() string

	// Parent returns the parent node or zero value if it's the root.
	Parent() T

	// NumChild returns the number of child nodes.
	NumChild() int

	// GetChildByIndex returns the child node at the given index.
	// It returns the zero value of T if the index is out of range.
	GetChildByIndex(i int) T
}

// NodeMarshaler extends Node with marshaling capability.
type NodeMarshaler[T comparable] interface {
	Node[T]

	// Marshal serializes the node into a byte slice.
	Marshal() ([]byte, error)
}

// encodeNode writes a node's data to the provided writer.
func encodeNode[T comparable](w io.Writer, id, parent uint32, node NodeMarshaler[T]) error {
	content, err := node.Marshal()
	if err != nil {
		return err
	}

	// Write header: {size: uint32, id: uint32, parent: uint32}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(content))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, id); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, parent); err != nil {
		return err
	}

	// Write content
	_, err = w.Write(content)
	return err
}

// Marshal serializes a tree structure into a byte slice.
func Marshal[T comparable](root NodeMarshaler[T]) ([]byte, error) {
	if root == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	var id uint32
	openSet := list.New()
	openSet.PushBack(pair.New(uint32(0), root))

	for openSet.Len() > 0 {
		front := openSet.Front()
		p := front.Value.(pair.Pair[uint32, NodeMarshaler[T]])
		openSet.Remove(front)

		id++
		if err := encodeNode(&buf, id, p.First, p.Second); err != nil {
			return nil, err
		}

		for i, n := 0, p.Second.NumChild(); i < n; i++ {
			child, ok := any(p.Second.GetChildByIndex(i)).(NodeMarshaler[T])
			if ok {
				openSet.PushBack(pair.New(id, child))
			}
		}
	}

	return buf.Bytes(), nil
}

// Options represents the configuration for stringifying a Node.
type Options struct {
	Prefix     string
	Parent     string // Default "│  "
	Space      string // Default "   "
	Branch     string // Default "├──"
	LastBranch string // Default "└──"
}

var defaultOptions = &Options{
	Parent:     "│   ",
	Space:      "    ",
	Branch:     "├── ",
	LastBranch: "└── ",
}

// Fix ensures all options have valid values, using defaults where necessary.
func (options *Options) Fix() {
	options.Parent = operator.Or(options.Parent, defaultOptions.Parent)
	if options.Branch == "" {
		options.Branch = defaultOptions.Branch
		options.LastBranch = operator.Or(options.LastBranch, defaultOptions.LastBranch)
	} else if options.LastBranch == "" {
		options.LastBranch = options.Branch
	}
	options.Space = operator.Or(options.Space, defaultOptions.Space)
}

// Stringify converts a node to a string representation.
func Stringify[T comparable](node Node[T], options *Options) string {
	if options == nil {
		options = defaultOptions
	} else {
		options.Fix()
	}

	if stringer, ok := node.(interface {
		Stringify(*Options) string
	}); ok {
		return stringer.Stringify(options)
	}

	var buf, stack bytes.Buffer
	if options.Prefix != "" {
		stack.WriteString(options.Prefix)
	}
	recursivelyPrintNode[T](node, &buf, &stack, "", 0, false, options)
	return buf.String()
}

// recursivelyPrintNode prints a node and its children recursively.
func recursivelyPrintNode[T comparable](
	x interface{},
	w io.Writer,
	stack *bytes.Buffer,
	prefix string,
	depth int,
	isLast bool,
	options *Options,
) {
	node, ok := x.(Node[T])
	if !ok {
		return
	}

	nprefix := stack.Len()
	value := node.String()
	parent := node.Parent()
	fmt.Fprintf(w, "%s%s%s\n", stack.String(), prefix, value)

	var zero T
	if parent != zero {
		if isLast {
			stack.WriteString(options.Space)
		} else {
			stack.WriteString(options.Parent)
		}
	}

	n := node.NumChild()
	for i := 0; i < n; i++ {
		isLast := i+1 == n
		appended := operator.Ternary(isLast, options.LastBranch, options.Branch)
		child := node.GetChildByIndex(i)
		if child == zero {
			continue
		}
		recursivelyPrintNode[T](child, w, stack, appended, depth+1, isLast, options)
	}

	if nprefix != stack.Len() {
		stack.Truncate(nprefix)
	}
}
