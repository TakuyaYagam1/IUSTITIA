package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	logkit "github.com/wahrwelt-kit/go-logkit"
	"golang.org/x/crypto/bcrypt"

	"github.com/TakuyaYagam1/iustitia/internal/apperr"
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/repo"
	sqlc "github.com/TakuyaYagam1/iustitia/internal/repo/sqlc"
	iustitiajwt "github.com/TakuyaYagam1/iustitia/pkg/jwt"
)

type LoginResult struct {
	User  *domain.User
	Token string
}

type User struct {
	store     repo.Store
	jwtSecret string
	jwtTTL    time.Duration
	logger    logkit.Logger
}

func NewUser(store repo.Store, jwtSecret string, jwtTTL time.Duration, logger logkit.Logger) *User {
	return &User{
		store:     store,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
		logger:    logger,
	}
}

func (u *User) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	if username == "" || password == "" {
		return nil, apperr.ErrInvalidCredentials
	}

	row, err := u.store.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, apperr.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password)); err != nil {
		return nil, apperr.ErrInvalidCredentials
	}

	token, err := iustitiajwt.Sign(u.jwtSecret, iustitiajwt.Claims{
		UserID: row.ID,
		Role:   row.Role,
		Dome:   row.Dome,
	}, u.jwtTTL)
	if err != nil {
		return nil, fmt.Errorf("User - Login - jwt.Sign: %w", err)
	}

	user, err := sqlcToUser(row)
	if err != nil {
		return nil, fmt.Errorf("User - Login - sqlcToUser: %w", err)
	}
	return &LoginResult{User: user, Token: token}, nil
}

func (u *User) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := u.store.GetUserByID(ctx, id.String())
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	return sqlcToUser(row)
}

func (u *User) ListByRole(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	if !role.Valid() {
		return nil, apperr.ErrBadRequest
	}
	rows, err := u.store.ListUsersByRole(ctx, string(role))
	if err != nil {
		return nil, fmt.Errorf("User - ListByRole - store.ListUsersByRole: %w", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, r := range rows {
		user, err := sqlcToUser(r)
		if err != nil {
			return nil, fmt.Errorf("User - ListByRole - sqlcToUser: %w", err)
		}
		out = append(out, user)
	}
	return out, nil
}

func sqlcToUser(row sqlc.User) (*domain.User, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, errors.New("usecase - sqlcToUser - malformed user id")
	}
	return &domain.User{
		ID:        id,
		Username:  row.Username,
		Password:  row.Password,
		Role:      domain.Role(row.Role),
		Dome:      row.Dome,
		CreatedAt: row.CreatedAt,
	}, nil
}
