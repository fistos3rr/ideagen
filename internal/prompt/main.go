package prompt

import (
	"embed"
	"text/template"
	"bytes"
)

//go:embed templates/system.md
//go:embed templates/prompts/*.md
var templateFS embed.FS

type PromptManager struct {
	SystemPrompt *template.Template
	Prompts      []*template.Template
}

func NewDefaultPromptManager() *PromptManager {
	manager := &PromptManager{}

	sysData, err := templateFS.ReadFile("templates/system.md")
	if err != nil {
		panic(err)
	}
	sysTmpl, err := template.New("system").Parse(string(sysData))
	if err != nil {
		panic(err)
	}
	manager.System = sysTmpl

	entries, err := templateFS.ReadDir("templates/prompts")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := templatesFS.ReadFile("templates/prompts/" + entry.Name())
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New(entry.Name()).Parse(string(data))
		if err != nil {
			return nil, err
		}
		manager.Prompts = append(manager.Prompts, tmpl)
	}

	return manager
}

type TemplateData struct {
	Type string
}

func executeTemplate(tmpl *template.Template, data TemplateData) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func writeSalt(&bytes.Buffer) error {

}
