//go:build dev

package swagger

import (
	"github.com/orkungursel/hey-taxi-location-api/internal/server"
)

func init() {
	server.Plug(func(s *server.Server, next server.Next) {
		if err := s.RegisterHttpApiAsRoot("/swagger/*", &SwaggerApi{Server: s}); err != nil {
			panic(err)
		}

		s.Logger().Info("swagger is enabled")

		next(nil)
	})
}
