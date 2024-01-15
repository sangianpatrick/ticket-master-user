package middleware

import (
	"context"
	"net/http"

	"github.com/sangianpatrick/tm-user/pkg/device"
)

func ClientDeviceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		d := device.Device{
			RemoteAddress: r.RemoteAddr,
			UserAgent:     r.UserAgent(),
		}

		d.XForwardedFor = r.Header.Get("X-Forwarded-For")

		ctx = context.WithValue(ctx, device.DeviceContextKey{}, d)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}
