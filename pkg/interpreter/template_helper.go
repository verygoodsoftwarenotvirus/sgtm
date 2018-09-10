package interpret

import (
	"bytes"
	"text/template"
)

func RenderTemplate(tmpl string, data interface{}, funcs template.FuncMap) (string, error) {
	t, err := template.New("t1").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
