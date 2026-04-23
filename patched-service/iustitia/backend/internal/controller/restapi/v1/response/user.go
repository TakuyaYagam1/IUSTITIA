package response

import (
	"github.com/TakuyaYagam1/iustitia/internal/domain"
	"github.com/TakuyaYagam1/iustitia/internal/openapi"
	"github.com/TakuyaYagam1/iustitia/internal/usecase"
)

func FromUserForMe(u *domain.User) openapi.User {
	return openapi.User{
		ID:        u.ID,
		Username:  u.Username,
		Role:      openapi.Role(u.Role),
		Dome:      u.Dome,
		CreatedAt: u.CreatedAt,
	}
}

func FromLoginResult(r *usecase.LoginResult) openapi.LoginResponse {
	return openapi.LoginResponse{
		Token:  r.Token,
		UserID: r.User.ID,
		Role:   openapi.Role(r.User.Role),
		Dome:   r.User.Dome,
	}
}

func FromUserListItem(u *domain.User) openapi.UserListItem {
	return openapi.UserListItem{
		ID:       u.ID,
		Username: u.Username,
		Role:     openapi.Role(u.Role),
		Dome:     u.Dome,
	}
}

func FromUserList(list []*domain.User) []openapi.UserListItem {
	out := make([]openapi.UserListItem, 0, len(list))
	for _, u := range list {
		out = append(out, FromUserListItem(u))
	}
	return out
}
