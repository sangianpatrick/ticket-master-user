package config

import (
	"encoding/json"
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
	appEnvironment string = "PRODUCTION"
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
	c.openTelemetry()
	c.logrus()
	c.jwt()
	c.basicAuth()
	c.postgres()

	return c
}

func (c *Config) app() {
	c.App.Name = os.Getenv("APP_NAME")
	c.App.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	c.App.Environment = appEnvironment
	c.App.Location, _ = time.LoadLocation("APP_LOCATION")

	environment := strings.ToUpper(strings.ReplaceAll(os.Getenv("APP_ENVIRONMENT"), " ", ""))
	if environment != "" {
		c.App.Environment = environment
	}

	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	c.App.Debug = debug
}

func (cfg *Config) openTelemetry() {
	cfg.OpenTelemetry.Collector.Endpoint = os.Getenv("OTEL_COLLECTOR_ENDPOINT")
}

func (c *Config) logrus() {
	c.Logrus.Formatter = &logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		DataKey:           "",
		FieldMap:          logrus.FieldMap{logrus.FieldKeyFunc: "caller", logrus.FieldKeyLevel: "severity", logrus.FieldKeyTime: "timestamp"},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return funcname, filename
		},
		PrettyPrint: false,
	}
}

func (c *Config) jwt() {
	jwtRsaPlain := os.Getenv("JWT_RSA")
	var jwtRsa = struct {
		PrivateKey string `json:"private"`
		PublicKey  string `json:"public"`
	}{}

	json.Unmarshal([]byte(jwtRsaPlain), &jwtRsa)

	c.JWT.PrivateKey = []byte(jwtRsa.PrivateKey)
	c.JWT.PublicKey = []byte(jwtRsa.PublicKey)
}

func (c *Config) basicAuth() {
	c.BasicAuth.Username = os.Getenv("BASIC_AUTH_USERNAME")
	c.BasicAuth.Password = os.Getenv("BASIC_AUTH_PASSWORD")
}

func (c *Config) postgres() {
	c.Postgres.Host = os.Getenv("POSTGRES_HOST")
	c.Postgres.Port, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	c.Postgres.User = os.Getenv("POSTGRES_USER")
	c.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Postgres.DBName = os.Getenv("POSTGRES_DBNAME")
	c.Postgres.SSLMode = os.Getenv("POSTGRES_SSLMODE")
	c.Postgres.MaxOpenConns, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	c.Postgres.MaxIdleConns, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
}

// Config contains collection of the properties for the application configurations.
type Config struct {
	// App is a set of properties for the application such as name, port, location, etc.
	App struct {
		Name        string
		Port        int
		Environment string
		Location    *time.Location
		Debug       bool
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
	// OpenTelemetry contains creds and other properties that will be needs to connect and run agent for application tracing and metrics collection.
	OpenTelemetry struct {
		Collector struct {
			Endpoint string
		}
	}
	Postgres struct {
		Host         string
		Port         int
		User         string
		Password     string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
	}
}
