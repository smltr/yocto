package main

import (
	"fmt"
)

// Expression represents the result of evaluating an AST node.
type Expression interface{}

// BuiltinFunction represents a built-in function.
type BuiltinFunction func(args []Expression) (Expression, error)

// Environment holds variable and function bindings.
type Environment struct {
	vars map[string]Expression
}

// NewEnvironment creates a new environment with built-in functions registered.
func NewEnvironment() *Environment {
	env := &Environment{
		vars: make(map[string]Expression),
	}
	RegisterBuiltins(env)
	return env
}

// RegisterBuiltins registers built-in functions to the environment.
func RegisterBuiltins(env *Environment) {
	env.vars["+"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		var sum float64
		for _, arg := range args {
			num, ok := arg.(Number)
			if !ok {
				return nil, fmt.Errorf("unsupported operand type for +: %T", arg)
			}
			sum += num.Value
		}
		return Number{Value: sum}, nil
	})

	env.vars["-"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("- expects at least one argument")
		}
		first, ok := args[0].(Number)
		if !ok {
			return nil, fmt.Errorf("unsupported operand type for -: %T", args[0])
		}
		if len(args) == 1 {
			return Number{Value: -first.Value}, nil
		}
		result := first.Value
		for _, arg := range args[1:] {
			num, ok := arg.(Number)
			if !ok {
				return nil, fmt.Errorf("unsupported operand type for -: %T", arg)
			}
			result -= num.Value
		}
		return Number{Value: result}, nil
	})

	env.vars["*"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		product := 1.0
		for _, arg := range args {
			num, ok := arg.(Number)
			if !ok {
				return nil, fmt.Errorf("unsupported operand type for *: %T", arg)
			}
			product *= num.Value
		}
		return Number{Value: product}, nil
	})

	env.vars["/"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("/ expects at least one argument")
		}
		first, ok := args[0].(Number)
		if !ok {
			return nil, fmt.Errorf("unsupported operand type for /: %T", args[0])
		}
		if len(args) == 1 {
			if first.Value == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return Number{Value: 1 / first.Value}, nil
		}
		result := first.Value
		for _, arg := range args[1:] {
			num, ok := arg.(Number)
			if !ok {
				return nil, fmt.Errorf("unsupported operand type for /: %T", arg)
			}
			if num.Value == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result /= num.Value
		}
		return Number{Value: result}, nil
	})

	env.vars["**"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("** expects exactly two arguments")
		}
		base, ok1 := args[0].(Number)
		exponent, ok2 := args[1].(Number)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("** expects number arguments")
		}
		result := pow(base.Value, exponent.Value)
		return Number{Value: result}, nil
	})

	env.vars["%"] = BuiltinFunction(func(args []Expression) (Expression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("%% expects exactly two arguments")
		}
		a, ok1 := args[0].(Number)
		b, ok2 := args[1].(Number)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("%% expects number arguments")
		}
		if b.Value == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		result := float64(int(a.Value) % int(b.Value))
		return Number{Value: result}, nil
	})
}

// pow is a helper function to perform exponentiation.
func pow(a, b float64) float64 {
	result := 1.0
	for i := 0; i < int(b); i++ {
		result *= a
	}
	return result
}

// Get retrieves a variable or function from the environment.
func (env *Environment) Get(name string) (Expression, bool) {
	value, exists := env.vars[name]
	return value, exists
}
