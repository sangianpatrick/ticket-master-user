package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sangianpatrick/tm-user/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	t.Run("should write object to json http response", func(t *testing.T) {
		res := httptest.NewRecorder()
		data := response.WebAPIEnvelope{
			Success: true,
			Status:  "OK",
			Message: "just for test",
			Data:    "data",
			Meta:    "meta",
			Errors:  nil,
		}
		dataBuff, _ := json.Marshal(data)

		err := response.JSON(res, http.StatusOK, data)
		assert.NoError(t, err, "should not be an error")
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "should be 200 OK")
		assert.JSONEq(t, string(dataBuff), res.Body.String(), "should match with web api envelope")
	})

}
