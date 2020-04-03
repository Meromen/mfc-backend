package httpserver

import (
	"context"
	"github.com/Meromen/mfc-backend/logger"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	server *http.Server
	logger logger.Logger
}

type Server interface {
	Start() error
}

func (s *server) Start() error {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		<-sigint

		if err := s.server.Shutdown(context.Background()); err != nil {
			s.logger.Println("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	s.logger.Println("Server successfully started at port %v\n", s.server.Addr)
	s.logger.Println(s.server.ListenAndServe())

	<-idleConnsClosed
	return nil
}

func NewServer(writeTimeout, readTimeout time.Duration, addr string, handler http.Handler, logger logger.Logger) Server {
	httpserver := http.Server{
		WriteTimeout: writeTimeout * time.Second,
		ReadTimeout:  readTimeout * time.Second,
		Addr:         addr,
		Handler:      handler,
	}

	return &server{
		server: &httpserver,
		logger: logger,
	}
}
