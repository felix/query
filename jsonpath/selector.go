package jsonpath

import (
	"src.userspace.com.au/query/json"
)

type Selector func(*json.Node) bool

// MatchAll returns a slice of the nodes that match the selector,
// from n and its children.
func (s Selector) MatchAll(n *json.Node) []*json.Node {
	return s.matchAllInto(n, nil)
}

func (s Selector) matchAllInto(n *json.Node, storage []*json.Node) []*json.Node {
	if s(n) {
		storage = append(storage, n)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		storage = s.matchAllInto(child, storage)
	}

	return storage
}

// Match returns true if the node matches the selector.
func (s Selector) Match(n *json.Node) bool {
	return s(n)
}

// MatchFirst returns the first node that matches s, from n and its children.
func (s Selector) MatchFirst(n *json.Node) *json.Node {
	if s.Match(n) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		m := s.MatchFirst(c)
		if m != nil {
			return m
		}
	}
	return nil
}

// Filter returns the nodes in nodes that match the selector.
func (s Selector) Filter(nodes []*json.Node) (result []*json.Node) {
	for _, n := range nodes {
		if s(n) {
			result = append(result, n)
		}
	}
	return result
}
