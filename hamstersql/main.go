package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct {
	Verbose        bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	SqlFolder      string `short:"i" long:"sql" description:"The folder of all the input sql files" value-name:"SQLFOLDER" required:"true"`
	OutputFolder   string `short:"o" long:"output" description:"The folder for the generated code" value-name:"OUTPUTFOLDER" required:"true"`
	TemplateFolder string `short:"t" long:"template" description:"The folder that contains the template that should get used" value-name:"TPLFOLDER" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println()
		fmt.Println("Example arguments: -i example-sql -t templates/python -o output -v")
		os.Exit(1)
	}

	fmt.Printf("SqlFolder: %s\n", opts.SqlFolder)
	fmt.Printf("OutputFolder: %s\n", opts.OutputFolder)
	fmt.Printf("TemplateFolder: %s\n", opts.TemplateFolder)

	snippets, err := parseSQLFunctionSnippetFolder(opts.SqlFolder, opts.Verbose)
	if err != nil {
		panic(err)
	}

	t, err := loadTemplate(opts.TemplateFolder, opts.Verbose)
	if err != nil {
		panic(err)
	}

	err = t.executeTemplate(opts.OutputFolder, snippets, opts.Verbose)
	if err != nil {
		panic(err)
	}
}
