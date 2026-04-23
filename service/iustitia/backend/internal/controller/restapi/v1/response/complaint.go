package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func FromComplaint(c *domain.Complaint) openapi.Complaint {
	return openapi.Complaint{
		ID:           c.ID,
		CaseID:       c.CaseID,
		AuthorID:     c.AuthorID,
		Text:         c.Text,
		EvidenceURL:  c.EvidenceURL,
		EvidenceData: c.EvidenceData,
		CreatedAt:    c.CreatedAt,
	}
}

func FromComplaintList(list []*domain.Complaint) []openapi.Complaint {
	out := make([]openapi.Complaint, 0, len(list))
	for _, c := range list {
		out = append(out, FromComplaint(c))
	}
	return out
}
