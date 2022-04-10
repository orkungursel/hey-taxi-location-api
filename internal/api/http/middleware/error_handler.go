package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
)

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				if e, ok := err.(*echo.HTTPError); ok {
					return e
				}

				if e, ok := err.(*app.Error); ok {
					code := e.Code()
					if code == http.StatusInternalServerError {
						return echo.NewHTTPError(http.StatusInternalServerError)
					}

					return echo.NewHTTPError(e.Code(), e.Error())
				}

				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return nil
		}
	}
}
