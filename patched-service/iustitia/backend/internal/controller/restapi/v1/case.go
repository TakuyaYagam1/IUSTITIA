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

// GET /api/cases.
func (h *Server) ListCases(w http.ResponseWriter, r *http.Request, params openapi.ListCasesParams) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	limit, offset := request.CaseListParams(params.Limit, params.Offset)

	cases, err := h.caseUC.CaseUC.List(r.Context(), limit, offset)
	if h.OnError(w, r, err, "ListCases", "CaseUC.List") {
		return
	}
	httputil.RenderOK(w, r, response.FromCaseList(cases, claims.Role))
}

// POST /api/cases. (citizen)
func (h *Server) CreateCase(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	var req openapi.CaseCreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "CreateCase", "DecodeJSON")
		return
	}
	c, err := h.caseUC.CaseUC.Create(r.Context(), claims.UserID, req.Defendant, req.Crime, req.Text)
	if h.OnError(w, r, err, "CreateCase", "CaseUC.Create") {
		return
	}
	httputil.RenderCreated(w, r, response.FromCase(c, claims.Role))
}

// POST /api/cases/{id}/accept. (registrar)
func (h *Server) AcceptCase(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	var req openapi.CaseAcceptRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "AcceptCase", "DecodeJSON")
		return
	}
	c, err := h.caseUC.CaseUC.Accept(r.Context(), id, req.ProsecutorID)
	if h.OnError(w, r, err, "AcceptCase", "CaseUC.Accept") {
		return
	}
	httputil.RenderOK(w, r, response.FromCase(c, claims.Role))
}

// POST /api/cases/{id}/dismiss. (registrar)
func (h *Server) DismissCase(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	if _, ok := helper.RequireClaims(w, r); !ok {
		return
	}
	var req openapi.CaseDismissRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "DismissCase", "DecodeJSON")
		return
	}
	if err := h.caseUC.CaseUC.Dismiss(r.Context(), id, req.Reason); err != nil {
		h.OnError(w, r, err, "DismissCase", "CaseUC.Dismiss")
		return
	}
	httputil.RenderOK(w, r, map[string]any{"ok": true})
}

// POST /api/cases/{id}/opinion. (prosecutor)
func (h *Server) FileOpinion(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	var req openapi.OpinionCreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "FileOpinion", "DecodeJSON")
		return
	}
	op, err := h.caseUC.CaseUC.FileOpinion(
		r.Context(), id, claims.UserID,
		domain.Verdict(req.PreliminaryVerdict), req.Reasoning,
	)
	if h.OnError(w, r, err, "FileOpinion", "CaseUC.FileOpinion") {
		return
	}
	httputil.RenderCreated(w, r, response.FromCaseOpinion(op))
}

// GET /api/cases/{id}/opinion. (judge, prosecutor)
func (h *Server) GetCaseOpinion(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	if _, ok := helper.RequireClaims(w, r); !ok {
		return
	}
	op, err := h.caseUC.CaseUC.GetOpinion(r.Context(), id)
	if h.OnError(w, r, err, "GetCaseOpinion", "CaseUC.GetOpinion") {
		return
	}
	httputil.RenderOK(w, r, response.FromCaseOpinion(op))
}

// GET /api/cases/{id}.
func (h *Server) GetCaseByID(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	c, err := h.caseUC.CaseUC.GetByID(r.Context(), id)
	if h.OnError(w, r, err, "GetCaseByID", "CaseUC.GetByID") {
		return
	}
	httputil.RenderOK(w, r, response.FromCase(c, claims.Role))
}

// POST /api/cases/search.
func (h *Server) SearchCases(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	var req openapi.SearchRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "SearchCases", "DecodeJSON")
		return
	}

	params := request.SearchRequestToParams(&req)

	cases, err := h.caseUC.CaseUC.Search(r.Context(), params)
	if h.OnError(w, r, err, "SearchCases", "CaseUC.Search") {
		return
	}
	httputil.RenderOK(w, r, response.FromCaseList(cases, claims.Role))
}

// POST /api/cases/{id}/verdict. (judge)
func (h *Server) SetCaseVerdict(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	var req openapi.VerdictRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "SetCaseVerdict", "DecodeJSON")
		return
	}
	verdict, sentence, reasoning := request.VerdictRequestFull(&req)

	res, err := h.caseUC.CaseUC.FinalizeVerdict(
		r.Context(), id, claims.UserID, verdict, sentence, reasoning, h.renderer,
	)
	if h.OnError(w, r, err, "SetCaseVerdict", "CaseUC.FinalizeVerdict") {
		return
	}
	httputil.RenderOK(w, r, response.FromVerdictResult(res))
}
