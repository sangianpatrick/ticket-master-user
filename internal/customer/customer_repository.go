package customer

import (
	"context"
	"database/sql"

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

type customerRepositoryImpl struct {
	logger *logrus.Logger
	db     *sql.DB
}

// NewCustomerRepository will instantiate new implementation of customer repository.
func NewCustomerRepository(logger *logrus.Logger, db *sql.DB) CustomerRepository {
	return &customerRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

// FindByEmail returns Customer object with that retrieved by the given email.
func (*customerRepositoryImpl) FindByEmail(ctx context.Context, email string, tx *sql.Tx) (c Customer, err error) {
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
func (*customerRepositoryImpl) Save(ctx context.Context, c Customer, tx *sql.Tx) (ID int64, err error) {
	panic("unimplemented")
}

// Update modify the existing Customer object with the new one based on the given ID.
func (*customerRepositoryImpl) Update(ctx context.Context, ID int64, c Customer) (err error) {
	panic("unimplemented")
}
