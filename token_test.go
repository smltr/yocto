package main

import "testing"

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Simple symbols and numbers",
			input: "(def (add a b) (+ a b))",
			expected: []Token{
				{Type: "PAREN", Value: "("},
				{Type: "SYMBOL", Value: "def"},
				{Type: "PAREN", Value: "("},
				{Type: "SYMBOL", Value: "add"},
				{Type: "SYMBOL", Value: "a"},
				{Type: "SYMBOL", Value: "b"},
				{Type: "PAREN", Value: ")"},
				{Type: "PAREN", Value: "("},
				{Type: "SYMBOL", Value: "+"},
				{Type: "SYMBOL", Value: "a"},
				{Type: "SYMBOL", Value: "b"},
				{Type: "PAREN", Value: ")"},
				{Type: "PAREN", Value: ")"},
			},
		},
		{
			name:  "String token",
			input: `("hello world")`,
			expected: []Token{
				{Type: "PAREN", Value: "("},
				{Type: "STRING", Value: "hello world"},
				{Type: "PAREN", Value: ")"},
			},
		},
		{
			name:  "Negative and floating numbers",
			input: "(-123 45.67)",
			expected: []Token{
				{Type: "PAREN", Value: "("},
				{Type: "NUMBER", Value: "-123"},
				{Type: "NUMBER", Value: "45.67"},
				{Type: "PAREN", Value: ")"},
			},
		},
		{
			name:  "Symbols with hyphens",
			input: "(my-variable-name 42)",
			expected: []Token{
				{Type: "PAREN", Value: "("},
				{Type: "SYMBOL", Value: "my-variable-name"},
				{Type: "NUMBER", Value: "42"},
				{Type: "PAREN", Value: ")"},
			},
		},
		{
			name:  "Escaped characters in string",
			input: `("Line1\nLine2\tTabbed")`,
			expected: []Token{
				{Type: "PAREN", Value: "("},
				{Type: "STRING", Value: "Line1\nLine2\tTabbed"},
				{Type: "PAREN", Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseString(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(result))
				return
			}
			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("Token %d mismatch: expected %+v, got %+v", i, tt.expected[i], token)
				}
			}
		})
	}
}
