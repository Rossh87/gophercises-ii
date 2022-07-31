package main

import (
	"flag"
	"fmt"
	"gophercises/linkParser/fetcher"
	"gophercises/linkParser/parser"
	"log"
)

func main() {
	fmt.Println("Starting...")

	fileFlag := flag.String("file", "", "file path to HTML to parse")
	urlFlag := flag.String("url", "", "URL of web page to parse")

	flag.Parse()

	t, err := fetcher.GetHTMLTarget(fileFlag, urlFlag)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	sourceInfo := t

	var reader fetcher.Reader

	if sourceInfo.SourceType == fetcher.Web {
		r, err := fetcher.WebReader(sourceInfo.String)

		if err != nil {
			log.Fatalf("%+v", err)
		}

		reader = r
	} else {
		r, err := fetcher.FileReader(sourceInfo.Path)

		if err != nil {
			log.Fatalf("%+v", err)
		}

		reader = r
	}

	results, err := parser.ParseHTML(reader)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	reader.Cleanup()

	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}
}
