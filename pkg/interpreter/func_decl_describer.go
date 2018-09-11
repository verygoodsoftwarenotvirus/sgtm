package interpret

import (
	"errors"
	"fmt"
	"go/ast"
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

func NewFuncDecl(f *ast.FuncDecl, v verbosity) *FuncDecl {
	funcDecl := &FuncDecl{
		original:           f,
		Verbosity:          v,
		Name:               f.Name.Name,
		ParameterArguments: parseArguments(f.Type.Params),
		ReturnArguments:    parseArguments(f.Type.Results),
	}

	return funcDecl
}

func parseArguments(in *ast.FieldList) []ArgDesc {
	var out []ArgDesc
	if in != nil {
		for _, t := range in.List {
			paramType, ok := t.Type.(*ast.Ident)
			if !ok {
				panic(errors.New("invalid param list?"))
			}

			var names []string
			for _, n := range t.Names {
				names = append(names, prepareName(n.Name))
			}

			out = append(out, ArgDesc{Type: prepareName(paramType.Name), Names: names})
		}
	}
	return out
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
			{{ if lt (len $arg.Names) 2 }}
				{{ if startsWithVowel $arg.Type }} an {{ else }} a {{ end }} {{ $arg.Type }}
			{{ else }}
				{{ len $arg.Names }} {{ $arg.Type }}s
				{{ if verbose $.Verbosity }} called
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

	return RenderTemplate(tmpl, f)
}

func (f FuncDecl) describeReturns() (string, error) {
	tmpl := `
	returning
	{{ if not .ReturnArguments }} nothing {{ else}}
		{{ range $i, $arg := .ReturnArguments}}
			{{ if ne $i 0 }} and {{ end }}
			{{ if lt (len $arg.Names) 2 }}
				{{ if startsWithVowel $arg.Type }} an {{ else }} a {{ end }} {{ $arg.Type }}
			{{ else }}
				{{ len $arg.Names }} {{ $arg.Type }}s
				{{ if verbose $.Verbosity }} called
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

	return RenderTemplate(tmpl, f)
}

func (f FuncDecl) describeBody() (string, error) {
	nbs := NewBlockStmt(f.original.Body.List, f.Verbosity)
	nbs.Describe()

	return "", nil
}
