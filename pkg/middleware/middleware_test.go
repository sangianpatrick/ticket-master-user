package middleware_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sangianpatrick/tm-user/pkg/middleware"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `json:"name"`
}

type LogData struct {
	Method       string          `json:"http.method"`
	RequestBody  json.RawMessage `json:"http.request.body"`
	ResponseBody string          `json:"http.response.body"`
}

func TestNewHTTPRequestLoggerMiddleware(t *testing.T) {
	t.Run("when debug mode is false", func(t *testing.T) {
		qs := map[string]string{
			"orderBy": "name",
		}
		user := User{
			Name: "TestUser1",
		}
		userBuff := new(bytes.Buffer)
		json.NewEncoder(userBuff).Encode(user)

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", userBuff)
		queryString := r.URL.Query()
		for k, v := range qs {
			queryString.Add(k, v)
		}
		r.URL.RawQuery = queryString.Encode()

		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
		})

		middleware.NewHTTPRequestLoggerMiddleware(logrus.New(), false).Middleware(handler).ServeHTTP(recorder, r)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("when debug mode is true", func(t *testing.T) {
		logDataBuff := new(bytes.Buffer)
		logger := logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetOutput(logDataBuff)

		qs := map[string]string{
			"orderBy": "name",
		}
		user := User{
			Name: "TestUser1",
		}
		userBuff := new(bytes.Buffer)
		json.NewEncoder(userBuff).Encode(user)

		r := httptest.NewRequest(http.MethodPost, "/just/for/testing", userBuff)
		r.SetBasicAuth("testuser", "testpassword")

		queryString := r.URL.Query()
		for k, v := range qs {
			queryString.Add(k, v)
		}
		r.URL.RawQuery = queryString.Encode()

		recorder := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
		})

		middleware.NewHTTPRequestLoggerMiddleware(logger, true).Middleware(handler).ServeHTTP(recorder, r)

		assert.Equal(t, http.StatusOK, recorder.Code)

		logData := LogData{}
		json.Unmarshal(logDataBuff.Bytes(), &logData)

		t.Log(logDataBuff.String())

		assert.Equal(t, http.MethodPost, logData.Method)
		assert.Equal(t, "OK", logData.ResponseBody)
	})
}
