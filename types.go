package main

import "fmt"

type Name string
type Number float64
type List []Expression
type String string
type Function func(args []Expression, env *Environment) (Expression, error)
type Boolean bool

func (b Boolean) Evaluate(env *Environment) (Expression, error) {
	return b, nil
}

func (f Function) Evaluate(env *Environment) (Expression, error) {
	return f, nil
}

func (n Name) Evaluate(env *Environment) (Expression, error) {
	value, ok := env.Get(n)
	if !ok {
		return nil, fmt.Errorf("undefined name: %s", n)
	}
	return value, nil
}

func (n Number) Evaluate(env *Environment) (Expression, error) {
	return n, nil
}

func (s String) Evaluate(env *Environment) (Expression, error) {
	return s, nil
}

func (l List) Evaluate(env *Environment) (Expression, error) {
	if len(l) == 0 {
		return nil, nil
	}

	// Macro expansion
	expanded, didExpand, err := MacroExpand(l, env)
	if err != nil {
		return nil, err
	}
	if didExpand {
		// If it was a macro, evaluate the expanded form
		return evalEval([]Expression{expanded}, env)
	}

	// If not a macro, proceed with normal evaluation
	switch first := l[0].(type) {
	case Name:
		switch string(first) {
		case "def":
			return evalDef(l[1:], env)
		case "defn":
			return evalDefn(l[1:], env)
		case "+":
			return evalAdd(l[1:], env)
		case "print":
			return evalPrint(l[1:], env)
		case "quote":
			return evalQuote(l[1:], env)
		case "quasiquote":
			return evalQuasiquote(l[1:], env)
		case "unquote":
			return evalUnquote(l[1:], env)
		case "defmacro":
			return evalDefMacro(l[1:], env)
		case "do":
			return evalDo(l[1:], env)
		case "and":
			return evalAnd(l[1:], env)
		case "or":
			return evalOr(l[1:], env)
		case "true":
			return Boolean(true), nil
		case "false":
			return Boolean(false), nil
		case "eval": // Add this line
			return evalEval(l[1:], env)
		}
	}

	// Function call
	fn, err := l[0].Evaluate(env)
	if err != nil {
		return nil, err
	}
	if f, ok := fn.(Function); ok {
		args := make([]Expression, len(l)-1)
		for i, arg := range l[1:] {
			evaluated, err := arg.Evaluate(env)
			if err != nil {
				return nil, err
			}
			args[i] = evaluated
		}
		return f(args, env)
	}
	return nil, fmt.Errorf("not a function: %v", l[0])
}

// macros ---

type Macro struct {
	params    List
	restParam Name
	body      List
}

func (m Macro) Evaluate(env *Environment) (Expression, error) {
	return m, nil
}

func quasiquoteExpand(expr Expression, env *Environment, depth int) (Expression, error) {
	switch e := expr.(type) {
	case List:
		if len(e) > 0 {
			if name, ok := e[0].(Name); ok {
				switch string(name) {
				case "unquote":
					if depth == 0 {
						if len(e) != 2 {
							return nil, fmt.Errorf("unquote requires exactly one argument")
						}
						return e[1].Evaluate(env)
					}
				case "unquote-splicing":
					if depth == 0 {
						if len(e) != 2 {
							return nil, fmt.Errorf("unquote-splicing requires exactly one argument")
						}
						spliced, err := e[1].Evaluate(env)
						if err != nil {
							return nil, err
						}
						if splicedList, ok := spliced.(List); ok {
							return splicedList, nil
						}
						return nil, fmt.Errorf("unquote-splicing argument must evaluate to a list")
					}
				}
			}
		}
		result := make(List, 0, len(e))
		for _, subExpr := range e {
			expanded, err := quasiquoteExpand(subExpr, env, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, expanded)
		}
		return result, nil
	case Name:
		if depth == 0 {
			value, ok := env.Get(e)
			if ok {
				return value, nil
			}
		}
		return e, nil
	default:
		return expr, nil
	}
}

func ExpandMacro(macro Macro, args []Expression, env *Environment) (Expression, error) {
	macroEnv := NewEnvironment(env)
	for i, param := range macro.params {
		macroEnv.Set(param.(Name), args[i])
	}

	if macro.restParam != "" {
		restArgs := args[len(macro.params):]
		macroEnv.Set(macro.restParam, List(restArgs))
	}

	var result Expression
	var err error
	for _, expr := range macro.body {
		result, err = quasiquoteExpand(expr, macroEnv, 0)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func IsMacro(expr Expression) bool {
	_, ok := expr.(Macro)
	return ok
}

func MacroExpand(expr Expression, env *Environment) (Expression, bool, error) {
	didExpand := false

	for {
		if list, ok := expr.(List); ok && len(list) > 0 {
			if name, ok := list[0].(Name); ok {
				value, found := env.Get(name)
				if found {
					if macro, ok := value.(Macro); ok {
						expanded, err := ExpandMacro(macro, list[1:], env)
						if err != nil {
							return nil, didExpand, err
						}
						expr = expanded
						didExpand = true
						continue
					}
				}
			}
		}
		break
	}

	return expr, didExpand, nil
}
