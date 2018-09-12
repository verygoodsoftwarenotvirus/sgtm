package main

import (
	"fmt"
)

type Person struct {
	Name string
	age  int
}

func (p *Person) SayHi() {
	fmt.Printf("My name is %s and I am %d years old", p.Name, p.age)
}
