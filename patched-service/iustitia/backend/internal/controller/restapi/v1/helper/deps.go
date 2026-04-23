package helper

import (
	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

type UserDeps struct {
	UserUC *usecase.User
}

type ComplaintDeps struct {
	ComplaintUC *usecase.Complaint
}

type CaseDeps struct {
	CaseUC *usecase.Case
}

type DocumentDeps struct {
	DocumentUC *usecase.Document
}

type ArchiveDeps struct {
	ArchiveUC *usecase.Archive
}

type InfraDeps struct {
	Logger logkit.Logger
}

type ServerDeps struct {
	User      UserDeps
	Complaint ComplaintDeps
	Case      CaseDeps
	Document  DocumentDeps
	Archive   ArchiveDeps
	Infra     InfraDeps
}
