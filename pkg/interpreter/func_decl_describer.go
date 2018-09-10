package interpret

import (
	"fmt"
	"strings"
	"text/template"
)

type ArgDesc struct {
	Type  string
	Names []string
}

type FuncDecl struct {
	Name               string
	ParameterArguments []ArgDesc
	ReturnArguments    []ArgDesc
}

func (f FuncDecl) Describe() (string, error) {
	argstmt, err := f.describeArguments()
	if err != nil {
		return "", err
	}

	retstmt, err := f.describeReturns()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`function declared called %s  %s  %s.`, f.Name, argstmt, retstmt), nil
}

func (f FuncDecl) describeArguments() (string, error) {
	var out = " accepting "
	if f.ParameterArguments == nil {
		out += "nothing. "
		return out, nil
	}

	for i, r := range f.ParameterArguments {
		if i != 0 {
			out += " and "
		}
		if startsWithVowel(r.Type) {
			out += fmt.Sprintf("an %s ", r.Type)
		} else {
			if len(r.Names) == 1 {
				out += fmt.Sprintf("a %s ", r.Type)
				if r.Names[0] != "" {
					out += fmt.Sprintf(" called %s", r.Names[0] )
				}
			} else {
				out += fmt.Sprintf(" %ss ", r.Type)
			}
		}

		if len(r.Names) > 2 {
			out += strings.Join(r.Names, ", and ")
		} else  if len(r.Names) == 2 {
			out += strings.Join(r.Names,  " and ")
		}
		out += ", "
	}
	return out, nil
}

func (f FuncDecl) describeReturns() (string, error) {
	var out = " returning "
	if f.ReturnArguments == nil {
		out += "nothing. "
		return out, nil
	}

	for _, r := range f.ReturnArguments {
		if startsWithVowel(r.Type) {
			out += fmt.Sprintf("an %s ", r.Type)
		} else {
			out += fmt.Sprintf("a %s ", r.Type)
		}

		out += strings.Join(r.Names, ", and ")
	}
	return out, nil
}

func (f FuncDecl) describeBody() (string, error) {
	return "", nil

}

func (f FuncDecl) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"dec": func(i int) int {
			i--
			return i
		},
		"returns": func() bool {
			answer := f.ReturnArguments != nil
			return answer
		},
		"accepts": func() bool {
			answer := f.ParameterArguments != nil
			return answer
		},
		"singular": func(ad ArgDesc) string {
			if len(ad.Names) <= 1 {
				if startsWithVowel(ad.Type) {
					return "an"
				}
				return "a"
			}
			return ""
		},
		"plural": func(ad ArgDesc) string {
			if len(ad.Names) >= 2 {
				return "s"
			}
			return ""
		},
		"joinargs": func(ad ArgDesc) string {
			if ad.Names == nil {
				return ""
			}

			if len(ad.Names) == 2 {
				return fmt.Sprintf("called %s", strings.Join(ad.Names, " and "))
			}

			return fmt.Sprintf("called %s", strings.Join(ad.Names, ", and "))
		},
	}
}
