package interpret

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var defaultFuncMap = sprig.TxtFuncMap()

func init() {
	defaultFuncMap["startsWithVowel"] = startsWithVowel
	defaultFuncMap["verbose"] = func() bool { return currentVerbosity == HighVerbosity }
	defaultFuncMap["prepare"] = prepareName
	defaultFuncMap["exported"] = func(s string) bool {
		if len(s) == 0 {
			return false
		}
		return s[0] == strings.ToUpper(s)[0]
	}
}

func RenderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("t").Funcs(defaultFuncMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

var _ Describer = (*noopDescriber)(nil)

type noopDescriber struct{}

func (n *noopDescriber) GetName() string {
	return ""
}

func (n *noopDescriber) Describe() (string, error) {
	return "", nil
}
