package interpret

import (
	"fmt"
	"go/ast"
)

type ArgDesc struct {
	Literal bool
	Pointer bool
	Type    string
	Name    string
	Names   []string
	Value   string
}

type FuncDecl struct {
	original           *ast.FuncDecl
	Name               string
	ReceiverData       *ReceiverData
	ParameterArguments []ArgDesc
	ReturnArguments    []ArgDesc
}

func NewFuncDecl(f *ast.FuncDecl) Describer {
	funcDecl := &FuncDecl{
		original:           f,
		Name:               f.Name.Name,
		ReceiverData:       parseReceiver(f.Recv),
		ParameterArguments: parseArguments(f.Type.Params),
		ReturnArguments:    parseArguments(f.Type.Results),
	}

	return funcDecl
}

func parseReceiver(in *ast.FieldList) *ReceiverData {
	if in == nil || in.List == nil || len(in.List) != 1 {
		return nil
	}

	r := in.List[0]
	rd := &ReceiverData{
		Name: r.Names[0].Name,
	}
	switch x := r.Type.(type) {
	case *ast.StarExpr:
		rd.Pointer = true
		if y, ok := x.X.(*ast.Ident); ok {
			rd.TypeName = y.Name
		}
	case *ast.Ident:
		rd.TypeName = x.Name
	}
	return rd
}

func parseArguments(in *ast.FieldList) []ArgDesc {
	var out []ArgDesc
	if in != nil {
		for _, t := range in.List {
			switch u := t.Type.(type) {
			case *ast.Ident:
				var names []string
				for _, n := range t.Names {
					names = append(names, prepareName(n.Name))
				}

				out = append(out, ArgDesc{Type: prepareName(u.Name), Names: names})
			case *ast.StarExpr:
				switch v := u.X.(type) {
				case *ast.Ident:
					name := fmt.Sprintf("%s dot %s", v.Name, v.Name)
					out = append(out, ArgDesc{Type: prepareName(name), Names: []string{"REPLACE ME"}})
				case *ast.SelectorExpr:
					if w, ok := v.X.(*ast.Ident); ok {
						name := fmt.Sprintf("%s dot %s", w.Name, v.Sel.Name)
						out = append(out, ArgDesc{Type: prepareName(name), Pointer: true, Names: []string{"REPLACE ME"}})
					}
				default:
					println() // TODO
				}
			case *ast.ArrayType:
				ad := ArgDesc{
					Type: fmt.Sprintf(" array of %s ", u.Elt.(*ast.Ident).Name),
					//Name: in.Name.Name,
					//Elt:  x.Elt.(*ast.Ident).Name,
				}
				out = append(out, ad)
			default:
				println() // TODO
			}
		}
	}
	return out
}

type ReceiverData struct {
	Pointer  bool
	Name     string
	TypeName string
}

func describeReceiver(in *ReceiverData) (string, error) {
	if in == nil {
		return "", nil
	}

	tmpl := `
	, which is attached to {{ if .Pointer }} a pointer to {{ else }} an instance of {{ end }} {{ prepare .TypeName }} {{ if verbose }} called {{ .Name }} {{ end }} .

	`

	return RenderTemplate(tmpl, in)
}

func describeArguments(in []ArgDesc) (string, error) {
	tmpl := `
	accepting
	{{ if not .Arguments }} nothing {{ else}}
		{{ range $i, $arg := .Arguments}}
			{{ if ne $i 0 }} and {{ end }}
			{{ if lt (len $arg.Names) 2 }}
				{{ if and (startsWithVowel $arg.Type) (not $arg.Pointer) }} an {{ else }} a {{ if $arg.Pointer }} pointer to {{ end }}{{ end }} {{ $arg.Type }}
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

	s := struct {
		Arguments []ArgDesc
	}{
		Arguments: in,
	}

	return RenderTemplate(tmpl, s)
}

func describeReturns(in []ArgDesc) (string, error) {
	tmpl := `
	returning
	{{ if not .Arguments }} nothing {{ else}}
		{{ range $i, $arg := .Arguments}}
			{{ if ne $i 0 }} and {{ end }}
			{{ if lt (len $arg.Names) 2 }}
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

	s := struct {
		Arguments []ArgDesc
	}{
		Arguments: in,
	}

	return RenderTemplate(tmpl, s)
}

func use(...interface{}) {}

func describeBody(in *ast.FuncDecl) (out string, err error) {
	if in.Body != nil {
		for _, s := range in.Body.List {
			switch x := s.(type) {
			case *ast.AssignStmt:
				o, err := parseAssignStmt(x)
				if err != nil {
					return out, err
				}
				out += o
			case *ast.ExprStmt:
				o, err := parseExprStmt(x)
				if err != nil {
					return out, err
				}
				out += o
			case *ast.ReturnStmt:
				o, err := NewReturnDescriber(x).Describe()
				if err != nil {
					return out, err
				}
				out += o
			default:
				print()
			}
		}
	}
	return
}

func (f FuncDecl) GetName() string {
	return f.Name
}

func (f FuncDecl) Describe() (string, error) {
	argstmt, err := describeArguments(f.ParameterArguments)
	if err != nil {
		return "", err
	}

	recvstmt, err := describeReceiver(f.ReceiverData)
	if err != nil {
		return "", err
	}

	retstmt, err := describeReturns(f.ReturnArguments)
	if err != nil {
		return "", err
	}

	bodystmt, err := describeBody(f.original)
	if err != nil {
		return "", err
	}

	tmpl := `function declared called {{ prepare .Name }} {{ .Receivers }} {{ .Args }} {{ .Returns }}. {{ if ne (len .Body) 0 }} The body contains {{ .Body }} {{ end }}`

	x := struct {
		Name      string
		Receivers string
		Args      string
		Returns   string
		Body      string
	}{
		Name:      f.Name,
		Receivers: recvstmt,
		Args:      argstmt,
		Returns:   retstmt,
		Body:      bodystmt,
	}

	return RenderTemplate(tmpl, x)
}
