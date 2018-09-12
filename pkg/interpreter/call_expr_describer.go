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
	{{ .PkgName }} dot {{ prepare .MethodName }} {{ if and verbose (ne (len .Arguments) 0) }} called with the argument {{ .Arguments }} {{ end }}
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

	return RenderTemplate(tmpl, y)
}
