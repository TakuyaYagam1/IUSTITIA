package domain

import (
	"time"

	"github.com/google/uuid"
)

type CaseStatus string

const (
	CaseStatusDraft    CaseStatus = "draft"
	CaseStatusOpen     CaseStatus = "open"
	CaseStatusAssigned CaseStatus = "assigned"
	CaseStatusHearing  CaseStatus = "hearing"
	CaseStatusClosed   CaseStatus = "closed"
)

func (s CaseStatus) Valid() bool {
	switch s {
	case CaseStatusDraft, CaseStatusOpen, CaseStatusAssigned, CaseStatusHearing, CaseStatusClosed:
		return true
	}
	return false
}

type Verdict string

const (
	VerdictGuilty    Verdict = "guilty"
	VerdictAcquitted Verdict = "acquitted"
	VerdictDismissed Verdict = "dismissed"
)

func (v Verdict) Valid() bool {
	switch v {
	case VerdictGuilty, VerdictAcquitted, VerdictDismissed:
		return true
	}
	return false
}

type Case struct {
	ID                   uuid.UUID
	SeqNum               int64
	Defendant            string
	Crime                string
	Status               CaseStatus
	Verdict              *Verdict
	ClassifiedNote       *string
	AuthorID             *uuid.UUID
	AssignedProsecutorID *uuid.UUID
	CreatedAt            time.Time
}

type CaseOpinion struct {
	ID                 uuid.UUID
	CaseID             uuid.UUID
	ProsecutorID       uuid.UUID
	PreliminaryVerdict Verdict
	Reasoning          string
	FiledAt            time.Time
}
