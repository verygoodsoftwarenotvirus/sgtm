package interpret

import (
	"go/ast"
)

type CompositeLiteralDescriber struct {
	AssignedName string
	original     *ast.CompositeLit
}

func NewCompositeLiteralDescriber(in *ast.CompositeLit, name string) Describer {
	return &CompositeLiteralDescriber{original: in, AssignedName: name}
}

func (c *CompositeLiteralDescriber) Describe() (out string, err error) {
	tmpl := `
	{{ if startsWithVowel .Type }}an {{ else }} a {{ end }} {{ prepare .Type }} literal {{ if verbose }} named {{ prepare .AssignedName }} {{ end }}  with the following constructing values: {{ $last := sub (len .Args) 1 }} {{ range $i, $desc := .Args }} {{ if eq $last $i }} and {{ end }} {{ $desc }} {{ if ne $last $i }} , {{ end }} {{ end }}.
	`

	var eltsDescriptions []string
	for _, a := range c.original.Elts {
		switch b := a.(type) {
		case *ast.BasicLit:
			s, err := NewBasicLitDescriber(b, nil).Describe()
			if err != nil {
				return "", err
			}
			eltsDescriptions = append(eltsDescriptions, s)
		}
	}

	var x struct {
		Type         string
		AssignedName string
		Args         []string
	}

	switch t := c.original.Type.(type) {
	case *ast.Ident:
		x = struct {
			Type         string
			AssignedName string
			Args         []string
		}{
			Type:         t.Name,
			AssignedName: c.AssignedName,
			Args:         eltsDescriptions,
		}
		return RenderTemplate(tmpl, x)
	default:
		print()
	}
	return "", nil
}

func (c *CompositeLiteralDescriber) GetName() string {
	return ""
}
