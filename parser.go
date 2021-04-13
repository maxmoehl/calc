package calc

import (
	"fmt"
)

// parserMap maps all available token types to the corresponding parsing function. Each function
// gets the current root operation, the full list of tokens and the position that it should parse.
// the returned operation is the new root operation. operation will be empty if err != nil.
var parserMap = map[string]func(root *operation, tokens []token, i int) (*operation, error){
	typeLiteral: parseLiteral,
	typeOperator: parseOperator,
}

// parse takes a list of tokens in the order they occur in the statement. It builds a abstract syntax tree
// by chaining together operations in a recursive structure.
func parse(tokens []token) (*operation, error) {
	// initialize values
	root, i, err := initOperation(tokens)
	if err != nil {
		return nil, err
	}

	for ; i < len(tokens); i++ {
		root, err = parserMap[tokens[i].Type()](root, tokens, i)
		if err != nil {
			return nil, err
		}
	}
	return root, nil
}

// initOperation takes care of all available special scenarios that might occur when initializing the
// first operation.
func initOperation(tokens []token) (*operation, int, error) {
	i := 0
	root := operation{}
	if tokens[i].Type() == typeLiteral {
		root.Left = tokens[i].Value().(float64)
		i++
		if tokens[i].Type() == typeOperator {
			root.Operator = tokens[i].Value().(rune)
		} else {
			return nil, 0, fmt.Errorf("expected literal after Operator but got Operator")
		}
	} else if tokens[i].Type() == typeOperator {
		root.Operator = tokens[i].Value().(rune)
	}
	i++
	return &root, i, nil
}

func parseOperator(root *operation, tokens []token, i int) (*operation, error) {
	if tokens[i-1].Type() == typeOperator {
		return nil, fmt.Errorf("expected literal after Operator but got Operator")
	}
	op := tokens[i].Value().(rune)
	if op == '+' || op == '-' {
		// If the operator is a addition or subtraction we shift the tree to the left.
		// This leaves the right side of the root node empty for the next literal.
		return &operation{
			Operator: op,
			Left:     root,
			Right:    nil,
		}, nil
	} else if op == '*' || op == '/' {
		// If the operator is a multiplication or division we shift the lowest value
		// on the right to a new position one operation lower to the left. The operator
		// of the created operation is given by the current token. This leaves the lowest
		// right leave empty for the next literal.
		right := getRightLeaf(root)
		
		right.Right = &operation{
			Operator: op,
			Left:     right.Right,
			Right:    nil,
		}
		return root, nil
	}
	return nil, fmt.Errorf("unknown operation '%s' at position %d", string(op), i)
}

func parseLiteral(root *operation, tokens []token, i int) (*operation, error) {
	if tokens[i-1].Type() == typeLiteral {
		return nil, fmt.Errorf("expected Operator after literal but got literal")
	}
	r := getRightLeaf(root)
	if r.Right != nil {
		return nil, fmt.Errorf("expected lowest right leave to be nil for this literal but got %v", r.Right)
	}
	r.Right = tokens[i].Value().(float64)
	return root, nil
}

// getRightLeaf returns the lowest operation by going down the right side of the
// operations tree. It returns the operation where the right leaf is not another
// operation.
func getRightLeaf(root *operation) *operation {
	if r, ok := root.Right.(*operation); ok {
		return getRightLeaf(r)
	} else {
		return root
	}
}