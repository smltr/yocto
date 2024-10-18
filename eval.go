package main

import (
	"fmt"
)

// EvaluateAST recursively evaluates an AST node within the given environment.
func EvaluateAST(node ASTNode, env *Environment) (Expression, error) {
	switch expr := node.(type) {
	case Number, StringLiteral:
		// Base cases: Numbers and Strings evaluate to themselves.
		return expr, nil
	case Symbol:
		// Look up the symbol in the environment.
		value, exists := env.Get(expr.Name)
		if !exists {
			return nil, fmt.Errorf("undefined symbol: %s", expr.Name)
		}
		return value, nil
	case List:
		if len(expr.Elements) == 0 {
			return nil, fmt.Errorf("empty expression")
		}

		// The first element should be a symbol representing the function/operator.
		first := expr.Elements[0]
		fnSymbol, ok := first.(Symbol)
		if !ok {
			return nil, fmt.Errorf("first element in list must be a symbol, got: %T", first)
		}

		// Retrieve the function from the environment.
		fnExpr, exists := env.Get(fnSymbol.Name)
		if !exists {
			return nil, fmt.Errorf("undefined function: %s", fnSymbol.Name)
		}

		// Assert that it's a BuiltinFunction.
		fn, ok := fnExpr.(BuiltinFunction)
		if !ok {
			return nil, fmt.Errorf("symbol %s is not a function", fnSymbol.Name)
		}

		// Evaluate all arguments.
		args := []Expression{}
		for _, argNode := range expr.Elements[1:] {
			evaluatedArg, err := EvaluateAST(argNode, env)
			if err != nil {
				return nil, err
			}
			args = append(args, evaluatedArg)
		}

		// Execute the function with evaluated arguments.
		return fn(args)
	default:
		return nil, fmt.Errorf("unknown AST node type: %T", expr)
	}
}
