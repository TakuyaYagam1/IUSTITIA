package apperr

import "errors"

var (
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("conflict")
	ErrInternal           = errors.New("internal error")
)
