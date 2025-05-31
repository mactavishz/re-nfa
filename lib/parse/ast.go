package parse

import "fmt"

// ModifierType enum
type ModifierType int

const (
	ModifierStar     ModifierType = iota // * (0 or more)
	ModifierPlus                         // + (1 or more)
	ModifierQuestion                     // ? (0 or 1)
	ModifierRange                        // {n,m} for future use
)

type ASTNode interface {
	String() string
	NodeType() string
}

// Modifier node - represents repetition modifiers
type Modifier struct {
	Child ASTNode
	Type  ModifierType
	Min   int // for future extension to {n,m} syntax
	Max   int // -1 for unlimited
}

func (q *Modifier) String() string {
	var str string
	switch q.Type {
	case ModifierStar:
		str = "*"
	case ModifierPlus:
		str = "+"
	case ModifierQuestion:
		str = "?"
	case ModifierRange:
		str = fmt.Sprintf("{%d,%d}", q.Min, q.Max)
	default:
		panic("invalid quantifiler")
	}
	return fmt.Sprintf("%s%s", q.Child.String(), str)
}

func (q *Modifier) NodeType() string {
	return "Quantifier"
}

// Top-level AstNode
type Expression struct {
	Child ASTNode
}

func (e *Expression) String() string {
	return e.Child.String()
}

func (e *Expression) NodeType() string {
	return "Expression"
}

type Term struct {
	Child       ASTNode
	HasModifier bool
}

func (t *Term) String() string {
	return t.Child.String()
}

func (t *Term) NodeType() string {
	if t.HasModifier {
		return "TermModified"
	} else {
		return "Term"
	}
}

// Alternation node - represents OR operation (a|b)
type Alternation struct {
	Left  ASTNode
	Right ASTNode
}

// Concatenation node - represents sequence of terms (ab)
type Concatenation struct {
	Left  ASTNode
	Right ASTNode
}

func (cat *Concatenation) String() string {
	return fmt.Sprintf("%s%s", cat.Left.String(), cat.Right.String())
}

func (cat *Concatenation) NodeType() string {
	return "Concatenation"
}

// Character node - represents a single character
type Character struct {
	Value rune
}

func (c *Character) String() string {
	return fmt.Sprintf("%c", c.Value)
}

func (c *Character) NodeType() string {
	return "Character"
}

// Group node - represents parenthesized expression
type Group struct {
	Expr ASTNode
}

func (g *Group) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (g *Group) NodeType() string {
	return "Group"
}

func (a *Alternation) String() string {
	return fmt.Sprintf("%s|%s", a.Left.String(), a.Right.String())
}

func (a *Alternation) NodeType() string {
	return "Alternation"
}
