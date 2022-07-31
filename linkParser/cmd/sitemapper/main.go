package main

import (
	"fmt"
	"gophercises/linkParser/fetcher"
	"log"
	"os"
)

func main() {
	entryPoint := os.Args[1]

	if len(entryPoint) == 0 {
		log.Fatalf("program requires a url or file entrypoint")
	}

	dummyFlag := ""

	target, err := fetcher.GetHTMLTarget(&dummyFlag, &entryPoint)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	siteMap := fetcher.Map(target)

	fmt.Println(siteMap.Urls())
}
