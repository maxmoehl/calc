package calc

import (
	"fmt"

	"github.com/maxmoehl/calc/types"
)

type Parser func(root types.Operation, tokens []Token, i int) (types.Operation, int, error)

var parser map[string]Parser

func init() {
	parser = make(map[string]Parser)
	parser[typeOperator] = parseOperator
	parser[typeLiteral] = parseLiteral
	parser[typeParentheses] = parseControl
	parser[typeIdentifier] = parseMacro
}

// parse takes a list of tokens in the order they occur in the statement. It builds a abstract syntax tree
// by chaining together operations in a recursive structure. Every Operation returned from parse is locked
// because those operations are considered done and should not be modified.
func parse(tokens []Token) (types.Operation, error) {
	var root types.Operation
	var err error

	for i := 0; i < len(tokens); i++ {
		root, i, err = parser[tokens[i].Type()](root, tokens, i)
		if err != nil {
			return nil, err
		}
	}

	if op, ok := root.(*operation); ok {
		op.locked = true
	}

	return root, nil
}

func parseOperator(root types.Operation, tokens []Token, i int) (types.Operation, int, error) {
	if i > 0 && tokens[i-1].Type() == typeOperator {
		return nil, i, fmt.Errorf("error: two consecutive operators at %d and %d", i-1, i)
	}
	op := tokens[i].Value().(rune)
	if op == '+' || op == '-' {
		// If the operator is a addition or subtraction we shift the tree to the left.
		// This leaves the right side of the root node empty for the next literal.
		// If this is the first Token of the expression, left will be nil and the
		// evaluation will assume it to be 0. This allows for negative signs
		// (and even unnecessary plus signs) at the beginning of an expression.
		return &operation{
			operator: op,
			left:     root,
			right:    nil,
		}, i, nil
	} else if op == '*' || op == '/' {
		if root == nil {
			return nil, i, fmt.Errorf("error: expression cannot start with %s", string(op))
		}
		if root.Locked() {
			return &operation{
				operator: op,
				left:     root,
			}, i, nil
		}
		// If the operator is a multiplication or division we shift the lowest value
		// on the right to a new position one Operation lower to the left. The operator
		// of the created Operation is given by the current Token. This leaves the lowest
		// right leave empty for the next literal or Operation.
		right := getRightLeaf(root)

		right.right = &operation{
			operator: op,
			left:     right.Right(),
			right:    nil,
		}
		return root, i, nil
	}
	return nil, i, fmt.Errorf("unknown Operation '%s' at position %d", string(op), i)
}

func parseLiteral(root types.Operation, tokens []Token, i int) (types.Operation, int, error) {
	if i > 0 && tokens[i-1].Type() != typeOperator {
		return nil, i, fmt.Errorf("expected operator before literal but got %s at position %d",
			tokens[i-1].Type(), i)
	}
	if root == nil {
		return &literalOperation{tokens[i].Value().(float64)}, i, nil
	}
	r := getRightLeaf(root)
	if r.Right() != nil {
		return nil, i, fmt.Errorf("expected lowest right leave to be nil for this literal but got %v", r.Right())
	}
	r.right = &literalOperation{tokens[i].Value().(float64)}
	return root, i, nil
}

func parseControl(root types.Operation, tokens []Token, i int) (types.Operation, int, error) {
	if i > 0 && tokens[i-1].Type() != typeOperator {
		return nil, i, fmt.Errorf("expected operator before parenthesis, but got %s, at position %d",
			tokens[i-1].Type(), i-1)
	}
	// store first value inside of parenthesis
	startIndex := i + 1
	// find corresponding closing parenthesis and check if one is found
	i = getClosingPart(tokens, i, '(', ')')
	if i == -1 {
		return nil, i, fmt.Errorf("missing closing parenthesis for opening parenthesis at position %d", startIndex)
	}
	// build the Operation for whatever was inside the parenthesis
	subOperation, err := parse(tokens[startIndex:i])
	if err != nil {
		return nil, i, err
	}

	if root == nil {
		return subOperation, i , nil
	}

	right := getRightLeaf(root)
	if right.Right() != nil {
		return nil, i, fmt.Errorf("expected right leaf to be empty but got %v", right.Right())
	}
	right.right = subOperation
	return root, i, nil
}

func parseMacro(root types.Operation, tokens []Token, i int) (types.Operation, int, error) {
	if i > 0 && tokens[i-1].Type() != typeOperator {
		return nil, i, fmt.Errorf("expected operator before identifiert but got %s", tokens[i-1].Type())
	}
	id := tokens[i].Value().(string)
	if macros[id] == nil {
		return nil, i, fmt.Errorf("unknown macro identifier %s", id)
	}
	i++
	if tokens[i].Type() != typeBraces || tokens[i].Value().(rune) != '{' {
		return nil, i, fmt.Errorf("expected opening brace after identifier but got type %s, value '%v'",
			tokens[i].Type(), tokens[i].Value())
	}
	startIndex := i + 1
	i = getClosingPart(tokens, i, '{', '}')
	if i == -1 {
		return nil, i, fmt.Errorf("unable to find closing brace for opening brace located at %d", startIndex)
	}

	parametersTokens, err := splitByComma(tokens[startIndex:i])
	if err != nil {
		return nil, i, err
	}

	var parameters []types.Operation
	var op types.Operation
	for _, parameterTokens := range parametersTokens {
		op, err = parse(parameterTokens)
		if err != nil {
			return nil, i, err
		}
		parameters = append(parameters, op)
	}
	macro, err := macros[id](parameters)
	if err != nil {
		return nil, i, err
	}

	if root == nil {
		return &macroOperation{macro}, i, nil
	}

	rightLeave := getRightLeaf(root)
	rightLeave.right = &macroOperation{macro}
	return root, i , nil
}

// getClosingPart tries to find the rune passed in as closing by ignoring all nested
// section delimited by opening and closed parameters.
func getClosingPart(tokens []Token, i int, opening, closing rune) int {
	parentheses := 1
	for i++; i < len(tokens); i++ {
		if tokens[i].Type() == typeParentheses || tokens[i].Type() == typeBraces {
			if tokens[i].Value().(rune) == opening {
				parentheses++
			} else if tokens[i].Value().(rune) == closing {
				parentheses--
			}
			if parentheses == 0 {
				return i
			}
		}
	}
	return -1
}

// getRightLeaf returns the lowest Operation by going down the right side of the
// operations tree. It returns the Operation where the right leaf is not another
// unlocked Operation (the right leaf is either a locked Operation or a number).
// We do not traverse into locked operations as they are considered a completed
// unit and should not be modified, see parse for more details.
func getRightLeaf(root types.Operation) *operation {
	if root.Right() == nil || root.Right().Locked() {
		return root.(*operation)
	} else {
		return getRightLeaf(root.Right())
	}
}

// splitByComma returns the passed in tokens as all tokens separated by commas.
// It ignores commas that are inside another macro.
func splitByComma(tokens []Token) ([][]Token, error) {
	var t Token
	var res [][]Token
	currentStart := 0
	for i := 0; i < len(tokens); i++ {
		t = tokens[i]
		if t.Type() == typeComma {
			// if the Token is a comma, take all tokens that passed since last comma or start and
			// store them in a slice.
			res = append(res, tokens[currentStart:i])
			currentStart = i + 1
		}
		if t.Type() == typeBraces {
			// If we find a opening brace we skip to the end of it to avoid the commas that this
			// macro might have.
			if t.Value().(rune) == '}' {
				return nil, fmt.Errorf("unexpected closing brace")
			}
			i = getClosingPart(tokens, i, '{', '}')
		}
		if i == len(tokens) - 1 {
			// if we are at the end of the section we group everything we have
			res = append(res, tokens[currentStart:i + 1])
		}
	}
	return res, nil
}
