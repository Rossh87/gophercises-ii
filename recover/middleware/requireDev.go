package middleware

import (
	"net/http"
	"os"
)

func RequireDev(next http.Handler) http.Handler {
	h := func(rw http.ResponseWriter, r *http.Request) {
		if genv, _ := os.LookupEnv("GO_ENV"); genv != "development" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(h)
}
