package middleware

import (
	"bytes"
	"encoding/json"
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

		requestHeader := r.Header

		buf, _ := io.ReadAll(r.Body)
		rcCopied1 := io.NopCloser(bytes.NewBuffer(buf))
		rcCopied2 := io.NopCloser(bytes.NewBuffer(buf))

		r.Body = rcCopied1

		now := time.Now()
		handler.ServeHTTP(wrappedResponseWriter, r)
		elapsed := time.Since(now)

		requestBodyData := make(map[string]interface{})
		json.NewDecoder(rcCopied2).Decode(&requestBodyData)

		result := recoreder.Result()

		defer result.Body.Close()

		responseBodyData := make(map[string]interface{})
		json.NewDecoder(result.Body).Decode(&responseBodyData)

		responseHeader := w.Header().Clone()

		captured := logrus.Fields{}
		captured["http.method"] = r.Method
		captured["http.url"] = r.RequestURI
		captured["http.request.body"] = requestBodyData
		captured["http.status_code"] = result.StatusCode
		for reqHeaderKey, reqHeaderCol := range requestHeader {
			captured[fmt.Sprintf("http.request.header.%s", strings.ReplaceAll(strings.ToLower(reqHeaderKey), " ", "_"))] = strings.Join(reqHeaderCol, ",")
		}
		for resHeaderKey, resHeaderCol := range responseHeader {
			captured[fmt.Sprintf("http.response.header.%s", strings.ReplaceAll(strings.ToLower(resHeaderKey), " ", "_"))] = strings.Join(resHeaderCol, ",")
		}
		captured["http.response.body"] = responseBodyData
		captured["time_consumption"] = elapsed.String()
		captured["custom_trace_id"] = traceID

		entry := hrlm.logger.WithContext(r.Context()).WithFields(captured)
		if result.StatusCode < 200 || result.StatusCode > 299 {
			entry.Error()
			return
		}
		entry.Info()
	})
}
