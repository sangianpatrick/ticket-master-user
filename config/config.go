package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
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
	c.logrus()

	return c
}

func (c *Config) app() {
	c.App.Name = os.Getenv("APP_NAME")
	c.App.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	c.App.Location, _ = time.LoadLocation("APP_LOCATION")
}

func (c *Config) logrus() {
	c.Logrus.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFunc:  "caller",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "timestamp",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return funcname, filename
		},
	}
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
	// Logrus is a logger package for logging. This object contains setup properties of logrus.
	Logrus struct {
		Formatter logrus.Formatter
	}
	// JWT contains both private and public key in buffered type.
	JWT struct {
		PrivateKey []byte
		PublicKey  []byte
	}
	// BasicAuth contains credentials for basic authentication.
	BasicAuth struct {
		Username string
		Password string
	}
}
