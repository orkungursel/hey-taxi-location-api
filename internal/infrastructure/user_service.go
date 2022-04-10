package infrastructure

import (
	"context"
	"errors"

	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
	userServiceGrpc "github.com/orkungursel/hey-taxi-location-api/proto"
)

type UserService struct {
	client userServiceGrpc.UserServiceClient
}

func NewUserService(config *config.Config, logger logger.ILogger, client userServiceGrpc.UserServiceClient) *UserService {
	return &UserService{
		client: client,
	}
}

func (us *UserService) GetUsersByIds(ctx context.Context, in []string) (map[string]model.User, error) {
	if len(in) == 0 {
		return nil, errors.New("empty user ids")
	}

	guir, err := us.client.GetUserInfo(ctx, &userServiceGrpc.GetUserInfoRequest{UserIds: in})
	if err != nil {
		return nil, err
	}

	var users map[string]model.User = make(map[string]model.User)
	for _, u := range guir.Users {
		users[u.Id] = model.User{
			Id:       u.Id,
			Name:     u.Name,
			Email:    u.Email,
			Nickname: u.Name,
			Picture:  u.Avatar,
		}
	}

	return users, nil
}
