package html

import (
	"strings"
	"testing"

	base "src.userspace.com.au/query"
)

func TestParse(t *testing.T) {
	src := `<html><body><p>One</p><p>Two</p></body></html>`

	doc, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}
	if doc == nil {
		t.Fatal("Expected node but got nil")
	}

	//doc.PrintTree(0)

	// document

	nt := doc.Type()
	nd := doc.Data()
	if nt != base.DocumentNode {
		t.Fatalf("Expected %q but got %q", "DocumentNode", nt)
	}
	if nd != "" {
		t.Fatalf("Expected %q but got %q", "", nd)
	}

	// get <html>
	n := doc.FirstChild()
	if n == nil {
		t.Fatal("Expected node but got nil")
	}

	nt = n.Type()
	nd = n.Data()
	if nt != base.ElementNode {
		t.Fatalf("Expected %q but got %q", "ElementNode", nd)
	}
	if nd != "html" {
		t.Fatalf("Expected %q but got %q", "html", nd)
	}

	// get <body>
	//n = n.FirstChild()
	// TODO why?
	n = n.LastChild()
	if n == nil {
		t.Fatal("Expected node but got nil")
	}

	nt = n.Type()
	nd = n.Data()
	if nt != base.ElementNode {
		t.Fatalf("Expected %q but got %q", "ElementNode", nd)
	}
	if nd != "body" {
		t.Fatalf("Expected %q but got %q", "body", nd)
	}

	// get first <p>
	n = n.LastChild()
	if n == nil {
		t.Fatal("Expected node but got nil")
	}

	nt = n.Type()
	nd = n.Data()
	if nt != base.ElementNode {
		t.Fatalf("Expected %q but got %q", "ElementNode", nd)
	}
	if nd != "p" {
		t.Fatalf("Expected %q but got %q", "p", nd)
	}
}
