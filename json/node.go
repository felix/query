package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"

	base "src.userspace.com.au/query"
)

// A Node consists of a NodeType and some data (tag name for
// element nodes, content for text) and are part of a tree of Nodes.
type Node struct {
	parent, prevSibling, nextSibling, firstChild, lastChild *Node

	nodeType base.NodeType
	data     string
	dataType string

	level int
}

func (n *Node) Parent() base.Node {
	if n.parent == nil {
		return nil
	}
	return n.parent
}
func (n *Node) NextSibling() base.Node {
	if n.nextSibling == nil {
		return nil
	}
	return n.nextSibling
}
func (n *Node) FirstChild() base.Node {
	if n.firstChild == nil {
		return nil
	}
	return n.firstChild
}
func (n *Node) PrevSibling() base.Node { return n.prevSibling }
func (n *Node) LastChild() base.Node   { return n.lastChild }
func (n *Node) Type() base.NodeType    { return base.NodeType(n.nodeType) }
func (n *Node) DataType() string       { return n.dataType }
func (n *Node) Attr() []base.Attribute { return nil }

func (n Node) String() string {
	return fmt.Sprintf("[%s] %s(%s)", base.NodeNames[n.nodeType], n.dataType, n.data)
}

func (n Node) PrintTree(level int) {
	for i := 1; i <= level; i++ {
		fmt.Printf("  ")
	}
	fmt.Println(n)
	for c := n.firstChild; c != nil; c = c.nextSibling {
		c.PrintTree(level + 1)
	}
}

// ChildNodes gets all child nodes of the node.
func (n *Node) ChildNodes() []*Node {
	var a []*Node
	for nn := n.firstChild; nn != nil; nn = nn.nextSibling {
		a = append(a, nn)
	}
	return a
}

// Data gets the value of the node and all its child nodes.
func (n *Node) Data() string {
	return n.data
}

// InnerText gets the value of the node and all its child nodes.
func (n *Node) InnerText() string {
	var buf bytes.Buffer
	if n.nodeType == base.TextNode {
		buf.WriteString(n.data)
	} else {
		for child := n.firstChild; child != nil; child = child.nextSibling {
			buf.WriteString(child.InnerText())
		}
	}
	return buf.String()
}

// SelectElement finds the first of child elements with the
// specified name.
func (n *Node) SelectElement(name string) *Node {
	for nn := n.firstChild; nn != nil; nn = nn.nextSibling {
		if nn.data == name {
			return nn
		}
	}
	return nil
}

func parseValue(x interface{}, top *Node, level int) {
	addNode := func(n *Node) {
		if n.level == top.level {
			top.nextSibling = n
			n.prevSibling = top
			n.parent = top.parent
			if top.parent != nil {
				top.parent.lastChild = n
			}
		} else if n.level > top.level {
			n.parent = top
			if top.firstChild == nil {
				top.firstChild = n
				top.lastChild = n
			} else {
				t := top.lastChild
				t.nextSibling = n
				n.prevSibling = t
				top.lastChild = n
			}
		}
	}

	// TODO check for null

	switch v := x.(type) {
	case []interface{}:
		for _, vv := range v {
			n := &Node{nodeType: base.ElementNode, level: level, dataType: "arrayitem"}
			addNode(n)
			parseValue(vv, n, level+1)
		}
	case map[string]interface{}:
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			n := &Node{data: key, nodeType: base.ElementNode, level: level, dataType: "object"}
			addNode(n)
			parseValue(v[key], n, level+1)
		}
	case string:
		n := &Node{data: v, nodeType: base.TextNode, level: level, dataType: "string"}
		addNode(n)
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		n := &Node{data: s, nodeType: base.TextNode, level: level, dataType: "number"}
		addNode(n)
	case bool:
		s := strconv.FormatBool(v)
		n := &Node{data: s, nodeType: base.TextNode, level: level, dataType: "boolean"}
		addNode(n)
	}
}

// Parse JSON document.
func Parse(r io.Reader) (*Node, error) {
	var v interface{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	doc := &Node{nodeType: base.DocumentNode}
	parseValue(v, doc, 1)
	return doc, nil
}
