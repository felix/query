package jsonpath

import (
	"fmt"
	"strconv"
	"strings"

	base "src.userspace.com.au/query"
	"src.userspace.com.au/query/lexer"
)

type Parser struct {
	l   *lexer.Lexer
	pos int
	tok *lexer.Token
}

func NewParser() (*Parser, error) {
	return &Parser{}, nil
}

func (p *Parser) next() (done bool) {
	p.tok, done = p.l.NextToken()
	if p.tok != nil {
		p.pos = p.tok.Position
	}
	return p.tok != nil && !done
}

func (p *Parser) Parse(input string) (Selector, error) {
	var sel, nr Selector
	var err error

	p.l = lexer.New(input, pathState)
	p.l.Start()

	// First token
	p.next()

	if p.tok.Type != TAbsolute {
		// TODO does jsonpath have relative searches
		return nil, fmt.Errorf("expected root, got %s", p.tok.Value)
	}

	p.next()

	if p.tok.Type == lexer.EOFToken {
		return rootSelector, nil
	}

	sel = rootSelector

	for {
		switch p.tok.Type {
		case TRecursive:
			p.next()
			p.next()
			if nr, err = p.parseStepSelector(); err != nil {
				return nil, err
			}
			sel = recursiveSelector(nr)

		case TChildDot:
			p.next()
			if nr, err = p.parseStepSelector(); err != nil {
				return nil, err
			}
			sel = childSelector(sel, nr)
		default:
			return sel, nil
		}
	}
	panic("unreachable")
}

func (p *Parser) parseStepSelector() (Selector, error) {
	var sel, nr Selector
	var err error

	sel = p.parseNodeTestSelector()
	p.next()

	switch p.tok.Type {
	case TPredicateStart:
		p.next()
		if nr, err = p.parsePredicateExprSelector(); err != nil {
			return nil, err
		}
		sel = childSelector(sel, nr)

	case lexer.EOFToken:
		return sel, nil

	default:
	}
	return sel, nil
}

func (p *Parser) parseNodeTestSelector() (sel Selector) {
	switch p.tok.Type {
	case TName:
		/*
			switch p.tok.Value {
			case "object", "array", "string", "number", "boolean", "null":
				// TODO
				//sel = typeSelector(p.tok.Value)
			default:
			}
		*/
		sel = nameSelector(p.tok.Value)

	case TQuotedName:
		sel = nameSelector(strings.Trim(p.tok.Value, `"'`))

	case TWildcard:
		sel = wildcardSelector

	default:
	}
	return sel
}

func (p *Parser) parsePredicateExprSelector() (Selector, error) {
	var err error

	if p.tok.Type != TNumber {
		return nil, fmt.Errorf("expecting number")
	}

	num, err := strconv.ParseInt(p.tok.Value, 10, 64)
	if err != nil {
		return nil, err
	}
	return arrayIndexSelector(num), nil
	/* TODO
	var els []int64
	if p.tok.Type == TPredicateEnd {
		return arrayIndexSelector(num), nil
	}

	els = append(els, num)

	p.next()

	if p.tok.Type == TRange {
		// FIXME
		p.next()
		num, err := strconv.ParseInt(p.tok.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		els = append(els, num)
		// We have start and finish range
	}
	if p.tok.Type == TUnion {
		// FIXME
	}
	}
	return childSelector(rootSelector, arrayIndexSelector(idx)), nil
	return sel, nil
	*/
}

// rootSelector checks node is root
func rootSelector(n base.Node) bool {
	result := false
	parent := n.Parent()
	if parent != nil && parent.Type() == base.DocumentNode {
		result = true
	} else {
		result = (n.Type() == base.DocumentNode)
	}
	return result
}

// wildcardSelector returns true
func wildcardSelector(n base.Node) bool {
	return true
}

// childSelector creates a selector for c being a child of p
func childSelector(p, c Selector) Selector {
	return func(n base.Node) bool {
		parent := n.Parent()
		result := (c(n) && parent != nil && p(parent))
		return result
	}
}

// nameSelector generates selector for object key == k
func nameSelector(k string) Selector {
	return func(n base.Node) bool {
		result := (n.Type() == base.ElementNode && n.Data() == k)
		return result
	}
}

// recursiveSelector matches any node below which matches a
func recursiveSelector(a Selector) Selector {
	return func(n base.Node) bool {
		if n.Type() != base.ElementNode {
			return false
		}
		return hasRecursiveMatch(n, a)
	}
}

func hasRecursiveMatch(n base.Node, a Selector) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if a(c) || (c.Type() == base.ElementNode && hasRecursiveMatch(c, a)) {
			return true
		}
	}
	return false
}

// arrayIndexSelector generates selector for node being idx index of parent
func arrayIndexSelector(idx int64) Selector {
	return func(n base.Node) bool {
		if n.DataType() != "arrayitem" {
			return false
		}
		parent := n.Parent()

		if parent == nil {
			return false
		}

		i := int64(0)
		for c := parent.FirstChild(); c != nil && i <= idx; c = c.NextSibling() {
			if i == idx && c == n {
				return true
			}
			i++
		}
		return false
	}
}

// typeSelector matches a node with type t
func typeSelector(t string) Selector {
	return func(n base.Node) bool {
		// FIXME
		if n.DataType() == t {
			return true
		}
		return false
	}
}
