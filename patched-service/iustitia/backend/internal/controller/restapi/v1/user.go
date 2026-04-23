package v1

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/request"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// POST /api/auth/login.
func (h *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req openapi.LoginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "Login", "DecodeJSON")
		return
	}

	username, password := request.LoginRequestToParams(&req)

	result, err := h.user.UserUC.Login(r.Context(), username, password)
	if h.OnError(w, r, err, "Login", "UserUC.Login") {
		return
	}

	httputil.RenderOK(w, r, response.FromLoginResult(result))
}

func (h *Server) Logout(w http.ResponseWriter, r *http.Request) {
	httputil.RenderNoContent(w, r)
}

// GET /api/auth/me.
func (h *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}

	user, err := h.user.UserUC.GetByID(r.Context(), claims.UserID)
	if h.OnError(w, r, err, "GetCurrentUser", "UserUC.GetByID") {
		return
	}

	httputil.RenderOK(w, r, response.FromUserForMe(user))
}

// GET /api/users?role=X. (registrar, judge)
func (h *Server) ListUsersByRole(w http.ResponseWriter, r *http.Request, params openapi.ListUsersByRoleParams) {
	if _, ok := helper.RequireClaims(w, r); !ok {
		return
	}
	users, err := h.user.UserUC.ListByRole(r.Context(), domain.Role(params.Role))
	if h.OnError(w, r, err, "ListUsersByRole", "UserUC.ListByRole") {
		return
	}
	httputil.RenderOK(w, r, response.FromUserList(users))
}
