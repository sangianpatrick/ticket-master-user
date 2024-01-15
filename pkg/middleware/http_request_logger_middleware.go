package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type HTTPRequestLoggerMiddleware interface {
	Middleware(http.Handler) http.Handler
}

type unimplementHTTPRequestLoggerMiddleware struct{}

func (uhrlm *unimplementHTTPRequestLoggerMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

type httpRequestLoggerMiddleware struct {
	logger *logrus.Logger
}

func NewHTTPRequestLoggerMiddleware(logger *logrus.Logger, debug bool) HTTPRequestLoggerMiddleware {
	if debug {
		return &httpRequestLoggerMiddleware{logger: logger}
	}

	return &unimplementHTTPRequestLoggerMiddleware{}
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	recorder http.ResponseWriter
}

func (wrw wrappedResponseWriter) WriteHeader(statusCode int) {
	wrw.recorder.WriteHeader(statusCode)
	wrw.ResponseWriter.WriteHeader(statusCode)
}

func (wrw wrappedResponseWriter) Write(b []byte) (n int, err error) {
	wrw.recorder.Write(b)
	return wrw.ResponseWriter.Write(b)
}

func (hrlm *httpRequestLoggerMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		span := trace.SpanFromContext(r.Context())
		traceID := span.SpanContext().TraceID().String()

		recoreder := httptest.NewRecorder()

		wrappedResponseWriter := wrappedResponseWriter{w, recoreder}
		wrappedResponseWriter.ResponseWriter.Header().Set("X-Custom-Trace-ID", traceID)
		wrappedResponseWriter.recorder.Header().Set("X-Custom-Trace-ID", traceID)

		now := time.Now()
		handler.ServeHTTP(wrappedResponseWriter, r)
		elapsed := time.Since(now)

		requestHeader := r.Header

		requestBody := new(bytes.Buffer)
		io.Copy(requestBody, r.Body)

		result := recoreder.Result()

		defer result.Body.Close()
		responseBodyBuff, _ := io.ReadAll(result.Body)

		responseHeader := w.Header().Clone()

		captured := logrus.Fields{}
		captured["http.method"] = r.Method
		captured["http.url"] = r.RequestURI
		captured["http.request.body"] = requestBody.String()
		captured["http.status_code"] = result.StatusCode
		for reqHeaderKey, reqHeaderCol := range requestHeader {
			captured[fmt.Sprintf("http.request.header.%s", strings.ReplaceAll(strings.ToLower(reqHeaderKey), " ", "_"))] = strings.Join(reqHeaderCol, ",")
		}
		for resHeaderKey, resHeaderCol := range responseHeader {
			captured[fmt.Sprintf("http.response.header.%s", strings.ReplaceAll(strings.ToLower(resHeaderKey), " ", "_"))] = strings.Join(resHeaderCol, ",")
		}
		captured["http.response.body"] = string(responseBodyBuff)
		captured["time_consumption"] = elapsed.String()
		captured["custom_trace_id"] = traceID

		hrlm.logger.WithContext(r.Context()).WithFields(captured).Info()
	})
}
