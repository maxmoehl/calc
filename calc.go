package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	typeOperator = "operator"
	typeLiteral  = "literal"
)

var validOperators = []rune{'+', '-'}
var validLiteralRunes = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.'}

func main() {
	input := os.Args[1:]
	tokens, err := tokenize(strings.Join(input, ""))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(100)
	}
	fmt.Println("the following instructions have been read by the lexer:")
	for _, t := range tokens {
		printToken(t)
	}
	var o operation
	o, err = buildAST(tokens)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(101)
	}
	fmt.Printf("Result: %f\n", o.eval())
}

func printToken(t token) {
	switch t.Type() {
	case typeOperator:
		fmt.Printf("%s\t%s\n", t.Type(), string(t.Value().(rune)))
	case typeLiteral:
		fmt.Printf("%s\t%f\n", t.Type(), t.Value().(float64))
	}
}
