package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
)

var DEBUG = false

func main() {
	res, err := run(strings.Join(os.Args[1:], ""))
	if _, found := os.LookupEnv("DEBUG"); found {
		DEBUG = true
	}
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%f\n", res)
}

func run(input string) (float64, error) {
	// run lexer
	tokens, err := tokenize(input)
	if err != nil {
		return math.NaN(), err
	}
	if DEBUG {
		fmt.Println("the following instructions have been read by the lexer:")
		for _, t := range tokens {
			printToken(t)
		}
	}

	// run parser
	var o *operation
	o, err = parse(tokens)
	if err != nil {
		return math.NaN(), err
	}
	if DEBUG {
		fmt.Printf("the following abstract syntax tree has been generated by the parser.\n" +
			"operator conversion:\n  %v -> %v\n  %v -> %v\n  %v -> %v\n  %v -> %v\n",
			'+', "+", '-', "-", '*', "*", '/', "/")
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(o)
	}

	// evaluate result
	var res float64
	res, err = o.eval()
	if err != nil {
		return math.NaN(), err
	}
	return res, nil
}

func printToken(t token) {
	switch t.Type() {
	case typeOperator:
		fmt.Printf("\t%s\t%s\n", t.Type(), string(t.Value().(rune)))
	case typeLiteral:
		fmt.Printf("\t%s\t%f\n", t.Type(), t.Value().(float64))
	}
}
