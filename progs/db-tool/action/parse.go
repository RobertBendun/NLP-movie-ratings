package action

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

func From(code string) (Action, error) {
	expr, err := parser.ParseExpr(code)
	if err != nil {
		return nil, err
	}
	return from(expr)
}

func from(expr ast.Expr) (Action, error) {
	switch v := expr.(type) {
	case *ast.Ident:
		switch v.Name {
		case "bag":
			return &BagOfWordsAction{}, nil
		case "id":
			return &IdAction{}, nil
		case "group":
			return &GroupAction{Keys: make(map[string]struct{})}, nil
		}

	case *ast.CallExpr:
		if ident, ok := v.Fun.(*ast.Ident); ok {
			switch ident.Name {
			case "bag":
				if err := expectArguments(v, "bag", 1); err != nil {
					return nil, err
				}
				return BagOfWordsAction{}, nil

			case "head":
				if err := expectArguments(v, "head", 1); err != nil {
					return nil, err
				}
				if action, err := expectAction(v.Args[0]); err == nil {
					return &HeadAction{
						Action: action,
						called: make(map[string]string),
					}, nil
				} else {
					return nil, err
				}
			case "join":
				if err := expectArguments(v, "join", 2); err != nil {
					return nil, err
				}
				delim, err := expectString(v.Args[0])
				if err != nil {
					return nil, fmt.Errorf("arg 0 of join: %v", err)
				}
				action, err := expectAction(v.Args[1])
				if err != nil {
					return nil, fmt.Errorf("arg 1 of join: %v", err)
				}
				return &JoinAction{
					Delim:  delim,
					Action: action,
					joined: make(map[string]*strings.Builder),
				}, nil
			case "limitedJoin":
				if err := expectArguments(v, "limitedJoin", 3); err != nil {
					return nil, err
				}
				delim, err := expectString(v.Args[0])
				if err != nil {
					return nil, fmt.Errorf("arg 0 of limitedJoin: %v", err)
				}
				action, err := expectAction(v.Args[1])
				if err != nil {
					return nil, fmt.Errorf("arg 1 of limitedJoin: %v", err)
				}
				limit, err := expectInt(v.Args[2])
				if err != nil {
					return nil, fmt.Errorf("arg 2 of limitedJoin: %v", err)
				}

				return &LimitedJoinAction{
					JoinAction: JoinAction{
						Delim:  delim,
						Action: action,
						joined: make(map[string]*strings.Builder),
					},
					counts: make(map[string]int),
					limit: limit,
				}, nil
			case "limitedCountedJoin":
				if err := expectArguments(v, "limitedCountedJoin", 3); err != nil {
					return nil, err
				}
				delim, err := expectString(v.Args[0])
				if err != nil {
					return nil, fmt.Errorf("arg 0 of limitedCountedJoin: %v", err)
				}
				action, err := expectAction(v.Args[1])
				if err != nil {
					return nil, fmt.Errorf("arg 1 of limitedCountedJoin: %v", err)
				}
				limit, err := expectInt(v.Args[2])
				if err != nil {
					return nil, fmt.Errorf("arg 2 of limitedCountedJoin: %v", err)
				}

				return &LimitedCountedJoinAction{
					LimitedJoinAction: LimitedJoinAction{
						JoinAction: JoinAction{
							Delim:  delim,
							Action: action,
							joined: make(map[string]*strings.Builder),
						},
						counts: make(map[string]int),
						limit: limit,
					},
				}, nil
			}
		} else {
			return nil, fmt.Errorf("Expected identifier, got %v", v.Fun)
		}
	}
	return nil, fmt.Errorf("Expected either identifier or function call, got %v", expr)
}

func expectArguments(expr *ast.CallExpr, name string, count int) error {
	suffix := "s"
	if len(expr.Args) == 1 {
		suffix = ""
	}
	if len(expr.Args) != count {
		return fmt.Errorf("`%s` action requires %d argument%s, got %d",
			name, count, suffix, len(expr.Args))
	}
	return nil
}

func expectString(expr ast.Expr) (string, error) {
	str, ok := expr.(*ast.BasicLit)
	if !ok || str.Kind != token.STRING {
		return "", errors.New("Expected string literal")
	}
	return str.Value[1 : len(str.Value)-1], nil
}

func expectInt(expr ast.Expr) (int, error) {
	str, ok := expr.(*ast.BasicLit)
	if !ok || str.Kind != token.INT {
		return 0, errors.New("Expected int literal")
	}
	return strconv.Atoi(str.Value)
}

func expectAction(expr ast.Expr) (Action, error) {
	return from(expr)
}
