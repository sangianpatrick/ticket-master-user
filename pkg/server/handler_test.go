package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sangianpatrick/tm-user/pkg/appstatus"
	"github.com/sangianpatrick/tm-user/pkg/response"
	"github.com/sangianpatrick/tm-user/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	t.Run("should return http response that written on not found handler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/not-found-route", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := server.NotFoundHandler(logrus.New())

		handler.ServeHTTP(rr, req)

		var respBody response.WebAPIEnvelope

		json.NewDecoder(rr.Body).Decode(&respBody)

		assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode, "http response status code should be 404 Not Found")
		assert.False(t, respBody.Success, "response body on field `success` should be FALSE")
		assert.Equal(t, respBody.Status, appstatus.NotFound, "response body on field `status` should be \"NotFound\"")
	})
}

func TestIndexHandler(t *testing.T) {
	t.Run("should return http response that written on index handler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/index", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := server.IndexHandler()

		handler.ServeHTTP(rr, req)

		var respBody response.WebAPIEnvelope

		json.NewDecoder(rr.Body).Decode(&respBody)

		assert.Equal(t, http.StatusOK, rr.Result().StatusCode, "http response status code should be 200 OK")
		assert.True(t, respBody.Success, "response body on field `success` should be TRUE")
		assert.Equal(t, respBody.Status, appstatus.OK, "response body on field `status` should be \"OK\"")
	})
}
