package jsonpath

import (
	"fmt"
)

type jsonpath struct {
	absolute bool
	steps    []step
}

func (jp jsonpath) String() string {
	out := ""
	if jp.absolute {
		out = "$"
	}
	for _, s := range jp.steps {
		out += s.String()
	}
	return out
}

type step struct {
	recursive bool
	selector  selector
	predicate *predicate
}

func (s step) String() string {
	if s.recursive {
		return ".."
	}
	if s.predicate != nil {
		return fmt.Sprintf("%s%s", s.selector.String(), s.predicate.String())
	}
	return s.selector.String()
}

type selector struct {
	wildcard bool
	value    string
}

func (s selector) String() string {
	if s.wildcard {
		return "[*]"
	}
	return fmt.Sprintf("[%s]", s.value)
}

type predicate struct {
	pType  string
	start  int
	end    int
	filter string
}

func (p predicate) String() string {
	switch p.pType {
	case "index":
		return fmt.Sprintf("[%d]", p.start)
	case "range":
		return fmt.Sprintf("[%d:%d]", p.start, p.end)
	case "union":
		// TODO
	}
	return ""
}
