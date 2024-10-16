package main

import "fmt"

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

func evalDefn(args []Expression, env *Environment) (Expression, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("defn requires 2 arguments")
	}

	signature, ok := args[0].(List)
	if !ok || len(signature) < 1 {
		return nil, fmt.Errorf("first argument to defn must be a list containing at least the function name")
	}

	name, ok := signature[0].(Name)
	if !ok {
		return nil, fmt.Errorf("function name must be a symbol")
	}

	params := signature[1:]
	body := args[1:]

	fn := Function(func(args []Expression, env *Environment) (Expression, error) {
		if len(args) != len(params) {
			return nil, fmt.Errorf("wrong number of arguments: expected %d, got %d", len(params), len(args))
		}
		newEnv := NewEnvironment(env)
		for i, param := range params {
			symbol, ok := param.(Name)
			if !ok {
				return nil, fmt.Errorf("parameter must be a symbol")
			}
			newEnv.Set(symbol, args[i])
		}
		var result Expression
		var err error
		for _, expr := range body {
			result, err = expr.Evaluate(newEnv)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	})

	env.Set(name, fn)
	return fn, nil
}

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
		if paramName, ok := param.(Name); ok && string(paramName) == "..." {
			if i+1 < len(signature[1:]) {
				restParam = signature[i+2].(Name)
				break
			} else {
				return nil, fmt.Errorf("... must be followed by a parameter name")
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
