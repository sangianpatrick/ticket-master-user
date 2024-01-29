package config_test

import (
	"testing"

	"github.com/sangianpatrick/tm-user/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Setenv("APP_ENVIRONMENT", "DEV")
	c1 := config.Get()
	c2 := config.Get()
	t.Run("should return the same config object event call in multiple times", func(t *testing.T) {
		assert.Equal(t, c1, c2, "both of c1 and c2 should have the same memory address")
	})

	t.Run("should cover the logurs formatter function", func(t *testing.T) {
		logger := logrus.New()
		logger.SetFormatter(c1.Logrus.Formatter)

		logger.Info("test logrus formatter")
	})

	t.Run("should load from env", func(t *testing.T) {
		assert.Equal(t, "DEV", c1.App.Environment, "should be \"DEV\"")
	})
}
