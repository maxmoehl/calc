package calc

import (
	"fmt"
	"math"

	"github.com/maxmoehl/calc/types"
)

// operation is a recursive struct that is the main building block of the abstract syntax tree.
type operation struct {
	// operator contains the operation that should be carried out on the left and right operand
	operator rune
	// left contains either a value (float64) or a pointer to a Node
	left types.Node
	// right contains either a value (float64) or a pointer to a Node
	right types.Node
	// locked stores whether or not this operation can be modified
	locked bool
}

func (o *operation) Locked() bool {
	return o.locked
}

// Eval evaluates an Operation by first evaluating all sub-operations and evaluating itself.
func (o *operation) Eval() (float64, error) {
	var l, r float64
	var err error
	if o.left == nil {
		l = 0
	} else {
		l, err = o.left.Eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	if o.right == nil {
		r = 0
	} else {
		r, err = o.right.Eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	return calc(o.operator, l, r)
}

type literal struct {
	value float64
}

func (l *literal) Locked() bool {
	return true
}

func (l *literal) Eval() (float64, error) {
	return l.value, nil
}

// calc carries out a Operation, indicated by operator, on the two operands, left and right.
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
		return math.NaN(), fmt.Errorf("unknown Operation: '%s'", string(operator))
	}
}
