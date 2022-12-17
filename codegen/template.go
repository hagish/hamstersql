package main

import (
	"bytes"
	"fmt"
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

func (templ *SqlSnippetTemplate) executeTemplateFile(file string, data any, verbose bool) (string, error) {
	if verbose {
		fmt.Println("Executing template", file)
	}
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

func (templ *SqlSnippetTemplate) writeTemplateOutput(outputFile string, content []byte, verbose bool) error {
	if verbose {
		fmt.Println("Writing template output", outputFile)
	}

	// check if file exists
	_, err := os.Stat(outputFile)
	if err == nil {
		// file exists, so we check for actual changes
		existingContent, err := os.ReadFile(outputFile)
		if err != nil {
			return err
		}
		if bytes.Equal(existingContent, content) {
			if verbose {
				fmt.Println("No changes, skipping", outputFile)
			}
			// no changes, so we can skip
			return nil
		}
	}

	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (templ *SqlSnippetTemplate) executeTemplate(outputFolder string, snippets []SqlFunctionSnippet, verbose bool) error {
	// generate all the static files
	for _, static := range templ.Config.StaticFiles {
		output, err := templ.executeTemplateFile(static, snippets, verbose)
		if err != nil {
			return err
		}
		err = templ.writeTemplateOutput(outputFolder+"/"+static, []byte(output), verbose)
		if err != nil {
			return err
		}
	}

	// collect all the groups
	groups := make([]string, 0)
	for _, snippet := range snippets {
		if !contains(groups, snippet.Group) {
			groups = append(groups, snippet.Group)
		}
	}

	// generate the group file for each group
	for _, group := range groups {
		// collect all the snippets for this group
		snippetsOfGroup := make([]SqlFunctionSnippet, 0)
		for _, snippet := range snippets {
			if snippet.Group == group {
				snippetsOfGroup = append(snippetsOfGroup, snippet)
			}
		}

		output, err := templ.executeTemplateFile(templ.Config.GroupFile, snippetsOfGroup, verbose)
		if err != nil {
			return err
		}

		outputFilename := strings.ReplaceAll(templ.Config.GroupFileName, "{{group}}", group)
		err = templ.writeTemplateOutput(outputFolder+"/"+outputFilename, []byte(output), verbose)
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

func loadTemplate(path string, verbose bool) (*SqlSnippetTemplate, error) {
	configFile := path + "/config.toml"
	if verbose {
		fmt.Println("Parsing template config", configFile)
	}
	var template SqlSnippetTemplate
	config, err := readTemplateConfig(configFile)
	if err != nil {
		return nil, err
	}
	template.Config = config
	template.Path = path
	return &template, nil
}
