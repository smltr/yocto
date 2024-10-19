package main

import (
	"bufio"
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
	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Println("Usage: yocto [filename.yoc]")
		os.Exit(1)
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	env := NewEnvironment(nil)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		tokens := tokenize(input)
		var result Expression
		var err error
		for len(tokens) > 0 {
			var expr Expression
			expr, tokens, err = parseExpr(tokens)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				break
			}
			result, err = expr.Evaluate(env)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				break
			}
		}
		if err == nil {
			fmt.Println(result)
		}
	}
}

func runFile(filename string) {
	if !strings.HasSuffix(filename, ".yoc") {
		fmt.Println("Error: File must have .yoc extension")
		os.Exit(1)
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	_, err = EvalString(string(contents))
	if err != nil {
		fmt.Printf("Error evaluating file: %v\n", err)
		os.Exit(1)
	}
}
