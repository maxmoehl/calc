package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input := os.Args[1:]
	tokens, err := tokenize(strings.Join(input, ""))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("the following instructions have been read by the lexer:")
	for _, t := range tokens {
		printToken(t)
	}
	var o operation
	o, err = parse(tokens)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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
