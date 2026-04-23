package domain

import (
	"time"

	"github.com/google/uuid"
)

type ArchiveEntry struct {
	ID             uuid.UUID
	CaseID         *uuid.UUID
	Defendant      string
	FinalVerdict   Verdict
	Sentence       *string
	ClassifiedNote *string
	ArchivedAt     time.Time
}
