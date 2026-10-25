package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type numArgs []string

func (n *numArgs) String() string {
	return strings.Join(*n, ", ")
}

func (n *numArgs) Set(value string) error {
	*n = append(*n, value)
	return nil
}

func main() {
	var nums numArgs
	flag.Var(&nums, "num", "A number")
	flag.Parse()

	input := make(chan int, len(nums))
	output := yieldPositive(input)

	for _, item := range nums {
		n, err := strconv.Atoi(item)
		if err != nil {
			fmt.Printf("Invalid number: %s", item)
		} else {
			input <- n
		}
	}

	close(input)

	for pn := range output {
		fmt.Printf("Positive: %d\n", pn)
	}
}
