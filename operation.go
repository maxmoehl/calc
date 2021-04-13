package main

import (
	"fmt"
	"math"
)

// operation is a recursive struct that is the main building block of the abstract syntax tree.
type operation struct {
	// Operator contains the operation that should be carried out on the Left and Right operand
	Operator rune
	// Left contains either a value (float64) or a pointer to another operation
	Left interface{}
	// Right contains either a value (float64) or a pointer to another operation
	Right interface{}
}

// eval evaluates an operation by first evaluating all sub-operations and evaluating itself.
func (o operation) eval() (float64, error) {
	var err error
	var left, right float64
	if o.Left == nil {
		left = 0
	} else if f, ok := o.Left.(float64); ok {
		left = f
	} else {
		left, err = o.Left.(*operation).eval()
	}
	if err != nil {
		return math.NaN(), err
	}

	if o.Right == nil {
		right = 0
	} else if f, ok := o.Right.(float64); ok {
		right = f
	} else {
		right, err = o.Right.(*operation).eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	return calc(o.Operator, left, right)
}

// calc carries out a operation, indicated by operator, on the two operands, left and right.
func calc(operator rune, left, right float64) (float64, error) {
	switch operator {
	case '+':
		return left + right, nil
	case '-':
		return left - right, nil
	case '*':
		return left * right, nil
	case '/':
		return left / right, nil
	default:
		return math.NaN(), fmt.Errorf("unknown operation: '%s'", string(operator))
	}
}
