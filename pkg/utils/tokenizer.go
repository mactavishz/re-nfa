package utils

import (
	"fmt"
	"unicode/utf8"
)

// Token types
type TokenType int

const (
	_              = iota
	CHAR TokenType = iota
	STAR
	PLUS
	OPTIONAL
	LPAREN
	RPAREN
	OR
	CARET
	DOLLAR
	EOF
)

// Token structure
type Token struct {
	Type  TokenType
	Value rune
}

// Tokenizer
type Tokenizer struct {
	input string
	pos   int
	buf   Token
}

func (t TokenType) String() string {
	switch t {
	case CHAR:
		return "CHAR"
	case STAR:
		return "STAR"
	case PLUS:
		return "PLUS"
	case OPTIONAL:
		return "OPTIONAL"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
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

func (t *Tokenizer) peek() Token {
	if t.buf.Type == 0 {
		t.buf = t.next()
	}
	return t.buf
}

func (t *Tokenizer) next() Token {
	if t.buf.Type != 0 {
		res := t.buf
		t.buf = Token{}
		return res
	}
	if t.pos >= len(t.input) {
		return Token{Type: EOF}
	}
	r, size := utf8.DecodeRuneInString(t.input[t.pos:])
	t.pos += size
	switch r {
	case '(':
		return Token{Type: LPAREN, Value: r}
	case ')':
		return Token{Type: RPAREN, Value: r}
	case '*':
		return Token{Type: STAR, Value: r}
	case '+':
		return Token{Type: PLUS, Value: r}
	case '?':
		return Token{Type: OPTIONAL, Value: r}
	case '|':
		return Token{Type: OR, Value: r}
	default:
		return Token{Type: CHAR, Value: r}
	}
}
