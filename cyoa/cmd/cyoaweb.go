package main

import (
	"flag"
	"fmt"
	"gophercises/cyoa/story"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filename := flag.String("file", "story.json", "file container story episodes")
	port := flag.String("port", "8080", "port to start web server on")

	flag.Parse()

	reader, err := os.Open(*filename)

	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	storyData, err := story.GetStory(reader)

	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	getStoryPath := story.WithPathFn(func(r *http.Request) string {
		path := strings.TrimSpace(r.URL.Path)

		if path == "/story" || path == "/story/" {
			path = "/story/intro"
		}

		path = path[len("/story/"):]

		return path
	})

	handler := story.GetStoryHandler(storyData, getStoryPath)

	mux := http.NewServeMux()

	mux.Handle("/story/", handler)

	fmt.Printf("starting server on port %s\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), mux))
}
