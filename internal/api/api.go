package api

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/api/http"
	"github.com/orkungursel/hey-taxi-location-api/internal/infrastructure"
	"github.com/orkungursel/hey-taxi-location-api/internal/server"
	userService "github.com/orkungursel/hey-taxi-location-api/proto"
)

func Api(s *server.Server, redisClient *redis.Client, userServiceGrpcClient userService.UserServiceClient) error {
	if s == nil {
		return errors.New("server is nil")
	}

	c := s.Config()
	if c == nil {
		return errors.New("config is nil")
	}

	if redisClient == nil {
		return errors.New("redis client is nil")
	}

	logger := s.Logger()

	// test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return err
	}
	logger.Info("connected to Redis")

	if userServiceGrpcClient == nil {
		return errors.New("user service client is nil")
	}

	userService := infrastructure.NewUserService(c, logger, userServiceGrpcClient)
	tokenService := infrastructure.NewTokenService(c, logger)

	locationRepo := infrastructure.NewLocationRepository(redisClient, logger)
	locationService := infrastructure.NewLocationService(locationRepo, logger, userService)

	ctrl := http.NewController(c, logger, locationService, tokenService)
	if err := s.RegisterHttpApi("/location", ctrl); err != nil {
		return err
	}

	return nil
}
