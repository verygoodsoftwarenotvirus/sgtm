package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
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

	func whatever(s, x string, fart bool) error {
		return nil	
	}

	func main() {
		fmt.Println("hello, world!")
	}
`
	x := interpret.NewInterpreter()
	p := parseChunkOfCode(codeSample)
	x.Interpret(p)
	y := x.RawOutput()
	use(y)
}
