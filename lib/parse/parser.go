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
*	char -> 'a'| 'b' | 'c' | ..., UTF8 character excluding '.', '|', '*', '+', '?', '(', ')'
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

func (p *Parser) match(expected TokenType) (*Token, error) {
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
	root := &Expression{}
	root.Child = p.parseExpr()
	return root
}

func (p *Parser) parseTerm() ASTNode {
	panic("TODO")
}

func (p *Parser) parseExpr() ASTNode {
	lTerm := p.parseTerm()
	if lTerm == nil {
		return nil
	}
	// an expression is finished only when the next token is EOF or ')'
	// (it is wrapped in a pair of parentheses)
	if p.check(EOF) || p.check(RPAREN) {
		return lTerm
	}
	var result ASTNode
	for {
		_, err := p.match(OR)
		if err != nil {
			panic(err)
		}
		rTerm := p.parseTerm()
		if rTerm == nil {
			panic("expected a term after |")
		}
		if result == nil {
			result = &Alternation{
				Left:  lTerm,
				Right: rTerm,
			}
		} else {
			result = &Alternation{
				Left:  result,
				Right: rTerm,
			}
		}
		if !p.check(OR) {
			return result
		}
	}
}

func (p *Parser) parseGroup() ASTNode {
	_, err := p.match(LPAREN)
	if err != nil {
		panic(err)
	}
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}
	_, err = p.match(LPAREN)
	if err != nil {
		panic(err)
	}
	return &Group{
		Expr: expr,
	}
}

func (p *Parser) parseItem() ASTNode {
	if p.check(EOF) {
		return nil
	}
	panic("TODO")
}

func (p *Parser) parseMTerm() ASTNode {
	if p.check(EOF) {
		return nil
	}
	panic("TODO")
}

func (p *Parser) parseChar() ASTNode {
	token, err := p.match(CHAR)
	if err != nil {
		panic(err)
	}
	return &Character{
		Value: token.Value,
	}
}
