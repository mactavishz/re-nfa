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
		wantErr  bool
	}{
		{
			"empty",
			"",
			&Expression{Child: nil},
			false,
		},
		{
			"single char",
			"a",
			&Expression{Child: &Character{Value: 'a'}},
			false,
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
			false,
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
			false,
		},
		{
			"simple group",
			"(a)",
			&Expression{
				Child: &Group{
					Expr: &Character{Value: 'a'},
				},
			},
			false,
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
			false,
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
			false,
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
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		wantErr  bool
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
			false,
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
			false,
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
			false,
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
			false,
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
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		wantErr  bool
	}{
		{
			"empty group",
			"()",
			&Expression{
				Child: &Group{Expr: nil},
			},
			false,
		},
		{
			"term after empty group",
			"()|a",
			&Expression{
				Child: &Alternation{
					Left: &Group{Expr: nil},
					Right: &Character{
						Value: 'a',
					},
				},
			},
			false,
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
			false,
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
			false,
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
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			ast, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		{"empty left term", "|a", true},
		{"unmatched left paren", "(", true},
		{"unmatched right paren", ")", true},
		{"invalid char after |", "a|", true},
		{"missing char between a||b", "a||b", true},
		{"invalid char after modifier", "a*|", true},
		{"invalid use of modifier", "a**", true},
		{"invalid char in group", "(a|)", true},
		{"invalid char after group", "(a)|", true},
		{"invalid char after nested group", "((a))|", true},
		{"invalid char after complex group", "((a|b)*)|", true},
		{"invalid char after empty group", "()|", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()
			fmt.Println(err)
			if (err != nil) != tt.expectError {
				t.Errorf("Parse() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
