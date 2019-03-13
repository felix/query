package json

import (
	"encoding/json"
	"io"
	"sort"
	"strconv"

	base "src.userspace.com.au/felix/query"
)

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
