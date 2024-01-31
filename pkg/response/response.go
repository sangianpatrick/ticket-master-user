package response

import (
	"encoding/json"
	"net/http"

	"github.com/sangianpatrick/tm-user/pkg/appstatus"
)

type WebAPIEnvelope struct {
	Success bool             `json:"success" xml:"success"`
	Status  appstatus.Status `json:"status" xml:"status"`
	Message string           `json:"message" xml:"message"`
	Data    interface{}      `json:"data,omitempty" xml:"data,omitempty"`
	Meta    interface{}      `json:"meta,omitempty" xml:"meta,omitempty"`
	Errors  []interface{}    `json:"errors,omitempty" xml:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, object interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(object)
}
