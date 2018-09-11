package interpret

import (
	"go/ast"
)

type TypeDescriber struct {
	original  *ast.TypeSpec
	Verbosity verbosity
	Type      string
	Name      string
	Fields    map[string][]string
}

func NewTypeDescriber(in *ast.TypeSpec, v verbosity) *TypeDescriber {
	switch x := in.Type.(type) {
	// NOTES: This is where we would handle aliases of any type, so if I write:
	// 		type Something int
	// it won't get passed here because we don't (and can't reliably) handle that case
	case *ast.StructType:
		td := &TypeDescriber{
			original:  in,
			Verbosity: v,
			Type:      "struct",
			Name:      in.Name.Name,
			Fields:    map[string][]string{},
		}

		for _, f := range x.Fields.List {
			tn, ok1 := f.Type.(*ast.Ident)
			if !ok1 {
				// I don't know what else to do when this happens, because I don't know if it can happen
				panic("aaaaaaaaaaaaaaaaa")
			}

			var names []string
			for _, n := range f.Names {
				names = append(names, n.Name)
			}

			td.Fields[tn.Name] = append(td.Fields[tn.Name], names...)
		}

		return td
	}
	return nil
}

func (td *TypeDescriber) Describe() (string, error) {
	tmpl := `type {{ .Name }} 
	{{ if exported .Name }}
		, which is exported,  
	{{ end }} 
	has the following fields: 
	{{ range $type, $vars := .Fields }} 
		{{ range $i, $var := $vars }} {{ $var }} {{ if and (gt (len $vars) 1) (eq (sub 1 $i) (len $vars)) }}, {{ end }}{{ end }} which 
		{{ if gt (len $vars) 1 }} are {{ else }} is a {{ end }}
		{{ $type }}{{ if gt (len $vars) 1 }}s{{ end }} 
	{{ end }}.`
	s, err := RenderTemplate(tmpl, td)
	return s, err
}
