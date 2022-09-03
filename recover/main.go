package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", greet)

	log.Fatal(http.ListenAndServe(":3000", &RecoveryMW{mux}))
}

func fPanic() {
	panic("Snap!!")
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

type RecoveryMW struct {
	handler http.Handler
}

func (rmw *RecoveryMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log err and stack trace
	// set status to 500
	// add 'something fucked up' msg for user
	// overwrite status header and remove any partially-written response
	// bytes
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			debug.PrintStack()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong!"))
		}
	}()

	rmw.handler.ServeHTTP(w, r)
}
