package server

import "github.com/labstack/echo/v4"

func (s *Server) mapHandlers() {
	s.echo.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"service": s.config.App.Name})
	})

	root := s.echo.Group("/api/v1")

	for _, api := range s.httpHandlers {
		if api.isRoot {
			api.handler.RegisterRoutes(s.echo.Group(api.prefix))
		} else {
			api.handler.RegisterRoutes(root.Group(api.prefix))
		}
	}

	for _, route := range s.echo.Routes() {
		s.logger.Debugf("route: %s [%s]", route.Path, route.Method)
	}
}
