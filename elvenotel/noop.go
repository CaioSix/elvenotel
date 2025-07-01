package elvenotel

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type NoopTelemetry struct {
	serviceName string
}

func NewNoopTelemetry(cfg Config) (*NoopTelemetry, error) {
	return &NoopTelemetry{serviceName: cfg.ServiceName}, nil
}

func (t *NoopTelemetry) GetServiceName() string         { return t.serviceName }
func (t *NoopTelemetry) LogInfo(args ...interface{})    {}
func (t *NoopTelemetry) LogErrorln(args ...interface{}) {}
func (t *NoopTelemetry) LogFatalln(args ...interface{}) { os.Exit(1) }
func (t *NoopTelemetry) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
func (t *NoopTelemetry) MeterRequestDuration() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
func (t *NoopTelemetry) MeterRequestsInFlight() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
func (t *NoopTelemetry) TraceStart(ctx context.Context, name string) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}
func (t *NoopTelemetry) MeterInt64Histogram(metric Metric) (metric.Int64Histogram, error) {
	return nil, nil
}
func (t *NoopTelemetry) MeterInt64UpDownCounter(metric Metric) (metric.Int64UpDownCounter, error) {
	return nil, nil
}
func (t *NoopTelemetry) Shutdown(ctx context.Context) {}
