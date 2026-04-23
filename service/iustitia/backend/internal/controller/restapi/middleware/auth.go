package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/wahrwelt-kit/go-httpkit/httperr"

	"github.com/TakuyaYagam1/iustitia/internal/domain"
	iustitiajwt "github.com/TakuyaYagam1/iustitia/pkg/jwt"
)

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if !strings.HasPrefix(authz, "Bearer ") {
				writeJSONError(w, httperr.ErrNotAuthenticated())
				return
			}
			tokenStr := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
			if tokenStr == "" {
				writeJSONError(w, httperr.ErrNotAuthenticated())
				return
			}

			parsed, err := iustitiajwt.Parse(secret, tokenStr)
			if err != nil {
				writeJSONError(w, httperr.ErrNotAuthenticated())
				return
			}

			uid, err := uuid.Parse(parsed.UserID)
			if err != nil {
				writeJSONError(w, httperr.ErrNotAuthenticated())
				return
			}

			claims := Claims{
				UserID: uid,
				Role:   domain.Role(parsed.Role),
				Dome:   parsed.Dome,
			}
			next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
		})
	}
}

func RequireRole(allowed ...domain.Role) func(http.Handler) http.Handler {
	set := make(map[domain.Role]struct{}, len(allowed))
	for _, v := range allowed {
		set[v] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, ok := GetClaims(r.Context())
			if !ok {
				writeJSONError(w, httperr.ErrNotAuthenticated())
				return
			}
			if _, ok := set[c.Role]; !ok {
				writeJSONError(w, httperr.ErrForbidden())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSONError(w http.ResponseWriter, e *httperr.HTTPError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.HTTPStatus())
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    e.GetCode(),
			"message": e.Error(),
		},
	})
}
