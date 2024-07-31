package utils

import (
	"fmt"

	"github.com/mactavishz/re-nfa/pkg/fsm"
)

/*
* Grammar for the regular expression:
* expr   -> term [ '|' term ]*
* term   -> factor*
* factor -> primary ['*' | '+' | '?']
* primary -> CHAR | '(' expr ')'
* CHAR   -> UTF8 character, excluding '|', '*', '+', '?', '(', ')'
* */

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

func (p *Parser) Parse() (*fsm.NFA, error) {
	nfa, err := p.expr()
	if err != nil {
		return nil, err
	}
	_, err = p.match(EOF)
	if err != nil {
		return nil, err
	}
	return nfa, nil
}

func (p *Parser) factor() (*fsm.NFA, error) {
	nfa, err := p.primary()
	if err != nil {
		return nil, err
	}
	if p.check(STAR) {
		_, err = p.match(STAR)
		if err != nil {
			return nil, err
		}
		nfa = fsm.Star(nfa)
	} else if p.check(PLUS) {
		_, err = p.match(PLUS)
		if err != nil {
			return nil, err
		}
		nfa = fsm.Plus(nfa)
	} else if p.check(OPTIONAL) {
		_, err = p.match(OPTIONAL)
		if err != nil {
			return nil, err
		}
		nfa = fsm.Optional(nfa)
	}
	return nfa, nil
}

func (p *Parser) term() (*fsm.NFA, error) {
	nfa, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.check(CHAR) || p.check(LPAREN) {
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		nfa = fsm.Concat(nfa, right)
	}
	return nfa, nil
}

func (p *Parser) expr() (*fsm.NFA, error) {
	nfa, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.check(OR) {
		_, err = p.match(OR)
		if err != nil {
			return nil, err
		}
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		nfa = fsm.Or(nfa, right)
	}
	return nfa, nil
}

func (p *Parser) primary() (*fsm.NFA, error) {
	if p.check(LPAREN) {
		_, err := p.match(LPAREN)
		if err != nil {
			return nil, err
		}
		nfa, err := p.expr()
		if err != nil {
			return nil, err
		}
		_, err = p.match(RPAREN)
		if err != nil {
			return nil, err
		}
		return nfa, nil
	} else {
		return p.char()
	}
}

func (p *Parser) char() (*fsm.NFA, error) {
	token, err := p.match(CHAR)
	if err != nil {
		return nil, err
	}
	start, accept := fsm.NewState(false), fsm.NewState(true)
	start.AddTransition(accept, token.Value)
	return fsm.NewNFA(start, accept), nil
}
