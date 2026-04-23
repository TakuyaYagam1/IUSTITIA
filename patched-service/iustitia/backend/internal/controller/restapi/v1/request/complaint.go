package request

import (
	"github.com/google/uuid"

	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func ComplaintCreateRequestToParams(req *openapi.ComplaintCreateRequest) (caseID uuid.UUID, text string) {
	return req.CaseID, req.Text
}

func EvidenceRequestToParams(req *openapi.EvidenceRequest) string {
	return req.URL
}
