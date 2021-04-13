package main

import (
	"fmt"
	"math"
)

type operation struct {
	// operator contains the operation that should be carried out on the left and right operand
	operator rune
	// left contains either a value (float64) or another Operation
	left interface{}
	// right contains either a value (float64) or another Operation
	right interface{}
}

func (o operation) eval() (float64, error) {
	var err error
	var left, right float64
	if o.left == nil {
		left = 0
	} else if f, ok := o.left.(float64); ok {
		left = f
	} else {
		left, err = o.left.(operation).eval()
	}
	if err != nil {
		return math.NaN(), err
	}

	if o.right == nil {
		right = 0
	} else if f, ok := o.right.(float64); ok {
		right = f
	} else {
		right, err = o.left.(operation).eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	return calc(o.operator, left, right)
}

func calc(operator rune, left, right float64) (float64, error) {
	switch operator {
	case '+':
		return left + right, nil
	case '-':
		return left - right, nil
	default:
		return math.NaN(), fmt.Errorf("unknown operation: '%s'", string(operator))
	}
}
