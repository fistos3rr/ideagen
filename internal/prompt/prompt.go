// Package prompt
package prompt

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	"io/fs"
	mathRand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	ErrSystemPromptNotExists = errors.New("system prompt not found, using default system prompt")
	ErrPromptsNotExists = errors.New("prompts not found, using default prompts")
	ErrPromptsFolderNotExists = errors.New("prompts folder not found, using default prompts")
)

func IsDefaultErr(err error) bool {
	if errors.Is(err, ErrPromptsFolderNotExists) ||
	errors.Is(err, ErrPromptsNotExists) ||
	errors.Is(err, ErrSystemPromptNotExists) {
		return true
	} else {
		return false
	}
}

//go:embed templates/system.md
//go:embed templates/prompts/*.md
var templateFS embed.FS

type PromptManager struct {
	SystemPrompt *template.Template
	Prompts      []*template.Template
}

func NewPromptManager(rootPath string) (*PromptManager, error) {
	var systemFS fs.FS
	var promptsFS fs.FS

	sysPath := filepath.Join(rootPath, "templates", "system.md")
	promptsPath := filepath.Join(rootPath, "templates", "prompts")
	var sysErr error
	var promptsErr error

	if _, err := os.Stat(sysPath); err == nil {
		systemFS = os.DirFS(rootPath)
	} else {
		sysErr = ErrSystemPromptNotExists
	}

	if _, err := os.Stat(promptsPath); err == nil {
		promptsFS = os.DirFS(rootPath)
	} else {
		promptsErr = ErrPromptsFolderNotExists
	}

	if sysErr != nil {
		systemFS = templateFS
	}

	if promptsErr != nil {
		promptsFS = templateFS
	}

	manager := &PromptManager{}

	sysData, err := fs.ReadFile(systemFS, "templates/system.md")
	if err != nil {
		return nil, err
	}
	sysTmpl, err := template.New("system").Parse(string(sysData))
	if err != nil {
		return nil, err
	}
	manager.SystemPrompt = sysTmpl

	entries, err := fs.ReadDir(promptsFS, "templates/prompts")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		data, err := fs.ReadFile(promptsFS, "templates/prompts/" + entry.Name())
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New(entry.Name()).Parse(string(data))
		if err != nil {
			return nil, err
		}
		manager.Prompts = append(manager.Prompts, tmpl)
	}

	err = nil
	switch {
	case sysErr != nil && promptsErr != nil:
		err = ErrPromptsNotExists
	case sysErr != nil:
		err = sysErr
	case promptsErr != nil:
		err = promptsErr
	}

	return manager, err
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
