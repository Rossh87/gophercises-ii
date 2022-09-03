package middleware

import (
	"bytes"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

type BufferedWriter struct {
	body    *bytes.Buffer
	headers http.Header
	status  int
}

func (bwr *BufferedWriter) Header() http.Header {
	return bwr.headers
}

func (bwr *BufferedWriter) Write(b []byte) {
	bwr.body.Write(b)
}

func (bwr *BufferedWriter) WriteHeader(code int) {
	bwr.status = code
}

func (bwr *BufferedWriter) Send(r http.ResponseWriter) {

	// aim is to copy FROM receiver object TO the responsewriter
	// param.

	// copy headers
	dest := r.Header()

	src := bwr.Header()

	for k, v := range src {
		for _, sglVal := range v {
			dest.Add(k, sglVal)
		}
	}

	// update status code with current MW status, if any
	r.WriteHeader(bwr.status)

}

func NewBufferedWriter() *BufferedWriter {
	var buf = &bytes.Buffer{}

	return &BufferedWriter{
		buf,
		make(http.Header),
		200,
	}
}
