package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/sangianpatrick/tm-user/config"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var db *sql.DB
var dbSyncOnce sync.Once

func GetDatabase() *sql.DB {
	dbSyncOnce.Do(func() {
		db = buildConnection()
	})

	return db
}

func buildConnection() *sql.DB {
	cfg := config.Get()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	conn, err := otelsql.Open(
		"postgres",
		dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("ticket-master"),
	)

	if err != nil {
		log.Println(err)
		return nil
	}

	conn.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)

	return conn
}
