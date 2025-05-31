package parse

import (
	"fmt"
)

/*
*
* Grammar for the regular expression (in EBNF):
*	expr := term ['|' term]*
*   term := mterm*
*   mterm := item [modifier]
*	modifier := '*' | '+' | '?'
*	item := char | group
*	group := '(' expr ')'
*	char -> 'a'| 'b' | 'c' | ..., UTF8 character excluding '|', '*', '+', '?', '(', ')'
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
