package request

import "github.com/TakuyaYagam1/iustitia/internal/openapi"

func LoginRequestToParams(req *openapi.LoginRequest) (username, password string) {
	return req.Username, req.Password
}
