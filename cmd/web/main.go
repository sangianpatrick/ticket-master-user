package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/internal/customer"
	"github.com/sangianpatrick/tm-user/pkg/apm"
	"github.com/sangianpatrick/tm-user/pkg/appvalidator"
	"github.com/sangianpatrick/tm-user/pkg/logrushook"
	globalMiddleware "github.com/sangianpatrick/tm-user/pkg/middleware"
	"github.com/sangianpatrick/tm-user/pkg/postgres"
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
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(logrus.AllLevels...)))
	logger.AddHook(logrushook.NewTraceIDLoggerHook())
	logger.AddHook(logrushook.NewStdoutLoggerHook(logrus.New(), cfg.Logrus.Formatter))

	pgdb := postgres.GetDatabase()

	appvalidator.Init()

	router := mux.NewRouter()
	router.Use(
		otelmux.Middleware(cfg.App.Name),
		globalMiddleware.ClientDeviceMiddleware,
		globalMiddleware.NewHTTPRequestLoggerMiddleware(logger, cfg.App.Debug).Middleware,
	)
	router.NotFoundHandler = server.NotFoundHandler(logger)
	router.StrictSlash(true)

	router.HandleFunc("/ticketmaster/user", server.IndexHandler()).Methods(http.MethodGet)

	customerRepository := customer.NewCustomerRepository(logger, pgdb, cfg.App.Location)
	customerUsecase := customer.NewCustomerUsecase(customer.CustomerUsecaseProps{
		Logger:             logger,
		Location:           cfg.App.Location,
		CustomerRepository: customerRepository,
	})
	customer.InitCustomerHTTPHandler(logger, router, customerUsecase)

	httpServer := server.NewServer(logger, router, cfg.App.Name, cfg.App.Port)
	httpServer.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	httpServer.Stop()
	pgdb.Close()
	otel.Stop(context.Background())

	logger.Info("shutdown with exit status 0")
	os.Exit(0)
}
