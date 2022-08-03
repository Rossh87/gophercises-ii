package main

import (
	"flag"
	"fmt"
	"gophercises/url-shortener/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := getHandler()

	fmt.Println("starting server on :8080")

	http.ListenAndServe(":8080", handler)
}

func getHandler() http.HandlerFunc {
	mux := getMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"yaml-godoc":      "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	yamlPath := readYAMLFlag()

	if len(yamlPath) > 0 {

		data, err := os.ReadFile(yamlPath)

		checkErr(err)

		fmt.Println(string(data))

		yamlHandler, err := handlers.YAMLHandler(data, mapHandler)

		checkErr(err)

		return yamlHandler
	}

	return mapHandler
}

func getMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func readYAMLFlag() string {
	var maybeFile string

	flag.StringVar(&maybeFile, "f", "default", "path to a YAML redirect file")

	flag.Parse()

	return maybeFile
}

func checkErr(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}
