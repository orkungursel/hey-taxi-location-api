package server

import (
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)

var ErrApiAlreadyExists = errors.New("api route already exists %s")

type HttpApiHandler interface {
	RegisterRoutes(group *echo.Group)
}

type HttpApiHandlerItem struct {
	prefix  string
	handler HttpApiHandler
	isRoot  bool
}

func (s *Server) addHttpApi(prefix string, h HttpApiHandler, isRoot bool) error {
	for _, ahi := range s.httpHandlers {
		if ahi.isRoot == isRoot && ahi.prefix == prefix {
			return errors.Errorf(ErrApiAlreadyExists.Error(), prefix)
		}
	}

	s.httpHandlers = append(s.httpHandlers, HttpApiHandlerItem{prefix: prefix, handler: h, isRoot: isRoot})

	return nil
}

func (s *Server) RegisterHttpApi(prefix string, h HttpApiHandler) error {
	return s.addHttpApi(prefix, h, false)
}

func (s *Server) RegisterHttpApiAsRoot(prefix string, h HttpApiHandler) error {
	return s.addHttpApi(prefix, h, true)
}
