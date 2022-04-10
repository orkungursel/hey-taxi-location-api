//go:generate mockgen -source user_service.go -destination mock/user_service_mock.go -package mock
package app

import (
	"context"

	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
)

type UserService interface {
	GetUsersByIds(ctx context.Context, in []string) (map[string]model.User, error)
}
