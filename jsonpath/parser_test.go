package jsonpath

import (
	"strings"
	"testing"

	"src.userspace.com.au/query/json"
)

func TestParse(t *testing.T) {
	tests := []struct {
		path, src, expect string
	}{
		{
			path:   "$",
			src:    `{"test":"one"}`,
			expect: "one",
		},
		{
			path:   "$.test",
			src:    `{"test":"one"}`,
			expect: "one",
		},
		{
			path:   "$.a.b",
			src:    `{"a":{"b":"two"}}`,
			expect: "two",
		},
		{
			path:   "$.a.b.c",
			src:    `{"a":{"b":{"c":"two"}}}`,
			expect: "two",
		},
		{
			path:   "$.a.b",
			src:    `{"fail":{"a":"one"},"a":{"b":"three"}}`,
			expect: "three",
		},
		{
			path:   "$.test.test",
			src:    `{"fail":{"test1":"one"},"test1":{"test3":"two"}}`,
			expect: "",
		},
		{
			path:   "$..test",
			src:    `{"fail":{"test1":{"test":"two"}}}`,
			expect: "two",
		},
		{
			path:   "$.a.b.c.d",
			src:    `{"a":{"b":{"c":{"d":"blah"}}}}`,
			expect: "blah",
		},
		{
			path:   "$.test[2]",
			src:    `{"test":[1,"two","three"]}`,
			expect: "three",
		},
		{
			path:   "$.*",
			src:    `{"test":"one"}`,
			expect: "one",
		},
	}

	p := Parser{}

	for _, tt := range tests {
		doc, err := json.Parse(strings.NewReader(tt.src))
		if err != nil {
			t.Fatalf("json.Parse(%q) failed: %s", tt.src, err)
		}
		//doc.PrintTree(0)

		sel, err := p.Parse(tt.path)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %s", tt.path, err)
		}
		actualText := ""
		actual := sel.MatchFirst(doc)
		if actual != nil {
			actualText = actual.InnerText()
		}

		if tt.expect == "" && actualText != "" {
			t.Fatalf("MatchFirst(%s, %s) => %s, expected nothing", tt.path, tt.src, actualText)
		} else if actualText != tt.expect {
			t.Fatalf("MatchFirst(%s, %s) => %s, expected %s", tt.path, tt.src, actualText, tt.expect)
		}
	}
}
