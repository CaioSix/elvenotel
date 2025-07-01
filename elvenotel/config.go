package elvenotel

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	ServiceName         string        `env:"OTEL_SERVICE_NAME" envDefault:"caio"`
	OTLPEndpoint        string        `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
	ResourceAttributes  string        `env:"OTEL_RESOURCE_ATTRIBUTES"`
	ServiceVersion      string        `env:"SERVICE_VERSION" envDefault:"0.0.1"`
	Enabled             bool          `env:"TELEMETRY_ENABLED" envDefault:"true"`
	OTELExporterTimeout time.Duration `env:"OTEL_EXPORTER_OTLP_TIMEOUT" envDefault:"10s"`
	OTELLogLevel        string        `env:"OTEL_LOG_LEVEL" envDefault:"info"`
	OTELPropagators     string        `env:"OTEL_PROPAGATORS" envDefault:"tracecontext,baggage"`
	LokiAppName         string        `env:"LOKI_APP_NAME" envDefault:"go-app"`
	LokiAuthToken       string        `env:"LOKI_AUTH_TOKEN"`
	LokiFlushTimeout    int           `env:"LOKI_FLUSH_TIMEOUT" envDefault:"2000"`
	LokiTenantID        string        `env:"LOKI_TENANT_ID" envDefault:"my-tenant"`
	LokiURL             string        `env:"LOKI_URL"`
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse telemetry config: %w", err)
	}
	return cfg, nil
}
