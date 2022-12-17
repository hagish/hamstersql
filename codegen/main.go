package main

import (
	"fmt"
)

func main() {
	snippets, err := parseSQLFunctionSnippetFolder("example-sql")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(snippets)
}
