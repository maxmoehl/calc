package calc

import (
	"fmt"
	"math"
	"strconv"
)

var (
	typeOperator    = "operator"
	typeParenthesis = "parenthesis"
	typeBrace       = "brace"
	typeComma       = "comma"
	typeWhitespace  = "whitespace"
	typeLiteral     = "literal"
	typeIdentifier  = "identifier"
)

// validRunes maps the type identifier for each allowed type to the runes it can consist of
var validRunes = map[string][]rune{
	typeOperator:    {'+', '-', '*', '/'},
	typeParenthesis: {'(', ')'},
	typeBrace:       {'{', '}'},
	typeComma:       {','},
	typeWhitespace:  {' ', '\n'},
	typeLiteral:     {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.'},
	typeIdentifier:  {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
}

// tokenize takes a string and creates a list of Token. In most cases each token
// consists of the type identifier and the rune that was detected. Literals and
// identifier have to be read by the external functions readIdentifier and
// readLiteral.
func tokenize(input string) ([]Token, error) {
	symbols := []rune(input)
	var t Token
	var err error
	var tokens []Token
	var s rune

	for i := 0; i < len(symbols); i++ {
		s = symbols[i]

		if isOfType(s, typeOperator) {
			tokens = append(tokens, token{typeOperator, s})
		} else if isOfType(s, typeParenthesis) {
			tokens = append(tokens, token{typeParenthesis, s})
		} else if isOfType(s, typeBrace) {
			tokens = append(tokens, token{typeBrace, s})
		} else if isOfType(s, typeComma) {
			tokens = append(tokens, token{typeComma, s})
		} else if isOfType(s, typeWhitespace) {
			// do nothing, if needed at some point these can also be read to give precise locations
			// of certain symbols, e.g. in case of an error
		} else if isOfType(s, typeLiteral) {
			t, i, err = readLiteral(symbols, i)
			tokens = append(tokens, t)
		} else if isOfType(s, typeIdentifier) {
			t, i = readIdentifier(symbols, i)
			tokens = append(tokens, t)
		} else {
			err = unknownSymbol(symbols[i], i)
		}
		if err != nil {
			return nil, err
		}
	}
	return tokens, nil
}

// readLiteral takes all symbols and the current position of the index. It then reads all
// symbols that belong to the current literal and returns the last index of the literal,
// a Token or an error.
func readLiteral(symbols []rune, i int) (Token, int, error) {
	start := i
	for ; i < len(symbols) && isOfType(symbols[i], typeLiteral); i++ {
	}
	t, err := convertLiteral(symbols[start:i])
	if err != nil {
		return nil, i, err
	}
	// decrease value of i, outer for loop will increase it again
	return t, i - 1, nil
}

// readLiteral takes all symbols and the current position of the index. It then reads all
// symbols that belong to the current identifier and returns the last index of the literal,
// a Token or an error.
func readIdentifier(symbols []rune, i int) (Token, int) {
	start := i
	for ; i < len(symbols) && isOfType(symbols[i], typeIdentifier); i++ {
	}
	// decrease value of i, outer for loop will increase it again
	return token{typeIdentifier, string(symbols[start:i])}, i - 1
}

// unknownSymbol generates an error message for some symbols that are not supported but known.
func unknownSymbol(symbol rune, position int) error {
	errString := fmt.Sprintf("unknown character '%s' at position %d\n", string(symbol), position+1)
	if symbol == '[' || symbol == ']' {
		errString += "\tdid u want to use parentheses or braces?\n"
	} else if symbol > 64 || symbol < 91 {
		// A - Z
		errString += "\tonly lowercase letters can be used as part of a macro identifier\n"
	}
	return fmt.Errorf(errString)
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
	return token{typeLiteral, v}, nil
}

// isOfType is a convenience function to check if a symbol is a valid symbol for
// the given type.
func isOfType(symbol rune, t string) bool {
	return runeSliceContains(validRunes[t], symbol)
}

// runeSliceContains checks if a runs slice contains a certain rune.
func runeSliceContains(s []rune, r rune) bool {
	for _, sr := range s {
		if sr == r {
			return true
		}
	}
	return false
}
