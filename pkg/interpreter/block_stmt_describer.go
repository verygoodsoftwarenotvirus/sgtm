package interpret

import (
	"go/ast"
)

type StatementParser struct {
	original []ast.Stmt
}

func NewStatementParser(bs []ast.Stmt) *StatementParser {
	sp := &StatementParser{
		original: bs,
	}

	return sp
}

func (i *interpreter) parseExprStmt(in *ast.ExprStmt) (string, error) {
	switch e := in.X.(type) {
	case *ast.CallExpr:
		// parse call expr
		return NewCallExpr(e).Describe()
	}
	return "", nil
}

func (i *interpreter) parseAssignStmt(in *ast.AssignStmt) (string, error) {
	var leftHandSide []string
	for i := range in.Lhs {
		if l, ok := in.Lhs[i].(*ast.Ident); ok {
			varName := l.Name
			leftHandSide = append(leftHandSide, varName)
		}
	}

	//var funcName string
	for j := range in.Rhs {
		switch in.Rhs[j].(type) {
		case *ast.CallExpr:
			/*
				Alright ya'll, here's where it gets real: this is a hackathon project, and while the Go programming
				language, which I love with all of my dear little heart, allows you to do some very interesting things
				with N operands on either side of an assign statement, we're going to assume that you only have one
				such function call for the sake of time and simplicity. Rob Pike, you won't read this, but I hope you
				appreciate the sentiments here.
			*/
			x := in.Rhs[0].(*ast.CallExpr)
			var args []ArgDesc
			for _, a := range x.Args {
				switch b := a.(type) {
				case *ast.Ident:
					ad := ArgDesc{
						Names: []string{b.Name},
					}
					args = append(args, ad)
				case *ast.BasicLit:
					ad := ArgDesc{
						Literal: true,
						Type:    typeTokenMap[b.Kind],
						Value:   prepareName(b.Value),
					}
					args = append(args, ad)
				}
			}

			return describeCallExpr(x.Fun.(*ast.Ident).Name, args, leftHandSide)
		}
	}
	return "", nil
}

func describeCallExpr(name string, args []ArgDesc, returns []string) (string, error) {
	tmpl := `
	a call to {{ .FuncName }}
		with arguments
		{{ $argCount := sub (len .Args) 1 }}
		{{ range $i, $arg := .Args }}
			{{ if $arg.Literal }}
				{{ if startsWithVowel $arg.Type }}an {{ else }}a {{ end }} {{ prepare $arg.Type }} literal with the value {{ prepare $arg.Value }}
			{{ else }}
				{{ range $_, $n := $arg.Names }} {{ $n }} {{ end }}
			{{ end }}
			{{ if ne $argCount $i }} and {{ end }}
		{{ end }}.
	assigning result{{ if gt (len $.Args) 1 }}s{{ end }} to
		{{ $varCount := sub (len .VarNames) 1 }}
		{{ range $i, $arg := .VarNames }}
			{{ $arg }} {{ if ne $varCount $i }} and {{ end }}
		{{ end }}`

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
