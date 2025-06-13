package parse

import (
	"testing"
)

func TestTokenizer_SimpleTokens(t *testing.T) {
	input := "a*+?(){}.,|"
	tok := NewTokenizer(input)
	expected := []struct {
		type_ TokenType
		val   rune
	}{
		{CHAR, 'a'},
		{STAR, '*'},
		{PLUS, '+'},
		{OPT, '?'},
		{LPAREN, '('},
		{RPAREN, ')'},
		{LBRACE, '{'},
		{RBRACE, '}'},
		{DOT, '.'},
		{COMMA, ','},
		{OR, '|'},
		{EOF, 0},
	}
	for i, exp := range expected {
		tokn := tok.consume()
		if tokn.Type != exp.type_ || tokn.Value != exp.val {
			t.Errorf("token %d: expected %s(%q), got %s(%q)", i, exp.type_.String(), exp.val, tokn.Type.String(), tokn.Value)
		}
	}
}

func TestTokenizer_EmptyInput(t *testing.T) {
	tok := NewTokenizer("")
	tokn := tok.consume()
	if tokn.Type != EOF {
		t.Errorf("expected EOF, got %s", tokn.Type.String())
	}
}

func TestTokenizer_UnicodeChar(t *testing.T) {
	input := "λ"
	tok := NewTokenizer(input)
	tokn := tok.consume()
	if tokn.Type != CHAR || tokn.Value != 'λ' {
		t.Errorf("expected CHAR(λ), got %s(%q)", tokn.Type.String(), tokn.Value)
	}
	if tok.consume().Type != EOF {
		t.Errorf("expected EOF after unicode char")
	}
}

func TestTokenizer_PeekAndNext(t *testing.T) {
	input := "ab*"
	tok := NewTokenizer(input)
	p1 := tok.peek()
	n1 := tok.consume()
	if p1.Type != CHAR || n1.Type != CHAR || p1.Value != 'a' || n1.Value != 'a' {
		t.Errorf("peek/next mismatch: peek %s(%q), next %s(%q)", p1.Type.String(), p1.Value, n1.Type.String(), n1.Value)
	}
	p2 := tok.peek()
	n2 := tok.consume()
	if p2.Type != CHAR || n2.Type != CHAR || p2.Value != 'b' || n2.Value != 'b' {
		t.Errorf("peek/next mismatch: peek %s(%q), next %s(%q)", p2.Type.String(), p2.Value, n2.Type.String(), n2.Value)
	}
	p3 := tok.peek()
	n3 := tok.consume()
	if p3.Type != STAR || n3.Type != STAR {
		t.Errorf("peek/next mismatch: peek %s, next %s", p3.Type.String(), n3.Type.String())
	}
	if tok.consume().Type != EOF {
		t.Errorf("expected EOF at end")
	}
}

func TestTokenizer_MultipleEOF(t *testing.T) {
	tok := NewTokenizer("")
	for range 3 {
		tokn := tok.consume()
		if tokn.Type != EOF {
			t.Errorf("expected EOF, got %s", tokn.Type.String())
		}
	}
}
