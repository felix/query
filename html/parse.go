package html

import (
	//"fmt"
	"io"

	x "golang.org/x/net/html"
)

func Parse(r io.Reader) (*Node, error) {
	xnode, err := x.Parse(r)
	if err != nil {
		return nil, err
	}
	/*
		if len(xnodes) > 1 {
			return nil, fmt.Errorf("found multiple HTML roots: %d", len(xnodes))
		}
	*/

	root := wrapNodes(xnode, 0)
	return root, nil
}

func wrapNodes(root *x.Node, l int) *Node {
	out := &Node{Node: root, level: l}

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		child := wrapNodes(c, l+1)
		out.appendChild(child)
	}
	/*
		if root.Parent != nil {
			out.parent = &Node{Node: root.Parent}
			if l > 0 {
				out.parent.level = l - 1
			}
		}

			if root.FirstChild != nil {
				out.firstChild = wrapNodes(root.FirstChild, l+1)
			}

			if root.NextSibling != nil {
				out.nextSibling = wrapNodes(root.NextSibling, l)
			}

			if root.LastChild != nil {
				out.lastChild = wrapNodes(root.LastChild, l+1)
			}
			if root.PrevSibling != nil {
				//out.prevSibling = wrapNodes(root.prevSibling, l)
				out.prevSibling = &Node{
					Node:   root.PrevSibling,
					level:  l,
					parent: out.parent,
				}
			}
	*/
	return out
}
