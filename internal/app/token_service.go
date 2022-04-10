//go:generate mockgen -source token_service.go -destination mock/token_service_mock.go -package mock
package app

import (
	"context"
	"net/http"
)

type Claims interface {
	GetSubject() string
	GetRole() string
	GetIssuer() string
}

type TokenService interface {
	ParseToken(ctx context.Context, token string) (Claims, error)
	ValidateAccessTokenFromRequest(ctx context.Context, r *http.Request) (Claims, error)
}
