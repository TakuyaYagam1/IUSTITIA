package v1

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/request"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// POST /api/complaints.
func (h *Server) CreateComplaint(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}

	var req openapi.ComplaintCreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "CreateComplaint", "DecodeJSON")
		return
	}

	caseID, text := request.ComplaintCreateRequestToParams(&req)

	complaint, err := h.complaint.ComplaintUC.Create(r.Context(), caseID, claims.UserID, text)
	if h.OnError(w, r, err, "CreateComplaint", "ComplaintUC.Create") {
		return
	}

	httputil.RenderCreated(w, r, response.FromComplaint(complaint))
}

// GET /api/complaints/{case_id}.
func (h *Server) ListComplaintsByCase(w http.ResponseWriter, r *http.Request, caseID openapi.UUID) {
	complaints, err := h.complaint.ComplaintUC.ListByCase(r.Context(), caseID)
	if h.OnError(w, r, err, "ListComplaintsByCase", "ComplaintUC.ListByCase") {
		return
	}
	httputil.RenderOK(w, r, response.FromComplaintList(complaints))
}

// POST /api/complaints/{id}/evidence.
func (h *Server) AttachEvidence(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	var req openapi.EvidenceRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "AttachEvidence", "DecodeJSON")
		return
	}

	url := request.EvidenceRequestToParams(&req)

	complaint, err := h.complaint.ComplaintUC.AttachEvidence(r.Context(), id, url)
	if h.OnError(w, r, err, "AttachEvidence", "ComplaintUC.AttachEvidence") {
		return
	}
	httputil.RenderOK(w, r, response.FromComplaint(complaint))
}
