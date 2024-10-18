package main

import (
	"os"
)

func main() {
	sourceString := readFile()
	tokens := ParseString(sourceString)
	_, _ = ParseTokens(tokens)
	// EvaluateAST(tree)
}

func readFile() string {
	if len(os.Args) < 2 {
		panic("No file name provided")
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	return string(data)
}
