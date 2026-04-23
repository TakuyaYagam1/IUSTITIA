package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
)

func FromDocument(d *domain.Document) openapi.Document {
	return openapi.Document{
		ID:        d.ID,
		CaseID:    d.CaseID,
		AuthorID:  d.AuthorID,
		Content:   d.Content,
		Template:  d.Template,
		CreatedAt: d.CreatedAt,
	}
}

func FromDocumentList(docs []*domain.Document) []openapi.Document {
	out := make([]openapi.Document, 0, len(docs))
	for _, d := range docs {
		out = append(out, FromDocument(d))
	}
	return out
}
