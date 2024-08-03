package encoding

import (
	"strconv"
	"text/scanner"
)

// NodeType represents type of node
type NodeType int

const (
	Invalid NodeType = iota
	Ident            // abc,true,false
	Int              // 1
	Float            // 1.2
	Char             // 'c'
	String           // "xyz"
	Object           // {}
	Array            // []
	KVPair           // e.g. x=y, x: y
)

func (t NodeType) String() string {
	if t >= 0 && t < NodeType(len(nodeTypes)) {
		return nodeTypes[t]
	}
	return "Unknown(" + strconv.Itoa(int(t)) + ")"
}

var nodeTypes = [...]string{
	Invalid: "Invalid",
	Ident:   "Ident",
	Int:     "Int",
	Float:   "Float",
	Char:    "Char",
	String:  "String",
	KVPair:  "KVPair",
	Object:  "Object",
	Array:   "Array",
}

type Node interface {
	// Pos returns position of node
	Pos() scanner.Position
	// Type returns type of node
	Type() NodeType
}

type Nodebase struct {
	pos scanner.Position
}

func CreateNodebase(pos scanner.Position) Nodebase { return Nodebase{pos: pos} }

func (n Nodebase) Pos() scanner.Position { return n.pos }

type LiteralNode struct {
	Nodebase
	nodeType NodeType
	Value    string
}

func NewLiteralNode(pos scanner.Position, nodeType NodeType) *LiteralNode {
	return &LiteralNode{
		Nodebase: CreateNodebase(pos),
		nodeType: nodeType,
	}
}

func (n LiteralNode) Type() NodeType { return n.nodeType }

type IdentNode struct {
	Nodebase
	Name string
}

func (n IdentNode) Type() NodeType { return Ident }

func NewIdentNode(pos scanner.Position, name string) *IdentNode {
	return &IdentNode{
		Nodebase: CreateNodebase(pos),
	}
}

type KVPairNode struct {
	K, V Node
}

func NewKVPairNode(k, v Node) *KVPairNode {
	return &KVPairNode{k, v}
}

func (n KVPairNode) Pos() scanner.Position { return n.K.Pos() }
func (n KVPairNode) Type() NodeType        { return KVPair }

type ObjectNode struct {
	Nodebase
	Elements []KVPairNode
}

func (n ObjectNode) Type() NodeType { return Object }

type ArrayNode struct {
	Nodebase
	Elements []Node
}

func (n ArrayNode) Type() NodeType { return Array }
