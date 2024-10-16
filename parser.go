package main

import (
	"fmt"
	"strconv"
	"strings"
)

func tokenize(input string) []string {
	specialTokens := []string{"(", ")", "'", "`", ",", "...", "\""}
	for _, token := range specialTokens {
		input = strings.ReplaceAll(input, token, " "+token+" ")
	}
	return strings.Fields(input)
}

func parseExpr(tokens []string) (Expression, []string, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("unexpected EOF")
	}
	token := tokens[0]
	tokens = tokens[1:]
	switch token {
	case "(":
		var list List
		for len(tokens) > 0 && tokens[0] != ")" {
			expr, remaining, err := parseExpr(tokens)
			if err != nil {
				return nil, tokens, err
			}
			list = append(list, expr)
			tokens = remaining
		}
		if len(tokens) == 0 {
			return nil, tokens, fmt.Errorf("missing closing parenthesis")
		}
		return list, tokens[1:], nil
	case ")":
		return nil, tokens, fmt.Errorf("unexpected closing parenthesis")
	case "\"":
		return parseString(token, tokens)
	case "'":
		expr, remaining, err := parseExpr(tokens)
		if err != nil {
			return nil, tokens, err
		}
		return List{Name("quote"), expr}, remaining, nil
	case "`":
		expr, remaining, err := parseExpr(tokens)
		if err != nil {
			return nil, tokens, err
		}
		return List{Name("quasiquote"), expr}, remaining, nil
	case ",":
		if len(tokens) > 0 && tokens[0] == "@" {
			tokens = tokens[1:]
			expr, remaining, err := parseExpr(tokens)
			if err != nil {
				return nil, tokens, err
			}
			return List{Name("unquote-splicing"), expr}, remaining, nil
		}
		expr, remaining, err := parseExpr(tokens)
		if err != nil {
			return nil, tokens, err
		}
		return List{Name("unquote"), expr}, remaining, nil
	case "...":
		return Name("..."), tokens, nil
	default:
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			return Number(num), tokens, nil
		}
		return Name(token), tokens, nil
	}
}

func parseString(firstToken string, remainingTokens []string) (Expression, []string, error) {
	if strings.HasSuffix(firstToken, "\"") && len(firstToken) > 1 {
		// Single-token string
		return String(firstToken[1 : len(firstToken)-1]), remainingTokens, nil
	}

	fullString := firstToken[1:] // Remove opening quote
	for len(remainingTokens) > 0 {
		token := remainingTokens[0]
		remainingTokens = remainingTokens[1:]
		if strings.HasSuffix(token, "\"") {
			// End of string found
			fullString += " " + token[:len(token)-1]
			return String(fullString), remainingTokens, nil
		}
		fullString += " " + token
	}
	return nil, remainingTokens, fmt.Errorf("unterminated string")
}

func Parse(input string) (Expression, error) {
	tokens := tokenize(input)
	expr, _, err := parseExpr(tokens)
	return expr, err
}
