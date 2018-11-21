package json

import (
	"bytes"
	"fmt"
	"io"

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
