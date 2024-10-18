package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Token represents a token with its type and value.
type Token struct {
	Type  string
	Value string
}

// tokenize converts the source string into a slice of tokens.
// Workflow:
//  1. Initialize an empty token list.
//  2. Iterate through each character in the source string.
//     a. Skip whitespace.
//     b. Identify and tokenize parentheses.
//     c. Tokenize strings enclosed in quotes.
//     d. Tokenize numbers as single tokens.
//     e. Tokenize symbols/identifiers.
//  3. Return the list of tokens.
func ParseString(str string) []Token {
	tokens := []Token{}

	// Step 1: Initialize
	input := str
	length := len(input)
	position := 0

	// Step 2: Iterate through each character
	for position < length {
		char := input[position]

		// Step 2a: Skip whitespace
		if isWhitespace(char) {
			position++
			continue
		}

		// Step 2b: Identify and tokenize parentheses
		if isParenthesis(char) {
			tokenType := "PAREN"
			tokenValue := string(char)
			tokens = append(tokens, Token{tokenType, tokenValue})
			position++
			continue
		}

		// Step 2c: Tokenize strings enclosed in quotes
		if char == '"' {
			strValue, newPos, err := parseString(input, position)
			if err != nil {
				// Handle error (e.g., unterminated string)
				panic(err)
			}
			tokens = append(tokens, Token{"STRING", strValue})
			position = newPos
			continue
		}

		// Step 2d: Tokenize numbers
		if isDigit(char) || (char == '-' && position+1 < length && isDigit(input[position+1])) {
			numValue, newPos := parseNumber(input, position)
			tokens = append(tokens, Token{"NUMBER", numValue})
			position = newPos
			continue
		}

		// Step 2e: Tokenize symbols/identifiers
		symbolValue, newPos := parseSymbol(input, position)
		tokens = append(tokens, Token{"SYMBOL", symbolValue})
		position = newPos
	}

	// Step 3: Return tokens
	return tokens
}

// Helper functions

func isWhitespace(char byte) bool {
	return unicode.IsSpace(rune(char))
}

func isParenthesis(char byte) bool {
	return char == '(' || char == ')'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func parseString(input string, start int) (string, int, error) {
	// Parse a string token starting at the given position
	end := start + 1
	var strBuilder strings.Builder
	for end < len(input) {
		if input[end] == '"' {
			return strBuilder.String(), end + 1, nil
		}
		// Handle escape characters if needed
		if input[end] == '\\' && end+1 < len(input) {
			end++
			switch input[end] {
			case 'n':
				strBuilder.WriteByte('\n')
			case 't':
				strBuilder.WriteByte('\t')
			case '"':
				strBuilder.WriteByte('"')
			default:
				strBuilder.WriteByte(input[end])
			}
		} else {
			strBuilder.WriteByte(input[end])
		}
		end++
	}
	return "", end, fmt.Errorf("unterminated string")
}

func parseNumber(input string, start int) (string, int) {
	end := start
	length := len(input)
	// Handle negative numbers
	if input[end] == '-' {
		end++
	}
	hasDot := false
	for end < length && (isDigit(input[end]) || (!hasDot && input[end] == '.')) {
		if input[end] == '.' {
			hasDot = true
		}
		end++
	}
	return input[start:end], end
}

func parseSymbol(input string, start int) (string, int) {
	end := start
	length := len(input)
	for end < length && !isWhitespace(input[end]) && !isParenthesis(input[end]) && input[end] != '"' {
		end++
	}
	return input[start:end], end
}
