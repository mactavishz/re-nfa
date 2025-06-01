package parse

import (
	"testing"
)

// TestASTCharacterNode tests a single character node (regex: a)
func TestASTCharacterNode(t *testing.T) {
	n := &Character{Value: 'a'}
	// fmt.Println(n.String())
	if n.ToRegex() != "a" {
		t.Errorf("expected 'a', got '%s'", n.ToRegex())
	}
	if n.NodeType() != "Character" {
		t.Errorf("expected 'Character', got '%s'", n.NodeType())
	}
}

// TestASTModifierNodes tests all modifier types (regex: a*, a+, a?, a{2,5})
func TestASTModifierNodes(t *testing.T) {
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
		// fmt.Println(m.mod.String())
		if m.mod.ToRegex() != m.expS {
			t.Errorf("expected '%s', got '%s'", m.expS, m.mod.ToRegex())
		}
		if m.mod.NodeType() != m.typeS {
			t.Errorf("expected '%s', got '%s'", m.typeS, m.mod.NodeType())
		}
	}
}

// TestASTConcatenationNode tests concatenation (regex: ab)
func TestASTConcatenationNode(t *testing.T) {
	n := &Concatenation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	// fmt.Println(n.String())
	if n.ToRegex() != "ab" {
		t.Errorf("expected 'ab', got '%s'", n.ToRegex())
	}
	if n.NodeType() != "Concatenation" {
		t.Errorf("expected 'Concatenation', got '%s'", n.NodeType())
	}
}

// TestASTAlternationNode tests alternation (regex: a|b)
func TestASTAlternationNode(t *testing.T) {
	n := &Alternation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	// fmt.Println(n.String())
	// Alternation has no String() method, so we check structure
	if n.Left.ToRegex() != "a" || n.Right.ToRegex() != "b" {
		t.Errorf("expected left 'a' and right 'b', got '%s' and '%s'", n.Left.ToRegex(), n.Right.ToRegex())
	}
}

// TestASTLongAlternationNode tests long alternation, alternation should be left-associative (regex: a|b|c|d)
func TestASTLongAlternationNode(t *testing.T) {
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
	// fmt.Println(n.String())
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
	if n.ToRegex() != "a|b|c|d" {
		t.Errorf("expected left 'a|b|c|d' , got '%s'", n.ToRegex())
	}
}

// TestASTGroupNode tests grouping (regex: (ab))
func TestASTGroupNode(t *testing.T) {
	cat := &Concatenation{
		Left:  &Character{Value: 'a'},
		Right: &Character{Value: 'b'},
	}
	// fmt.Println(cat.String())
	g := &Group{Expr: cat}
	if g.ToRegex() != "(ab)" {
		t.Errorf("expected '(ab)', got '%s'", g.ToRegex())
	}
	if g.NodeType() != "Group" {
		t.Errorf("expected 'Group', got '%s'", g.NodeType())
	}
}

// TestASTExpressionNode tests the Expression wrapper (regex: a)
func TestASTExpressionNode(t *testing.T) {
	c := &Character{Value: 'a'}
	e := &Expression{Child: c}
	// fmt.Println(e.String())
	if e.ToRegex() != "a" {
		t.Errorf("expected 'a', got '%s'", e.ToRegex())
	}
	if e.NodeType() != "Expression" {
		t.Errorf("expected 'Expression', got '%s'", e.NodeType())
	}
}

// TestASTTermNode tests Term with and without modifier (regex: a, a*)
func TestASTTermNode(t *testing.T) {
	c := &Character{Value: 'a'}
	term := &Term{Child: c, HasModifier: false}
	// fmt.Println(term.String())
	if term.ToRegex() != "a" {
		t.Errorf("expected 'a', got '%s'", term.ToRegex())
	}
	if term.NodeType() != "Term" {
		t.Errorf("expected 'Term', got '%s'", term.NodeType())
	}
	m := &Modifier{Child: c, Type: ModifierStar}
	tm := &Term{Child: m, HasModifier: true}
	// fmt.Println(tm.String())
	if tm.ToRegex() != "a*" {
		t.Errorf("expected 'a*', got '%s'", tm.ToRegex())
	}
	if tm.NodeType() != "TermModified" {
		t.Errorf("expected 'TermModified', got '%s'", tm.NodeType())
	}
}

// TestASTNestedGroupsAndModifiers tests complex regex: ((a|b)*c+)?
func TestASTNestedGroupsAndModifiers(t *testing.T) {
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
	// fmt.Println(mod2.String())
	if mod2.ToRegex() != "((a|b)*c+)?" {
		t.Errorf("expected '((a|b)*c+)?', got '%s'", mod2.ToRegex())
	}
}

// TestASTEdgeCases tests empty group, deeply nested, and invalid modifier range
func TestASTEdgeCases(t *testing.T) {
	// Empty group: ()
	emptyGroup := &Group{Expr: &Concatenation{Left: &Character{Value: '\u0000'}, Right: &Character{Value: '\u0000'}}}
	// This is a hack, as there's no explicit Empty node
	if emptyGroup.ToRegex() != "(\u0000\u0000)" {
		t.Errorf("expected '(\u0000\u0000)', got '%s'", emptyGroup.ToRegex())
	}
	// fmt.Println(emptyGroup.String())
	// Deeply nested: (((a)))
	deep := &Group{Expr: &Group{Expr: &Group{Expr: &Character{Value: 'a'}}}}
	// fmt.Println(deep.String())
	if deep.ToRegex() != "(((a)))" {
		t.Errorf("expected '(((a)))', got '%s'", deep.ToRegex())
	}
}
