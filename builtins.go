package main

import (
	"fmt"
	"math"
)

// core ----------------------------------------------------------------------------------

func evalDef(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def requires 2 arguments")
	}
	symbol, ok := args[0].(Name)
	if !ok {
		return nil, fmt.Errorf("first argument to def must be a symbol")
	}
	value, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}
	env.Set(symbol, value)
	return value, nil
}

func evalPrint(args []Expression, env *Environment) (Expression, error) {
	for _, arg := range args {
		value, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		fmt.Print(value)
	}
	fmt.Println()
	return nil, nil
}

func evalDo(args []Expression, env *Environment) (Expression, error) {
	var result Expression
	var err error
	for _, arg := range args {
		result, err = arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalQuote(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quote requires exactly one argument")
	}
	return args[0], nil
}

func evalQuasiquote(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quasiquote requires exactly one argument")
	}
	return quasiquoteExpand(args[0], env, 0)
}

func evalUnquote(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("unquote requires exactly one argument")
	}
	return args[0].Evaluate(env)
}

func evalDefMacro(args []Expression, env *Environment) (Expression, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("defmacro requires at least 2 arguments")
	}

	signature, ok := args[0].(List)
	if !ok || len(signature) < 1 {
		return nil, fmt.Errorf("first argument to defmacro must be a list containing at least the macro name")
	}

	name, ok := signature[0].(Name)
	if !ok {
		return nil, fmt.Errorf("macro name must be a symbol")
	}

	params := make(List, 0)
	var restParam Name
	for i, param := range signature[1:] {
		if paramName, ok := param.(Name); ok && string(paramName) == "&" {
			if i+1 < len(signature[1:]) {
				restParam = signature[i+2].(Name)
				break
			} else {
				return nil, fmt.Errorf("& must be followed by a parameter name")
			}
		}
		params = append(params, param)
	}

	body := args[1:]

	macro := Macro{params: params, restParam: restParam, body: body}
	env.Set(name, macro)

	return macro, nil
}

func evalEval(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("eval expects exactly one argument")
	}

	// Unwrap the quote if present
	var exprToEval Expression
	if list, ok := args[0].(List); ok && len(list) > 0 {
		if name, ok := list[0].(Name); ok && (string(name) == "quote" || string(name) == "quasiquote") {
			if len(list) != 2 {
				return nil, fmt.Errorf("quote expects exactly one argument")
			}
			exprToEval = list[1]
		} else {
			exprToEval = args[0]
		}
	} else {
		exprToEval = args[0]
	}

	// Evaluate the unwrapped expression
	return exprToEval.Evaluate(env)
}

func evalLambda(args []Expression, env *Environment) (Expression, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("lambda requires at least 2 arguments")
	}
	signature, ok := args[0].(List)
	if !ok || len(signature) == 0 {
		return nil, fmt.Errorf("first argument to lambda must be a non-empty list")
	}
	params := signature[1:]
	body := args[1:]
	return &Function{params: params, body: body, env: env}, nil
}

func evalDefn(args []Expression, env *Environment) (Expression, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("defn requires at least 2 arguments")
	}
	signature, ok := args[0].(List)
	if !ok || len(signature) < 2 {
		return nil, fmt.Errorf("first argument to defn must be a list containing at least the function name and parameters")
	}
	name, ok := signature[0].(Name)
	if !ok {
		return nil, fmt.Errorf("first element of the signature list must be a symbol")
	}
	params := signature[1:]
	body := args[1:]
	fn := &Function{params: params, body: body, env: env}
	env.Set(name, fn)
	return fn, nil
}

func evalIf(args []Expression, env *Environment) (Expression, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("if requires 2 or 3 arguments")
	}
	condition, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}
	if condition != nil && condition != Boolean(false) {
		return args[1].Evaluate(env)
	} else if len(args) == 3 {
		return args[2].Evaluate(env)
	}
	return nil, nil
}

// logic --------------------------------------------------------------------------------

func evalAnd(args []Expression, env *Environment) (Expression, error) {
	if len(args) == 0 {
		return Boolean(true), nil
	}

	var result Expression
	for _, arg := range args {
		evaluated, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if evaluated == nil || evaluated == Boolean(false) {
			return Boolean(false), nil
		}
		result = evaluated
	}
	return result, nil
}

func evalOr(args []Expression, env *Environment) (Expression, error) {
	if len(args) == 0 {
		return Boolean(false), nil // Empty 'or' is false
	}

	for _, arg := range args {
		evaluated, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if evaluated != nil && evaluated != Boolean(false) {
			return evaluated, nil // Short-circuit: return first truthy value
		}
	}
	return Boolean(false), nil
}

func evalNot(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("not requires exactly one argument")
	}

	value, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	return Boolean(value == nil || value == Boolean(false)), nil
}

// math ----------------------------------------------------------------------------------

func evalAdd(args []Expression, env *Environment) (Expression, error) {
	var result float64
	for _, arg := range args {
		value, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		num, ok := value.(Number)
		if !ok {
			return nil, fmt.Errorf("+ expects numbers, got %T", value)
		}
		result += float64(num)
	}
	return Number(result), nil
}

func evalSubtract(args []Expression, env *Environment) (Expression, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("- requires at least one argument")
	}

	first, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	firstNum, ok := first.(Number)
	if !ok {
		return nil, fmt.Errorf("- expects numbers, got %T", first)
	}

	if len(args) == 1 {
		return Number(-float64(firstNum)), nil
	}

	result := float64(firstNum)
	for _, arg := range args[1:] {
		value, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		num, ok := value.(Number)
		if !ok {
			return nil, fmt.Errorf("- expects numbers, got %T", value)
		}
		result -= float64(num)
	}
	return Number(result), nil
}

func evalMultiply(args []Expression, env *Environment) (Expression, error) {
	result := 1.0
	for _, arg := range args {
		value, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		num, ok := value.(Number)
		if !ok {
			return nil, fmt.Errorf("* expects numbers, got %T", value)
		}
		result *= float64(num)
	}
	return Number(result), nil
}

func evalDivide(args []Expression, env *Environment) (Expression, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("/ requires at least one argument")
	}

	first, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	firstNum, ok := first.(Number)
	if !ok {
		return nil, fmt.Errorf("/ expects numbers, got %T", first)
	}

	if len(args) == 1 {
		return Number(1 / float64(firstNum)), nil
	}

	result := float64(firstNum)
	for _, arg := range args[1:] {
		value, err := arg.Evaluate(env)
		if err != nil {
			return nil, err
		}
		num, ok := value.(Number)
		if !ok {
			return nil, fmt.Errorf("/ expects numbers, got %T", value)
		}
		if num == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result /= float64(num)
	}
	return Number(result), nil
}

func evalPower(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("** requires exactly two arguments")
	}

	base, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	exponent, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	baseNum, ok := base.(Number)
	if !ok {
		return nil, fmt.Errorf("** expects numbers, got %T for base", base)
	}

	exponentNum, ok := exponent.(Number)
	if !ok {
		return nil, fmt.Errorf("** expects numbers, got %T for exponent", exponent)
	}

	result := math.Pow(float64(baseNum), float64(exponentNum))
	return Number(result), nil
}

// compare -------------------------------------------------------------------------------

func evalEqual(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("= requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	return Boolean(left == right), nil
}

func evalNotEqual(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("!= requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	return Boolean(left != right), nil
}

func evalLessThan(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("< requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	leftNum, leftOk := left.(Number)
	rightNum, rightOk := right.(Number)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("< expects numbers, got %T and %T", left, right)
	}

	return Boolean(leftNum < rightNum), nil
}

func evalLessThanOrEqual(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("<= requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	leftNum, leftOk := left.(Number)
	rightNum, rightOk := right.(Number)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("<= expects numbers, got %T and %T", left, right)
	}

	return Boolean(leftNum <= rightNum), nil
}

func evalGreaterThan(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("> requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	leftNum, leftOk := left.(Number)
	rightNum, rightOk := right.(Number)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf("> expects numbers, got %T and %T", left, right)
	}

	return Boolean(leftNum > rightNum), nil
}

func evalGreaterThanOrEqual(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(">= requires exactly two arguments")
	}

	left, err := args[0].Evaluate(env)
	if err != nil {
		return nil, err
	}

	right, err := args[1].Evaluate(env)
	if err != nil {
		return nil, err
	}

	leftNum, leftOk := left.(Number)
	rightNum, rightOk := right.(Number)

	if !leftOk || !rightOk {
		return nil, fmt.Errorf(">= expects numbers, got %T and %T", left, right)
	}

	return Boolean(leftNum >= rightNum), nil
}
