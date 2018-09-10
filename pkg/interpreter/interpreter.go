package interpret

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"
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

func (i *interpreter) prepareImport(spec *ast.ImportSpec) string {
	return i.replacer.Replace(strings.Replace(spec.Path.Value, `"`, ``, -1))
}

func (i *interpreter) handleImport(d *ast.GenDecl) {
	i.addToOutput("importing")
	for ix, spec := range d.Specs {
		if is, ok := spec.(*ast.ImportSpec); ok {
			i.addToOutput(i.prepareImport(is))
			if ix != len(d.Specs)-1 {
				i.addToOutput(" and ")
			}
		}
	}
	i.addToOutput(".")
}

func (i *interpreter) handleFunction(f *ast.FuncDecl) error {
	var err error
	funcDecl := &FuncDecl{
		Name:               f.Name.Name,
		ParameterArguments: []ArgDesc{},
		ReturnArguments:    []ArgDesc{},
	}

	funcDecl.ParameterArguments, err = i.parseArguments(f.Type.Params)
	if err != nil {
		return err
	}

	funcDecl.ReturnArguments, err = i.parseArguments(f.Type.Results)
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

func (i *interpreter) parseArguments(in *ast.FieldList) ([]ArgDesc, error) {
	var out []ArgDesc
	if in != nil {
		for _, t := range in.List {
			paramType, ok := t.Type.(*ast.Ident)
			if !ok {
				return nil, errors.New("invalid param list?")
			}

			var names []string
			for _, n := range t.Names {
				names = append(names, n.Name)
			}

			out = append(out, ArgDesc{Type: paramType.Name, Names: names})
		}
	}
	return out, nil
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
