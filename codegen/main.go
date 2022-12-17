package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type sqlFunctionSnippet struct {
	properties map[string]string
	parameters []string
	sql        string
	group      string
}

func (self *sqlFunctionSnippet) finish(inGroup string) {
	self.sql = strings.TrimSpace(self.sql)
	self.group = inGroup
}

func parseSQLFunctionSnippetFile(file string, group string) ([]sqlFunctionSnippet, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// convert the bytes to a string
	str := string(bytes)

	// split the string into lines
	lines := strings.Split(str, "\n")

	var snippet *sqlFunctionSnippet

	snippets := make([]sqlFunctionSnippet, 0)

	reMatchWord := regexp.MustCompile(`:\w+`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "--") {
			if snippet != nil {
				snippet.finish(group)
				snippets = append(snippets, *snippet)
				snippet = nil
			}

			snippet = &sqlFunctionSnippet{sql: "", properties: make(map[string]string), parameters: make([]string, 0)}

			err := parseSqlSnippetPropertyLine(line, snippet)
			if err != nil {
				return nil, err
			}
		} else {
			if snippet == nil {
				return nil, errors.New("sql snippet not started")
			}

			paramMatches := reMatchWord.FindAllStringSubmatch(line, -1)
			for _, paramMatch := range paramMatches {
				param := paramMatch[0]
				snippet.parameters = append(snippet.parameters, param)
			}

			snippet.sql += line + "\n"
		}
	}

	if snippet != nil {
		snippet.finish(group)
		snippets = append(snippets, *snippet)
		snippet = nil
	}

	return snippets, nil
}

func splitIntoWords(str string) []string {
	words := make([]string, 0)
	word := ""

	for i := 0; i < len(str); i++ {
		letter := str[i]
		if letter == ' ' {
			if word != "" {
				words = append(words, word)
			}

			word = ""
		} else {
			word += string(letter)
		}
	}

	words = append(words, word)

	return words
}

func parseSqlSnippetPropertyLine(line string, snippet *sqlFunctionSnippet) error {
	line = strings.TrimPrefix(line, "--")
	line = strings.TrimSpace(line)
	words := splitIntoWords(line)

	for i := 0; i < len(words); i++ {
		word := words[i]
		if word[0] == ':' {
			key := word[1:]
			value := ""

			if i+1 < len(words) {
				value = words[i+1]
				i++
			}

			snippet.properties[key] = value
		}
	}

	return nil
}

func parseSQLFunctionSnippetFolder(folder string) ([]sqlFunctionSnippet, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	snippets := make([]sqlFunctionSnippet, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.HasSuffix(fileName, ".sql") {
			filePath := folder + "/" + fileName
			group := strings.TrimSuffix(fileName, ".sql")
			fileSnippets, err := parseSQLFunctionSnippetFile(filePath, group)

			if err != nil {
				return nil, err
			}

			snippets = append(snippets, fileSnippets...)
		}
	}

	return snippets, nil
}

func main() {
	snippets, err := parseSQLFunctionSnippetFolder("example-sql")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(snippets)
}