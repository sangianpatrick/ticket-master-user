package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Server is a set of actions both for start and stop the http server. It wraps the built-in http server package.
type Server interface {
	// Start will start the http server and begins the request from http protocol within the given port. This is a non-blocking process, do not invocate it in go routine.
	Start()
	// Stop will block the incomming requset and shutdown the http server. But it will wait until the on-going requests to be finished before shutdown.
	Stop()
}

type server struct {
	logger     *logrus.Logger
	name       string
	port       int
	httpServer *http.Server
}

func NewServer(logger *logrus.Logger, httpHandler http.Handler, name string, port int) Server {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: httpHandler,
	}

	s := &server{
		logger:     logger,
		name:       name,
		port:       port,
		httpServer: httpServer,
	}

	return s
}

// // Start will start the http server and begins the request from http protocol within the given port. This is a non-blocking process, do not invocate it in go routine.
func (s *server) Start() {
	go func() {
		s.httpServer.ListenAndServe()
	}()

	s.logger.Infof("http server starts listen to :%d", s.port)
}

// Stop will block the incomming requset and shutdown the http server. But it will wait until the on-going requests to be finished before shutdown.
func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s.httpServer.Shutdown(ctx)

	s.logger.Info("http server is gracefully shutdown")
}
