package customer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sangianpatrick/tm-user/pkg/apperror"
	"github.com/sangianpatrick/tm-user/pkg/appstatus"
	"github.com/sangianpatrick/tm-user/pkg/appvalidator"
	"github.com/sangianpatrick/tm-user/pkg/response"
	"github.com/sirupsen/logrus"
)

type CustomerHTTPHandler struct {
	Logger          *logrus.Logger
	CustomerUsecase CustomerUsecase
}

func InitCustomerHTTPHandler(logger *logrus.Logger, router *mux.Router, customerUsecase CustomerUsecase) {
	handler := &CustomerHTTPHandler{
		Logger:          logger,
		CustomerUsecase: customerUsecase,
	}

	router.HandleFunc("/ticket-master/v1/customerapp/customers/register", handler.Register).Methods(http.MethodPost)
}

func (handler *CustomerHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := RegisterParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		response.JSON(w, http.StatusUnprocessableEntity, response.WebAPIEnvelope{
			Success: false,
			Status:  appstatus.UnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	if err := appvalidator.ValidateStruct(ctx, params); err != nil {
		response.JSON(w, http.StatusBadRequest, response.WebAPIEnvelope{
			Success: false,
			Status:  appstatus.BadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := handler.CustomerUsecase.Register(ctx, params); err != nil {
		apperr := apperror.Destruct(err)
		response.JSON(w, apperr.HTTPStatusCode, response.WebAPIEnvelope{
			Success: false,
			Status:  apperr.Status,
			Message: apperr.Error(),
		})
		return
	}

	response.JSON(w, http.StatusCreated, response.WebAPIEnvelope{
		Success: true,
		Status:  appstatus.Created,
		Message: "customer is successfuly registered and need to verify by email",
	})
}
