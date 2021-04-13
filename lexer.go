package main

import (
	"fmt"
	"math"
	"strconv"
)

var (
	typeOperator = "operator"
	typeLiteral  = "literal"

	validOperators = []rune{'+', '-', '*', '/'}
	validLiteralRunes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.'}
)

func tokenize(input string) ([]token, error) {
	symbols := []rune(input)
	var t token
	var ok bool
	var err error
	var tokens []token
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

		if isWhitespace(symbols[i]) {
			continue
		}

		return nil, unknownSymbol(symbols[i], i)
	}

	return tokens, nil
}

// isOperator checks if a rune is an operator and if to returns the proper operator struct
func isOperator(symbol rune) (bool, token) {
	for _, o := range validOperators {
		if o == symbol {
			return true, operator{symbol}
		}
	}
	return false, nil
}

// isLiteral checks if the next symbol is part of a literal. If so it tries to find the end of the literal
// and tries to parse the literal. If the first return value is true if a literal was found and parsed successfully,
// otherwise it is false. The returned index is the same as the one passed in if no literal was found, or the position
// of the last element of the literal if one was found. If parsing the literal fails, the error is returned.
// If bool is true, error is always nil.
func isLiteral(symbols []rune, i int) (bool, token, int, error) {
	if isValidLiteralRune(symbols[i]) {
		start := i
		for ; i < len(symbols) && isValidLiteralRune(symbols[i]); i++ {
		}
		t, err := convertLiteral(symbols[start:i])
		if err != nil {
			return false, nil, i, err
		}
		// reset value, for loop will increase it again
		i--
		return true, t, i, nil
	}
	return false, nil, i, nil
}

// isValidLiteralRune checks if a symbol is a valid part of a literal
func isValidLiteralRune(symbol rune) bool {
	for _, l := range validLiteralRunes {
		if l == symbol {
			return true
		}
	}
	return false
}

// isWhitespace checks if a symbol is a valid whitespace
func isWhitespace(symbol rune) bool {
	return symbol == ' ' || symbol == '\n' || symbol == '\r'
}

// convertLiteral takes a list of runes, parses it as a float64 and stores it in a token
func convertLiteral(symbols []rune) (token, error) {
	v, err := strconv.ParseFloat(string(symbols), 64)
	if err != nil {
		return nil, err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return nil, fmt.Errorf("unable to parse literal: the parsed value is not a valid number: '%f'", v)
	}
	return literal{v}, nil
}

// unknownSymbol generates an error message for some symbols that are not supported but known.
func unknownSymbol(symbol rune, position int) error {
	switch symbol {
	case ',':
		return fmt.Errorf("unknown character '%s' at position %d\n" +
			"\tyou should use a '.' for decimal places and omit any ','", string(symbol), position+1)
	default:
		return fmt.Errorf("unknown character '%s' at position %d", string(symbol), position+1)
	}
}
