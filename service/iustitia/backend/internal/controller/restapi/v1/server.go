package v1

import (
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httputil"
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/errmap"
	"github.com/TakuyaYagam1/iustitia/internal/controller/restapi/v1/helper"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

var _ openapi.ServerInterface = (*Server)(nil)

type Server struct {
	openapi.Unimplemented

	user      helper.UserDeps
	complaint helper.ComplaintDeps
	caseUC    helper.CaseDeps
	document  helper.DocumentDeps
	archive   helper.ArchiveDeps
	infra     helper.InfraDeps
	renderer  *usecase.VerdictRenderer
}

func NewServer(deps *helper.ServerDeps) *Server {
	if deps == nil {
		return nil
	}
	return &Server{
		user:      deps.User,
		complaint: deps.Complaint,
		caseUC:    deps.Case,
		document:  deps.Document,
		archive:   deps.Archive,
		infra:     deps.Infra,
		renderer:  usecase.NewVerdictRenderer(),
	}
}

func (h *Server) OnError(w http.ResponseWriter, r *http.Request, err error, op, step string) bool {
	if err == nil {
		return false
	}
	httpErr := errmap.MapAppError(err)
	log := logkit.FromContext(r.Context())
	if httpErr.HTTPStatus() >= 500 {
		log.Error("restapi - v1 - "+op+" - "+step, logkit.Error(err))
	} else {
		log.Info("restapi - v1 - "+op+" - "+step, logkit.Fields{"err": err.Error()})
	}
	httputil.HandleError(w, r, httpErr)
	return true
}

// GET /api/health
func (h *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	httputil.RenderOK(w, r, openapi.HealthStatus{
		Status:  "operational",
		Version: "3.1.0",
	})
}
