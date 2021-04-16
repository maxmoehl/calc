package calc

import (
	"fmt"

	"github.com/maxmoehl/calc/types"
)

type Parser func(root types.Node, tokens []Token, i int) (types.Node, int, error)

var parser map[string]Parser

func init() {
	parser = make(map[string]Parser)
	parser[typeOperator] = parseOperator
	parser[typeLiteral] = parseLiteral
	parser[typeParenthesis] = parseControl
	parser[typeIdentifier] = parseMacro
}

// parse takes a list of tokens in the order they occur in the statement. It builds a abstract syntax tree
// by chaining together operations in a recursive structure. Every Operation returned from parse is locked
// because those operations are considered done and should not be modified.
func parse(tokens []Token) (types.Node, error) {
	var root types.Node
	var err error

	for i := 0; i < len(tokens); i++ {
		root, i, err = parser[tokens[i].Type()](root, tokens, i)
		if err != nil {
			return nil, err
		}
	}

	if op, ok := root.(*operation); ok {
		_, err = getRightOperationNil(root)
		if err == nil {
			return nil, fmt.Errorf("expression has trailing operand")
		}
		op.locked = true
	}

	return root, nil
}

// parseOperator handles tokens that are of typeOperator. The handling is
// defined in parsePlusMinus and parseMulDiv
func parseOperator(root types.Node, tokens []Token, i int) (types.Node, int, error) {
	var err error
	op := tokens[i].Value().(rune)
	if op == '+' || op == '-' {
		root, err = parsePlusMinus(root, op)
		return root, i, err
	} else if op == '*' || op == '/' {
		root, err = parseMulDiv(root, op)
		return root, i, err
	}
	return nil, i, fmt.Errorf("unknown Operation '%s' at position %d", string(op), i)
}

// parsePlusMinus parses the operators '+' and '-'. The tree gets shifted to the
// left by replacing the root element with a new element and placing the root
// element on the left side of the new element.
// This leaves the right side of the root node empty for the next Node.
// If this is the first Token of the expression, left will be nil since the
// passed in root node is nil and the evaluation will assume it to be 0. This
// allows for negative signs (and even unnecessary plus signs) at the
// beginning of an expression.
func parsePlusMinus(root types.Node, operator rune) (types.Node, error) {
	// This covers a special case: If there is a root element and the root element
	// is a operation, the right side needs to have a node. Otherwise we have two
	// operators without a operand in between them.
	if o, ok := root.(*operation); root != nil && ok && o.right == nil {
		return nil, fmt.Errorf("expected right side of root node to be non-nil but got nil")
	}
	return &operation{
		operator: operator,
		left:     root,
		right:    nil,
	}, nil
}

// parseMulDiv parses the operators '*' and '/'. The lowest value on the right
// that is Locked get shifted to a new position one Operation lower to the left.
// The operator of the created Operation is given by the current Token. This
// leaves the lowest right leave empty for the next Node. If the root element
// is already locked create a new root node of type operation and place the
// passed in root node on the left side of the newly created operation.
func parseMulDiv(root types.Node, operator rune) (types.Node, error) {
	if root == nil {
		return nil, fmt.Errorf("error: expression cannot start with %s", string(operator))
	}
	if root.Locked() {
		return &operation{
			operator: operator,
			left:     root,
		}, nil
	}

	right, err := getRightOperationNonNil(root)
	if err != nil {
		return nil, err
	}

	right.right = &operation{
		operator: operator,
		left:     right.Right(),
		right:    nil,
	}
	return root, nil
}

func parseLiteral(root types.Node, tokens []Token, i int) (types.Node, int, error) {
	if root == nil {
		return &literal{tokens[i].Value().(float64)}, i, nil
	}
	r, err := getRightOperationNil(root)
	if err != nil {
		return nil, i, err
	}
	r.right = &literal{tokens[i].Value().(float64)}
	return root, i, nil
}

func parseControl(root types.Node, tokens []Token, i int) (types.Node, int, error) {
	// store first value inside of parenthesis
	startIndex := i + 1
	// find corresponding closing parenthesis and check if one is found
	i, err := getClosingPart(tokens, i)
	if i == -1 {
		return nil, i, fmt.Errorf("missing closing parenthesis for opening parenthesis at position %d", startIndex)
	}
	// build the Operation for whatever was inside the parenthesis
	op, err := parse(tokens[startIndex:i])
	if err != nil {
		return nil, i, err
	}

	if root == nil {
		return op, i , nil
	}

	right, err := getRightOperationNil(root)
	if err != nil {
		return nil, i, err
	}
	right.right = op
	return root, i, nil
}

func parseMacro(root types.Node, tokens []Token, i int) (types.Node, int, error) {
	id := tokens[i].Value().(string)
	if macroIndex[id] == nil {
		return nil, i, fmt.Errorf("unknown macro identifier %s", id)
	}
	i++
	if tokens[i].Type() != typeBrace || tokens[i].Value().(rune) != '{' {
		return nil, i, fmt.Errorf("expected opening brace after identifier but got type %s, value '%v'",
			tokens[i].Type(), tokens[i].Value())
	}
	startIndex := i + 1
	i, err := getClosingPart(tokens, i)
	if err != nil {
		return nil, i, err
	}

	parametersTokens, err := splitByComma(tokens[startIndex:i])
	if err != nil {
		return nil, i, err
	}

	var parameters []types.Node
	var op types.Node
	for _, parameterTokens := range parametersTokens {
		op, err = parse(parameterTokens)
		if err != nil {
			return nil, i, err
		}
		parameters = append(parameters, op)
	}
	macro, err := macroIndex[id](parameters)
	if err != nil {
		return nil, i, err
	}

	if root == nil {
		return &macroOperation{macro}, i, nil
	}

	rightLeave, err := getRightOperationNil(root)
	if err != nil {
		return nil, i, err
	}
	rightLeave.right = &macroOperation{macro}
	return root, i , nil
}

// getClosingPart tries to find the rune passed in as closing by ignoring all nested
// section delimited by opening and closed parameters.
func getClosingPart(tokens []Token, i int) (int, error) {
	var opening, closing rune
	t := tokens[i].Type()
	if t == typeParenthesis || t == typeBrace {
		opening = validRunes[t][0]
		closing = validRunes[t][1]
	} else {
		return -1, fmt.Errorf("unknown control combination %s", t)
	}
	parentheses := 1
	for i++; i < len(tokens); i++ {
		if tokens[i].Type() == typeParenthesis || tokens[i].Type() == typeBrace {
			if tokens[i].Value().(rune) == opening {
				parentheses++
			} else if tokens[i].Value().(rune) == closing {
				parentheses--
			}
			if parentheses == 0 {
				return i, nil
			}
		}
	}
	return -1, fmt.Errorf("missing at least one closing part of control %s", t)
}

// getRightOperationNil searches the lowest operation on the right and expects that
// operation is nil on the right side.
func getRightOperationNil(n types.Node) (*operation, error) {
	o, err := getRightOperation(n)
	if err != nil {
		return nil, err
	}
	if o.right != nil {
		return nil, fmt.Errorf("lowest operation on the right side has no nil right side: %+v", o.right)
	}
	return o, nil
}

// getRightOperationNonNil is the corresponding function to getRightOperationNil. It
// searches for the lowest operation on the right and expects the right side to be
// nil. If that is not the case an error is returned
func getRightOperationNonNil(n types.Node) (*operation, error) {
	o, err := getRightOperation(n)
	if err != nil {
		return nil, err
	}
	if o.right == nil {
		return nil, fmt.Errorf("lowest operation on the right side is nil")
	}
	return o, nil
}

// getRightOperation returns the lowest operation by going down the right side of the
// operations tree. It returns the operation where the right leaf is not another
// unlocked operation (the right leaf is either a locked operation, a number or nil).
// Before passing a node it should be checked if that node is locked. If a locked node
// is passed in or the node is not of type operation, nil is returned.
// We do not traverse into locked operations as they are considered a completed
// unit and should not be modified, see parse for more details.
func getRightOperation(n types.Node) (*operation, error) {
	if n.Locked() {
		return nil, fmt.Errorf("received locked Node but expected unlocked Node")
	}
	o, ok := n.(*operation)
	if !ok {
		return nil, fmt.Errorf("expected Node of type operation but got %T", n)
	}
	if o.Right() == nil || o.Right().Locked() {
		return o, nil
	} else {
		return getRightOperation(o.Right())
	}
}

// splitByComma returns the passed in tokens as all tokens separated by commas.
// It ignores commas that are inside another macro.
func splitByComma(tokens []Token) ([][]Token, error) {
	var t Token
	var res [][]Token
	var err error
	currentStart := 0
	for i := 0; i < len(tokens); i++ {
		t = tokens[i]
		if t.Type() == typeComma {
			// if the Token is a comma, take all tokens that passed since last comma or start and
			// store them in a slice.
			res = append(res, tokens[currentStart:i])
			currentStart = i + 1
		}
		if t.Type() == typeBrace {
			// If we find a opening brace we skip to the end of it to avoid the commas that any
			// nested macro might have.
			if t.Value().(rune) == '}' {
				return nil, fmt.Errorf("unexpected closing brace")
			}
			i, err = getClosingPart(tokens, i)
			if err != nil {
				return nil, err
			}
		}
		if i == len(tokens) - 1 {
			// if we are at the end of the section we group everything we have
			res = append(res, tokens[currentStart:i + 1])
		}
	}
	return res, nil
}
