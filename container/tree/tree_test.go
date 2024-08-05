package tree_test

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/gopherd/core/container/tree"
)

// MockNode implements the Node and NodeMarshaler interfaces for testing
type MockNode struct {
	id       int
	parent   *MockNode
	children []*MockNode
	data     string
}

func (n *MockNode) String() string    { return n.data }
func (n *MockNode) Parent() *MockNode { return n.parent }
func (n *MockNode) NumChild() int     { return len(n.children) }
func (n *MockNode) GetChildByIndex(i int) *MockNode {
	if i < 0 || i >= len(n.children) {
		return nil
	}
	return n.children[i]
}
func (n *MockNode) Marshal() ([]byte, error) { return []byte(n.data), nil }

func TestStringify(t *testing.T) {
	root := &MockNode{id: 1, data: "Root"}
	child1 := &MockNode{id: 2, parent: root, data: "Child 1"}
	child2 := &MockNode{id: 3, parent: root, data: "Child 2"}
	grandchild := &MockNode{id: 4, parent: child1, data: "Grandchild"}

	root.children = []*MockNode{child1, child2}
	child1.children = []*MockNode{grandchild}

	tests := []struct {
		name     string
		node     tree.Node[*MockNode]
		options  *tree.Options
		expected string
	}{
		{
			name:     "Default options",
			node:     root,
			options:  nil,
			expected: "Root\n├── Child 1\n│   └── Grandchild\n└── Child 2\n",
		},
		{
			name: "Custom options",
			node: root,
			options: &tree.Options{
				Prefix:     "  ",
				Parent:     "| ",
				Space:      "  ",
				Branch:     "+- ",
				LastBranch: "\\- ",
			},
			expected: "  Root\n  +- Child 1\n  | \\- Grandchild\n  \\- Child 2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tree.Stringify(tt.node, tt.options)
			if result != tt.expected {
				t.Errorf("Stringify() =\n%v\nwant\n%v", result, tt.expected)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	root := &MockNode{id: 1, data: "Root"}
	child := &MockNode{id: 2, parent: root, data: "Child"}
	root.children = []*MockNode{child}

	data, err := tree.Marshal[*MockNode](root)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify the structure of the marshaled data
	reader := bytes.NewReader(data)

	// Check root node
	var size, id, parentID uint32
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		t.Fatalf("Failed to read size: %v", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &id); err != nil {
		t.Fatalf("Failed to read id: %v", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &parentID); err != nil {
		t.Fatalf("Failed to read parentID: %v", err)
	}

	content := make([]byte, size)
	if _, err := reader.Read(content); err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	if id != 1 || parentID != 0 || string(content) != "Root" {
		t.Errorf("Unexpected root node data: id=%d, parentID=%d, content=%s", id, parentID, content)
	}

	// Check child node
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		t.Fatalf("Failed to read size: %v", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &id); err != nil {
		t.Fatalf("Failed to read id: %v", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &parentID); err != nil {
		t.Fatalf("Failed to read parentID: %v", err)
	}

	content = make([]byte, size)
	if _, err := reader.Read(content); err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	if id != 2 || parentID != 1 || string(content) != "Child" {
		t.Errorf("Unexpected child node data: id=%d, parentID=%d, content=%s", id, parentID, content)
	}

	// Ensure we've read all the data
	if reader.Len() != 0 {
		t.Errorf("Unexpected extra data in marshaled output")
	}
}

func TestMarshalNilRoot(t *testing.T) {
	data, err := tree.Marshal[*MockNode](nil)
	if err != nil {
		t.Fatalf("Marshal(nil) error = %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Marshal(nil) returned non-empty data: %v", data)
	}
}

func TestOptionsFix(t *testing.T) {
	tests := []struct {
		name     string
		input    tree.Options
		expected tree.Options
	}{
		{
			name:     "Empty options",
			input:    tree.Options{},
			expected: tree.Options{Parent: "│   ", Space: "    ", Branch: "├── ", LastBranch: "└── "},
		},
		{
			name:     "Custom Branch without LastBranch",
			input:    tree.Options{Branch: ">> "},
			expected: tree.Options{Parent: "│   ", Space: "    ", Branch: ">> ", LastBranch: ">> "},
		},
		{
			name:     "Custom everything",
			input:    tree.Options{Parent: "P ", Space: "S ", Branch: "B ", LastBranch: "L "},
			expected: tree.Options{Parent: "P ", Space: "S ", Branch: "B ", LastBranch: "L "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := tt.input
			options.Fix()
			if !reflect.DeepEqual(options, tt.expected) {
				t.Errorf("fix() = %v, want %v", options, tt.expected)
			}
		})
	}
}
