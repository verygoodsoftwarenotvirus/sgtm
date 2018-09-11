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

func main() {
	codeSample := `
	package main

	import(
		"fmt"
		"log"
	)

	type Structure interface {
		DoSomeWork(arg string) (uintptr, error) 
	}

	func (f fart) whatever(a, url string, arbitrary bool) error {
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

	y := x.RawOutput()
	fmt.Println(y)
}
