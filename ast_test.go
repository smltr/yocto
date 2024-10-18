package main

import (
	"reflect"
	"testing"
)

func TestParseAST(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected ASTNode
	}{
		{
			name: "Simple symbols and numbers",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "def"},
				{"PAREN", "("},
				{"SYMBOL", "add"},
				{"SYMBOL", "a"},
				{"SYMBOL", "b"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "+"},
				{"SYMBOL", "a"},
				{"SYMBOL", "b"},
				{"PAREN", ")"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					Symbol{Name: "def"},
					List{
						Elements: []ASTNode{
							Symbol{Name: "add"},
							Symbol{Name: "a"},
							Symbol{Name: "b"},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "+"},
							Symbol{Name: "a"},
							Symbol{Name: "b"},
						},
					},
				},
			},
		},
		{
			name: "String token",
			tokens: []Token{
				{"PAREN", "("},
				{"STRING", "hello world"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					StringLiteral{Value: "hello world"},
				},
			},
		},
		{
			name: "Negative and floating numbers",
			tokens: []Token{
				{"PAREN", "("},
				{"NUMBER", "-123"},
				{"NUMBER", "45.67"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					Number{Value: -123},
					Number{Value: 45.67},
				},
			},
		},
		{
			name: "Symbols with hyphens",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "my-variable-name"},
				{"NUMBER", "42"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					Symbol{Name: "my-variable-name"},
					Number{Value: 42},
				},
			},
		},
		{
			name: "Escaped characters in string",
			tokens: []Token{
				{"PAREN", "("},
				{"STRING", "Line1\nLine2\tTabbed"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					StringLiteral{Value: "Line1\nLine2\tTabbed"},
				},
			},
		},
		{
			name: "Multiple variable definitions",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "def"},
				{"SYMBOL", "x"},
				{"NUMBER", "10"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "def"},
				{"SYMBOL", "y"},
				{"NUMBER", "20"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					List{
						Elements: []ASTNode{
							Symbol{Name: "def"},
							Symbol{Name: "x"},
							Number{Value: 10},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "def"},
							Symbol{Name: "y"},
							Number{Value: 20},
						},
					},
				},
			},
		},
		{
			name: "Variable definition followed by an expression",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "def"},
				{"SYMBOL", "z"},
				{"NUMBER", "30"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "+"},
				{"SYMBOL", "x"},
				{"SYMBOL", "y"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					List{
						Elements: []ASTNode{
							Symbol{Name: "def"},
							Symbol{Name: "z"},
							Number{Value: 30},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "+"},
							Symbol{Name: "x"},
							Symbol{Name: "y"},
						},
					},
				},
			},
		},
		{
			name: "Multiple print statements",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "print"},
				{"STRING", "Hello"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "print"},
				{"STRING", "World"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					List{
						Elements: []ASTNode{
							Symbol{Name: "print"},
							StringLiteral{Value: "Hello"},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "print"},
							StringLiteral{Value: "World"},
						},
					},
				},
			},
		},
		{
			name: "Nested lists with multiple top-level expressions",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "defn"},
				{"SYMBOL", "multiply"},
				{"PAREN", "("},
				{"SYMBOL", "a"},
				{"SYMBOL", "b"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "*"},
				{"SYMBOL", "a"},
				{"SYMBOL", "b"},
				{"PAREN", ")"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "multiply"},
				{"NUMBER", "5"},
				{"NUMBER", "6"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					List{
						Elements: []ASTNode{
							Symbol{Name: "defn"},
							Symbol{Name: "multiply"},
							List{
								Elements: []ASTNode{
									Symbol{Name: "a"},
									Symbol{Name: "b"},
								},
							},
							List{
								Elements: []ASTNode{
									Symbol{Name: "*"},
									Symbol{Name: "a"},
									Symbol{Name: "b"},
								},
							},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "multiply"},
							Number{Value: 5},
							Number{Value: 6},
						},
					},
				},
			},
		},
		{
			name:   "Empty input",
			tokens: []Token{},
			expected: List{
				Elements: []ASTNode{},
			},
		},
		{
			name: "Single top-level expression",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "print"},
				{"STRING", "Solo Expression"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					Symbol{Name: "print"},
					StringLiteral{Value: "Solo Expression"},
				},
			},
		},
		{
			name: "Complex multiple expressions",
			tokens: []Token{
				{"PAREN", "("},
				{"SYMBOL", "defn"},
				{"SYMBOL", "square"},
				{"PAREN", "("},
				{"SYMBOL", "n"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "*"},
				{"SYMBOL", "n"},
				{"SYMBOL", "n"},
				{"PAREN", ")"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "print"},
				{"SYMBOL", "square"},
				{"NUMBER", "4"},
				{"PAREN", ")"},
				{"PAREN", "("},
				{"SYMBOL", "print"},
				{"NUMBER", "16"},
				{"PAREN", ")"},
			},
			expected: List{
				Elements: []ASTNode{
					List{
						Elements: []ASTNode{
							Symbol{Name: "defn"},
							Symbol{Name: "square"},
							List{
								Elements: []ASTNode{
									Symbol{Name: "n"},
								},
							},
							List{
								Elements: []ASTNode{
									Symbol{Name: "*"},
									Symbol{Name: "n"},
									Symbol{Name: "n"},
								},
							},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "print"},
							Symbol{Name: "square"},
							Number{Value: 4},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "print"},
							Number{Value: 16},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := ParseTokens(tt.tokens)
			if err != nil {
				t.Fatalf("ParseTokens returned an error: %v", err)
			}

			if !reflect.DeepEqual(ast, tt.expected) {
				t.Errorf("Expected AST %+v, got %+v", tt.expected, ast)
			}
		})
	}
}
