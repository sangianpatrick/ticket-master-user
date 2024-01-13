package server

import (
	"fmt"
	"net/http"

	"github.com/sangianpatrick/tm-user/pkg/response"
	"github.com/sirupsen/logrus"
)

func NotFoundHandler(logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.WithContext(ctx).WithFields(logrus.Fields{
			"status": "InvalidURI",
			"path":   r.RequestURI,
		}).Info("attempt to unregistered path")

		errorMessage := fmt.Sprintf("this URI '%s' is not registered", r.RequestURI)

		resp := response.WebAPIEnvelope{
			Success: false,
			Status:  "NotFound",
			Message: "attempt to invalid / unregistered path",
			Errors: []interface{}{
				errorMessage,
			},
		}

		response.JSON(w, http.StatusOK, resp)
	})
}
