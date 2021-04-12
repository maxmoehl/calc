package main

import (
	"fmt"
)

func buildAST(tokens []token) (operation, error) {
	// initialize values
	o := operation{}
	i := 0

	if tokens[i].Type() == typeLiteral {
		o.left = tokens[i].Value().(float64)
		i++
		if tokens[i].Type() == typeOperator {
			o.operator = tokens[i].Value().(rune)
		} else {
			return operation{}, fmt.Errorf("expected literal after operator but got operator")
		}
	} else if tokens[i].Type() == typeOperator {
		o.operator = tokens[i].Value().(rune)
	}
	i++

	for ; i < len(tokens); i++ {
		if tokens[i].Type() == typeOperator {
			if tokens[i-1].Type() == typeOperator {
				return operation{}, fmt.Errorf("expected literal after operator but got operator")
			}
			o = operation{
				operator: tokens[i].Value().(rune),
				left:     o,
				right:    nil,
			}
			continue
		}
		if tokens[i].Type() == typeLiteral {
			if tokens[i-1].Type() == typeLiteral {
				return operation{}, fmt.Errorf("expected operator after literal but got literal")
			}
			o.right = tokens[i].Value().(float64)
		}
	}
	return o, nil
}
