package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
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
		color.Red(err.Error())
		os.Exit(1)
	}
	var res float64
	res, err = o.eval()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Result: %f\n", res)
}

func printToken(t token) {
	switch t.Type() {
	case typeOperator:
		fmt.Printf("\t%s\t%s\n", t.Type(), string(t.Value().(rune)))
	case typeLiteral:
		fmt.Printf("\t%s\t%f\n", t.Type(), t.Value().(float64))
	}
}
