package main

import (
	"bytes"
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
	"text/template"
)

type SqlSnippetTemplateConfig struct {
	StaticFiles   []string
	GroupFile     string
	GroupFileName string
	Name          string
}

type SqlSnippetTemplate struct {
	Config *SqlSnippetTemplateConfig
	Path   string
}

func (templ *SqlSnippetTemplate) executeTemplateFile(file string, data any) (string, error) {
	text, err := templ.getFileContent(file)
	if err != nil {
		return "", err
	}

	// Create a new template and parse the letter into it.
	tpl, err := template.New(file).Parse(text)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = tpl.Execute(&output, data)
	return string(output.Bytes()), nil
}

func (templ *SqlSnippetTemplate) getFileContent(file string) (string, error) {
	bytes, err := os.ReadFile(templ.Path + "/" + file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}
	return false
}

func (templ *SqlSnippetTemplate) executeTemplate(outputFolder string, snippets []SqlFunctionSnippet) error {
	// generate all the static files
	for _, static := range templ.Config.StaticFiles {
		output, err := templ.executeTemplateFile(static, snippets)
		if err != nil {
			return err
		}
		err = os.WriteFile(outputFolder+"/"+static, []byte(output), 0644)
		if err != nil {
			return err
		}
	}

	// collect all the groups
	groups := make([]string, 0)
	for _, snippet := range snippets {
		if !contains(groups, snippet.Groupt) {
			groups = append(groups, snippet.Groupt)
		}
	}

	// generate the group file for each group
	for _, group := range groups {
		// collect all the snippets for this group
		snippetsOfGroup := make([]SqlFunctionSnippet, 0)
		for _, snippet := range snippets {
			if snippet.Groupt == group {
				snippetsOfGroup = append(snippetsOfGroup, snippet)
			}
		}

		output, err := templ.executeTemplateFile(templ.Config.GroupFile, snippetsOfGroup)
		if err != nil {
			return err
		}

		outputFilename := strings.ReplaceAll(templ.Config.GroupFileName, "{{group}}", group)
		err = os.WriteFile(outputFolder+"/"+outputFilename, []byte(output), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func readTemplateConfig(file string) (*SqlSnippetTemplateConfig, error) {
	// read file as bytes
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// parse the bytes as a TOML file
	var cfg SqlSnippetTemplateConfig
	err = toml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadTemplate(path string) (*SqlSnippetTemplate, error) {
	configFile := path + "/config.toml"
	var template SqlSnippetTemplate
	config, err := readTemplateConfig(configFile)
	if err != nil {
		return nil, err
	}
	template.Config = config
	template.Path = path
	return &template, nil
}
