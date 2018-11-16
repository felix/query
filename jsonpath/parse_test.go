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
			path:   "$.test",
			src:    `{"test":"one"}`,
			expect: "one",
		},
	}

	p := Parser{}

	for _, tt := range tests {
		doc, err := json.Parse(strings.NewReader(tt.src))
		if err != nil {
			t.Errorf("json.Parse(%q) failed: %s", tt.src, err)
		}

		sel, err := p.Parse(tt.path)
		if err != nil {
			t.Errorf("Parse(%q) failed: %s", tt.path, err)
		}
		actual := sel.MatchFirst(doc)
		actualText := actual.InnerText()
		if actualText != tt.expect {
			t.Errorf("MatchFirst(%s) => %s, expected %s", tt.src, actualText, tt.expect)
		}
	}
}
