package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
)

func Auth(ts app.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := ts.ValidateAccessTokenFromRequest(c.Request().Context(), c.Request())

			if err != nil {
				return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
			}

			c.Set("claims", claims)

			return next(c)
		}
	}
}
