package main

import (
	"fmt"
)

type Vector []float64

func main() {
	a, b := Vector{1, 2, 3}, Vector{4, 5, 6}
	a, b, _ = swap(a, b)
	fmt.Printf("Swapped vectors a: %v, b: %v\n", a, b)

	c := Vector{7, 8, 9}
	sum, _ := add([]Vector{a, b, c}...)
	fmt.Printf("Sum of all vectors: %v\n", sum)

	multiplier := 3
	scaled := scale(sum, multiplier)

	fmt.Printf("Scaled up by %d the sum is: %v\n", multiplier, scaled)
}
