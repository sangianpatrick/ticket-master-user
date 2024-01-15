package device

// DeviceContextKey is context key for tracing. It will be injected to the given context as the key.
type DeviceContextKey struct{}

// Device contains information from requester device.
type Device struct {
	UserAgent     string
	RemoteAddress string
	XForwardedFor string
}
