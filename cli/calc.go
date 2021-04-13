// This application is the cli version of calc. It serves as a basic wrapper to
// enable the use of calc as a package in an application.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/maxmoehl/calc"

	"github.com/fatih/color"
)

func main() {
	if _, found := os.LookupEnv("DEBUG"); found {
		calc.SetDebug(true)
	}
	res, err := calc.Eval(strings.Join(os.Args[1:], ""))
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%f\n", res)
}
