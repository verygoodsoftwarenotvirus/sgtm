package main

import (
	"fmt"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"regexp"
)

func clean(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}

func parseChunkOfCode(code string) *ast.File {
	p, err := parser.ParseFile(token.NewFileSet(), "example.go", code, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func parseExampleFile(filename string) *ast.File {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return parseChunkOfCode(string(b))
}

const (
	codeSample = `
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

`
)

func main() {
	var (
		filename string         //= "example_packages/quine/main.go"
		chunks   []string = nil //[]string{"SomeInterface"}
		p        *ast.File
	)

	x := interpret.NewInterpreter(nil)

	if filename != "" {
		p = parseExampleFile(filename)
	} else {
		p = parseChunkOfCode(codeSample)
	}

	if err := x.InterpretFile(p, chunks); err != nil {
		log.Fatal(err)
	}

	y := x.RawOutput()
	fmt.Println(y)
}
