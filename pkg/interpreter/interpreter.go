package interpret

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"
)

type verbosity int

const (
	Normal verbosity = iota
	HighVerbosity
)

type Interpreter interface {
	Interpret(input *ast.File) error
	RawOutput() string
}

type interpreter struct {
	fileset      *token.FileSet
	debug        bool
	outputString string
	logger       *log.Logger
	replacer     *strings.Replacer
}

func NewInterpreter() Interpreter {
	return &interpreter{
		logger:   log.New(os.Stdout, "", log.LstdFlags),
		replacer: defaultStringReplacer,
	}
}

func (i *interpreter) addToOutput(s string) {
	i.outputString += " " + s
}

func (i *interpreter) RawOutput() string {
	return i.outputString
}

func (i *interpreter) handleImport(d *ast.GenDecl) {
	i.addToOutput("importing")
	for ix, spec := range d.Specs {
		if is, ok := spec.(*ast.ImportSpec); ok {
			i.addToOutput(prepareName(is.Path.Value))
			if is.Name != nil && is.Name.Name != "" && is.Name.Name != is.Path.Value {
				i.addToOutput(fmt.Sprintf(" as %s ", is.Name.Name))
			}
			if ix != len(d.Specs)-1 {
				i.addToOutput(" and ")
			}
		}
	}
	i.addToOutput(".")
}

func (i *interpreter) handleFunction(f *ast.FuncDecl) error {
	funcDecl, err := NewFuncDecl(f)
	if err != nil {
		return err
	}

	if s, err := funcDecl.Describe(); err != nil {
		return err
	} else {
		i.addToOutput(s)
	}
	return nil
}

func (i *interpreter) Interpret(input *ast.File) error {
	i.addToOutput(fmt.Sprintf("package %s. ", input.Name.Name))
	for _, decl := range input.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			if x.Tok == token.IMPORT {
				i.handleImport(x)
			}
		case *ast.FuncDecl:
			if err := i.handleFunction(x); err != nil {
				return err
			}
		}
	}
	return nil
}
