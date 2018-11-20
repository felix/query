package jsonpath

import (
	//"fmt"
	"strings"

	"src.userspace.com.au/query/lexer"
)

const (
	_ lexer.TokenType = iota

	TAbsolute       // $
	TCurrent        // @
	TChildDot       // .
	TRecursive      // ..
	TChildStart     // [
	TChildEnd       // ]
	TQuotedName     // ' or "
	TName           // A word-like
	TWildcard       // *
	TFilterStart    // ?(
	TFilterEnd      // )
	TPredicateStart // [
	TPredicateEnd   // ]
	TNumber         // [-]0-9 and hex also
	TRange          // :
	TIndex          // [234]
	TUnion          // [,]
	TSlice          // [start,end,step]
	TScriptStart    // (
	TScriptEnd      // )
)

// helpful for debugging
var tokenNames = map[lexer.TokenType]string{
	lexer.ErrorToken: "TError",
	lexer.EOFToken:   "EOF",
	TAbsolute:        "TAbsolute",
	TCurrent:         "TCurrent",
	TChildDot:        "TChildDot",
	TRecursive:       "TRecursive",
	TChildStart:      "TChildStart",
	TChildEnd:        "TChildEnd",
	TQuotedName:      "TQuoteName",
	TName:            "TName",
	TWildcard:        "TWildcard",
	TFilterStart:     "TFilterStart",
	TFilterEnd:       "TFilterEnd",
	TPredicateStart:  "TPredicateStart",
	TPredicateEnd:    "TPredicateEnd",
	TNumber:          "TNumber",
	TRange:           "TRange",
	TIndex:           "TIndex",
	TUnion:           "TUnion",
	TSlice:           "TSlice",
	TScriptStart:     "TScriptStart",
	TScriptEnd:       "TScriptEnd",
}

const (
	alpha         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric       = "0123456789"
	alphanumeric  = alpha + numeric
	signednumeric = "-" + numeric
)

func pathState(l *lexer.Lexer) lexer.StateFunc {
	l.SkipWhitespace()
	if l.Accept("$") {
		l.Emit(TAbsolute)
		return stepState
	}
	return nil
}

// stepState starts with a . or bracket
func stepState(l *lexer.Lexer) lexer.StateFunc {
	for {
		switch t := l.Next(); {
		case t == '.':
			// Don't emit dot as it is used for subsequent child
			if l.Peek() == '.' {
				l.Emit(TRecursive)
				return stepState
			}
			l.Emit(TChildDot)
			return childState

		case t == '[':
			l.Emit(TChildStart)
			return childState

		case t == lexer.EOFRune:
			l.Emit(lexer.EOFToken)
			return nil

		default:
			return l.ErrorState("abrupt end to stepState")
		}
	}
}

// childState has started with . or [
func childState(l *lexer.Lexer) lexer.StateFunc {
	for {
		switch t := l.Next(); {

		case t == '*':
			l.Emit(TWildcard)
			return childState

		case t == '\'' || t == '"':
			// FIXME what other characters?
			l.AcceptRun(alphanumeric + "-_")
			l.Accept(string(t))
			l.Emit(TQuotedName)
			return childState

		case strings.IndexRune(alpha, t) != -1:
			// Key lookup
			l.AcceptRun(alphanumeric)
			l.Emit(TName)
			return childState

		case t == '.':
			l.Backup()
			return stepState

		case t == '[':
			// FIXME predicate or another child
			l.Emit(TPredicateStart)
			return predicateState

		case t == ']':
			l.Emit(TChildEnd)
			return stepState

		case t == lexer.EOFRune:
			l.Emit(lexer.EOFToken)
			return nil

		default:
			return l.ErrorState("abrupt end to childState")
		}
	}
}

func predicateState(l *lexer.Lexer) lexer.StateFunc {
	for {
		switch t := l.Next(); {
		case t == '?':
			if l.Accept("(") {
				l.Emit(TFilterStart)
				return filterState
			}
			return l.ErrorState("expecting (")

		case t == '(':
			l.Emit(TScriptStart)
			return scriptState

		case t == ']':
			l.Emit(TPredicateEnd)
			return stepState

		default:
			l.Backup()
			return predicateExprState
		}
	}
}

func predicateExprState(l *lexer.Lexer) lexer.StateFunc {
	open := true
	inRange := false
	for {
		switch t := l.Next(); {

		case t == '-' || '0' <= t && t <= '9':
			l.AcceptRun(numeric)
			l.Emit(TNumber)
			open = false

		case t == ':':
			if inRange {
				return l.ErrorState("invalid range")
			}
			if open {
				l.Backup()
				l.Emit(TNumber)
				l.Accept(":")
			}
			l.Emit(TRange)
			inRange = true
			open = true

		case t == ',':
			l.Emit(TUnion)

		case t == ']':
			l.Backup()
			if open && inRange {
				l.Emit(TNumber)
			}
			return predicateState

		default:
			return l.ErrorState("invalid predicate expression")
		}
	}
}

// filterState TODO
func filterState(l *lexer.Lexer) lexer.StateFunc {
	for {
		switch t := l.Next(); {
		case t == ')':
			l.Emit(TFilterEnd)
			return predicateState
		default:
			return l.ErrorState("invalid filter expression")
		}
	}
}

// scriptState TODO
func scriptState(l *lexer.Lexer) lexer.StateFunc {
	for {
		switch t := l.Next(); {
		case t == ')':
			l.Emit(TScriptEnd)
			return predicateState
		default:
			return l.ErrorState("invalid script expression")
		}
	}
}
