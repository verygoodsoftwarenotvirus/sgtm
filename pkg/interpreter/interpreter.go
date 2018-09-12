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

var currentVerbosity verbosity

const (
	NormalVerbosity verbosity = iota
	HighVerbosity
)

var defaultInterpreter *interpreter

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

func NewInterpreter(thingsToRead []string, verbose bool) Interpreter {
	partsToRead := map[string]struct{}{}
	for _, p := range thingsToRead {
		partsToRead[p] = struct{}{}
	}

	defaultInterpreter = &interpreter{
		partsToRead: partsToRead,
		logger:      log.New(os.Stdout, "", log.LstdFlags),
		replacer:    defaultStringReplacer,
	}

	currentVerbosity = NormalVerbosity
	if verbose {
		currentVerbosity = HighVerbosity
	}

	return defaultInterpreter
}

func (i *interpreter) RawOutput() string {
	return replace(strings.Join(i.output, ".\n"))
}

func (i *interpreter) addToOutput(s string) {
	i.output = append(i.output, clean(s))
}

func (i *interpreter) clear() {
	i.output = []string{}
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
	i.addToOutput(fmt.Sprintf("package %s", input.Name.Name))

	chunksFound := map[string]string{}
	for _, decl := range input.Decls {
		switch x := decl.(type) {
		case *ast.GenDecl:
			switch x.Tok {
			case token.IMPORT:
				is := NewImportSpec(x)
				desc, err := is.Describe()
				if err != nil {
					return err
				}
				i.addToOutput(desc)
				chunksFound[is.GetName()] = desc
			case token.TYPE:
				for _, spec := range x.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						td := NewTypeDescriber(ts)
						desc, err := td.Describe()

						if err != nil {
							return err
						}
						i.addToOutput(desc)
						chunksFound[td.GetName()] = desc
					}
				}
			}
		case *ast.FuncDecl:
			fd := NewFuncDecl(x)
			desc, err := fd.Describe()
			if err != nil {
				return err
			}
			i.addToOutput(desc)
			chunksFound[fd.GetName()] = desc
		}
	}

	if chunks != nil && len(chunks) > 0 {
		i.clear()
	}

	for _, chunk := range chunks {
		if s, ok := chunksFound[chunk]; ok {
			i.addToOutput(s)
		}
	}

	return nil
}
