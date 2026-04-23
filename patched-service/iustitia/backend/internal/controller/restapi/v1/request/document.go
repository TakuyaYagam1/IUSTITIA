package request

import (
	"github.com/google/uuid"

	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func DocumentGenerateRequestToParams(req *openapi.DocumentGenerateRequest) (caseID uuid.UUID, tpl string) {
	return req.CaseID, req.Template
}
