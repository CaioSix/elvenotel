package elvenotel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// newLoggerProvider cria um logger provider com suporte a exportação
func newLoggerProvider(ctx context.Context, cfg Config) (*log.LoggerProvider, error) {
	res := newResource(cfg.ServiceName, cfg.ServiceVersion)

	// Cria o exportador OTLP com configurações adequadas
	exporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlploggrpc.WithTimeout(cfg.OTELExporterTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Configura o processador em lote
	processor := log.NewBatchProcessor(exporter)

	lp := log.NewLoggerProvider(
		log.WithProcessor(processor),
		log.WithResource(res),
	)

	return lp, nil
}

// newMeterProvider cria um meter provider com suporte a exportação
func newMeterProvider(ctx context.Context, cfg Config) (*metric.MeterProvider, error) {
	res := newResource(cfg.ServiceName, cfg.ServiceVersion)

	// Cria o exportador OTLP com configurações adequadas
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlpmetricgrpc.WithTimeout(cfg.OTELExporterTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return mp, nil
}

// newTracerProvider cria um tracer provider com suporte a exportação
func newTracerProvider(ctx context.Context, cfg Config) (*trace.TracerProvider, error) {
	res := newResource(cfg.ServiceName, cfg.ServiceVersion)

	// Cria o exportador OTLP com configurações adequadas
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithTimeout(cfg.OTELExporterTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

// newResource cria um resource OTEL
func newResource(serviceName string, serviceVersion string) *resource.Resource {
	hostName, _ := os.Hostname()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.HostName(hostName),
	)
}
