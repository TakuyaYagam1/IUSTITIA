package docgen

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type DocumentContext struct {
	CaseID    string
	Defendant string
	Verdict   string
	Details   string
}

func (d *DocumentContext) ReadClassified(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}

func Generate(userTemplate string, ctx *DocumentContext) (string, error) {
	tmpl, err := template.New("doc").Parse(userTemplate)
	if err != nil {
		return "", fmt.Errorf("docgen - Generate - Parse: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("docgen - Generate - Execute: %w", err)
	}
	return buf.String(), nil
}
