package interpret

import (
	"fmt"
	"go/ast"
)

func parseExprStmt(in *ast.ExprStmt) (string, error) {
	switch e := in.X.(type) {
	case *ast.CallExpr:
		// parse call expr
		return NewCallExpr(e).Describe()
	}
	return "", nil
}

func parseAssignStmt(in *ast.AssignStmt) (out string, err error) {
	var leftHandSide []string
	for _, v := range in.Lhs {
		switch l := v.(type) {
		case *ast.Ident:
			varName := l.Name
			leftHandSide = append(leftHandSide, varName)
		case *ast.SelectorExpr:
			if m, ok := l.X.(*ast.Ident); ok {
				varName := fmt.Sprintf("%s dot %s", prepareName(m.Name), prepareName(l.Sel.Name))
				leftHandSide = append(leftHandSide, varName)
			} else {
				print()
			}
		default:
			print()
		}
	}

	for j := range in.Rhs {
		switch k := in.Rhs[j].(type) {
		case *ast.CallExpr:
			/*
				Alright ya'll, here's where it gets real: this is a hackathon project, and while the Go programming
				language, which I love with all of my dear little heart, allows you to do some very interesting things
				with N operands on either side of an assign statement, we're going to assume that you only have one
				such function call for the sake of time and simplicity. Rob Pike, you won't read this, but I hope you
				appreciate the sentiments here.
			*/
			if len(in.Lhs) != len(in.Rhs) {
				println()
			}

			var (
				args     []ArgDesc
				funcName string
			)
			for i, x := range in.Rhs {
				switch y := x.(type) {
				case *ast.CallExpr:
					switch z := y.Fun.(type) {
					case *ast.Ident:
						funcName = z.Name
						for _, a := range y.Args {
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
								if len(in.Rhs) == len(in.Lhs) {
									ad.Names = append(ad.Names, leftHandSide[i])
								}

								args = append(args, ad)
							}
						}
					}
				}
			}
			s, err := describeCallExpr(funcName, args, leftHandSide)
			if err != nil {
				return "", err
			}
			out += s

		case *ast.CompositeLit:
			var varName string
			if leftHandSide != nil && len(in.Rhs) == len(in.Lhs) {
				if len(leftHandSide)-1 < j {
					println()
				}

				varName = leftHandSide[j]
			}

			s, err := NewCompositeLiteralDescriber(k, varName).Describe()
			if err != nil {
				return "", err
			}
			out += s
		}
	}
	return
}
