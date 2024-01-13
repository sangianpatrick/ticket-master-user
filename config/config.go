package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	configSyncOnce sync.Once
	cfg            *Config
)

// Get returns config concrete object. It's done in singleton and thread-safe.
func Get() *Config {
	configSyncOnce.Do(func() {
		cfg = load()
	})

	return cfg
}

func load() *Config {
	c := new(Config)
	c.app()

	return c
}

func (c *Config) app() {
	c.App.Name = os.Getenv("APP_NAME")
	c.App.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	c.App.Location, _ = time.LoadLocation("APP_LOCATION")
}

// Config contains collection of the properties for the application configurations.
type Config struct {
	// App is a set of properties for the application such as name, port, location, etc.
	App struct {
		Name     string
		Port     int
		Location *time.Location
	}
	// CORS is a set of properties for Cross Origins Resource Sharing. It should be set to allow direct request from the external web application (from browser).
	CORS struct {
		AllowedOrigins []string
		AllowedMethods []string
		ExposedHeaders []string
	}
}
