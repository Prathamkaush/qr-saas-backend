package templates

import (
	"bytes"
	"html/template"
)

func RenderVCard(t *Template, inst *TemplateInstance) (string, error) {
	tpl, err := template.ParseFiles("internal/templates/views/vcard.html")
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, inst.Data); err != nil {
		return "", err
	}

	return out.String(), nil
}

func RenderGeneric(t *Template, inst *TemplateInstance) (string, error) {
	return "<h1>Template</h1>", nil
}

func RenderSocial(t *Template, inst *TemplateInstance) (string, error) {
	return "<h1>Social Template</h1>", nil
}

func RenderEvent(t *Template, inst *TemplateInstance) (string, error) {
	return "<h1>Event Template</h1>", nil
}
