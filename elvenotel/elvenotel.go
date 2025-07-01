package elvenotel

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TelemetryProvider interface {
	GetServiceName() string
	LogInfo(args ...interface{})
	LogErrorln(args ...interface{})
	LogFatalln(args ...interface{})
	MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error)
	MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error)
	TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span)
	LogRequest() gin.HandlerFunc
	MeterRequestDuration() gin.HandlerFunc
	MeterRequestsInFlight() gin.HandlerFunc
	Shutdown(ctx context.Context)
}

type Telemetry struct {
	lp     *log.LoggerProvider
	mp     *metric.MeterProvider
	tp     *sdktrace.TracerProvider
	log    *zap.SugaredLogger
	meter  otelmetric.Meter
	tracer oteltrace.Tracer
	cfg    Config
}

func NewTelemetry(ctx context.Context, cfg Config) (*Telemetry, error) {
	lp, err := newLoggerProvider(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger provider: %w", err)
	}

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
			otelzap.NewCore(cfg.ServiceName, otelzap.WithLoggerProvider(lp)),
		),
	)

	mp, err := newMeterProvider(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider: %w", err)
	}

	tp, err := newTracerProvider(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	return &Telemetry{
		lp:     lp,
		mp:     mp,
		tp:     tp,
		log:    logger.Sugar(),
		meter:  mp.Meter(cfg.ServiceName),
		tracer: tp.Tracer(cfg.ServiceName),
		cfg:    cfg,
	}, nil
}

func (t *Telemetry) GetServiceName() string {
	return t.cfg.ServiceName
}

func (t *Telemetry) LogInfo(args ...interface{}) {
	t.log.Info(args...)
}

func (t *Telemetry) LogErrorln(args ...interface{}) {
	t.log.Errorln(args...)
}

func (t *Telemetry) LogFatalln(args ...interface{}) {
	t.log.Fatalln(args...)
}

func (t *Telemetry) MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error) {
	histogram, err := t.meter.Int64Histogram(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}
	return histogram, nil
}

func (t *Telemetry) MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error) {
	counter, err := t.meter.Int64UpDownCounter(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}
	return counter, nil
}

func (t *Telemetry) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return t.tracer.Start(ctx, name)
}

func (t *Telemetry) Shutdown(ctx context.Context) {
	if t.lp != nil {
		_ = t.lp.Shutdown(ctx)
	}
	if t.mp != nil {
		_ = t.mp.Shutdown(ctx)
	}
	if t.tp != nil {
		_ = t.tp.Shutdown(ctx)
	}
}
