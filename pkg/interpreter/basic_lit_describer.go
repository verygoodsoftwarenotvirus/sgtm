package interpret

import (
	"go/ast"
	"go/token"
)

var typeTokenMap = map[token.Token]string{
	token.IDENT:  "identifier",
	token.INT:    "integer",
	token.FLOAT:  "float",
	token.IMAG:   "imaginary number",
	token.CHAR:   "character",
	token.STRING: "string",
}

type BasicLitDescriber struct {
	original *ast.BasicLit
	varNames []string
}

func NewBasicLitDescriber(bs *ast.BasicLit, varNames []string) *BasicLitDescriber {
	sp := &BasicLitDescriber{original: bs, varNames: varNames}
	return sp
}

func (d *BasicLitDescriber) GetName() string {
	return ""
}

func (d *BasicLitDescriber) Describe() (string, error) {
	tmpl := `
	{{ if verbose }}
		{{ if startsWithVowel .Type }} an {{ else }} a {{ end }} 
		{{ .Type }} literal with a value of 
	{{ end }} {{ .Value }}
	`

	x := ArgDesc{
		Literal: true,
		Type:    typeTokenMap[d.original.Kind],
		Value:   d.original.Value,
	}

	s, err := RenderTemplate(tmpl, x)

	return s, err
}
