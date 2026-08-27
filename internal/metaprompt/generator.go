// Package metaprompt provides generator for meta prompts
package metaprompt

import (
	"embed"
	"math/rand"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*.yaml templates/dicts/*.yaml
var fs embed.FS

var (
	templateStrings []string
	dicts map[string][]string
	rng *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	loadTemplates()
	loadDicts()
}

func loadTemplates() {
	data, err := fs.ReadFile("templates/templates.yaml")
	if err != nil {
		panic("failed to read templates.yaml: " + err.Error())
	}
	if err := yaml.Unmarshal(data, &templateStrings); err != nil {
		panic("failed to parse templates.yaml: " + err.Error())
	}
	if len(templateStrings) == 0 {
		panic("templates.yaml is empty")
	}
}

func loadDicts() {
	dicts = make(map[string][]string)
	dictNames := []string{"verbs", "adjectives", "styles"}

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
		dicts[name] = list
	}
}

func randomFrom(dictName string) string {
	list, ok := dicts[dictName]
	if !ok || len(list) == 0 {
		return "undefined"
	}

	return list[rng.Intn(len(list))]
}

func GenerateMetaPrompt(topic, requirements) string {
	idx := rng.Intn(len(templateStrings))
	tmplStr := templateStrings[idx]

	data := map[string]string{
		"Verb": randomFrom("verbs"),
		"Adjective": randomFrom("adjectives"),
		"Style": randomFrom("styles"),
		"Topic": topic,
	}

	tmpl, err := template.New("meta").Parse(tmplStr)
	if err != nil {
		return tmplStr + "\n\n" + requirements
	}
	var result string
	if err := tmpl.Execute(&result, data); err != nil {
		return tmplStr + "\n\n" + requirements
	}
	return result + "\n\n" + requirements
}
