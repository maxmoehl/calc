// This application is the cli version of calc. It serves as a basic wrapper to
// enable the use of calc as a package in an application.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/maxmoehl/calc"
)

func main() {
	if _, found := os.LookupEnv("DEBUG"); found {
		calc.SetDebug(true)
	}

	interactive := flag.Bool("interactive", false, "start interactive mode")
	flag.Parse()

	if *interactive {
		runInteractive()
	}

	if len(os.Args) == 1 {

		fmt.Println("Usage:")
		fmt.Println("  either execute a single calculation:")
		fmt.Println("    calc <mathematical expression>")
		fmt.Println("  or start the interactive mode:")
		fmt.Println("    calc -interactive")
		fmt.Println()
		fmt.Println("Loaded macros:")
		fmt.Println("  " + calc.GetLoadedMacros())
		return
	}

	res, err := calc.Eval(strings.Join(os.Args[1:], ""))
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	fmt.Printf("%g\n", res)
}

func runInteractive() {
	s := bufio.NewScanner(os.Stdin)
	var err error
	var in string
	var f float64
	for {
		fmt.Print("> ")
		s.Scan()
		in = s.Text()
		if in == "exit" {
			fmt.Println("bye")
			os.Exit(0)
		}
		f, err = calc.Eval(in)
		if err != nil {
			printError(err)
			continue
		}
		fmt.Printf("%g\n", f)
	}
}

func printError(err error) {
	fmt.Println("\x1b[31m" + err.Error() + "\x1b[0m")
}
