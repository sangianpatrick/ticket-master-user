package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/internal/container"
	"github.com/sangianpatrick/tm-user/pkg/applogger"
)

func main() {
	_ = config.Get()
	logger := applogger.GetLogrus()

	ctx := context.Background()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)

	<-container.Run(ctx, sigterm)

	logger.Info("shutdown with exit status 0")
	os.Exit(0)
}
