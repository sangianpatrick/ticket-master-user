package customer

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/sangianpatrick/tm-user/pkg/apperror"
	"github.com/sangianpatrick/tm-user/pkg/appstatus"
	"github.com/sirupsen/logrus"
)

// CustomerRepository contains bevahior of customer repository such as creation, data retrieval and modification.
type CustomerRepository interface {
	Save(ctx context.Context, c Customer, tx *sql.Tx) (ID int64, err error)
	FindByID(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error)
	FindByIDForUpdate(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error)
	FindByEmail(ctx context.Context, email string, tx *sql.Tx) (c Customer, err error)
	Update(ctx context.Context, ID int64, c Customer) (err error)
}

type sqlCommand interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type customerRepositoryImpl struct {
	logger *logrus.Logger
	db     *sql.DB
	loc    *time.Location
}

// NewCustomerRepository will instantiate new implementation of customer repository.
func NewCustomerRepository(logger *logrus.Logger, db *sql.DB, loc *time.Location) CustomerRepository {
	return &customerRepositoryImpl{
		logger: logger,
		db:     db,
		loc:    loc,
	}
}

func (cr *customerRepositoryImpl) exec(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (result sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = cmd.PrepareContext(ctx, query); err != nil {
		cr.logger.Error(query, err)
		return
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			cr.logger.Error(query, err)
		}
	}()

	if result, err = stmt.ExecContext(ctx, args...); err != nil {
		cr.logger.Error(query, err)
	}

	return
}

// FindByEmail returns Customer object with that retrieved by the given email.
func (cr *customerRepositoryImpl) FindByEmail(ctx context.Context, email string, tx *sql.Tx) (c Customer, err error) {
	panic("unimplemented")
}

// FindByID returns Customer object with that retrieved by the given ID.
func (*customerRepositoryImpl) FindByID(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error) {
	panic("unimplemented")
}

// FindByIDForUpdate behaves like `FindByID()` but it locked the row for modification in transactional mode.
func (*customerRepositoryImpl) FindByIDForUpdate(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error) {
	panic("unimplemented")
}

// Save save Customer object and returns the generated sequential ID.
func (cr *customerRepositoryImpl) Save(ctx context.Context, c Customer, tx *sql.Tx) (ID int64, err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := "INSERT INTO customers (email, password, name, is_verified, created_at) VALUES ($1, $2, $3, $4, $5)"

	stmt, err := cmd.PrepareContext(ctx, query)
	if err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, c.Email, c.Password, c.Name, c.IsVerified, c.CreatedAt)
	if err = row.Scan(&ID); err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	return
}

// Update modify the existing Customer object with the new one based on the given ID.
func (*customerRepositoryImpl) Update(ctx context.Context, ID int64, c Customer) (err error) {
	panic("unimplemented")
}
