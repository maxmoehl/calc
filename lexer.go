package calc

import (
	"fmt"
	"math"
	"strconv"
)

var (
	typeOperator    = "operator   :"
	typeLiteral     = "literal    :"
	typeParentheses = "parenthesis:"
	typeIdentifier  = "identifier :"
	typeBraces      = "brace      :"
	typeComma       = "comma      :"

	validOperators       = []rune{'+', '-', '*', '/'}
	validLiteralRunes    = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.'}
	validParentheses     = []rune{'(', ')'}
	validIdentifierRunes = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	validBraces          = []rune{'{', '}'}
	validCommas          = []rune{','}
	whitespaces          = []rune{' ', '\n'}
)

func tokenize(input string) ([]Token, error) {
	symbols := []rune(input)
	var t Token
	var ok bool
	var err error
	var tokens []Token
	var newI int

	for i := 0; i < len(symbols); i++ {

		if ok, t = isOperator(symbols[i]); ok {
			tokens = append(tokens, t)
			continue
		}

		if ok, t, newI, err = isLiteral(symbols, i); ok {
			// if value could be parsed, everything is fine
			tokens = append(tokens, t)
			// use newI for further tokenization
			i = newI
			continue
		} else if err != nil {
			// otherwise check if an error occurred or if no symbol was found
			return nil, err
		}

		if ok, t = isParenthesis(symbols[i]); ok {
			tokens = append(tokens, t)
			continue
		}

		if ok, t = isBrace(symbols[i]); ok {
			tokens = append(tokens, t)
			continue
		}

		if ok, t, newI = isIdentifier(symbols, i); ok {
			tokens = append(tokens, t)
			i = newI
			continue
		}

		if ok, t = isComma(symbols[i]); ok {
			tokens = append(tokens, t)
			continue
		}

		if isWhitespace(symbols[i]) {
			continue
		}

		return nil, unknownSymbol(symbols[i], i)
	}

	return tokens, nil
}

// isOperator checks if a rune is an operator and if to returns the proper operator struct
func isOperator(symbol rune) (bool, Token) {
	if runeSliceContains(validOperators, symbol) {
		return true, operator{symbol}
	}
	return false, nil
}

// isLiteral checks if the next symbol is part of a literal. If so it tries to find the end of the literal
// and tries to parse the literal. If the first return value is true if a literal was found and parsed successfully,
// otherwise it is false. The returned index is the same as the one passed in if no literal was found, or the position
// of the last element of the literal if one was found. If parsing the literal fails, the error is returned.
// If bool is true, error is always nil.
func isLiteral(symbols []rune, i int) (bool, Token, int, error) {
	if isValidLiteralRune(symbols[i]) {
		start := i
		for ; i < len(symbols) && isValidLiteralRune(symbols[i]); i++ {
		}
		t, err := convertLiteral(symbols[start:i])
		if err != nil {
			return false, nil, i, err
		}
		// decrease value of i, outer for loop will increase it again
		return true, t, i - 1, nil
	}
	return false, nil, i, nil
}

func isParenthesis(symbol rune) (bool, Token) {
	if runeSliceContains(validParentheses, symbol) {
		return true, parenthesis{symbol}
	}
	return false, nil
}

// isWhitespace checks if a symbol is a valid whitespace
func isWhitespace(symbol rune) bool {
	return runeSliceContains(whitespaces, symbol)
}

// isValidLiteralRune checks if a symbol is a valid part of a literal
func isValidLiteralRune(symbol rune) bool {
	return runeSliceContains(validLiteralRunes, symbol)
}

// convertLiteral takes a list of runes, parses it as a float64 and stores it in a Token
func convertLiteral(symbols []rune) (Token, error) {
	v, err := strconv.ParseFloat(string(symbols), 64)
	if err != nil {
		return nil, err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return nil, fmt.Errorf("unable to parse literal: the parsed value is not a valid number: '%f'", v)
	}
	return literal{v}, nil
}

func isBrace(symbol rune) (bool, Token) {
	if runeSliceContains(validBraces, symbol) {
		return true, brace{symbol}
	}
	return false, nil
}

func isIdentifier(symbols []rune, i int) (bool, Token, int) {
	if isValidIdentifierRune(symbols[i]) {
		start := i
		for ; i < len(symbols) && isValidIdentifierRune(symbols[i]); i++ {
		}
		// decrease value of i, outer for loop will increase it again
		return true, identifier{string(symbols[start:i])}, i - 1
	}
	return false, nil, i
}

func isComma(symbol rune) (bool, Token) {
	if runeSliceContains(validCommas, symbol) {
		return true, comma{symbol}
	}
	return false, nil
}

func isValidIdentifierRune(symbol rune) bool {
	return runeSliceContains(validIdentifierRunes, symbol)
}

// unknownSymbol generates an error message for some symbols that are not supported but known.
func unknownSymbol(symbol rune, position int) error {
	switch symbol {
	case '[':
		fallthrough
	case ']':
		return fmt.Errorf("unknown character '%s' at position %d\n"+
			"\tdid u want to use parantheses or braces?", string(symbol), position+1)
	default:
		return fmt.Errorf("unknown character '%s' at position %d", string(symbol), position+1)
	}
}

func runeSliceContains(s []rune, r rune) bool {
	for _, sr := range s {
		if sr == r {
			return true
		}
	}
	return false
}
