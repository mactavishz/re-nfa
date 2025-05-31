package parse

import (
	"fmt"
)

/*
*
* Grammar for the regular expression (in EBNF):
*	expr ::= seq_expr alt_tail | seq_expr
*   alt_tail ::= '|' seq_expr alt_tail | '|' seq_expr
*   seq_expr ::= term seq_expr | term
*	term ::= item modifier | item
*	modifier ::= '*' | '+' | '?'
*	item ::= char | group
*	group ::= '(' expr ')'
*	char -> UTF8 character, excluding '|', '*', '+', '?', '(', ')'
*
 */

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(input string) *Parser {
	tokenizer := NewTokenizer(input)
	parser := &Parser{tokenizer: tokenizer}
	return parser
}

func (p *Parser) match(expected TokenType) (Token, error) {
	t := p.tokenizer.next()
	if t.Type != expected {
		return t, fmt.Errorf("unexpected token, expected %s but got %s", expected.String(), t.String())
	} else {
		return t, nil
	}
}

func (p *Parser) check(expected TokenType) bool {
	t := p.tokenizer.peek()
	if t.Type == expected {
		return true
	} else {
		return false
	}
}

func (p *Parser) Parse() ASTNode {
	panic("TODO")
}
