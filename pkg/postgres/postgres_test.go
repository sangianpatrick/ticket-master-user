package postgres_test

import (
	"testing"

	"github.com/sangianpatrick/tm-user/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabase(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "root")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DBNAME", "mydb")
	t.Setenv("POSTGRES_SSLMODE", "disable")

	t.Run("try to build connection and get the db object", func(t *testing.T) {
		db := postgres.GetDatabase()
		assert.NotNil(t, db, "db object should not be null")
	})
	t.Run("db object is singleton", func(t *testing.T) {
		db1 := postgres.GetDatabase()
		db2 := postgres.GetDatabase()

		assert.Equal(t, db1, db2, "both db1 and db2 should have the same reference")
	})
}
