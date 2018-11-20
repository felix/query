package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
)

// A NodeType is the type of a Node.
type NodeType uint

const (
	// DocumentNode the root of the document tree.
	DocumentNode NodeType = iota
	// ElementNode is an element.
	ElementNode
	// TextNode is the text content of a node.
	TextNode
)

var NodeNames = map[NodeType]string{
	DocumentNode: "DocumentNode",
	ElementNode:  "ElementNode",
	TextNode:     "TextNode",
}

// A Node consists of a NodeType and some Data (tag name for
// element nodes, content for text) and are part of a tree of Nodes.
type Node struct {
	Parent, PrevSibling, NextSibling, FirstChild, LastChild *Node

	Type     NodeType
	Data     string
	DataType string

	level int
}

func (n Node) String() string {
	return fmt.Sprintf("[%s] %s(%s)", NodeNames[n.Type], n.DataType, n.Data)
}

func (n Node) PrintTree(level int) {
	for i := 1; i <= level; i++ {
		fmt.Printf("  ")
	}
	fmt.Println(n)
	for _, c := range n.ChildNodes() {
		c.PrintTree(level + 1)
	}
}

// ChildNodes gets all child nodes of the node.
func (n *Node) ChildNodes() []*Node {
	var a []*Node
	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		a = append(a, nn)
	}
	return a
}

// InnerText gets the value of the node and all its child nodes.
func (n *Node) InnerText() string {
	var output func(*bytes.Buffer, *Node)
	output = func(buf *bytes.Buffer, n *Node) {
		if n.Type == TextNode {
			buf.WriteString(n.Data)
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}
	var buf bytes.Buffer
	output(&buf, n)
	return buf.String()
}

// SelectElement finds the first of child elements with the
// specified name.
func (n *Node) SelectElement(name string) *Node {
	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		if nn.Data == name {
			return nn
		}
	}
	return nil
}

func parseValue(x interface{}, top *Node, level int) {
	addNode := func(n *Node) {
		if n.level == top.level {
			top.NextSibling = n
			n.PrevSibling = top
			n.Parent = top.Parent
			if top.Parent != nil {
				top.Parent.LastChild = n
			}
		} else if n.level > top.level {
			n.Parent = top
			if top.FirstChild == nil {
				top.FirstChild = n
				top.LastChild = n
			} else {
				t := top.LastChild
				t.NextSibling = n
				n.PrevSibling = t
				top.LastChild = n
			}
		}
	}

	// TODO check for null

	switch v := x.(type) {
	case []interface{}:
		for _, vv := range v {
			n := &Node{Type: ElementNode, level: level, DataType: "arrayitem"}
			addNode(n)
			parseValue(vv, n, level+1)
		}
	case map[string]interface{}:
		// The Go’s map iteration order is random.
		// (https://blog.golang.org/go-maps-in-action#Iteration-order)
		var keys []string
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			n := &Node{Data: key, Type: ElementNode, level: level, DataType: "object"}
			addNode(n)
			parseValue(v[key], n, level+1)
		}
	case string:
		n := &Node{Data: v, Type: TextNode, level: level, DataType: "string"}
		addNode(n)
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		n := &Node{Data: s, Type: TextNode, level: level, DataType: "number"}
		addNode(n)
	case bool:
		s := strconv.FormatBool(v)
		n := &Node{Data: s, Type: TextNode, level: level, DataType: "boolean"}
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
	doc := &Node{Type: DocumentNode}
	parseValue(v, doc, 1)
	return doc, nil
}
