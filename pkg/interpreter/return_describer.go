package interpret

import "go/ast"

type ReturnDescriber struct {
	original *ast.ReturnStmt
}

func NewReturnDescriber(in *ast.ReturnStmt) Describer {
	return &ReturnDescriber{original: in}
}

func (rd *ReturnDescriber) Describe() (string, error) {
	var returnNames []string
	for _, r := range rd.original.Results {
		switch s := r.(type) {
		case *ast.Ident:
			returnNames = append(returnNames, s.Name)
		}
	}

	tmpl := `
	returning {{ range $_, $rv := . }} {{ . }} {{ end }}  
	`

	return RenderTemplate(tmpl, returnNames)
}

func (rd *ReturnDescriber) GetName() string {
	return ""
}
