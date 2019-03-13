package json

import (
	base "src.userspace.com.au/felix/query"
)

// NodeNavigator implements the Nav interface for navigating JSON nodes.
type NodeNavigator struct {
	root, cur *Node
}

func (a *NodeNavigator) Current() *Node {
	return a.cur
}

func (a *NodeNavigator) NodeType() base.NodeType {
	return a.cur.nodeType
}

func (a *NodeNavigator) LocalName() string {
	return a.cur.data

}

func (a *NodeNavigator) Prefix() string {
	return ""
}

func (a *NodeNavigator) Value() string {
	switch a.cur.nodeType {
	case base.ElementNode:
		return a.cur.data
	case base.TextNode:
		return a.cur.data
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
	if n := a.cur.parent; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (x *NodeNavigator) MoveToNextAttribute() bool {
	return false
}

func (a *NodeNavigator) MoveToChild() bool {
	if n := a.cur.firstChild; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (a *NodeNavigator) MoveToFirst() bool {
	for n := a.cur.prevSibling; n != nil; n = n.prevSibling {
		a.cur = n
	}
	return true
}

func (a *NodeNavigator) String() string {
	return a.Value()
}

func (a *NodeNavigator) MoveToNext() bool {
	if n := a.cur.nextSibling; n != nil {
		a.cur = n
		return true
	}
	return false
}

func (a *NodeNavigator) MoveToPrevious() bool {
	if n := a.cur.prevSibling; n != nil {
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
