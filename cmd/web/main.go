package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Get()

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{})

	router := mux.NewRouter()
	router.NotFoundHandler = server.NotFoundHandler(logger)
	router.HandleFunc("/ticketmaster/user", server.IndexHandler(logger)).Methods(http.MethodGet)

	httpServer := server.NewServer(logger, router, cfg.App.Name, cfg.App.Port)
	httpServer.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	httpServer.Stop()
}
