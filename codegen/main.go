package main

func main() {
	snippets, err := parseSQLFunctionSnippetFolder("example-sql")
	if err != nil {
		panic(err)
	}

	t, err := loadTemplate("templates/python")
	if err != nil {
		panic(err)
	}

	err = t.executeTemplate("output", snippets)
	if err != nil {
		panic(err)
	}
}
