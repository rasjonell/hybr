package nginx

import (
	"bytes"
	"fmt"
	"text/template"
)

func parseTemplate(name string, data any) ([]byte, error) {
	content, err := templatesFS.ReadFile("templates/" + name + ".conf.templ")
	if err != nil {
		return nil, fmt.Errorf("Failed to read template %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("Failed to execure template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}
