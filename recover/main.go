package main

import (
	"bytes"
	"errors"
	"fmt"
	"gophercises/recover/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", greet)

	withRecovery := middleware.Compose(middleware.PanicRecover)
	withDev := middleware.Compose(middleware.RequireDev)

	mux.Handle("/source/", withDev(http.HandlerFunc(renderSource)))

	log.Fatal(http.ListenAndServe(":3000", withRecovery(mux)))
}

func shiftPath(path string) (head, tail string) {
	path = filepath.Clean("/" + path)

	i := strings.Index(path[1:], "/") + 1

	if i <= 0 {
		return path[1:], "/"
	}

	return path[1:i], path[i:]
}

func renderSource(rw http.ResponseWriter, r *http.Request) {
	_, tail := shiftPath(r.URL.Path)
	fmt.Printf("Attempting to open %s", tail)

	sourceBytes, err := os.ReadFile(tail)

	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	highlighted := bytes.Buffer{}

	formatter := html.New()

	if lineNumber, ok := r.URL.Query()["lineNumber"]; ok {
		ln, _ := strconv.ParseInt(lineNumber[0], 10, 32)
		lineRange := [][2]int{{int(ln), int(ln)}}
		formatter = html.New(html.HighlightLines(lineRange), html.WithLineNumbers(true))
	}

	lexer := lexers.Get("go")

	highlightStyle := styles.Get("monokai")

	iterator, _ := lexer.Tokenise(nil, string(sourceBytes))

	formatter.Format(&highlighted, highlightStyle, iterator)

	tpl := `
	<h1>
	Source:
	</h1>	
	</br>
	<div style="margin: auto; width: 85%; overflow: auto; padding: 2rem;">
	[[highlighted]]
	</div>
	`

	tpl = strings.Replace(tpl, "[[highlighted]]", highlighted.String(), 1)

	rw.Write([]byte(tpl))
}

func fPanic() {
	panic(errors.New("snap"))
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	fPanic()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Greetings!</h1>")
	fPanic()
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Greetings!</h1>")
}
