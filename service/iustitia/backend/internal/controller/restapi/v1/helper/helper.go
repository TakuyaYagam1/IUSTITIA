package helper

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httperr"
	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/middleware"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
)

func RequireClaims(w http.ResponseWriter, r *http.Request) (middleware.Claims, bool) {
	c, ok := middleware.GetClaims(r.Context())
	if !ok {
		httputil.HandleError(w, r, httperr.ErrNotAuthenticated())
		return middleware.Claims{}, false
	}
	return c, true
}

func IsJudge(c middleware.Claims) bool { return c.Role == domain.RoleJudge }

func IsProsecutor(c middleware.Claims) bool { return c.Role == domain.RoleProsecutor }

func IsCitizen(c middleware.Claims) bool { return c.Role == domain.RoleCitizen }

func IsRegistrar(c middleware.Claims) bool { return c.Role == domain.RoleRegistrar }
