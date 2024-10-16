package main

import (
	"testing"
)

func TestYoctoLisp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Define and call function",
			input:    "(defn (add a b) (+ a b)) (add 1 2)",
			expected: "3",
		},
		{
			name:     "Define variable",
			input:    "(def one 1) one",
			expected: "1",
		},
		{
			name:     "Simple arithmetic",
			input:    "(+ 2 3)",
			expected: "5",
		},
		{
			name:     "Define and use a macro with quasiquote",
			input:    "(defmacro (add a b) `(+ ,a ,b)) (add 1 1)",
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EvalString(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	_, err := EvalString("(+ 1 2")
	if err == nil {
		t.Error("Expected parse error, got nil")
	}
}

func TestEvalError(t *testing.T) {
	_, err := EvalString("(/ 1 0)")
	if err == nil {
		t.Error("Expected evaluation error (division by zero), got nil")
	}
}
