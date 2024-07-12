package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is the global tracer
var Tracer = otel.GetTracerProvider().Tracer("gin-server")

// TracingMiddleware is a Gin middleware for tracing
func TracingMiddleware(c *gin.Context) {
	// Create a carrier to extract context
	carrier := propagation.HeaderCarrier(c.Request.Header)
	ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), carrier)

	// Define span options
	spanOptions := []trace.SpanStartOption{
		trace.WithAttributes(semconv.HTTPMethodKey.String(c.Request.Method)),
		trace.WithAttributes(semconv.HTTPTargetKey.String(c.FullPath())),
		trace.WithAttributes(semconv.HTTPRouteKey.String(c.FullPath())),
		trace.WithAttributes(semconv.HTTPURLKey.String(fmt.Sprintf("%s://%s%s", c.Request.URL.Scheme, c.Request.Host, c.Request.RequestURI))),
		trace.WithAttributes(semconv.UserAgentOriginalKey.String(c.Request.UserAgent())),
		trace.WithAttributes(semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength)),
		trace.WithAttributes(semconv.HTTPSchemeKey.String(c.Request.URL.Scheme)),
		trace.WithAttributes(semconv.NetTransportTCP),
		trace.WithSpanKind(trace.SpanKindServer),
	}

	// Start a new span with the extracted context and options
	ctx, span := Tracer.Start(ctx, c.Request.Method+" "+c.FullPath(), spanOptions...)
	defer span.End()

	// Set the context with the new span in the request
	c.Request = c.Request.WithContext(ctx)

	// Proceed to the next middleware/handler
	c.Next()

	// Inject headers back into response
	propagator := otel.GetTextMapPropagator()
	carrier = propagation.HeaderCarrier{}
	propagator.Inject(ctx, carrier)

	for k, v := range carrier {
		for _, value := range v {
			c.Writer.Header().Set(k, value)
		}
	}

	// Set HTTP status code in span attributes
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Writer.Status()))
}
