package server_test

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/sangianpatrick/tm-user/pkg/server"
	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	t.Run("it should run the server and subsequently shut it down gracefully.", func(t *testing.T) {
		buf := new(bytes.Buffer)

		logger := logrus.New()
		logger.SetOutput(buf)
		logger.SetFormatter(&logrus.JSONFormatter{})

		router := http.NewServeMux()
		port := 8080
		name := "test-server"

		httpServer := server.NewServer(logger, router, name, port)
		httpServer.Start()

		time.Sleep(time.Second * 3)

		httpServer.Stop()
	})
}
