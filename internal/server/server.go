package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

type Server struct {
	echo         *echo.Echo
	config       *config.Config
	logger       logger.ILogger
	ctx          context.Context
	httpHandlers []HttpApiHandlerItem
	done         chan struct{}
}

func New(ctx context.Context, config *config.Config, logger logger.ILogger) *Server {
	return &Server{
		ctx:    ctx,
		echo:   echo.New(),
		config: config,
		logger: logger,
		done:   make(chan struct{}),
	}
}

// Run starts the server
func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(s.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := runPlugs(s); err != nil {
		return err
	}

	go s.startHttpServer(ctx, cancel)
	go s.waitForSignal(ctx)

	<-s.done

	s.logger.Info("shutting down...")

	if err := s.shutdownHttpServer(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) Config() *config.Config {
	return s.config
}

func (s *Server) Logger() logger.ILogger {
	return s.logger
}

func (s *Server) Context() context.Context {
	return s.ctx
}

func (s *Server) Wait() chan struct{} {
	return s.done
}

// waitForSignal waits for the cancellation token
func (s *Server) waitForSignal(ctx context.Context) {
	defer close(s.done)
	<-ctx.Done()
	s.done <- struct{}{}
}
