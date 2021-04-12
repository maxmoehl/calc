package main

type operation struct {
	// operator contains the operation that should be carried out on the left and right operand
	operator rune
	// left contains either a value (float64) or another Operation
	left interface{}
	// right contains either a value (float64) or another Operation
	right interface{}
}

func (o operation) eval() float64 {
	var left, right float64
	if o.left == nil {
		left = 0
	} else if f, ok := o.left.(float64); ok {
		left = f
	} else {
		left = o.left.(operation).eval()
	}

	if o.right == nil {
		right = 0
	} else if f, ok := o.right.(float64); ok {
		right = f
	} else {
		right = o.left.(operation).eval()
	}
	return calc(o.operator, left, right)
}

func calc(operator rune, left, right float64) float64 {
	switch operator {
	case '+':
		return left + right
	case '-':
		return left - right
	default:
		panic("unknown operation: " + string(operator))
	}
}
