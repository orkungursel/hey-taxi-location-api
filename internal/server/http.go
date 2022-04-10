package server

import (
	"context"
	"net/http"
	"time"
)

// startHttpServer starts the server
func (s *Server) startHttpServer(ctx context.Context, cancel context.CancelFunc) {
	s.configure()
	s.mapHandlers()

	// create http server
	httpServer := &http.Server{
		Addr: s.config.Server.Http.Host + ":" + s.config.Server.Http.Port,
	}

	s.logger.Infof("starting server on %s", httpServer.Addr)

	// start the server
	if err := s.echo.StartServer(httpServer); err != nil && err != http.ErrServerClosed {
		s.logger.Error(err)
		cancel()
		return
	}
}

// shutdownHttpServer stops the server
func (s *Server) shutdownHttpServer(ctx context.Context) error {
	s.logger.Info("stopping server...")
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Server.Http.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Errorf("error while shutting down server: %s", err)
		return err
	}

	<-ctx.Done()
	s.logger.Info("stopped server...")

	return nil
}
