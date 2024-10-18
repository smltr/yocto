package main

import (
	"fmt"
	"strconv"
)

// ASTNode represents a node in the abstract syntax tree.
type ASTNode interface{}

// Symbol represents an identifier or variable name.
type Symbol struct {
	Name string
}

// Number represents a numerical literal.
type Number struct {
	Value float64
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Value string
}

// List represents a list of AST nodes (e.g., function calls, expressions).
type List struct {
	Elements []ASTNode
}

// ParseTokens converts a slice of Tokens into an ASTNode.
// Steps:
//  1. Initialize an empty stack to manage nested lists.
//  2. Create a root list to start building the AST.
//  3. Iterate over each token in the tokens slice:
//     a. If token is a "PAREN" with value "(", push the current list onto the stack and start a new list.
//     b. If token is a "PAREN" with value ")", finalize the current list and append it to the previous list from the stack.
//     c. If token is a "SYMBOL", "NUMBER", or "STRING", create the corresponding ASTNode and append it to the current list.
//     d. Handle errors such as unexpected closing parentheses or invalid tokens.
//  4. After processing all tokens, ensure that the stack is empty (all opened parentheses are closed).
//  5. Return the root ASTNode representing the entire parsed structure.
func ParseTokens(tokens []Token) (ASTNode, error) {
	// Step 1: Initialize an empty stack to manage nested lists.
	stack := []List{}

	// Step 2: Create a root list to start building the AST.
	currentList := List{Elements: []ASTNode{}}

	// Iterate over each token in the tokens slice.
	for i, token := range tokens {
		switch token.Type {
		case "PAREN":
			if token.Value == "(" {
				// Step 3a: Start a new list.
				// Push the current list onto the stack.
				stack = append(stack, currentList)
				// Initialize a new current list.
				currentList = List{Elements: []ASTNode{}}
			} else if token.Value == ")" {
				// Step 3b: Finalize the current list.
				if len(stack) == 0 {
					// Error: Unexpected closing parenthesis.
					return nil, fmt.Errorf("unexpected closing parenthesis at token index %d", i)
				}
				// Pop the last list from the stack.
				parentList := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				// Append the current list to the parent list.
				parentList.Elements = append(parentList.Elements, currentList)
				// Set the current list to the parent list.
				currentList = parentList
			}
		case "SYMBOL":
			// Step 3c: Create a Symbol ASTNode and append to current list.
			symbol := Symbol{Name: token.Value}
			currentList.Elements = append(currentList.Elements, symbol)
		case "NUMBER":
			// Step 3c: Convert the string to float64 and create a Number ASTNode.
			numValue, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number '%s' at token index %d: %v", token.Value, i, err)
			}
			number := Number{Value: numValue}
			currentList.Elements = append(currentList.Elements, number)
		case "STRING":
			// Step 3c: Create a StringLiteral ASTNode and append to current list.
			stringLit := StringLiteral{Value: token.Value}
			currentList.Elements = append(currentList.Elements, stringLit)
		default:
			// Step 3d: Handle unknown token types.
			return nil, fmt.Errorf("unknown token type '%s' at token index %d", token.Type, i)
		}
	}

	// Step 4: Ensure that all opened parentheses are closed.
	if len(stack) != 0 {
		return nil, fmt.Errorf("unexpected end of tokens: missing closing parentheses")
	}

	// Step 5: Return the root ASTNode.
	// If the entire input is a single list, return it directly.
	if len(currentList.Elements) == 1 {
		return currentList.Elements[0], nil
	}

	// Otherwise, return the root list containing all top-level elements.
	return currentList, nil
}
