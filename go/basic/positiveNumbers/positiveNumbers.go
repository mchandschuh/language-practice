package main

import (
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
	*n = append(*n, num)
	return nil
}

func yieldPositive(input <-chan int) <-chan int {
	output := make(chan int, cap(input))

	go func() {
		defer close(output)
		for item := range input {
			if item > 0 {
				output <- item
			}
		}
	}()

	return output
}
