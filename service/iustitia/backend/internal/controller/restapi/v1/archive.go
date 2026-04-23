package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httperr"
	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/request"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// 256 KiB is plenty for a JSON patch body.
const maxArchivePatchBytes = 256 << 10

// GET /api/archive.
func (h *Server) ListArchive(w http.ResponseWriter, r *http.Request, params openapi.ListArchiveParams) {
	limit, offset := request.ArchiveListParams(params.Limit, params.Offset)

	entries, err := h.archive.ArchiveUC.List(r.Context(), limit, offset)
	if h.OnError(w, r, err, "ListArchive", "ArchiveUC.List") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveList(entries))
}

// GET /api/archive/{id}.
func (h *Server) GetArchiveByID(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	entry, err := h.archive.ArchiveUC.GetByID(r.Context(), id)
	if h.OnError(w, r, err, "GetArchiveByID", "ArchiveUC.GetByID") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveEntry(entry))
}

// PATCH /api/archive/{id}.
func (h *Server) PatchArchive(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	if r.Body == nil {
		h.OnError(w, r,
			httperr.New(io.EOF, http.StatusBadRequest, "BAD_REQUEST"),
			"PatchArchive", "empty body")
		return
	}
	defer func() { _ = r.Body.Close() }()

	var fields map[string]any
	dec := json.NewDecoder(io.LimitReader(r.Body, maxArchivePatchBytes))
	if err := dec.Decode(&fields); err != nil {
		h.OnError(w, r,
			httperr.New(err, http.StatusBadRequest, "BAD_REQUEST"),
			"PatchArchive", "DecodeJSON")
		return
	}

	entry, err := h.archive.ArchiveUC.Update(r.Context(), id, fields)
	if h.OnError(w, r, err, "PatchArchive", "ArchiveUC.Update") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveEntry(entry))
}
