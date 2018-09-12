package interpret

import (
	"go/ast"
	"strings"
)

type TypeDescriber struct {
	original *ast.TypeSpec
	Type     string
	Name     string
	Fields   map[string][]string
}

type InterfaceDescriber struct {
	Name    string
	Methods []InterfaceMethodDescriber
}

type InterfaceMethodDescriber struct {
	BelongsTo string
	Name      string
	Args      []ArgDesc
	Returns   []ArgDesc
}

type ArrayDescriber struct {
	Type string
	Name string
	Elt  string
}

func NewTypeDescriber(in *ast.TypeSpec) Describer {
	switch x := in.Type.(type) {
	// NOTES: This is where we would handle aliases of any type, so if I write:
	// 		type Something int
	// it won't get passed here because we don't (and can't reliably) handle that case
	case *ast.StructType:
		td := &TypeDescriber{
			original: in,
			Type:     "struct",
			Name:     in.Name.Name,
			Fields:   map[string][]string{},
		}

		for _, f := range x.Fields.List {
			tn, ok := f.Type.(*ast.Ident)
			if !ok {
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
	case *ast.InterfaceType:
		out := &InterfaceDescriber{
			Name: in.Name.Name,
		}
		for _, method := range x.Methods.List {
			if ft, ok := method.Type.(*ast.FuncType); ok && len(method.Names) > 0 {
				out.Methods = append(out.Methods, InterfaceMethodDescriber{
					BelongsTo: out.Name,
					Name:      method.Names[0].Name,
					Returns:   parseArguments(ft.Results),
					Args:      parseArguments(ft.Params),
				})
			}
		}
		return out
	case *ast.ArrayType:
		ad := &ArrayDescriber{
			Type: "array",
			Name: in.Name.Name,
			Elt:  x.Elt.(*ast.Ident).Name,
		}
		return ad
	}
	return &noopDescriber{}
}

func (td *TypeDescriber) GetName() string {
	return td.Name
}

func (td *TypeDescriber) Describe() (string, error) { /////////////////////////////////////////////////////// {{ else if eq (len $vars) 2 }} and
	tmpl := `type {{ prepare .Name }}
	{{ if exported .Name }}
		, which is exported,
	{{ end }}
	has the following fields:
	{{ range $type, $vars := .Fields }}
		{{ range $i, $var := $vars }} {{ $var }} {{ if and (gt (len $vars) 1) (eq (sub 1 $i) (len $vars)) }}, {{ end }}{{ end }} which
		{{ if gt (len $vars) 1 }} are {{ else }} is a {{ end }}
		{{ prepare $type }}{{ if gt (len $vars) 1 }}s{{ end }}
	{{ end }}.`
	s, err := RenderTemplate(tmpl, td)
	return s, err
}

func (id *InterfaceDescriber) Describe() (string, error) {
	tmpl := `interface called {{ prepare .Name }} {{ if verbose }} which has {{ len .Methods }} {{ if eq (len .Methods) 1 }} method {{ else }} methods {{ end }} {{ end }} . `

	var mds []string
	for _, m := range id.Methods {
		s, err := m.Describe()
		if err != nil {
			return "", err
		}
		mds = append(mds, s)
	}

	s, err := RenderTemplate(tmpl, id)
	if err != nil {
		return "", err
	}

	return s + strings.Join(mds, " "), nil
}

func (id *InterfaceDescriber) GetName() string {
	return id.Name
}

func (id *InterfaceMethodDescriber) Describe() (string, error) {

	tmpl := `method {{ if verbose }} declared {{ end }} called {{ prepare .Name }} {{ if verbose }} belonging to the {{ prepare .BelongsTo }} interface {{ end }} {{ .Args }} {{ .Returns }} .`

	args, err := describeArguments(id.Args)
	if err != nil {
		return "", err
	}

	returns, err := describeReturns(id.Returns)
	if err != nil {
		return "", err
	}

	x := struct {
		Name      string
		BelongsTo string
		Args      string
		Returns   string
	}{
		Name:      id.Name,
		BelongsTo: id.BelongsTo,
		Args:      args,
		Returns:   returns,
	}

	return RenderTemplate(tmpl, x)
}

func (ad *ArrayDescriber) GetName() string {
	return ad.Name
}

func (ad *ArrayDescriber) Describe() (string, error) {
	tmpl := `Declaring type {{ prepare .Name }} which is a {{ prepare .Elt }} {{ .Type }}`
	s, err := RenderTemplate(tmpl, ad)
	if err != nil {
		return "", err
	}
	return s, nil
}
