package interpret

import (
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"
)

type verbosity int

var currentVerbosity = HighVerbosity

const (
	NormalVerbosity verbosity = iota
	HighVerbosity
)

type Interpreter interface {
	InterpretFile(input *ast.File, chunks []string) error
	RawOutput() string
}

type interpreter struct {
	partsToRead map[string]struct{}
	output      []string
	verbosity   verbosity
	fileset     *token.FileSet
	debug       bool
	logger      *log.Logger
	replacer    *strings.Replacer
}

func NewInterpreter(thingsToRead []string) Interpreter {
	partsToRead := map[string]struct{}{}
	for _, p := range thingsToRead {
		partsToRead[p] = struct{}{}
	}

	return &interpreter{
		partsToRead: partsToRead,
		verbosity:   NormalVerbosity,
		logger:      log.New(os.Stdout, "", log.LstdFlags),
		replacer:    defaultStringReplacer,
	}
}

func (i *interpreter) RawOutput() string {
	return strings.Join(i.output, ".\n")
}

func (i *interpreter) addToOutput(s string) {
	i.output = append(i.output, clean(s))
}

func (i *interpreter) handleImport(d *ast.GenDecl) {
	s, _ := NewImportSpec(d).Describe()
	i.addToOutput(s)
}

func (i *interpreter) handleFunction(f *ast.FuncDecl) error {
	s, err := NewFuncDecl(f).Describe()
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
			desc := NewTypeDescriber(ts)
			if desc != nil {
				s, err := desc.Describe()
				if err != nil {
					panic(err)
				}
				i.addToOutput(s)
			}
		}
	}
}

func (i *interpreter) InterpretFile(input *ast.File, chunks []string) error {
	var err error
	chunksFound := map[string]string{}
	for _, decl := range input.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			switch x.Tok {
			case token.IMPORT:
				i := NewImportSpec(x)
				chunksFound[i.GetName()], err = i.Describe()
				if err != nil {
					return err
				}
			case token.TYPE:
				for _, spec := range x.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						i := NewTypeDescriber(ts)
						chunksFound[i.GetName()], err = i.Describe()
						if err != nil {
							return err
						}
					}
				}
			}
		case *ast.FuncDecl:
			fd := NewFuncDecl(x)
			chunksFound[fd.GetName()], err = fd.Describe()
			if err != nil {
				return err
			}
		}
	}

	for _, chunk := range chunks {
		if s, ok := chunksFound[chunk]; ok {
			i.addToOutput(s)
		}
	}

	return nil
}
