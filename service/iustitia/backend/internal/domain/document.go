package domain

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID
	CaseID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	Template  string
	CreatedAt time.Time
}
