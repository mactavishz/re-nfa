package parse

import (
	"fmt"
	"reflect"
	"testing"
)

func astToString(n ASTNode) string {
	if n == nil {
		return "nil"
	}
	switch v := n.(type) {
	case *Expression:
		return fmt.Sprintf("Expression{Child: %s}", astToString(v.Child))
	case *Character:
		return fmt.Sprintf("Character{Value: %c}", v.Value)
	case *Alternation:
		return fmt.Sprintf("Alternation{Left: %s, Right: %s}", astToString(v.Left), astToString(v.Right))
	case *Concatenation:
		return fmt.Sprintf("Concatenation{Left: %s, Right: %s}", astToString(v.Left), astToString(v.Right))
	case *Group:
		return fmt.Sprintf("Group{Expr: %s}", astToString(v.Expr))
	case *Modifier:
		return fmt.Sprintf("Modifier{Child: %s, Type: %v}", astToString(v.Child), v.Type)
	default:
		return fmt.Sprintf("%T", n)
	}
}

func TestParserBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ASTNode
	}{
		{
			"empty",
			"",
			&Expression{Child: nil},
		},
		{
			"single char",
			"a",
			&Expression{Child: &Character{Value: 'a'}},
		},
		{
			"simple alternation",
			"a|b",
			&Expression{
				Child: &Alternation{
					Left:  &Character{Value: 'a'},
					Right: &Character{Value: 'b'},
				},
			},
		},
		{
			"simple concatenation",
			"ab",
			&Expression{
				Child: &Concatenation{
					Left:  &Character{Value: 'a'},
					Right: &Character{Value: 'b'},
				},
			},
		},
		{
			"simple group",
			"(a)",
			&Expression{
				Child: &Group{
					Expr: &Character{Value: 'a'},
				},
			},
		},
		{
			"simple modifier",
			"a*",
			&Expression{
				Child: &Modifier{
					Child: &Character{Value: 'a'},
					Type:  ModifierStar,
				},
			},
		},
		{
			"simple plus",
			"a+",
			&Expression{
				Child: &Modifier{
					Child: &Character{Value: 'a'},
					Type:  ModifierPlus,
				},
			},
		},
		{
			"simple optional",
			"a?",
			&Expression{
				Child: &Modifier{
					Child: &Character{Value: 'a'},
					Type:  ModifierQuestion,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast := parser.Parse()
			if !reflect.DeepEqual(ast, tt.expected) {
				t.Errorf("expected %s, got %s", astToString(tt.expected), astToString(ast))
			}
		})
	}
}

func TestParserComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ASTNode
	}{
		{
			"complex alternation",
			"a|b|c|d",
			&Expression{
				Child: &Alternation{
					Left: &Alternation{
						Left: &Alternation{
							Left:  &Character{Value: 'a'},
							Right: &Character{Value: 'b'},
						},
						Right: &Character{Value: 'c'},
					},
					Right: &Character{Value: 'd'},
				},
			},
		},
		{
			"complex concatenation",
			"abcd",
			&Expression{
				Child: &Concatenation{
					Left: &Concatenation{
						Left: &Concatenation{
							Left:  &Character{Value: 'a'},
							Right: &Character{Value: 'b'},
						},
						Right: &Character{Value: 'c'},
					},
					Right: &Character{Value: 'd'},
				},
			},
		},
		{
			"nested groups",
			"(a(b(c)))",
			&Expression{
				Child: &Group{
					Expr: &Concatenation{
						Left: &Character{Value: 'a'},
						Right: &Group{
							Expr: &Concatenation{
								Left:  &Character{Value: 'b'},
								Right: &Group{Expr: &Character{Value: 'c'}},
							},
						},
					},
				},
			},
		},
		{
			"complex modifiers",
			"a*b+c?",
			&Expression{
				Child: &Concatenation{
					Left: &Concatenation{
						Left: &Modifier{
							Child: &Character{Value: 'a'},
							Type:  ModifierStar,
						},
						Right: &Modifier{
							Child: &Character{Value: 'b'},
							Type:  ModifierPlus,
						},
					},
					Right: &Modifier{
						Child: &Character{Value: 'c'},
						Type:  ModifierQuestion,
					},
				},
			},
		},
		{
			"mixed operations",
			"a|b*|(c+)",
			&Expression{
				Child: &Alternation{
					Left: &Alternation{
						Left: &Character{Value: 'a'},
						Right: &Modifier{
							Child: &Character{Value: 'b'},
							Type:  ModifierStar,
						},
					},
					Right: &Group{
						Expr: &Modifier{
							Child: &Character{Value: 'c'},
							Type:  ModifierPlus,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast := parser.Parse()
			if !reflect.DeepEqual(ast, tt.expected) {
				t.Errorf("expected %s, got %s", astToString(tt.expected), astToString(ast))
			}
		})
	}
}

func TestParserEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ASTNode
	}{
		{
			"empty group",
			"()",
			&Expression{
				Child: &Group{Expr: nil},
			},
		},
		{
			"nested empty groups",
			"(()())",
			&Expression{
				Child: &Group{
					Expr: &Concatenation{
						Left:  &Group{Expr: nil},
						Right: &Group{Expr: nil},
					},
				},
			},
		},
		{
			"complex empty groups",
			"(()|())",
			&Expression{
				Child: &Group{
					Expr: &Alternation{
						Left:  &Group{Expr: nil},
						Right: &Group{Expr: nil},
					},
				},
			},
		},
		{
			"nested modifiers",
			"(a*)*",
			&Expression{
				Child: &Modifier{
					Child: &Group{
						Expr: &Modifier{
							Child: &Character{Value: 'a'},
							Type:  ModifierStar,
						},
					},
					Type: ModifierStar,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast := parser.Parse()
			if !reflect.DeepEqual(ast, tt.expected) {
				t.Errorf("expected %s, got %s", astToString(tt.expected), astToString(ast))
			}
		})
	}
}

func TestParserErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"unmatched left paren", "(", true},
		{"unmatched right paren", ")", true},
		{"invalid char after |", "a|", true},
		{"invalid char after modifier", "a*|", true},
		{"invalid use of modifier", "a**", true},
		{"invalid char in group", "(a|)", true},
		{"invalid char after group", "(a)|", true},
		{"invalid char after nested group", "((a))|", true},
		{"invalid char after complex group", "((a|b)*)|", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("unexpected panic: %v", r)
					}
				} else if tt.expectError {
					t.Error("expected panic but got none")
				}
			}()

			parser := NewParser(tt.input)
			parser.Parse()
		})
	}
}
