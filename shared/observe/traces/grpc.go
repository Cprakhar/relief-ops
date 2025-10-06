package traces

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// WithTracingInterceptors returns gRPC server options with OpenTelemetry tracing interceptors.
func WithTracingInterceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(newServerHandler()),
	}
}

// DialOptionsWithTracing returns gRPC dial options with OpenTelemetry tracing interceptors.
func DialOptionsWithTracing() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(newClientHandler()),
	}
}

// newServerHandler creates a new gRPC server stats handler with OpenTelemetry tracing.
func newServerHandler() stats.Handler {
	return otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
}

// newClientHandler creates a new gRPC client stats handler with OpenTelemetry tracing.
func newClientHandler() stats.Handler {
	return otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}