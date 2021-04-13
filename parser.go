package main

import (
	"fmt"
)

// parserMap maps all available token types to the corresponding parsing function. Each function
// gets the current root operation, the full list of tokens and the position that it should parse.
// the returned operation is the new root operation. operation will be empty if err != nil.
var parserMap = map[string]func(root operation, tokens []token, i int) (operation, error){
	typeLiteral: parseLiteral,
	typeOperator: parseOperator,
}

// parse takes a list of tokens in the order they occur in the statement. It builds a abstract syntax tree
// by chaining together operations in a recursive structure.
func parse(tokens []token) (operation, error) {
	// initialize values
	root, i, err := initOperation(tokens)
	if err != nil {
		return operation{}, err
	}

	for ; i < len(tokens); i++ {
		root, err = parserMap[tokens[i].Type()](root, tokens, i)
		if err != nil {
			return operation{}, nil
		}
	}
	return root, nil
}

// initOperation takes care of all available special scenarios that might occur when initializing the
// first operation.
func initOperation(tokens []token) (operation, int, error) {
	i := 0
	root := operation{}
	if tokens[i].Type() == typeLiteral {
		root.left = tokens[i].Value().(float64)
		i++
		if tokens[i].Type() == typeOperator {
			root.operator = tokens[i].Value().(rune)
		} else {
			return operation{}, 0, fmt.Errorf("expected literal after operator but got operator")
		}
	} else if tokens[i].Type() == typeOperator {
		root.operator = tokens[i].Value().(rune)
	}
	i++
	return root, i, nil
}

func parseOperator(root operation, tokens []token, i int) (operation, error) {
	if tokens[i-1].Type() == typeOperator {
		return operation{}, fmt.Errorf("expected literal after operator but got operator")
	}
	return operation{
		operator: tokens[i].Value().(rune),
		left:     root,
		right:    nil,
	}, nil
}

func parseLiteral(root operation, tokens []token, i int) (operation, error) {
	if tokens[i-1].Type() == typeLiteral {
		return operation{}, fmt.Errorf("expected operator after literal but got literal")
	}
	root.right = tokens[i].Value().(float64)
	return root, nil
}