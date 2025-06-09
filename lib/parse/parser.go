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
	// we only consume EOF at the top level
	_, err := p.match(EOF)
	if err != nil {
		panic(err)
	}
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
	if p.check(CHAR) {
		return p.parseChar()
	} else if p.check(LPAREN) {
		return p.parseGroup()
	} else {
		panic("expected a Character or a Group")
	}
}

func (p *Parser) parseMTerm() ASTNode {
	item := p.parseItem()
	var t ModifierType
	var err error
	if p.check(STAR) {
		_, err = p.match(STAR)
		t = ModifierStar
	} else if p.check(PLUS) {
		_, err = p.match(PLUS)
		t = ModifierPlus
	} else if p.check(OPT) {
		_, err = p.match(OPT)
		t = ModifierQuestion
	} else {
		panic("expected a '*', '+' or '?'")
	}
	if err != nil {
		panic(err)
	}
	return &Modifier{
		Child: item,
		Type:  t,
	}
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
