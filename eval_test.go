package main

import (
	"reflect"
	"testing"
)

// TestEvaluateAST tests the EvaluateAST function to ensure it correctly
// evaluates AST nodes representing mathematical expressions.
func TestEvaluateAST(t *testing.T) {
	tests := []struct {
		name     string
		ast      ASTNode
		expected Expression
	}{
		{
			name: "Simple addition",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "+"},
					Number{Value: 1},
					Number{Value: 2},
				},
			},
			expected: Number{Value: 3},
		},
		{
			name: "Simple subtraction",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "-"},
					Number{Value: 5},
					Number{Value: 3},
				},
			},
			expected: Number{Value: 2},
		},
		{
			name: "Simple multiplication",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "*"},
					Number{Value: 4},
					Number{Value: 3},
				},
			},
			expected: Number{Value: 12},
		},
		{
			name: "Simple division",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "/"},
					Number{Value: 10},
					Number{Value: 2},
				},
			},
			expected: Number{Value: 5},
		},
		{
			name: "Exponentiation",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "**"},
					Number{Value: 2},
					Number{Value: 3},
				},
			},
			expected: Number{Value: 8},
		},
		{
			name: "Modulo operation",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "%"},
					Number{Value: 10},
					Number{Value: 3},
				},
			},
			expected: Number{Value: 1},
		},
		{
			name: "Nested arithmetic expressions",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "+"},
					List{
						Elements: []ASTNode{
							Symbol{Name: "*"},
							Number{Value: 2},
							Number{Value: 3},
						},
					},
					List{
						Elements: []ASTNode{
							Symbol{Name: "/"},
							Number{Value: 10},
							Number{Value: 5},
						},
					},
				},
			},
			expected: Number{Value: 8}, // (2*3) + (10/5) = 6 + 2 = 8
		},
		{
			name: "Multiple operations",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "-"},
					Number{Value: 20},
					Number{Value: 5},
					Number{Value: 3},
				},
			},
			expected: Number{Value: 12}, // 20 - 5 - 3 = 12
		},
		{
			name: "Division by zero",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "/"},
					Number{Value: 10},
					Number{Value: 0},
				},
			},
			expected: nil, // Expect an error
		},
		{
			name: "Invalid number format",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "+"},
					Symbol{Name: "a"}, // 'a' is not a number
					Number{Value: 2},
				},
			},
			expected: nil, // Expect an error
		},
		{
			name: "Unary negation",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "-"},
					Number{Value: 5},
				},
			},
			expected: Number{Value: -5}, // Assuming unary minus is supported
		},
		{
			name: "Zero arguments to addition",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "+"},
				},
			},
			expected: Number{Value: 0}, // Assuming + with no args returns 0
		},
		{
			name: "Single argument to addition",
			ast: List{
				Elements: []ASTNode{
					Symbol{Name: "+"},
					Number{Value: 5},
				},
			},
			expected: Number{Value: 5}, // Assuming + with one arg returns the arg
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the environment with built-in functions
			env := NewEnvironment()
			RegisterBuiltins(env)

			// Evaluate the AST
			result, err := EvaluateAST(tt.ast, env)

			if tt.expected == nil {
				// Expecting an error
				if err == nil {
					t.Errorf("Expected an error, but got result: %+v", result)
				}
				return
			}

			// No error expected
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare the expected and actual results
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
