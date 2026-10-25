package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	calc "github.com/mchandschuh/language-practice/go/basic/calc/math"
)

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: calc <expression>")
		return 2
	}

	expression := strings.Join(args, " ")
	parsed, err := calc.Parse(expression)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	evaluator, err := calc.NewEvaluator(ctx, calc.DefaultQueueCapacity)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer evaluator.Close()

	results, err := evaluator.Submit(ctx, parsed, calc.Environment{})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	result := <-results
	if result.Err != nil {
		fmt.Fprintln(stderr, result.Err)
		return 1
	}

	fmt.Fprintln(stdout, result.Value)
	return 0
}
