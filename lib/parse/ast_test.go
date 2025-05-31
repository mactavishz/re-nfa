package parse

import (
	"testing"
)

// TestCharacterNode tests a single character node (regex: a)
func TestCharacterNode(t *testing.T) {
	n := &Character{Value: 'a'}
	if n.String() != "a" {
		t.Errorf("expected 'a', got '%s'", n.String())
	}
	if n.NodeType() != "Character" {
		t.Errorf("expected 'Character', got '%s'", n.NodeType())
	}
}

// TestModifierNodes tests all modifier types (regex: a*, a+, a?, a{2,5})
func TestModifierNodes(t *testing.T) {
	base := &Character{Value: 'a'}
	mods := []struct {
		mod   *Modifier
		expS  string
		typeS string
	}{
		{&Modifier{Child: base, Type: ModifierStar}, "a*", "Quantifier"},
		{&Modifier{Child: base, Type: ModifierPlus}, "a+", "Quantifier"},
		{&Modifier{Child: base, Type: ModifierQuestion}, "a?", "Quantifier"},
		{&Modifier{Child: base, Type: ModifierRange, Min: 2, Max: 5}, "a{2,5}", "Quantifier"},
	}
	for _, m := range mods {
		if m.mod.String() != m.expS {
			t.Errorf("expected '%s', got '%s'", m.expS, m.mod.String())
		}
		if m.mod.NodeType() != m.typeS {
			t.Errorf("expected '%s', got '%s'", m.typeS, m.mod.NodeType())
		}
	}
}

// TestConcatenationNode tests concatenation (regex: ab)
func TestConcatenationNode(t *testing.T) {
	n := &Concatenation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	if n.String() != "ab" {
		t.Errorf("expected 'ab', got '%s'", n.String())
	}
	if n.NodeType() != "Concatenation" {
		t.Errorf("expected 'Concatenation', got '%s'", n.NodeType())
	}
}

// TestAlternationNode tests alternation (regex: a|b)
func TestAlternationNode(t *testing.T) {
	n := &Alternation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	// Alternation has no String() method, so we check structure
	if n.Left.String() != "a" || n.Right.String() != "b" {
		t.Errorf("expected left 'a' and right 'b', got '%s' and '%s'", n.Left.String(), n.Right.String())
	}
}

// TestLongAlternationNode tests long alternation, alternation should be left-associative (regex: a|b|c|d)
func TestLongAlternationNode(t *testing.T) {
	// Build AST for a|b|c|d: Alternation(Alternation(Alternation(a, b), c), d)
	n := &Alternation{
		Left: &Alternation{
			Left: &Alternation{
				Left:  &Character{Value: 'a'},
				Right: &Character{Value: 'b'},
			},
			Right: &Character{Value: 'c'},
		},
		Right: &Character{Value: 'd'},
	}

	// Check the structure recursively
	alt1, ok := n.Left.(*Alternation)
	if !ok {
		t.Fatalf("expected n.Left to be Alternation")
	}
	alt2, ok := alt1.Left.(*Alternation)
	if !ok {
		t.Fatalf("expected n.Left.Left to be Alternation")
	}

	// Check leaf nodes
	a, ok := alt2.Left.(*Character)
	if !ok || a.Value != 'a' {
		t.Errorf("expected leftmost to be 'a', got %#v", alt2.Left)
	}
	b, ok := alt2.Right.(*Character)
	if !ok || b.Value != 'b' {
		t.Errorf("expected second to be 'b', got %#v", alt2.Right)
	}
	c, ok := alt1.Right.(*Character)
	if !ok || c.Value != 'c' {
		t.Errorf("expected third to be 'c', got %#v", alt1.Right)
	}
	d, ok := n.Right.(*Character)
	if !ok || d.Value != 'd' {
		t.Errorf("expected rightmost to be 'd', got %#v", n.Right)
	}
	// Alternation has no String() method, so we check structure
	if n.String() != "a|b|c|d" {
		t.Errorf("expected left 'a|b|c|d' , got '%s'", n.String())
	}
}

// TestGroupNode tests grouping (regex: (ab))
func TestGroupNode(t *testing.T) {
	cat := &Concatenation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	g := &Group{Expr: cat}
	if g.String() != "(ab)" {
		t.Errorf("expected '(ab)', got '%s'", g.String())
	}
	if g.NodeType() != "Group" {
		t.Errorf("expected 'Group', got '%s'", g.NodeType())
	}
}

// TestExpressionNode tests the Expression wrapper (regex: a)
func TestExpressionNode(t *testing.T) {
	c := &Character{Value: 'a'}
	e := &Expression{Child: c}
	if e.String() != "a" {
		t.Errorf("expected 'a', got '%s'", e.String())
	}
	if e.NodeType() != "Expression" {
		t.Errorf("expected 'Expression', got '%s'", e.NodeType())
	}
}

// TestTermNode tests Term with and without modifier (regex: a, a*)
func TestTermNode(t *testing.T) {
	c := &Character{Value: 'a'}
	term := &Term{Child: c, HasModifier: false}
	if term.String() != "a" {
		t.Errorf("expected 'a', got '%s'", term.String())
	}
	if term.NodeType() != "Term" {
		t.Errorf("expected 'Term', got '%s'", term.NodeType())
	}
	m := &Modifier{Child: c, Type: ModifierStar}
	tm := &Term{Child: m, HasModifier: true}
	if tm.String() != "a*" {
		t.Errorf("expected 'a*', got '%s'", tm.String())
	}
	if tm.NodeType() != "TermModified" {
		t.Errorf("expected 'TermModified', got '%s'", tm.NodeType())
	}
}

// TestNestedGroupsAndModifiers tests complex regex: ((a|b)*c+)?
func TestNestedGroupsAndModifiers(t *testing.T) {
	alt := &Alternation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	group1 := &Group{Expr: alt}                          // (a|b)
	mod1 := &Modifier{Child: group1, Type: ModifierStar} // (a|b)*
	cat := &Concatenation{
		Left:  mod1,
		Right: &Modifier{Child: &Character{Value: 'c'}, Type: ModifierPlus},
	} // (a|b)*c+
	group2 := &Group{Expr: cat}                              // ((a|b)*c+)
	mod2 := &Modifier{Child: group2, Type: ModifierQuestion} // ((a|b)*c+)?
	if mod2.String() != "((a|b)*c+)?" {
		t.Errorf("expected '((a|b)*c+)?', got '%s'", mod2.String())
	}
}

// TestEdgeCases tests empty group, deeply nested, and invalid modifier range
func TestEdgeCases(t *testing.T) {
	// Empty group: ()
	emptyGroup := &Group{Expr: &Concatenation{Left: &Character{Value: '\u0000'}, Right: &Character{Value: '\u0000'}}}
	// This is a hack, as there's no explicit Empty node
	if emptyGroup.String() != "(\u0000\u0000)" {
		t.Errorf("expected '(\u0000\u0000)', got '%s'", emptyGroup.String())
	}

	// Deeply nested: (((a)))
	deep := &Group{Expr: &Group{Expr: &Group{Expr: &Character{Value: 'a'}}}}
	if deep.String() != "(((a)))" {
		t.Errorf("expected '(((a)))', got '%s'", deep.String())
	}
}
