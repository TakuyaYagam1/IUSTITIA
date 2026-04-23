package v1

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/request"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// POST /api/documents/generate.
func (h *Server) GenerateDocument(w http.ResponseWriter, r *http.Request) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}

	var req openapi.DocumentGenerateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		h.OnError(w, r, err, "GenerateDocument", "DecodeJSON")
		return
	}

	caseID, tpl := request.DocumentGenerateRequestToParams(&req)

	doc, err := h.document.DocumentUC.Generate(r.Context(), caseID, claims.UserID, tpl)
	if h.OnError(w, r, err, "GenerateDocument", "DocumentUC.Generate") {
		return
	}

	httputil.RenderOK(w, r, response.FromDocument(doc))
}

// GET /api/cases/{id}/documents.
func (h *Server) ListCaseDocuments(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	docs, err := h.document.DocumentUC.ListByCase(r.Context(), id)
	if h.OnError(w, r, err, "ListCaseDocuments", "DocumentUC.ListByCase") {
		return
	}
	httputil.RenderOK(w, r, response.FromDocumentList(docs))
}

// GET /api/documents/{id}.
func (h *Server) GetDocumentByID(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	doc, err := h.document.DocumentUC.GetByID(r.Context(), id)
	if h.OnError(w, r, err, "GetDocumentByID", "DocumentUC.GetByID") {
		return
	}
	httputil.RenderOK(w, r, response.FromDocument(doc))
}
