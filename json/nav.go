package json

import (
	"fmt"

	base "src.userspace.com.au/query"
)

// NodeNavigator implements the Nav interface for navigating JSON nodes.
type NodeNavigator struct {
	root, cur *Node
}

func (a *NodeNavigator) Current() *Node {
	return a.cur
}

func (a *NodeNavigator) NodeType() base.NodeType {
	switch a.cur.Type {
	case TextNode:
		return base.TextNode
	case DocumentNode:
		return base.DocumentNode
	case ElementNode:
		return base.ElementNode
	default:
		panic(fmt.Sprintf("unknown node type %v", a.cur.Type))
	}
}

func (a *NodeNavigator) LocalName() string {
	return a.cur.Data

}

func (a *NodeNavigator) Prefix() string {
	return ""
}

func (a *NodeNavigator) Value() string {
	switch a.cur.Type {
	case ElementNode:
		return a.cur.InnerText()
	case TextNode:
		return a.cur.Data
	}
	return ""
}

func (a *NodeNavigator) Copy() base.Nav {
	n := *a
	return &n
}

func (a *NodeNavigator) MoveToRoot() {
	a.cur = a.root
}

func (a *NodeNavigator) MoveToParent() bool {
	if n := a.cur.Parent; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (x *NodeNavigator) MoveToNextAttribute() bool {
	return false
}

func (a *NodeNavigator) MoveToChild() bool {
	if n := a.cur.FirstChild; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (a *NodeNavigator) MoveToFirst() bool {
	for n := a.cur.PrevSibling; n != nil; n = n.PrevSibling {
		a.cur = n
	}
	return true
}

func (a *NodeNavigator) String() string {
	return a.Value()
}

func (a *NodeNavigator) MoveToNext() bool {
	if n := a.cur.NextSibling; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (a *NodeNavigator) MoveToPrevious() bool {
	if n := a.cur.PrevSibling; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (a *NodeNavigator) MoveTo(other base.Nav) bool {
	node, ok := other.(*NodeNavigator)
	if !ok || node.root != a.root {
		return false
	}
	a.cur = node.cur
	return true
}
