package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// TelemetryMiddleware adds OpenTelemetry metrics
func TelemetryMiddleware() gin.HandlerFunc {
	meter := otel.Meter("api-gateway")
	requestCounter, _ := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)

	return func(c *gin.Context) {
		requestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("method", c.Request.Method),
			attribute.String("path", c.Request.URL.Path),
		))

		c.Next()
	}
}