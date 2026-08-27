// Package metaprompt provides generator for meta prompts
package metaprompt

import (
	"embed"
	"math/rand"
	"text/template"
	"time"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/luxfi/go-bip39"
)

//go:embed templates/*.yaml templates/dicts/*.yaml
var fs embed.FS

type MetaPromptGenerator struct {
	rng *rand.Rand
	dicts map[string][]string
	templateStrings []string
	requirements []string
}

func NewMetaPromptGenerator() *MetaPromptGenerator {
	generator := new(MetaPromptGenerator)

	generator.init()

	return generator
}


func (g *MetaPromptGenerator) init() {
	g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.loadTemplates()
	g.loadDicts()
}

func (g *MetaPromptGenerator) loadTemplates() {
	data, err := fs.ReadFile("templates/templates.yaml")
	if err != nil {
		panic("failed to read templates.yaml: " + err.Error())
	}
	if err := yaml.Unmarshal(data, &g.templateStrings); err != nil {
		panic("failed to parse templates.yaml: " + err.Error())
	}
	if len(g.templateStrings) == 0 {
		panic("templates.yaml is empty")
	}

	data, err = fs.ReadFile("templates/requirements.yaml")
	if err != nil {
		g.requirements = []string{}
		return
	}
	if err := yaml.Unmarshal(data, &g.requirements); err != nil {
		panic("failed to parse requirements.yaml: " + err.Error())
	}
}

func (g *MetaPromptGenerator) loadDicts() {
	g.dicts = make(map[string][]string)
	dictNames := []string{}

	for _, name := range dictNames {
		path := "templates/dicts/" + name + ".yaml"
		data, err := fs.ReadFile(path)
		if err != nil {
			panic("failed to read " + path + ": " + err.Error())
		}
		var list []string
		if err := yaml.Unmarshal(data, &list); err != nil {
			panic("failed to parse " + path + ": " + err.Error())
		}
		g.dicts[name] = list
	}
}

func (g *MetaPromptGenerator) randomFrom(dictName string) string {
	list, ok := g.dicts[dictName]
	if !ok || len(list) == 0 {
		return "undefined"
	}

	return list[g.rng.Intn(len(list))]
}

func (g *MetaPromptGenerator) GenerateMetaPrompt(topic string) string {
	idx := g.rng.Intn(len(g.templateStrings))
	tmplStr := g.templateStrings[idx]

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}

	data := map[string]string{
		"Verb": g.randomFrom("verbs"),
		"Adjective": g.randomFrom("adjectives"),
		"Style": g.randomFrom("styles"),
		"Topic": topic,
		"Seed": mnemonic,
	}

	tmpl, err := template.New("meta").Parse(tmplStr)
	if err != nil {
		return g.withMetadata(tmplStr)
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return g.withMetadata(tmplStr)
	}
	return g.withMetadata(sb.String())
}

func (g *MetaPromptGenerator) withMetadata(tmpl string) string {
	if len(g.requirements) == 0 {
		return tmpl
	}

	var sb strings.Builder
	sb.WriteString(tmpl)
	sb.WriteString("\n\n")
	sb.WriteString("Important:\n")
	for _, r := range g.requirements {
		sb.WriteString(r)
		sb.WriteString("\n")
	}

	return sb.String()
}
