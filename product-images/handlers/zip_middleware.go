package handlers

import (
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			rw.Write([]byte("hello"))
		}
	})
}
