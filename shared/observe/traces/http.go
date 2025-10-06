package traces

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func GinTracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName,
		otelgin.WithTracerProvider(otel.GetTracerProvider()),
		otelgin.WithPropagators(otel.GetTextMapPropagator()),
	)
}
