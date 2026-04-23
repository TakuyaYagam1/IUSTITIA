package v1

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
)

// GET /api/hearings. (judge)
// Список дел со статусом hearing с приклеенным opinion прокурора.
func (h *Server) ListHearings(w http.ResponseWriter, r *http.Request) {
	if _, ok := helper.RequireClaims(w, r); !ok {
		return
	}
	items, err := h.caseUC.CaseUC.ListHearings(r.Context())
	if h.OnError(w, r, err, "ListHearings", "CaseUC.ListHearings") {
		return
	}
	httputil.RenderOK(w, r, response.FromHearings(items))
}
