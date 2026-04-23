package middleware

import (
	"context"

	"github.com/google/uuid"

	"github.com/TakuyaYagam1/iustitia/internal/domain"
)

type Claims struct {
	UserID uuid.UUID
	Role   domain.Role
	Dome   string
}

type ctxKey int

const (
	ctxClaims ctxKey = iota
)

func GetClaims(ctx context.Context) (Claims, bool) {
	v, ok := ctx.Value(ctxClaims).(Claims)
	return v, ok
}

func withClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, ctxClaims, c)
}
