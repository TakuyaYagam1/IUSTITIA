package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

// PATCH - Vuln 5 (classified_note leak к не-judge ролям) для архива.
// Было: ClassifiedNote (SECRET_MARKER_A) выдавался любому auth-юзеру
// через GET /api/archive и /api/archive/{id}.
// Стало: converter принимает role; для всех, кроме judge, поле зануляется.
func FromArchiveEntry(e *domain.ArchiveEntry, role domain.Role) openapi.ArchiveEntry {
	out := openapi.ArchiveEntry{
		ID:           e.ID,
		Defendant:    e.Defendant,
		FinalVerdict: openapi.Verdict(e.FinalVerdict),
		Sentence:     e.Sentence,
		ArchivedAt:   e.ArchivedAt,
	}
	if role == domain.RoleJudge {
		out.ClassifiedNote = e.ClassifiedNote
	}
	if e.CaseID != nil {
		v := *e.CaseID
		out.CaseID = &v
	}
	return out
}

func FromArchiveList(list []*domain.ArchiveEntry, role domain.Role) []openapi.ArchiveEntry {
	out := make([]openapi.ArchiveEntry, 0, len(list))
	for _, e := range list {
		out = append(out, FromArchiveEntry(e, role))
	}
	return out
}
