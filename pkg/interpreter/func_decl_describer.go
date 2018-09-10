package interpret

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"text/template"
)

type ArgDesc struct {
	Type  string
	Names []string
}

type FuncDecl struct {
	original           *ast.FuncDecl
	Verbosity          verbosity
	Name               string
	ParameterArguments []ArgDesc
	ReturnArguments    []ArgDesc
}

func NewFuncDecl(f *ast.FuncDecl) (*FuncDecl, error) {
	var err error
	funcDecl := &FuncDecl{
		original:           f,
		Name:               f.Name.Name,
		ParameterArguments: []ArgDesc{},
		ReturnArguments:    []ArgDesc{},
	}

	funcDecl.ParameterArguments, err = parseArguments(f.Type.Params)
	if err != nil {
		return nil, err
	}

	funcDecl.ReturnArguments, err = parseArguments(f.Type.Results)
	if err != nil {
		return nil, err
	}

	return funcDecl, nil
}

func parseArguments(in *ast.FieldList) ([]ArgDesc, error) {
	var out []ArgDesc
	if in != nil {
		for _, t := range in.List {
			paramType, ok := t.Type.(*ast.Ident)
			if !ok {
				return nil, errors.New("invalid param list?")
			}

			var names []string
			for _, n := range t.Names {
				names = append(names, prepareName(n.Name))
			}

			out = append(out, ArgDesc{Type: paramType.Name, Names: names})
		}
	}
	return out, nil
}

func (f FuncDecl) Describe() (string, error) {
	argstmt, err := f.describeArguments()
	if err != nil {
		return "", err
	}

	retstmt, err := f.describeReturns()
	if err != nil {
		return "", err
	}

	if _, err = f.describeBody(); err != nil {
		return "", err
	}

	return fmt.Sprintf(`function declared called %s  %s  %s.`, f.Name, argstmt, retstmt), nil
}

func (f FuncDecl) describeArguments() (string, error) {
	tmpl := `
	accepting
	{{ if not .ParameterArguments }} nothing {{ else}}
		{{ range $i, $arg := .ParameterArguments}}
			{{ if ne $i 0 }} and {{ end }}
			{{ if eq (len $arg.Names) 1 }}
				{{ if startsWithVowel $arg.Type }} an {{ else }} a {{ end }} {{ $arg.Type }}
			{{ else }}
				{{ len $arg.Names }} {{ $arg.Type }}s
				{{ if verbose }} called
					{{ range $i, $x := $arg.Names }}
						{{ if (ne (len $x) 0) }}
							{{ if ne $i 0 }} and {{ end }} {{ $x }}
						{{ end }}
					{{ end }}
				{{ end }}
			{{ end }}
		{{ end }}
	{{ end }}

	`

	return RenderTemplate(tmpl, f, f.TemplateFuncs())
}

func (f FuncDecl) describeReturns() (string, error) {
	var out = " returning "
	if f.ReturnArguments == nil {
		out += "nothing. "
		return out, nil
	}

	for _, r := range f.ReturnArguments {
		if startsWithVowel(r.Type) {
			out += fmt.Sprintf("an %s ", r.Type)
		} else {
			out += fmt.Sprintf("a %s ", r.Type)
		}

		out += strings.Join(r.Names, ", and ")
	}
	return out, nil
}

func (f FuncDecl) describeBody() (string, error) {
	for _, b := range f.original.Body.List {
		println(b)
	}

	return "", nil
}

func (f FuncDecl) TemplateFuncs() template.FuncMap {
	fm := defaultFuncMap
	fm["startsWithVowel"] = startsWithVowel
	fm["verbose"] = func() bool { return f.Verbosity == HighVerbosity }
	return fm
}
