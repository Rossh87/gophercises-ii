package fetcher

import (
	"errors"
	"fmt"
	"gophercises/linkParser/parser"
	"log"
	"net/http"
	"net/url"
	"os"
)

type HtmlSource int
type Mode int

const (
	File HtmlSource = iota
	Web
)

type Reader interface {
	Read(p []byte) (int, error)
	Cleanup()
}

type fileReader struct {
	readable *os.File
}

func (f *fileReader) Cleanup() {
	f.readable.Close()
}

func (f *fileReader) Read(p []byte) (int, error) {
	return f.readable.Read(p)
}

type webReader struct {
	readable *http.Response
}

func (f *webReader) Cleanup() {
	f.readable.Body.Close()
}

func (f *webReader) Read(p []byte) (int, error) {
	return f.readable.Body.Read(p)
}

func FileReader(filePath string) (*fileReader, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	fr := fileReader{file}

	return &fr, nil
}

func WebReader(url string) (*webReader, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	wr := webReader{resp}

	return &wr, nil
}

type htmlTarget struct {
	SourceType HtmlSource
	Path       string
	Host       string
	Scheme     string
	String     string
}

func parseEntryPoint(pathString string) (*url.URL, error) {
	parsed, err := url.Parse(pathString)

	if err != nil {
		return nil, err
	}

	if len(parsed.Host) == 0 || len(parsed.Scheme) == 0 {
		return nil, errors.New("provided URL must specify at minimum a scheme and a host")
	}

	return parsed, nil
}

func GetHTMLTarget(fileFlag *string, urlFlag *string) (htmlTarget, error) {
	hasFilepath := len(*fileFlag) > 0
	hasUrl := len(*urlFlag) > 0

	if hasUrl && hasFilepath {
		return htmlTarget{}, errors.New("url and filepath options may not both be set")
	}

	if hasUrl {
		parsed, err := parseEntryPoint(*urlFlag)

		if err != nil {
			return htmlTarget{}, err
		}

		return htmlTarget{Path: parsed.Path, SourceType: Web, Host: parsed.Host, Scheme: parsed.Scheme, String: parsed.String()}, nil
	}

	if hasFilepath {
		return htmlTarget{Path: *fileFlag, SourceType: File}, nil
	}

	return htmlTarget{}, errors.New("must specify either web address or file path to provide HTML")
}

type siteMap struct {
	target    htmlTarget
	failures  []string
	usedPaths map[string]struct{}
	urls      []string
}

type addResult struct {
	fullUrl   string
	unvisited bool
}

func (s *siteMap) Add(maybeUrl string) addResult {
	parsed, err := url.Parse(maybeUrl)

	if err != nil {
		s.failures = append(s.failures, maybeUrl)
		return addResult{"", false}
	}

	// if link points to different host, it doesn't belong in sitemap
	if len(parsed.Host) > 0 && parsed.Host != s.target.Host {
		return addResult{"", false}
	}

	// otherwise, if host and/or scheme are empty, populate them with
	// target host and scheme
	parsed.Host = s.target.Host
	parsed.Scheme = s.target.Scheme

	// this will give us back the query string too if one exists, which I *think* is what we want??
	fullUrl := parsed.String()
	// fmt.Println(fullUrl)

	if _, used := s.usedPaths[fullUrl]; used {
		// we've already seen this url, so don't add it
		return addResult{"", false}
	}

	s.urls = append(s.urls, fullUrl)

	s.usedPaths[fullUrl] = struct{}{}

	return addResult{fullUrl, true}
}

func (s *siteMap) Urls() []string {
	return s.urls
}

func pop(xs *[]string) string {
	ln := len(*xs)

	if ln == 0 {
		return ""
	}

	t := (*xs)[ln-1]

	*xs = (*xs)[:ln-1]

	return t
}

func Map(target htmlTarget) siteMap {
	siteMap := siteMap{target: target, usedPaths: make(map[string]struct{})}

	links := []string{target.String}

	for len(links) > 0 {
		next := pop(&links)

		r, err := http.Get(next)

		if err != nil {
			log.Fatalf("request to %s failed\nreason:\t%v", next, err)
		}

		nextLinks, err := parser.ParseHTML(r.Body)

		r.Body.Close()

		if err != nil {
			log.Fatalf("failed to parse links for document revceived from %s\nreason:\t%+v", next, err)
		}

		for _, l := range nextLinks {
			if addResult := siteMap.Add(l.Href); addResult.unvisited {
				fmt.Printf("appending new link:\n%s\n", addResult.fullUrl)
				links = append(links, addResult.fullUrl)
			}
		}
	}

	return siteMap
}
