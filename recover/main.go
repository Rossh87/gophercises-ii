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
	"runtime/debug"
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

	withRecovery := middleware.Compose(panicRecover)
	withDev := middleware.Compose(requireDev)

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

type stackPath struct {
	lineNo int
	path   string
}

func parseStackPath(stack string) stackPath {
	fields := strings.Fields(stack)
	data := strings.Split(fields[0], ":")

	ln, err := strconv.ParseInt(data[1], 10, 32)

	if err != nil {
		panic(err)
	}

	return stackPath{
		int(ln),
		data[0],
	}
}

func formatStack(stack string) string {
	line := strings.Builder{}
	out := strings.Builder{}

	for _, ch := range stack {
		if ch == '\n' {
			maybePath := strings.Trim(line.String(), " \t")

			if maybePath[0] == '/' {
				stackPath := parseStackPath(maybePath)
				maybePath = fmt.Sprintf("\t<a href=\"/source%s?lineNumber=%d\">%s</a>", stackPath.path, stackPath.lineNo, maybePath)
			}

			out.Write([]byte(maybePath + "\n"))
			line.Reset()
		} else {
			line.WriteRune(ch)
		}
	}

	if line.Len() > 0 {
		maybePath := strings.Trim(line.String(), " \t")

		if maybePath[0] == '/' {
			stackPath := parseStackPath(maybePath)
			maybePath = fmt.Sprintf("<a href=\"/source/%s?lineNumber=%d\">%s</a>", stackPath.path, stackPath.lineNo, maybePath)
		}

		out.Write([]byte(maybePath))
		line.Reset()
	}

	return out.String()
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

// It is important that this handler be at the bottom of the middleware stack:
// otherwise, subsequent middlewares could overwrite the error message and/or
// status code in the buffered writer.
func panicRecover(next http.Handler) http.Handler {
	h := func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				var panicMsg string

				if err, ok := r.(error); ok {
					panicMsg = err.Error()
				} else {
					panicMsg = "Panic: unknown reason"
				}

				log.Print(panicMsg)

				stack := string(debug.Stack())

				log.Print(stack)

				var usrMsg bytes.Buffer

				genv, _ := os.LookupEnv("GO_ENV")

				if genv == "development" {
					formattedStack := formatStack(stack)
					usrMsg.Write([]byte(fmt.Sprintf("<h1>Panic: %s</h1></br><pre>%s</pre>", panicMsg, formattedStack)))
				} else {
					usrMsg.Write([]byte("something went wrong!"))
				}

				rw.Write(usrMsg.Bytes())
			}
		}()

		bufWriter := middleware.NewBufferedWriter()

		next.ServeHTTP(bufWriter, r)

		bufWriter.Send(rw)
	}

	return http.HandlerFunc(h)
}

func requireDev(next http.Handler) http.Handler {
	h := func(rw http.ResponseWriter, r *http.Request) {
		if genv, _ := os.LookupEnv("GO_ENV"); genv != "development" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(h)
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
