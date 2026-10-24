package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type positiveNumbers []int

func (n *positiveNumbers) String() string {
	values := make([]string, len(*n))
	for i, value := range *n {
		values[i] = strconv.Itoa(value)
	}

	return strings.Join(values, ", ")
}

func (n *positiveNumbers) Set(value string) error {
	var num, _ = strconv.Atoi(value)
	if num > 0 {
		*n = append(*n, num)
	}
	return nil
}

func main() {
	var nums positiveNumbers
	flag.Var(&nums, "num", "A number")
	flag.Parse()

	fmt.Printf("Positive Numbers: %s!\n", nums.String())
}
