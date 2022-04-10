package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
)

func GetUserId(c echo.Context) (string, error) {
	claims := c.Get("claims")
	if claims == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "claims is nil")
	}

	return claims.(app.Claims).GetSubject(), nil
}
