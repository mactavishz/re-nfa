package parse

import (
	"fmt"
	"strings"
)

/*
*
* Grammar for the regular expression (in EBNF):
*	expr := [term ['|' expr]]
*   term := mterm+
*   mterm := item [modifier]
*	modifier := '*' | '+' | '?'
*	item := char | group
*	group := '(' expr ')'
*	char -> 'a'| 'b' | 'c' | ..., UTF8 character excluding '.', '|', '*', '+', '?', '(', ')'
*
 */

type SyntaxError struct {
	Column  int
	Message string
	Line    string
}

func (s SyntaxError) Error() string {
	res := fmt.Sprintf("SyntaxError (column %d): %s\n", s.Column, s.Message)
	res += fmt.Sprintf(">>>%s\n", s.Line)
	res += fmt.Sprintf(">>>%s^\n", strings.Repeat("-", s.Column))
	return res
}

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(input string) *Parser {
	tokenizer := NewTokenizer(input)
	parser := &Parser{tokenizer: tokenizer}
	return parser
}

func (p *Parser) match(expected TokenType) (*Token, error) {
	t := p.tokenizer.consume()
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

func (p *Parser) Parse() (ASTNode, error) {
	root := &Expression{}
	child, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	root.Child = child
	// we only consume EOF at the top level
	token, err := p.match(EOF)
	if err != nil {
		return nil, &SyntaxError{
			Column:  token.Column,
			Message: "expected EOF",
			Line:    p.tokenizer.input,
		}
	}
	return root, nil
}

func (p *Parser) parseTerm() (ASTNode, error) {
	if p.check(EOF) || p.check(OR) {
		return nil, &SyntaxError{
			Column:  p.tokenizer.buf.Column,
			Message: "expected a term",
			Line:    p.tokenizer.input,
		}
	}
	lMTerm, err := p.parseMTerm()
	if err != nil {
		return nil, err
	}

	for !p.check(EOF) && !p.check(OR) && !p.check(RPAREN) {
		// this must be correct
		rMTerm, err := p.parseMTerm()
		if err != nil {
			return nil, err
		}
		lMTerm = &Concatenation{
			Left:  lMTerm,
			Right: rMTerm,
		}
	}
	return lMTerm, nil
}

func (p *Parser) parseExpr() (ASTNode, error) {
	// handle empty expression, when the next token is EOF or ')' (it is wrapped in a pair of parentheses)
	if p.check(EOF) || p.check(RPAREN) {
		return nil, nil
	}
	lTerm, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	if lTerm == nil {
		return nil, &SyntaxError{
			Column:  p.tokenizer.buf.Column,
			Message: "expected a term",
			Line:    p.tokenizer.input,
		}
	}
	for p.check(OR) {
		token, err := p.match(OR)
		if err != nil {
			return nil, &SyntaxError{
				Column:  token.Column,
				Message: "expected a | or EOF",
				Line:    p.tokenizer.input,
			}
		}
		rTerm, err := p.parseTerm()
		if err != nil {
			return nil, &SyntaxError{
				Column:  p.tokenizer.buf.Column,
				Message: "expected a term after |",
				Line:    p.tokenizer.input,
			}
		}
		if rTerm == nil {
			return nil, &SyntaxError{
				Column:  p.tokenizer.buf.Column,
				Message: "expected a term after |",
				Line:    p.tokenizer.input,
			}
		}
		lTerm = &Alternation{
			Left:  lTerm,
			Right: rTerm,
		}
	}
	return lTerm, nil
}

func (p *Parser) parseGroup() (ASTNode, error) {
	token, err := p.match(LPAREN)
	if err != nil {
		return nil, &SyntaxError{
			Column:  token.Column,
			Message: "expected a (",
			Line:    p.tokenizer.input,
		}
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	token, err = p.match(RPAREN)
	if err != nil {
		return nil, &SyntaxError{
			Column:  token.Column,
			Message: "expected a )",
			Line:    p.tokenizer.input,
		}
	}
	return &Group{
		Expr: expr,
	}, nil
}

func (p *Parser) parseItem() (ASTNode, error) {
	if p.check(CHAR) {
		return p.parseChar()
	} else if p.check(LPAREN) {
		return p.parseGroup()
	} else {
		return nil, &SyntaxError{
			Column:  p.tokenizer.buf.Column,
			Message: "expected a Character or a Group",
			Line:    p.tokenizer.input,
		}
	}
}

func (p *Parser) parseMTerm() (ASTNode, error) {
	item, err := p.parseItem()
	if err != nil {
		return nil, err
	}
	var t ModifierType
	if p.check(STAR) {
		_, err = p.match(STAR)
		t = ModifierStar
	} else if p.check(PLUS) {
		_, err = p.match(PLUS)
		t = ModifierPlus
	} else if p.check(OPT) {
		_, err = p.match(OPT)
		t = ModifierQuestion
	} else if p.check(CHAR) || p.check(LPAREN) || p.check(RPAREN) || p.check(OR) || p.check(EOF) {
		// if this is not a modified term
		return item, nil
	} else {
		return nil, &SyntaxError{
			Column:  p.tokenizer.buf.Column,
			Message: "expected a '*', '+' or '?'",
			Line:    p.tokenizer.input,
		}
	}
	if err != nil {
		return nil, err
	}
	return &Modifier{
		Child: item,
		Type:  t,
	}, nil
}

func (p *Parser) parseChar() (ASTNode, error) {
	token, err := p.match(CHAR)
	if err != nil {
		return nil, &SyntaxError{
			Column:  token.Column,
			Message: "expected a character",
			Line:    p.tokenizer.input,
		}
	}
	return &Character{
		Value: token.Value,
	}, nil
}
