package interpret

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

var defaultFuncMap = sprig.TxtFuncMap()

func RenderTemplate(tmpl string, data interface{}, funcs template.FuncMap) (string, error) {
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
