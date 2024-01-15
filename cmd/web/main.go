package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/pkg/apm"
	globalMiddleware "github.com/sangianpatrick/tm-user/pkg/middleware"
	"github.com/sangianpatrick/tm-user/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	cfg := config.Get()

	otel := apm.GetOpenTelemetry()
	otel.Start(context.Background())

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(cfg.Logrus.Formatter)
	logger.Hooks.Add(otellogrus.NewHook(otellogrus.WithLevels(logrus.AllLevels...)))

	router := mux.NewRouter()
	router.Use(
		otelmux.Middleware(cfg.App.Name),
		globalMiddleware.ClientDeviceMiddleware,
		globalMiddleware.NewHTTPRequestLoggerMiddleware(logger, cfg.App.Debug).Middleware,
	)
	router.NotFoundHandler = server.NotFoundHandler(logger)
	router.StrictSlash(true)

	router.HandleFunc("/ticketmaster/user", server.IndexHandler(logger)).Methods(http.MethodGet)

	httpServer := server.NewServer(logger, router, cfg.App.Name, cfg.App.Port)
	httpServer.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	httpServer.Stop()
	otel.Stop(context.Background())
}
