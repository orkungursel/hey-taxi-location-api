package main

import (
	"context"

	"github.com/orkungursel/hey-taxi-location-api/config"
	_ "github.com/orkungursel/hey-taxi-location-api/internal/api"
	"github.com/orkungursel/hey-taxi-location-api/internal/server"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
	_ "github.com/orkungursel/hey-taxi-location-api/pkg/swagger"
)

// @title                       Hey Taxi Location API
// @version                     1.0
// @BasePath                    /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @desc                        Add Bearer token to the request header
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.New()
	logger := logger.New(config)
	defer logger.Sync()

	logger.Infof("current profile: %s", config.GetProfile())

	if err := server.New(ctx, config, logger).Run(); err != nil {
		logger.Fatalf("error while starting server: %s", err)
	}
}
