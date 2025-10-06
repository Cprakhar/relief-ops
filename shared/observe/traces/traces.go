package traces

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerConfig struct {
	ServiceName      string
	Environment      string
	Secure           bool
	ExporterEndpoint string
}

var once sync.Once

// InitTrace initializes OpenTelemetry tracing with the given configuration, sets global tracer provider, and returns a shutdown function.
func InitTrace(ctx context.Context, cfg *TracerConfig) (func(context.Context) error, error) {
	// Ensure propagator is set only once
	once.Do(func() {
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
	})
	// Create trace exporter (OTLP over HTTP)
	traceExporter := newTraceExporter(ctx, cfg.ExporterEndpoint, cfg.Secure)

	// Create tracer provider
	tp, err := newTracProvider(ctx, cfg, traceExporter)
	if err != nil {
		return nil, err
	}

	// Set global tracer provider
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// GetTracer returns a tracer for the given name.
func GetTracer(name string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name)
}

// newTraceExporter creates a new OTLP trace exporter.
func newTraceExporter(ctx context.Context, exporterEndpoint string, secure bool) sdktrace.SpanExporter {
	var opts []otlptracehttp.Option
	if !secure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	opts = append(opts,
		otlptracehttp.WithEndpoint(exporterEndpoint),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(opts...),
	)
	if err != nil {
		return nil
	}
	return exporter
}

// newTracProvider creates a new tracer provider with the given configuration and exporter.
func newTracProvider(ctx context.Context, cfg *TracerConfig, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	// Create resource with service attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentNameKey.String(cfg.Environment),
			semconv.ServiceNamespaceKey.String("relief-ops"),
			semconv.HostArchAMD64,
			semconv.OSTypeLinux,
			semconv.TelemetrySDKLanguageGo,
		),
		resource.WithHost(),
		resource.WithOSType(),
		resource.WithTelemetrySDK(),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)

	return tp, nil
}
