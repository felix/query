package jsonpath

import (
	"testing"

	"src.userspace.com.au/query/lexer"
)

func TestValidStates(t *testing.T) {
	tests := []struct {
		path   string
		tokens []lexer.Token
	}{
		{
			path: "$.test",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "test"},
			},
		},
		{
			path: "$[test]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TName, Value: "test"},
				lexer.Token{Type: TChildEnd, Value: "]"},
			},
		},
		{
			path: "$[one][two]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TChildEnd, Value: "]"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TName, Value: "two"},
				lexer.Token{Type: TChildEnd, Value: "]"},
			},
		},
		{
			path: "$.one.two",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TChildEnd, Value: ""},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "two"},
			},
		},
		{
			path: "$[*]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TWildcard, Value: "*"},
				lexer.Token{Type: TChildEnd, Value: "]"},
			},
		},
		{
			path: "$[one][*]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TChildEnd, Value: "]"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TWildcard, Value: "*"},
				lexer.Token{Type: TChildEnd, Value: "]"},
			},
		},
		{
			path: "$.one[1,2,3]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TPredicateStart, Value: "["},
				lexer.Token{Type: TNumber, Value: "1"},
				lexer.Token{Type: TUnion, Value: ","},
				lexer.Token{Type: TNumber, Value: "2"},
				lexer.Token{Type: TUnion, Value: ","},
				lexer.Token{Type: TNumber, Value: "3"},
				lexer.Token{Type: TPredicateEnd, Value: "]"},
			},
		},
		{
			path: "$.one[1:3]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TPredicateStart, Value: "["},
				lexer.Token{Type: TNumber, Value: "1"},
				lexer.Token{Type: TRange, Value: ":"},
				lexer.Token{Type: TNumber, Value: "3"},
				lexer.Token{Type: TPredicateEnd, Value: "]"},
			},
		},
		{
			path: "$.one[:3]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TPredicateStart, Value: "["},
				lexer.Token{Type: TNumber, Value: ""},
				lexer.Token{Type: TRange, Value: ":"},
				lexer.Token{Type: TNumber, Value: "3"},
				lexer.Token{Type: TPredicateEnd, Value: "]"},
			},
		},
		{
			path: "$.one[3:]",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
				lexer.Token{Type: TPredicateStart, Value: "["},
				lexer.Token{Type: TNumber, Value: "3"},
				lexer.Token{Type: TRange, Value: ":"},
				lexer.Token{Type: TNumber, Value: ""},
				lexer.Token{Type: TPredicateEnd, Value: "]"},
			},
		},
		{
			path: "$..one",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TRecursive, Value: "."},
				lexer.Token{Type: TChildDot, Value: "."},
				lexer.Token{Type: TName, Value: "one"},
			},
		},
		{
			path: "$['one']",
			tokens: []lexer.Token{
				lexer.Token{Type: TAbsolute, Value: "$"},
				lexer.Token{Type: TChildStart, Value: "["},
				lexer.Token{Type: TQuotedName, Value: "'one'"},
				lexer.Token{Type: TChildEnd, Value: "]"},
			},
		},
	}

	for _, tt := range tests {
		t.Log("testing path: ", tt.path)
		l := lexer.New(tt.path, pathState)
		l.Start()

		func() {
			for i, expected := range tt.tokens {
				actual, done := l.NextToken()
				if done || actual == nil {
					t.Errorf("Lexer(%q) finished early, expecting %v", tt.path, expected)
					return
				}
				if actual.Type != expected.Type {
					t.Errorf("Lexer(%q) token %d => %s, expected %s", tt.path, i, tokenNames[actual.Type], tokenNames[expected.Type])
					return
				}
				if actual.Value != expected.Value {
					t.Errorf("Lexer(%q) token %d =>  %v, expected %v", tt.path, i, actual, expected)
					return
				}
			}
		}()
	}
}
