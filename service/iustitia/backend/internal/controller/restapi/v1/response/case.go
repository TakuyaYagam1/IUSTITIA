package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

func FromCase(c *domain.Case) openapi.Case {
	out := openapi.Case{
		ID:                   c.ID,
		SeqNum:               int(c.SeqNum),
		Defendant:            c.Defendant,
		Crime:                c.Crime,
		Status:               openapi.CaseStatus(c.Status),
		ClassifiedNote:       c.ClassifiedNote,
		AuthorID:             c.AuthorID,
		AssignedProsecutorID: c.AssignedProsecutorID,
		CreatedAt:            c.CreatedAt,
	}
	if c.Verdict != nil {
		v := openapi.CaseVerdict(*c.Verdict)
		out.Verdict = &v
	}
	return out
}

func FromCaseOpinion(op *domain.CaseOpinion) openapi.CaseOpinion {
	return openapi.CaseOpinion{
		ID:                 op.ID,
		CaseID:             op.CaseID,
		ProsecutorID:       op.ProsecutorID,
		PreliminaryVerdict: openapi.PreliminaryVerdict(op.PreliminaryVerdict),
		Reasoning:          op.Reasoning,
		FiledAt:            op.FiledAt,
	}
}

func FromVerdictResult(res *usecase.VerdictResult) openapi.VerdictResult {
	out := openapi.VerdictResult{
		Case:         FromCase(res.Case),
		ArchiveEntry: FromArchiveEntry(res.ArchiveEntry),
	}
	if res.Document != nil {
		out.Document = FromDocument(res.Document)
	}
	return out
}

func FromHearings(list []*usecase.HearingItem) []openapi.HearingItem {
	out := make([]openapi.HearingItem, 0, len(list))
	for _, h := range list {
		if h == nil || h.Case == nil || h.Opinion == nil {
			continue
		}
		out = append(out, openapi.HearingItem{
			Case:    FromCase(h.Case),
			Opinion: FromCaseOpinion(h.Opinion),
		})
	}
	return out
}

func FromCaseList(list []*domain.Case) []openapi.Case {
	out := make([]openapi.Case, 0, len(list))
	for _, c := range list {
		out = append(out, FromCase(c))
	}
	return out
}
