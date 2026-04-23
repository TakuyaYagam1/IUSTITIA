package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

// PATCH - Vuln 5 (classified_note leak к не-judge ролям).
// Было: ClassifiedNote всегда пробрасывался в ответ; prosecutor и registrar
// (и даже через IDOR любой auth-юзер) читали SECRET_MARKER_J прямо из
// GET /api/cases и GET /api/cases/{id}.
// Стало: converter принимает role; для всех, кроме judge, поле зануляется.
// Ответ остаётся в той же JSON-форме (поле `classified_note` есть и равно
// null) - клиент-сайд ничего чинить не надо.
func FromCase(c *domain.Case, role domain.Role) openapi.Case {
	out := openapi.Case{
		ID:                   c.ID,
		SeqNum:               int(c.SeqNum),
		Defendant:            c.Defendant,
		Crime:                c.Crime,
		Status:               openapi.CaseStatus(c.Status),
		CreatedAt:            c.CreatedAt,
		AuthorID:             c.AuthorID,
		AssignedProsecutorID: c.AssignedProsecutorID,
	}
	if role == domain.RoleJudge {
		out.ClassifiedNote = c.ClassifiedNote
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
		Case:         FromCase(res.Case, domain.RoleJudge),
		ArchiveEntry: FromArchiveEntry(res.ArchiveEntry, domain.RoleJudge),
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
			Case:    FromCase(h.Case, domain.RoleJudge),
			Opinion: FromCaseOpinion(h.Opinion),
		})
	}
	return out
}

func FromCaseList(list []*domain.Case, role domain.Role) []openapi.Case {
	out := make([]openapi.Case, 0, len(list))
	for _, c := range list {
		out = append(out, FromCase(c, role))
	}
	return out
}
