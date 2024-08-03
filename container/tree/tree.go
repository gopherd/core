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

// Node represents a generic printable node
type Node[T comparable] interface {
	String() string          // String returns node self information
	Parent() T               // Parent returns parent node or nil
	NumChild() int           // NumChild returns number of child
	GetChildByIndex(i int) T // GetChildByIndex gets child by index
}

type NodeMarshaler[T comparable] interface {
	Node[T]
	Marshal() ([]byte, error)
}

func encodeNode[T comparable](w io.Writer, id, parent uint32, node NodeMarshaler[T]) error {
	content, err := node.Marshal()
	if err != nil {
		return err
	}
	// header: {size: uint32, id: uint32, parent: uint32}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(content))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(id)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(parent)); err != nil {
		return err
	}
	// content
	_, err = w.Write(content)
	return err
}

func Marshal[T comparable](root NodeMarshaler[T]) ([]byte, error) {
	var buf bytes.Buffer
	if root == nil {
		return nil, nil
	}
	var id uint32
	var openSet = list.New()
	openSet.PushBack(pair.Make(uint32(0), root))
	for openSet.Len() > 0 {
		var front = openSet.Front()
		var p = front.Value.(pair.Pair[uint32, NodeMarshaler[T]])
		openSet.Remove(front)
		id++
		if err := encodeNode(&buf, id, p.First, p.Second); err != nil {
			return nil, err
		}
		for i, n := 0, p.Second.NumChild(); i < n; i++ {
			var child = p.Second.GetChildByIndex(i)
			openSet.PushBack(pair.Make(id, child))
		}
	}
	return buf.Bytes(), nil
}

// Options represents a options for stringify Node
type Options struct {
	Prefix     string
	Parent     string // default "│  "
	Space      string // default "   "
	Branch     string // default "├──"
	LastBranch string // default "└──"
}

var defaultOptions = &Options{
	Parent:     "│   ",
	Space:      "    ",
	Branch:     "├── ",
	LastBranch: "└── ",
}

func (options *Options) fix() {
	options.Parent = operator.Or(options.Parent, defaultOptions.Parent)
	if options.Branch == "" {
		options.Branch = defaultOptions.Branch
		options.LastBranch = operator.Or(options.LastBranch, defaultOptions.LastBranch)
	} else if options.LastBranch == "" {
		options.LastBranch = options.Branch
	}
	options.Space = operator.Or(options.Space, defaultOptions.Space)
}

// Stringify converts node to string
func Stringify[T comparable](node Node[T], options *Options) string {
	if options == nil {
		options = defaultOptions
	} else {
		options.fix()
	}
	if stringer, ok := node.(interface {
		Stringify(*Options) string
	}); ok {
		return stringer.Stringify(options)
	}
	var (
		buf   bytes.Buffer
		stack bytes.Buffer
	)
	if options.Prefix != "" {
		stack.WriteString(options.Prefix)
	}
	recursivelyPrintNode[T](node, &buf, &stack, "", 0, false, options)
	return buf.String()
}

func recursivelyPrintNode[T comparable](
	x interface{},
	w io.Writer,
	stack *bytes.Buffer,
	prefix string,
	depth int,
	isLast bool,
	options *Options,
) {
	var node, ok = x.(Node[T])
	if !ok {
		return
	}
	var nprefix = stack.Len()
	var value = node.String()
	var parent = node.Parent()
	fmt.Fprintf(w, "%s%s%s\n", stack.String(), prefix, value)
	var zero T

	if parent != zero {
		if isLast {
			stack.WriteString(options.Space)
		} else {
			stack.WriteString(options.Parent)
		}
	}
	var n = node.NumChild()
	for i := 0; i < n; i++ {
		isLast = i+1 == n
		var appended = operator.Ternary(isLast, options.LastBranch, options.Branch)
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
