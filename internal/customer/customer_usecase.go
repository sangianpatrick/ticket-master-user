package customer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sangianpatrick/tm-user/pkg/apperror"
	"github.com/sangianpatrick/tm-user/pkg/appstatus"
	"github.com/sangianpatrick/tm-user/pkg/utils"
	"github.com/sirupsen/logrus"
)

type CustomerUsecaseProps struct {
	Logger             *logrus.Logger
	Location           *time.Location
	CustomerRepository CustomerRepository
}

type CustomerUsecase interface {
	Register(ctx context.Context, params RegisterParams) (err error)
}

type customerUsecaseImpl struct {
	logger             *logrus.Logger
	location           *time.Location
	customerRepository CustomerRepository
}

func NewCustomerUsecase(props CustomerUsecaseProps) CustomerUsecase {
	return &customerUsecaseImpl{
		logger:             props.Logger,
		location:           props.Location,
		customerRepository: props.CustomerRepository,
	}
}

func (cr *customerUsecaseImpl) checkCustomer(ctx context.Context, params RegisterParams) (err error) {
	existingCustomer, err := cr.customerRepository.FindByEmail(ctx, params.Email, nil)
	if err == nil {
		if existingCustomer.IsDeleted {
			return nil
		}
		err = apperror.SetError(fmt.Errorf("customer exists"), appstatus.Conflict, http.StatusConflict)
		return
	}

	apperr := apperror.Destruct(err)

	if apperr.Status != appstatus.NotFound {
		return err
	}

	return nil
}

// Register will register the customer into the application and the customer data will be saved.
// The registered customers should verify the their account by click the url that has been sent to their own email.
func (cr *customerUsecaseImpl) Register(ctx context.Context, params RegisterParams) (err error) {
	if err = cr.checkCustomer(ctx, params); err != nil {
		return
	}

	now := time.Now()
	passwordSalt := utils.GenerateUniqueID(utils.AlphaNumeric, 16)
	cryptedPassword := utils.Encrypt(params.Password, passwordSalt)

	newCustomer := Customer{
		ID:           0,
		Email:        params.Email,
		Password:     cryptedPassword,
		PasswordSalt: passwordSalt,
		Name:         params.Name,
		IsVerified:   false,
		IsDeleted:    false,
		CreatedAt:    now,
	}

	ID, err := cr.customerRepository.Save(ctx, newCustomer, nil)
	if err != nil {
		return
	}

	newCustomer.ID = ID

	registerEvent := RegisterEvent{}
	registerEvent.Customer = newCustomer
	registerEvent.Verification.ExpiredAt = now.Add(5 * time.Minute)

	return
}
