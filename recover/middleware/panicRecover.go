package middleware

import (
	"bytes"
	"fmt"
	"gophercises/recover/stackFormat"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

// It is important that this handler be at the bottom of the middleware stack:
// otherwise, subsequent middlewares could overwrite the error message and/or
// status code in the buffered writer.
func PanicRecover(next http.Handler) http.Handler {
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
					formattedStack := stackFormat.FormatStack(stack)
					usrMsg.Write([]byte(fmt.Sprintf("<h1>Panic: %s</h1></br><pre>%s</pre>", panicMsg, formattedStack)))
				} else {
					usrMsg.Write([]byte("something went wrong!"))
				}

				rw.Write(usrMsg.Bytes())
			}
		}()

		bufWriter := NewBufferedWriter()

		next.ServeHTTP(bufWriter, r)

		bufWriter.Send(rw)
	}

	return http.HandlerFunc(h)
}
