package domain

import (
	"time"

	"github.com/google/uuid"
)

type Complaint struct {
	ID           uuid.UUID
	CaseID       uuid.UUID
	AuthorID     uuid.UUID
	Text         string
	EvidenceURL  *string
	EvidenceData *string
	CreatedAt    time.Time
}
