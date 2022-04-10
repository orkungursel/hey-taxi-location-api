package server

import (
	"time"

	emw "github.com/labstack/echo/v4/middleware"
	"github.com/orkungursel/hey-taxi-location-api/internal/server/middleware"
)

// configure the echo server
func (s *Server) configure() {
	s.echo.HidePort = true
	s.echo.HideBanner = true

	// add pre middlewares
	s.echo.Pre(middleware.AddTrailingSlash())
	s.echo.Pre(middleware.Logger(s.logger))

	// add middlewares
	s.echo.Use(emw.Recover())
	s.echo.Use(middleware.CORS(s.config))
	s.echo.Use(emw.Secure())
	s.echo.Use(emw.BodyLimit(s.config.Server.Http.BodyLimit))
	s.echo.Use(emw.Gzip())
	s.echo.Use(emw.RequestID())
	s.echo.Use(emw.TimeoutWithConfig(emw.TimeoutConfig{
		Timeout:      time.Duration(s.config.Server.Http.RequestTimeout) * time.Second,
		ErrorMessage: "{\"error\":\"request timeout\"}",
	}))
}
