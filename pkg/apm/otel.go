package apm

import (
	"context"
	"sync"

	"github.com/sangianpatrick/tm-user/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var otelImpl *openTelemetryImpl
var syncOnce sync.Once

// OpenTelemetry is wrapper of OpenTelemetry object. It contains behavior to start and stop the agent from tracing and collecting metrics.
type OpenTelemetry interface {
	// Start will start the OpenTelemetry to capture the trace from any request and collecting the metrics.
	Start(ctx context.Context) (err error)
	// Stop will stop the OpenTelemetry agent from capturing and collecting the traces/metrics.
	Stop(ctx context.Context) (err error)
}

type openTelemetryImpl struct {
	serviceName   string
	environment   string
	endpoint      string
	resource      *resource.Resource
	traceExporter *otlptrace.Exporter
}

func newOpenTelemetry(serviceName, environment, endpoint string) *openTelemetryImpl {
	ctx := context.Background()
	res, _ := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		),
	)

	return &openTelemetryImpl{
		serviceName: serviceName,
		environment: environment,
		endpoint:    endpoint,
		resource:    res,
	}
}

// GetOpenTelemetry returns object the implements the OpenTelemetry wrapper. It's done in singleton pattern and thread-safe. The object is instantiated once.
func GetOpenTelemetry() OpenTelemetry {
	cfg := config.Get()
	syncOnce.Do(func() {
		otelImpl = newOpenTelemetry(cfg.App.Name, cfg.App.Environment, cfg.OpenTelemetry.Collector.Endpoint)
	})

	return otelImpl
}

// Start will start the OpenTelemetry to capture the trace from any request and collecting the metrics.
func (ot *openTelemetryImpl) Start(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(ot.endpoint),
		otlptracegrpc.WithDialOption(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		otel.Handle(err)
		return
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(ot.resource),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	ot.traceExporter = exporter

	return
}

// Stop will stop the OpenTelemetry agent from capturing and collecting the traces/metrics.
func (ot *openTelemetryImpl) Stop(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	err = ot.traceExporter.Shutdown(ctx)
	if err != nil {
		otel.Handle(err)
		return
	}

	return
}
