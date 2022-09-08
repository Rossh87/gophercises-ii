package middleware

import (
	"bytes"
	"fmt"
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

func (bwr *BufferedWriter) Write(b []byte) (int, error) {
	return bwr.body.Write(b)
}

func (bwr *BufferedWriter) WriteHeader(code int) {
	bwr.status = code
}

func (bwr *BufferedWriter) Send(r http.ResponseWriter) error {
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

	bodyBytes := bwr.body.Bytes()
	contentLen := fmt.Sprintf("%d", len(bodyBytes))
	dest.Set("Content-Length", contentLen)

	// copy body
	_, err := r.Write(bodyBytes)

	if flusher, ok := r.(http.Flusher); ok {
		flusher.Flush()
	}

	return err
}

func NewBufferedWriter() *BufferedWriter {
	var buf = &bytes.Buffer{}

	return &BufferedWriter{
		buf,
		make(http.Header),
		200,
	}
}

func Compose(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		ln := len(mws)

		currHandler := final

		for i := ln - 1; i >= 0; i-- {
			currHandler = mws[i](currHandler)
		}

		return currHandler
	}
}
