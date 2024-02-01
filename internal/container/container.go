package container

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/internal/customer"
	"github.com/sangianpatrick/tm-user/pkg/apm"
	"github.com/sangianpatrick/tm-user/pkg/applogger"
	"github.com/sangianpatrick/tm-user/pkg/appvalidator"
	globalMiddleware "github.com/sangianpatrick/tm-user/pkg/middleware"
	"github.com/sangianpatrick/tm-user/pkg/postgres"
	"github.com/sangianpatrick/tm-user/pkg/server"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func run(ctx context.Context, shutdownSignal <-chan os.Signal, hasStop chan struct{}) {
	cfg := config.Get()

	otel := apm.GetOpenTelemetry()
	otel.Start(ctx)

	logger := applogger.GetLogrus()

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

	<-shutdownSignal

	httpServer.Stop()
	pgdb.Close()
	otel.Stop(ctx)

	hasStop <- struct{}{}
}

func Run(ctx context.Context, shutdownSignal <-chan os.Signal) (hasStop <-chan struct{}) {
	stopChan := make(chan struct{}, 1)
	go run(ctx, shutdownSignal, stopChan)

	return stopChan
}
