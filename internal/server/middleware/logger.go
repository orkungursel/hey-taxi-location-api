package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

func Logger(log logger.ILogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			log.EchoCtx(c, start)

			return nil
		}
	}
}
