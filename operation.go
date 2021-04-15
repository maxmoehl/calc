package calc

import (
	"fmt"
	"math"
)

type Operation interface {
	Operator() rune
	Left() Operation
	Right() Operation
	// Locked indicates whether or not it is possible to go down into the Operation.
	// If false, Operator, Left and Right might return unexpected values.
	Locked() bool
	// Eval returns the value this operation resolves to, or an error if one occurs.
	Eval() (float64, error)
}

// operation is a recursive struct that is the main building block of the abstract syntax tree.
type operation struct {
	// operator contains the operation that should be carried out on the Left and Right operand
	operator rune
	// left contains either a value (float64) or a pointer to another Operation
	left Operation
	// right contains either a value (float64) or a pointer to another Operation
	right Operation
	// locked stores whether or not this Operation can be modified
	locked bool
}

func (o *operation) Operator() rune {
	return o.operator
}

func (o *operation) Left() Operation {
	return o.left
}

func (o *operation) Right() Operation {
	return o.right
}

func (o *operation) Locked() bool {
	return o.locked
}

// Eval evaluates an Operation by first evaluating all sub-operations and evaluating itself.
func (o *operation) Eval() (float64, error) {
	var l, r float64
	var err error
	if o.Left() == nil {
		l = 0
	} else {
		l, err = o.Left().Eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	if o.Right() == nil {
		r = 0
	} else {
		r, err = o.Right().Eval()
	}
	if err != nil {
		return math.NaN(), err
	}
	return calc(o.Operator(), l, r)
}

type literalOperation struct {
	value float64
}

func (l *literalOperation) Operator() rune {
	return 'l'
}

func (l *literalOperation) Left() Operation {
	return nil
}

func (l *literalOperation) Right() Operation {
	return nil
}

func (l *literalOperation) Locked() bool {
	return true
}

func (l *literalOperation) Eval() (float64, error) {
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
