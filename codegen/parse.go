package main

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type SqlSnippetReturnType int

const (
	SqlSnippetReturnTypeOneRow SqlSnippetReturnType = iota
	SqlSnippetReturnTypeManyRows
	SqlSnippetReturnTypeAffectedRows
	SqlSnippetReturnTypeScalar
	SqlSnippetReturnTypeInsertID
	SqlSnippetReturnTypeNone
)

func (t SqlSnippetReturnType) String() string {
	return [...]string{"OneRow", "ManyRows", "AffectedRows", "Scalar", "InsertID", "None"}[t]
}

type SqlFunctionSnippet struct {
	Properties    map[string]string
	Parameters    []string
	SqlQuery      string
	Groupt        string
	ReturnType    SqlSnippetReturnType
	ReturnTypeStr string
	Doc           string
}

func (snippet *SqlFunctionSnippet) finish(inGroup string) {
	snippet.SqlQuery = strings.TrimSpace(snippet.SqlQuery)
	snippet.Groupt = inGroup

	if _, ok := snippet.Properties["one"]; ok {
		snippet.ReturnType = SqlSnippetReturnTypeOneRow
	} else if _, ok := snippet.Properties["many"]; ok {
		snippet.ReturnType = SqlSnippetReturnTypeManyRows
	} else if _, ok := snippet.Properties["affected"]; ok {
		snippet.ReturnType = SqlSnippetReturnTypeAffectedRows
	} else if _, ok := snippet.Properties["scalar"]; ok {
		snippet.ReturnType = SqlSnippetReturnTypeScalar
	} else if _, ok := snippet.Properties["insert"]; ok {
		snippet.ReturnType = SqlSnippetReturnTypeInsertID
	} else {
		snippet.ReturnType = SqlSnippetReturnTypeNone
	}

	snippet.ReturnTypeStr = snippet.ReturnType.String()
}

func parseSQLFunctionSnippetFile(file string, group string) ([]SqlFunctionSnippet, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// convert the bytes to a string
	str := string(bytes)

	// split the string into lines
	lines := strings.Split(str, "\n")

	var snippet *SqlFunctionSnippet

	snippets := make([]SqlFunctionSnippet, 0)

	reMatchWord := regexp.MustCompile(`:\w+`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "-- :name") {
			if snippet != nil {
				snippet.finish(group)
				snippets = append(snippets, *snippet)
				snippet = nil
			}

			snippet = &SqlFunctionSnippet{SqlQuery: "", Properties: make(map[string]string), Parameters: make([]string, 0)}

			err := parseSqlSnippetPropertyLine(line, snippet)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "-- :doc") {
			if snippet == nil {
				return nil, errors.New("coc line without snippet")
			}
			doc := strings.TrimPrefix(line, "-- :doc")
			doc = strings.TrimSpace(doc)
			snippet.Doc = doc
		} else {
			if snippet == nil {
				return nil, errors.New("sql snippet not started")
			}

			paramMatches := reMatchWord.FindAllStringSubmatch(line, -1)
			for _, paramMatch := range paramMatches {
				param := paramMatch[0]
				param = strings.TrimPrefix(param, ":")
				snippet.Parameters = append(snippet.Parameters, param)
			}

			snippet.SqlQuery += line + "\n"
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

func parseSqlSnippetPropertyLine(line string, snippet *SqlFunctionSnippet) error {
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

			snippet.Properties[key] = value
		}
	}

	return nil
}

func parseSQLFunctionSnippetFolder(folder string) ([]SqlFunctionSnippet, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	snippets := make([]SqlFunctionSnippet, 0)

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
