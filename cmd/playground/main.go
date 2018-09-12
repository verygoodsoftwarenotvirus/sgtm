package main

import (
	"fmt"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/interpreter"
	"github.com/verygoodsoftwarenotvirus/sgtm/pkg/speakers/say"
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

	//import(
	//	"fmt"
	//	"log"
	//)

	func main() {
		something, err := someFunction(someArg, 9)
		fmt.Println("hello, world!")
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

	speaker := say.SaySpeaker{}
	speaker.GenerateSpeech(y, "")
}
