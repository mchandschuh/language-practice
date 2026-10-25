package main

import (
	"fmt"
	"os"
	"strings"

	calc "github.com/mchandschuh/language-practice/go/basic/calc/math"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: calc <expression>")
		os.Exit(2)
	}

	expression := strings.Join(os.Args[1:], " ")
	parsed, err := calc.Parse(expression)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result, err := calc.Evaluate(parsed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}
