package templating

import (
	"bytes"
	"html/template"
	"io"
	"strings"
)

type TemplateParameters map[string]string

func PrepareCommandTemplate(payloadTemplate string, params *TemplateParameters) (io.Reader, error) {

	t := template.Must(template.New("payload").Parse(payloadTemplate))
	var b bytes.Buffer

	err := t.Execute(&b, params)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(b.String()), nil
}
