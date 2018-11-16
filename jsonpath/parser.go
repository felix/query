package jsonpath

import (
	"fmt"
	"strings"

	//base "src.userspace.com.au/query"
	"src.userspace.com.au/query/json"
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
	fmt.Printf("%s(%d): '%s'\n", tokenNames[p.tok.Type], p.tok.Type, p.tok.Value)
	return p.tok != nil && !done
}

func (p *Parser) Parse(input string) (Selector, error) {
	p.l = lexer.New(input, pathState)
	p.l.Start()

	// First token
	p.next()
	if p.tok.Type != TAbsolute {
		return nil, fmt.Errorf("expected root, got %s", p.tok.Value)
	}

	result, err := p.parseQualifiedSelector()
	if err != nil {
		return nil, err
	}
	return childSelector(rootSelector, result), nil
}

// parseQualifiedSelector
func (p *Parser) parseQualifiedSelector() (result Selector, err error) {
	p.next()

	switch p.tok.Type {
	case TRecursive:
		nr, _ := p.parseStepSelector()
		result = recursiveSelector(nr)

	case TChildDot, TChildStart:
		result, err = p.parseStepSelector()

	default:
		return nil, fmt.Errorf("expected . or .. or something, got %s", p.tok.Value)
	}
	return result, nil
}

func (p *Parser) parseStepSelector() (result Selector, err error) {
	p.next()
	result, err = p.parseNodeTestSelector()
	if err != nil {
		return nil, err
	}
	p.next()
	if p.tok.Type == TPredicateStart {
		// TODO
	}
	return result, nil
}

func (p *Parser) parseNodeTestSelector() (result Selector, err error) {
	switch p.tok.Type {
	case TName:
		/*
			switch p.tok.Value {
			case "object", "array", "string", "number", "boolean", "null":
				// TODO
				//result = typeSelector(p.tok.Value)
			default:
			}
		*/
		result = nameSelector(p.tok.Value)
	case TWildcard:
		result = wildcardSelector
	default:
		fmt.Println("here: ", tokenNames[p.tok.Type])
	}
	return result, err
}

func (p *Parser) parseChildSelector() Selector {
	var result Selector
	p.next()
	switch p.tok.Type {
	case TQuotedName:
		result = nameSelector(strings.Trim(p.tok.Value, `"'`))
	case TName:
		result = nameSelector(p.tok.Value)
	}
	p.next()
	return result
}

// rootSelector checks node is root
func rootSelector(n *json.Node) bool {
	result := (n.Type == json.DocumentNode)
	fmt.Printf("rootSelector => type: %s, val: %s, result: %t\n", json.NodeNames[n.Type], n.Data, result)
	return result
}

// wildcardSelector returns true
func wildcardSelector(n *json.Node) bool {
	return true
}

// childSelector creates a selector for c being a child of p
func childSelector(p, c Selector) Selector {
	return func(n *json.Node) bool {
		fmt.Printf("childSelector => type: %s, val: %s\n", json.NodeNames[n.Type], n.Data)
		result := (c(n) && n.Parent != nil && p(n.Parent))
		fmt.Printf("childSelector => type: %s, val: %s, result: %t\n", json.NodeNames[n.Type], n.Data, result)
		return result
	}
}

// nameSelector generates selector for object key == k
func nameSelector(k string) Selector {
	return func(n *json.Node) bool {
		result := (n.Type == json.ElementNode && n.Data == k)
		fmt.Printf("nameSelector => type: %s, val: %s, result: %t\n", json.NodeNames[n.Type], n.Data, result)
		return result
	}
}

// recursiveSelector matches any node below which matches a
func recursiveSelector(a Selector) Selector {
	return func(n *json.Node) bool {
		if n.Type != json.ElementNode {
			return false
		}
		return hasRecursiveMatch(n, a)
	}
}

func hasRecursiveMatch(n *json.Node, a Selector) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if a(c) || (c.Type == json.ElementNode && hasRecursiveMatch(c, a)) {
			return true
		}
	}
	return false
}

// typeSelector matches a node with type t
func typeSelector(t string) Selector {
	return func(n *json.Node) bool {
		if n.DataType == t {
			return true
		}
		return false
	}
}
