package jsonpath

import (
	base "src.userspace.com.au/felix/query"
)

type Selector func(base.Node) bool

// MatchAll returns a slice of the nodes that match the selector,
// from n and its children.
func (s Selector) MatchAll(n base.Node) []base.Node {
	return s.matchAllInto(n, nil)
}

func (s Selector) matchAllInto(n base.Node, storage []base.Node) []base.Node {
	if s(n) {
		storage = append(storage, n)
	}

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		storage = s.matchAllInto(child, storage)
	}

	return storage
}

// Match returns true if the node matches the selector.
func (s Selector) Match(n base.Node) bool {
	return s(n)
}

// MatchFirst returns the first node that matches s, from n and its children.
func (s Selector) MatchFirst(n base.Node) base.Node {
	if s.Match(n) {
		return n
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		m := s.MatchFirst(c)
		if m != nil {
			return m
		}
	}
	return nil
}

// Filter returns the nodes in nodes that match the selector.
func (s Selector) Filter(nodes []base.Node) (result []base.Node) {
	for _, n := range nodes {
		if s(n) {
			result = append(result, n)
		}
	}
	return result
}
