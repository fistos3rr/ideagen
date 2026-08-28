// Package prompt
package prompt

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	mathRand "math/rand"
	"text/template"
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
	manager.SystemPrompt = sysTmpl

	entries, err := templateFS.ReadDir("templates/prompts")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := templateFS.ReadFile("templates/prompts/" + entry.Name())
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New(entry.Name()).Parse(string(data))
		if err != nil {
			panic(err)
		}
		manager.Prompts = append(manager.Prompts, tmpl)
	}

	return manager
}

func (m *PromptManager) GetSystemPrompt(t string) (string, error) {
	data := TemplateData{
		Type: t,
	}

	buf, err := executeTemplate(m.SystemPrompt, data)
	if err != nil {
		return "", err
	}

	err = writeSalt(buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *PromptManager) GetRandPrompt(t string) (string, error) {
	if len(m.Prompts) == 0 {
		return "", errors.New("prompts slice empty")
	}

	idx := mathRand.Intn(len(m.Prompts))
	pr := m.Prompts[idx]

	data := TemplateData{
		Type: t,
	}

	buf, err := executeTemplate(pr, data)
	if err != nil {
		return "", err
	}

	err = writeSalt(buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *PromptManager) GetPrompts(t string) (string, string, error) {
	sysPr, err := m.GetSystemPrompt(t)
	if err != nil {
		return "", "", err
	}
	pr, err := m.GetRandPrompt(t)
	if err != nil {
		return sysPr, "", err
	}

	return sysPr, pr, nil
}

type TemplateData struct {
	Type string
}

func executeTemplate(tmpl *template.Template, data TemplateData) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func writeSalt(buf *bytes.Buffer) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	buf.WriteString("\n[")
	buf.WriteString(hex.EncodeToString(salt))
	buf.WriteString("]")

	return nil
}
