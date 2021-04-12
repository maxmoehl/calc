package main

import (
	"fmt"
	"strconv"
)

func tokenize(input string) (tokens []token, err error) {
	symbols := []rune(input)
	var t token
	var ok bool

	for i := 0; i < len(symbols); i++ {

		if ok, t = isOperator(symbols[i]); ok {
			tokens = append(tokens, t)
			continue
		}

		if isValidLiteralRune(symbols[i]) {
			start := i
			for ; i < len(symbols) && isValidLiteralRune(symbols[i]); i++ {
			}
			t, err = parseLiteral(symbols[start:i])
			if err != nil {
				return tokens, fmt.Errorf("unable to parse literal '%s' at position %d - %d\n\t%s", string(symbols[start:i]), start, i-1, err.Error())
			}
			tokens = append(tokens, t)
			// reset value, for loop will increase it again
			i--
			continue
		}

		if isWhitespace(symbols[i]) {
			continue
		}

		return tokens, fmt.Errorf("unknown character '%s' at position %d", string(symbols[i]), i)
	}

	return
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

// isValidLiteralRune checks if a symbol is a valid part of a literal
func isValidLiteralRune(symbol rune) bool {
	for _, l := range validLiteralRunes {
		if l == symbol {
			return true
		}
	}
	return false
}

func isWhitespace(symbol rune) bool {
	return symbol == ' '
}

func parseLiteral(symbols []rune) (token, error) {
	v, err := strconv.ParseFloat(string(symbols), 64)
	if err != nil {
		return nil, err
	}
	return literal{v}, nil
}
