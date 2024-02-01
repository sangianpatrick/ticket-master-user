package container_test

import (
	"context"
	"os"
	"testing"

	"github.com/sangianpatrick/tm-user/internal/container"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Setenv("OTEL_COLLECTOR_ENDPOINT", "localhost:4444")

	ctx := context.Background()
	shutdownChan := make(chan os.Signal, 1)
	shutdownChan <- os.Interrupt

	hasStop := container.Run(ctx, shutdownChan)

	assert.Equal(t, struct{}{}, <-hasStop)
}
