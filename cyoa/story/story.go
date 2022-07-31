package story

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed story-template.html
var storyTemplate string

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("tpl").Parse(storyTemplate))
}

type handler struct {
	s        Story
	t        *template.Template
	pathFunc func(r *http.Request) string
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	return path
}

func WithPathFn(f func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFunc = f
	}
}

func GetStoryHandler(s Story, opts ...HandlerOption) http.Handler {
	handler := handler{s, tpl, defaultPathFunc}

	for _, optFn := range opts {
		optFn(&handler)
	}

	return handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)

		if err != nil {
			log.Printf("%+v", err)
			http.Error(w, "Problems!!", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func GetStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)

	var s Story

	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}

type StoryOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Episode struct {
	Title   string        `json:"title"`
	Story   []string      `json:"story"`
	Options []StoryOption `json:"options"`
}

type Story map[string]Episode
