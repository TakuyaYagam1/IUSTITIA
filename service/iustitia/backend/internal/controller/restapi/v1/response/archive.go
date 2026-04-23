package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func FromArchiveEntry(e *domain.ArchiveEntry) openapi.ArchiveEntry {
	out := openapi.ArchiveEntry{
		ID:             e.ID,
		Defendant:      e.Defendant,
		FinalVerdict:   openapi.Verdict(e.FinalVerdict),
		Sentence:       e.Sentence,
		ClassifiedNote: e.ClassifiedNote,
		ArchivedAt:     e.ArchivedAt,
	}
	if e.CaseID != nil {
		v := *e.CaseID
		out.CaseID = &v
	}
	return out
}

func FromArchiveList(list []*domain.ArchiveEntry) []openapi.ArchiveEntry {
	out := make([]openapi.ArchiveEntry, 0, len(list))
	for _, e := range list {
		out = append(out, FromArchiveEntry(e))
	}
	return out
}
