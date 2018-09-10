package main

import (
	"fmt"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
)

func use(...interface{}) {}

func parseChunkOfCode(code string) *ast.File {
	p, err := parser.ParseFile(token.NewFileSet(), "example.go", code, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func main() {
	codeSample := `
	package main

	import(
		"fmt"
		"log"
	)

	func whatever(a, url string, arbitrary bool) error {
		return nil
	}

	func main() {
		fmt.Println("hello, world!")
	}
`
	x := interpret.NewInterpreter()
	p := parseChunkOfCode(codeSample)
	if err := x.Interpret(p); err != nil {
		log.Fatal(err)
	}

	fmt.Println(regexp.MustCompile(`\s+`).ReplaceAllString(x.RawOutput(), " "))
}
