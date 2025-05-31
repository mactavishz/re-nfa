package parse

import (
	"fmt"
	"strings"
)

// ModifierType enum
type ModifierType int

const (
	ModifierStar     ModifierType = iota // * (0 or more)
	ModifierPlus                         // + (1 or more)
	ModifierQuestion                     // ? (0 or 1)
	ModifierRange                        // {n,m} for future use
)

type ASTNode interface {
	ToRegex() string
	NodeType() string
	String() string
}

// Modifier node - represents repetition modifiers
type Modifier struct {
	Child ASTNode
	Type  ModifierType
	Min   int // for future extension to {n,m} syntax
	Max   int // -1 for unlimited
}

func (q *Modifier) ToRegex() string {
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
	return fmt.Sprintf("%s%s", q.Child.ToRegex(), str)
}

func (q *Modifier) NodeType() string {
	return "Quantifier"
}

func (q *Modifier) String() string {
	return fmt.Sprintf("Modifier {\n  Child: %s,\n  Type: %d\n}", indentString(q.Child.String(), 2), q.Type)
}

// Top-level AstNode
type Expression struct {
	Child ASTNode
}

func (e *Expression) ToRegex() string {
	return e.Child.ToRegex()
}

func (e *Expression) NodeType() string {
	return "Expression"
}

func (e *Expression) String() string {
	return fmt.Sprintf("Expression {\n  Child: %s\n}", indentString(e.Child.String(), 2))
}

type Term struct {
	Child       ASTNode
	HasModifier bool
}

func (t *Term) ToRegex() string {
	return t.Child.ToRegex()
}

func (t *Term) NodeType() string {
	if t.HasModifier {
		return "TermModified"
	} else {
		return "Term"
	}
}

func (t *Term) String() string {
	modifier := "ModifierNone"
	if t.HasModifier {
		modifier = "Modified"
	}
	return fmt.Sprintf("Term {\n  Child: %s,\n  Modifier: %s\n}", indentString(t.Child.String(), 2), modifier)
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

func (cat *Concatenation) ToRegex() string {
	return fmt.Sprintf("%s%s", cat.Left.ToRegex(), cat.Right.ToRegex())
}

func (cat *Concatenation) NodeType() string {
	return "Concatenation"
}

func (cat *Concatenation) String() string {
	return fmt.Sprintf("Concatenation {\n  Left: %s,\n  Right: %s\n}", indentString(cat.Left.String(), 2), indentString(cat.Right.String(), 2))
}

// Character node - represents a single character
type Character struct {
	Value rune
}

func (c *Character) ToRegex() string {
	return fmt.Sprintf("%c", c.Value)
}

func (c *Character) NodeType() string {
	return "Character"
}

func (ch *Character) String() string {
	return fmt.Sprintf("Character {Value: '%c'}", ch.Value)
}

// Group node - represents parenthesized expression
type Group struct {
	Expr ASTNode
}

func (g *Group) ToRegex() string {
	return fmt.Sprintf("(%s)", g.Expr.ToRegex())
}

func (g *Group) NodeType() string {
	return "Group"
}

func (g *Group) String() string {
	return fmt.Sprintf("Group {\n  Expr: %s\n}", indentString(g.Expr.String(), 2))
}

func (a *Alternation) ToRegex() string {
	return fmt.Sprintf("%s|%s", a.Left.ToRegex(), a.Right.ToRegex())
}

func (a *Alternation) NodeType() string {
	return "Alternation"
}

func (a *Alternation) String() string {
	return fmt.Sprintf("Alternation {\n  Left: %s,\n  Right: %s\n}", indentString(a.Left.String(), 2), indentString(a.Right.String(), 2))
}

// indentString indents each line of s by n spaces
func indentString(s string, n int) string {
	if n <= 0 {
		return s
	}

	prefix := strings.Repeat(" ", n)
	lines := strings.Split(s, "\n")

	for i := range lines {
		if i > 0 { // Don't indent the first line
			lines[i] = prefix + lines[i]
		}
	}

	return strings.Join(lines, "\n")
}
