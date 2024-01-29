package customer_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sangianpatrick/tm-user/internal/customer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepositorySave(t *testing.T) {
	t.Run("success with transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		cr := customer.NewCustomerRepository(logrus.New(), db, time.UTC)
		ctx := context.Background()

		newCustomer := customer.Customer{
			Email:      "patrick@test.com",
			Password:   "passwordtest",
			Name:       "Patrick Test",
			IsVerified: false,
			IsDeleted:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}

		expectedQuery := "INSERT INTO customers"

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQuery).
			ExpectQuery().
			WithArgs(newCustomer.Email, newCustomer.Password, newCustomer.Name, newCustomer.IsVerified, newCustomer.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(int64(1)))
		mock.ExpectCommit()

		tx, err := db.BeginTx(ctx, nil)
		assert.NoError(t, err, "should not be an error when begin transaction")
		id, err := cr.Save(ctx, newCustomer, tx)
		assert.NoError(t, err, "should not be an error when save customer data")
		tx.Commit()

		assert.Equal(t, int64(1), id, "should be int64 with value '1'")
	})
	t.Run("error caused by network connection to db", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		cr := customer.NewCustomerRepository(logrus.New(), db, time.UTC)
		ctx := context.Background()

		newCustomer := customer.Customer{
			Email:      "patrick@test.com",
			Password:   "passwordtest",
			Name:       "Patrick Test",
			IsVerified: false,
			IsDeleted:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}

		expectedQuery := "INSERT INTO customers"

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQuery).
			ExpectQuery().
			WithArgs(newCustomer.Email, newCustomer.Password, newCustomer.Name, newCustomer.IsVerified, newCustomer.CreatedAt).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		tx, err := db.BeginTx(ctx, nil)
		assert.NoError(t, err, "should not be an error when begin transaction")
		_, err = cr.Save(ctx, newCustomer, tx)
		assert.Error(t, err, "should be an error when saving customer data caused by network connection")
		tx.Rollback()
	})
	t.Run("error when preparing query", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		cr := customer.NewCustomerRepository(logrus.New(), db, time.UTC)
		ctx := context.Background()

		newCustomer := customer.Customer{
			Email:      "patrick@test.com",
			Password:   "passwordtest",
			Name:       "Patrick Test",
			IsVerified: false,
			IsDeleted:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}

		expectedQuery := "INSERT INTO customers"

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQuery).
			WillReturnError(fmt.Errorf("sql: invalid statement just for test"))
		mock.ExpectRollback()

		tx, err := db.BeginTx(ctx, nil)
		assert.NoError(t, err, "should not be an error when begin transaction")
		_, err = cr.Save(ctx, newCustomer, tx)
		assert.Error(t, err, "should be an error when saving customer data caused by query preparation failed")
		tx.Rollback()
	})
}
