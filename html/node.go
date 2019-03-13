package html

import (
	"fmt"

	x "golang.org/x/net/html"
	base "src.userspace.com.au/felix/query"
)

// A Node consists of a NodeType and some data (tag name for
// element nodes, content for text) and are part of a tree of Nodes.
type Node struct {
	parent, prevSibling, nextSibling, firstChild, lastChild *Node

	level int

	*x.Node
}

func (n *Node) Parent() base.Node {
	if n.parent == nil {
		return nil
	}
	return &Node{Node: n.Node.Parent}
}
func (n *Node) NextSibling() base.Node {
	if n.nextSibling == nil {
		return nil
	}
	return &Node{Node: n.Node.NextSibling}
}
func (n *Node) FirstChild() base.Node {
	if n.firstChild == nil {
		return nil
	}
	return &Node{Node: n.Node.FirstChild}
}
func (n *Node) PrevSibling() base.Node { return &Node{Node: n.Node.PrevSibling} }
func (n *Node) LastChild() base.Node   { return &Node{Node: n.Node.LastChild} }
func (n *Node) Type() base.NodeType    { return base.NodeType(n.Node.Type) }
func (n *Node) DataType() string       { return "string" }
func (n *Node) Attr() []base.Attribute {
	out := make([]base.Attribute, len(n.Node.Attr))
	for _, a := range n.Node.Attr {
		out = append(out, base.Attribute(a))
	}
	return out
}

// Data gets the value of the node and all its child nodes.
func (n *Node) Data() string { return n.Node.Data }

// InnerText gets the value of the node and all its child nodes.
func (n *Node) InnerText() string {
	// FIXME
	return n.Node.Data
}

func (n Node) String() string {
	return fmt.Sprintf("[%s] %s(%s)", base.NodeNames[n.Type()], n.DataType(), n.Data())
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

func (n *Node) appendChild(c *Node) {
	if c.parent != nil || c.prevSibling != nil || c.nextSibling != nil {
		panic("html: appendChild called for an attached child Node")
	}
	last := n.lastChild
	if last != nil {
		last.nextSibling = c
	} else {
		n.firstChild = c
	}
	n.lastChild = c
	c.parent = n
	c.prevSibling = last
}
