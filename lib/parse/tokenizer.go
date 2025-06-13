package parse

import (
	"fmt"
	"unicode/utf8"
)

// Token types
type TokenType int

const (
	CHAR TokenType = iota
	STAR
	PLUS
	OPT
	DOT
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	COMMA
	OR
	EOF
)

// Token structure
type Token struct {
	Type   TokenType
	Column int
	Value  rune
}

// Tokenizer, only caches one token in advance
type Tokenizer struct {
	input string
	pos   int
	buf   *Token
}

func (t TokenType) String() string {
	switch t {
	case CHAR:
		return "CHAR"
	case STAR:
		return "STAR"
	case PLUS:
		return "PLUS"
	case OPT:
		return "OPT"
	case DOT:
		return "DOT"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case COMMA:
		return "COMMA"
	case OR:
		return "OR"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Type.String(), string(t.Value))
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: input}
}

func (t *Tokenizer) peek() *Token {
	if t.buf == nil {
		t.buf = t.consume()
	}
	return t.buf
}

func (t *Tokenizer) consume() *Token {
	if t.buf != nil {
		res := t.buf
		t.buf = nil
		return res
	}
	if t.pos >= len(t.input) {
		return &Token{Type: EOF, Column: len(t.input) - 1}
	}
	r, size := utf8.DecodeRuneInString(t.input[t.pos:])
	oldPos := t.pos
	t.pos += size
	switch r {
	case '(':
		return &Token{Type: LPAREN, Value: r, Column: oldPos}
	case ')':
		return &Token{Type: RPAREN, Value: r, Column: oldPos}
	case '{':
		return &Token{Type: LBRACE, Value: r, Column: oldPos}
	case '}':
		return &Token{Type: RBRACE, Value: r, Column: oldPos}
	case '.':
		return &Token{Type: DOT, Value: r, Column: oldPos}
	case ',':
		return &Token{Type: COMMA, Value: r, Column: oldPos}
	case '*':
		return &Token{Type: STAR, Value: r, Column: oldPos}
	case '+':
		return &Token{Type: PLUS, Value: r, Column: oldPos}
	case '?':
		return &Token{Type: OPT, Value: r, Column: oldPos}
	case '|':
		return &Token{Type: OR, Value: r, Column: oldPos}
	default:
		return &Token{Type: CHAR, Value: r, Column: oldPos}
	}
}
