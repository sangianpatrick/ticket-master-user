package customer

import (
	"context"
	"database/sql"
	"fmt"
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
	Update(ctx context.Context, ID int64, c Customer, tx *sql.Tx) (err error)
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

func (cr *customerRepositoryImpl) fetch(ctx context.Context, cmd sqlCommand, query string, args ...interface{}) (bunchOfCustomer []Customer, err error) {
	var stmt *sql.Stmt
	if stmt, err = cmd.PrepareContext(ctx, query); err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var c Customer
		var updatedAt sql.NullTime

		err = rows.Scan(
			&c.ID, &c.Email, &c.Password, &c.PasswordSalt, &c.Name, &c.IsVerified, &c.IsDeleted, &c.CreatedAt, &updatedAt,
		)

		if err != nil {
			cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
			err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
			return
		}

		if updatedAt.Valid {
			c.UpdatedAt = &updatedAt.Time
		}

		bunchOfCustomer = append(bunchOfCustomer, c)
	}

	return
}

// FindByEmail returns Customer object with that retrieved by the given email.
func (cr *customerRepositoryImpl) FindByEmail(ctx context.Context, email string, tx *sql.Tx) (c Customer, err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := `
		SELECT ids, email, name, password, password_salt, is_verified, is_deleted, created_at, updated_at FROM customers WHERE email = $1 ORDER BY id DESC LIMIT 1 
	`

	bunchOfCustomer, err := cr.fetch(ctx, cmd, query, email)
	if err != nil {
		return
	}

	dataLength := len(bunchOfCustomer)

	if dataLength < 1 {
		err = apperror.SetError(fmt.Errorf("customer is not found"), appstatus.NotFound, http.StatusNotFound)
		return
	}

	c = bunchOfCustomer[dataLength-1]

	return
}

// FindByID returns Customer object with that retrieved by the given ID.
func (cr *customerRepositoryImpl) FindByID(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := `
		SELECT id, email, password, password_salt, is_verified, is_deleted, created_at, updated_at FROM customers WHERE id = $1 FOR UPDATE
	`

	bunchOfCustomer, err := cr.fetch(ctx, cmd, query, ID)
	if err != nil {
		return
	}

	dataLength := len(bunchOfCustomer)

	if dataLength < 1 {
		err = apperror.SetError(fmt.Errorf("customer is not found"), appstatus.NotFound, http.StatusNotFound)
		return
	}

	c = bunchOfCustomer[dataLength-1]

	return
}

// FindByIDForUpdate behaves like `FindByID()` but it locked the row for modification in transactional mode.
func (cr *customerRepositoryImpl) FindByIDForUpdate(ctx context.Context, ID int64, tx *sql.Tx) (c Customer, err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := `
		SELECT id, email, password, password_salt, is_verified, is_deleted, created_at, updated_at FROM customers WHERE id = $1 FOR UPDATE 
	`

	bunchOfCustomer, err := cr.fetch(ctx, cmd, query, ID)
	if err != nil {
		return
	}

	dataLength := len(bunchOfCustomer)

	if dataLength < 1 {
		err = apperror.SetError(fmt.Errorf("customer is not found"), appstatus.NotFound, http.StatusNotFound)
		return
	}

	c = bunchOfCustomer[dataLength-1]

	return
}

// Save save Customer object and returns the generated sequential ID.
func (cr *customerRepositoryImpl) Save(ctx context.Context, c Customer, tx *sql.Tx) (ID int64, err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := "INSERT INTO customers (email, password, password_salt, name, is_verified, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	stmt, err := cmd.PrepareContext(ctx, query)
	if err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, c.Email, c.Password, c.PasswordSalt, c.Name, c.IsVerified, c.CreatedAt)
	if err = row.Scan(&ID); err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	return
}

// Update modify the existing Customer object with the new one based on the given ID.
func (cr *customerRepositoryImpl) Update(ctx context.Context, ID int64, c Customer, tx *sql.Tx) (err error) {
	var cmd sqlCommand = cr.db
	if tx != nil {
		cmd = tx
	}

	query := "UPDATE customers SET email = $1, password = $2, password_salt = $3, name = $4, is_verified = $5, is_deleted = $6, updated_at = $7 WHERE id = $8"

	stmt, err := cmd.PrepareContext(ctx, query)
	if err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, c.Email, c.Password, c.Name, c.IsVerified, c.IsDeleted, c.UpdatedAt, ID); err != nil {
		cr.logger.WithContext(ctx).WithField("query", query).WithError(err).Error()
		err = apperror.SetError(err, appstatus.InternalServerError, http.StatusInternalServerError)
		return
	}

	return
}
