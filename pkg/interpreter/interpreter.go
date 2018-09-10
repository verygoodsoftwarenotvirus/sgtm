package interpret

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"
)

type Interpreter interface {
	Interpret(input *ast.File)
	RawOutput() string
}

type interpreter struct {
	fileset *token.FileSet
	debug   bool
	outputString string
	logger  *log.Logger
	replacer *strings.Replacer
}

func NewInterpreter() Interpreter {
	return &interpreter{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		replacer: defaultStringReplacer,
	}
}

func (i *interpreter) addToOutput(s string) {
	i.outputString += " " + s
}

func (i *interpreter) RawOutput() string {
	return i.outputString
}

func (i *interpreter) prepareImport(spec *ast.ImportSpec) string {
	return i.replacer.Replace(strings.Replace(spec.Path.Value, `"`, ``, -1))
}

func (i *interpreter) handleImport(d *ast.GenDecl) {
	i.addToOutput("importing")
	for ix, spec := range d.Specs {
		if is, ok := spec.(*ast.ImportSpec); ok  {
			i.addToOutput(i.prepareImport(is))
			if ix != len(d.Specs) - 1 {
				i.addToOutput("and")
			}
		}
	}
}

func (i *interpreter) Interpret(input *ast.File) {
	for _,  decl := range input.Decls {
		switch d := decl.(type)  {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				i.handleImport(d)
			}
		case *ast.FuncDecl:
			funcName := d.Name.Name
			i.addToOutput(fmt.Sprintf("function declared called %s", funcName))

			for _, t := range d.Type.Params.List {
				paramType, ok := t.Type.(*ast.Ident)
				if !ok {
					panic("invalid param list?")
				}
				summary := fmt.Sprintf(" accepts %ss called ", paramType)

				for ix, n := range t.Names {
					summary += n.Name
					if ix != len(t.Names)-1 {
						summary += " and "
					}
				}

				i.addToOutput(summary)
			}
		}
	}
}
