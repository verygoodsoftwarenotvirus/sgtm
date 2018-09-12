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
	return "", nil
}
