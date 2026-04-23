package middleware

import (
	"net/http"
	"slices"
	"strings"
)

func CORS(origins []string) func(http.Handler) http.Handler {
	wildcard := slices.Contains(origins, "*")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqOrigin := r.Header.Get("Origin")

			allow := ""
			switch {
			case reqOrigin == "":
			case wildcard:
				allow = reqOrigin
			case slices.Contains(origins, reqOrigin):
				allow = reqOrigin
			}

			if allow != "" {
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", allow)
				h.Set("Access-Control-Allow-Credentials", "true")
				h.Set("Vary", "Origin")
				if r.Method == http.MethodOptions {
					h.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
					h.Set("Access-Control-Allow-Headers",
						strings.Join([]string{"Authorization", "Content-Type", "X-Request-ID"}, ", "),
					)
					h.Set("Access-Control-Max-Age", "600")
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
