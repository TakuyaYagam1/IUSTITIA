package errmap

import (
	"errors"
	"net/http"

	"github.com/wahrwelt-kit/go-httpkit/httperr"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
)

func MapAppError(err error) *httperr.HTTPError {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, apperr.ErrInvalidCredentials),
		errors.Is(err, apperr.ErrNotAuthenticated):
		return httperr.ErrNotAuthenticated()
	case errors.Is(err, apperr.ErrForbidden):
		return httperr.ErrForbidden()
	case errors.Is(err, apperr.ErrNotFound):
		return httperr.ErrNotFound()
	case errors.Is(err, apperr.ErrConflict):
		return httperr.ErrConflict()
	case errors.Is(err, apperr.ErrBadRequest):
		return httperr.New(err, http.StatusBadRequest, "BAD_REQUEST")
	default:
		return httperr.New(err, http.StatusInternalServerError, "INTERNAL_ERROR")
	}
}
