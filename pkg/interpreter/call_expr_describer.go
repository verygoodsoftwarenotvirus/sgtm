package interpret

import (
	"go/ast"
)

type CallExpr struct {
	original *ast.CallExpr
}

func NewCallExpr(og *ast.CallExpr) *CallExpr {
	sp := &CallExpr{
		original: og,
	}

	return sp
}

func (d *CallExpr) GetName() string {
	return ""
}

func (d *CallExpr) Describe() (string, error) {
	tmpl := `
	{{ prepare .PkgName }} dot {{ prepare .MethodName }}
		{{ if verbose }} called with the following argument{{ if gt (len .Arguments) 1 }}s {{ end }} {{ end }}
	{{ range $i, $arg := .Arguments }} 
		{{ if startsWithVowel $arg.Type }}an {{ else }} a {{ end }} 
		{{ $arg.Type }}, {{ $arg.Value }}, 
	{{ end }}
	`

	var y struct {
		PkgName    string
		MethodName string
		Arguments  []ArgDesc
	}

	switch x := d.original.Fun.(type) {
	case *ast.SelectorExpr:
		pkgName := x.X.(*ast.Ident).Name
		methodName := x.Sel.Name

		y = struct {
			PkgName    string
			MethodName string
			Arguments  []ArgDesc
		}{
			PkgName:    pkgName,
			MethodName: methodName,
			Arguments:  nil,
		}
	}

	for _, arg := range d.original.Args {
		switch a := arg.(type) {
		case *ast.BasicLit:
			ad := ArgDesc{
				Type:  typeTokenMap[a.Kind],
				Value: a.Value,
			}
			y.Arguments = append(y.Arguments, ad)
		}
	}

	s, err := RenderTemplate(tmpl, y)
	return s, err
}

func describeCallExpr(name string, args []ArgDesc, returns []string) (string, error) {
	tmpl := `
	a call to {{ .FuncName }}
		with arguments
		{{ $argCount := sub (len .Args) 1 }}
		{{ range $i, $arg := .Args }}
			{{ if $arg.Literal }}
				{{ if startsWithVowel $arg.Type }}an {{ else }}a {{ end }}
				{{ prepare $arg.Type }} literal with the value {{ prepare $arg.Value }}
			{{ else }}
				{{ range $_, $n := $arg.Names }} {{ $n }} {{ end }}
			{{ end }}
			{{ if ne $argCount $i }} and {{ end }}
		{{ end }}.
	assigning result{{ if gt (len $.Args) 1 }}s{{ end }} to
		{{ $varCount := sub (len .VarNames) 1 }}
		{{ range $i, $arg := .VarNames }}
			{{ $arg }} {{ if ne $varCount $i }} and {{ end }}
		{{ end }}.`

	x := struct {
		FuncName string
		Args     []ArgDesc
		VarNames []string
	}{
		FuncName: name,
		Args:     args,
		VarNames: returns,
	}

	s, err := RenderTemplate(tmpl, x)
	return s, err
}
