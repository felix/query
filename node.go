package query

import (
	"golang.org/x/net/html"
)

type NodeType uint32

const (
	ErrorNode     = NodeType(html.ErrorNode)
	TextNode      = NodeType(html.TextNode)
	DocumentNode  = NodeType(html.DocumentNode)
	ElementNode   = NodeType(html.ElementNode)
	CommentNode   = NodeType(html.CommentNode)
	DoctypeNode   = NodeType(html.DoctypeNode)
	AttributeNode = 100
	AnyNode       = 101
	// Add json node types
)

type Node interface {
	Parent() Node
	FirstChild() Node
	LastChild() Node
	PrevSibling() Node
	NextSibling() Node
	Type() NodeType
	Data() string
	Attr() []Attribute
}

type Attribute struct {
	Namespace, Key, Val string
}
