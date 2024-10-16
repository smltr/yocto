package main

import (
	"fmt"
	"os"
	"strings"
)

type Expression interface {
	Evaluate(env *Environment) (Expression, error)
}

func EvalString(input string) (string, error) {
	env := NewEnvironment(nil)
	tokens := tokenize(input)
	var result Expression
	var err error

	for len(tokens) > 0 {
		var expr Expression
		expr, tokens, err = parseExpr(tokens)
		if err != nil {
			return "", err
		}
		result, err = expr.Evaluate(env)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%v", result), nil
}

func main() {
	// Example usage: go run main.go example.yoc
	if len(os.Args) != 2 {
		fmt.Println("Usage: yocto <filename.yoc>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".yoc") {
		fmt.Println("Error: File must have .yoc extension")
		os.Exit(1)
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Add built-in functions to the environment
	_, err = EvalString(string(contents))
	if err != nil {
		fmt.Printf("Error evaluating file: %v\n", err)
		os.Exit(1)
	}

}
