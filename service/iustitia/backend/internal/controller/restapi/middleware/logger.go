package middleware

import (
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	logkit "github.com/wahrwelt-kit/go-logkit"
)

func RequestLogger(logger logkit.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := chimw.GetReqID(r.Context())

			scoped := logger.WithFields(logkit.RequestID(reqID))
			ctx := logkit.IntoContext(r.Context(), scoped)

			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r.WithContext(ctx))

			scoped.Info("http_request", logkit.Fields{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      ww.Status(),
				"bytes":       ww.BytesWritten(),
				"duration_ms": time.Since(start).Milliseconds(),
				"remote":      r.RemoteAddr,
			})
		})
	}
}
