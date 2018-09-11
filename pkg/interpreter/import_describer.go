package interpret

import (
	"go/ast"
	"strings"
)

type ImportSpec struct {
	verbosity verbosity
	original  *ast.GenDecl
	Imports   map[string]string // map[importAlias]importPath
}

func NewImportSpec(in *ast.GenDecl, v verbosity) *ImportSpec {
	imports := map[string]string{}
	for _, spec := range in.Specs {
		if is, ok := spec.(*ast.ImportSpec); ok {
			imp := prepareName(is.Path.Value)
			if is.Name != nil && is.Name.Name != "" && is.Name.Name != is.Path.Value {
				imports[is.Name.Name] = imp
			} else {
				imports[imp] = imp
			}
		}
	}

	return &ImportSpec{
		original:  in,
		verbosity: v,
		Imports:   imports,
	}
}

func (i *ImportSpec) Describe() (string, error) {
	tmpl := `importing {{ range $alias, $path := $.Imports }} {{ $path }} {{ if and (ne $alias "") (ne $alias $path) }} as {{ $path }}{{ end}}{{ if gt (len $.Imports) 2 }}, {{ else if eq (len $.Imports) 2 }} and {{ end }} {{ end}}. `
	s, err := RenderTemplate(tmpl, i)
	s = strings.Replace(strings.TrimSpace(s), "and  .", ".", 1) // gotta get rid of the excess `and` at the end
	return s, err
}
