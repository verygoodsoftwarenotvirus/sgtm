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
	NormalVerbosity verbosity = iota
	HighVerbosity
)

type Interpreter interface {
	Interpret(input *ast.File) error
	RawOutput() string
}

type interpreter struct {
	verbosity    verbosity
	fileset      *token.FileSet
	debug        bool
	outputString string
	logger       *log.Logger
	replacer     *strings.Replacer
}

func NewInterpreter() Interpreter {
	return &interpreter{
		verbosity: NormalVerbosity,
		logger:    log.New(os.Stdout, "", log.LstdFlags),
		replacer:  defaultStringReplacer,
	}
}

func (i *interpreter) RawOutput() string {
	return i.outputString
}

func (i *interpreter) addToOutput(s string) {
	i.outputString += " " + s
}

func (i *interpreter) handleImport(d *ast.GenDecl) {
	s, _ := NewImportSpec(d, i.verbosity).Describe()
	i.addToOutput(s)
}

func (i *interpreter) handleFunction(f *ast.FuncDecl) error {
	s, err := NewFuncDecl(f, i.verbosity).Describe()
	if err != nil {
		return err
	} else {
		i.addToOutput(s)
	}
	return nil
}

func (i *interpreter) handleType(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		if ts, ok := spec.(*ast.TypeSpec); ok {
			s, _ := NewTypeDescriber(ts, i.verbosity).Describe()
			i.addToOutput(s)
		}
	}
}

func (i *interpreter) Interpret(input *ast.File) error {
	i.addToOutput(fmt.Sprintf("package %s. ", input.Name.Name))
	for _, decl := range input.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			switch x.Tok {
			case token.IMPORT:
				i.handleImport(x)
			case token.TYPE:
				i.handleType(x)
			}
		case *ast.FuncDecl:
			if err := i.handleFunction(x); err != nil {
				return err
			}
		}
	}
	return nil
}
