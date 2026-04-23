package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httperr"
	"github.com/wahrwelt-kit/go-httpkit/httputil"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/request"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/response"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// 256 KiB is plenty for a JSON patch body.
const maxArchivePatchBytes = 256 << 10

// GET /api/archive.
func (h *Server) ListArchive(w http.ResponseWriter, r *http.Request, params openapi.ListArchiveParams) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	limit, offset := request.ArchiveListParams(params.Limit, params.Offset)

	entries, err := h.archive.ArchiveUC.List(r.Context(), limit, offset)
	if h.OnError(w, r, err, "ListArchive", "ArchiveUC.List") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveList(entries, claims.Role))
}

// GET /api/archive/{id}.
func (h *Server) GetArchiveByID(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	entry, err := h.archive.ArchiveUC.GetByID(r.Context(), id)
	if h.OnError(w, r, err, "GetArchiveByID", "ArchiveUC.GetByID") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveEntry(entry, claims.Role))
}

// PATCH /api/archive/{id}.
//
// PATCH - Vuln 4 (Mass Assignment + IDOR on classified_note).
// Было: тело декодировалось в map[string]any и проксировалось как есть
// в usecase и дальше в SQLStore.UpdateArchiveDynamic, где в whitelist'е
// колонок находился classified_note. Любой атакующий с ролью judge мог
// передать `{"classified_note": "my-text"}` и перезаписать SECRET_MARKER_A.
// Стало:
//  1. Декодируем в ТИПИЗИРОВАННУЮ openapi.ArchivePatchRequest.
//  2. Строим map явно только из sentence / final_verdict - classified_note
//     из тела просто игнорируется, даже если он присутствует.
//  3. Store-whitelist тоже лишён classified_note (см. sqlstore.go) -
//     defense-in-depth: SQL никогда не выставит эту колонку, даже если
//     кто-то обойдёт handler.
func (h *Server) PatchArchive(w http.ResponseWriter, r *http.Request, id openapi.UUID) {
	claims, ok := helper.RequireClaims(w, r)
	if !ok {
		return
	}
	if r.Body == nil {
		h.OnError(w, r,
			httperr.New(io.EOF, http.StatusBadRequest, "BAD_REQUEST"),
			"PatchArchive", "empty body")
		return
	}
	defer func() { _ = r.Body.Close() }()

	var req openapi.ArchivePatchRequest
	dec := json.NewDecoder(io.LimitReader(r.Body, maxArchivePatchBytes))
	if err := dec.Decode(&req); err != nil {
		h.OnError(w, r,
			httperr.New(err, http.StatusBadRequest, "BAD_REQUEST"),
			"PatchArchive", "DecodeJSON")
		return
	}

	fields := map[string]any{}
	if req.Sentence != nil {
		fields["sentence"] = *req.Sentence
	}
	if req.FinalVerdict != nil {
		fields["final_verdict"] = string(*req.FinalVerdict)
	}

	entry, err := h.archive.ArchiveUC.Update(r.Context(), id, fields)
	if h.OnError(w, r, err, "PatchArchive", "ArchiveUC.Update") {
		return
	}
	httputil.RenderOK(w, r, response.FromArchiveEntry(entry, claims.Role))
}
