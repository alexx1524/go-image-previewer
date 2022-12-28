package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alexx1524/go-image-previewer/internal/cache"
	"github.com/alexx1524/go-image-previewer/internal/config"
	"github.com/alexx1524/go-image-previewer/internal/logger"
	"github.com/alexx1524/go-image-previewer/internal/server/http/middlewares"
	"github.com/alexx1524/go-image-previewer/internal/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	logger  logger.Logger
	server  *http.Server
	cache   cache.Cache
	storage storage.Storage
	client  *http.Client
	Router  *mux.Router
}

func NewServer(l logger.Logger, c config.Config, cache cache.Cache, storage storage.Storage) *Server {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:              net.JoinHostPort(c.HTTPServer.Address, strconv.Itoa(c.HTTPServer.Port)),
		ReadHeaderTimeout: time.Duration(c.HTTPServer.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	server := &Server{
		logger:  l,
		server:  httpServer,
		cache:   cache,
		storage: storage,
		Router:  router,
		client: &http.Client{
			Timeout: time.Duration(c.HTTPServer.ReadHeaderTimeout) * time.Second,
		},
	}

	logMiddleware := middlewares.LoggingMiddleware{
		Logger: l,
	}

	server.InitializeImagesRoutes()

	router.Use(logMiddleware.Middleware)

	return server
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(err.Error())
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("Stopping HTTP server error %s", err.Error()))
		return err
	}
	<-ctx.Done()
	return nil
}
