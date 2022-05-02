package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/orkungursel/hey-taxi-location-api/config"
)

func CORS(c *config.Config) echo.MiddlewareFunc {
	return emw.CORSWithConfig(emw.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     c.Server.Http.CorsOrigins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXRealIP,
			echo.HeaderAuthorization,
			echo.HeaderXRequestID,
			echo.HeaderXCorrelationID,
			echo.HeaderXCSRFToken,
			"X-Envoy-External-Address",
		},
		ExposeHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	})
}
