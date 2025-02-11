package tracing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Config struct {
	Enabled            bool          `env:"TRACE_ENABLED"`
	Host               string        `env:"TRACE_HOST,notEmpty"`
	Port               string        `env:"TRACE_PORT,notEmpty"`
	BatchTimeout       time.Duration `env:"TRACE_BATCH_TIMEOUT" envDefault:"5000ms"`
	MaxExportBatchSize int           `env:"TRACE_MAX_EXPORT_BATCH_SIZE" envDefault:"512"`
	HttpInsecure       bool          `env:"TRACE_HTTP_INSECURE"`
}

// Provide provides configuration for tracing
func (cfg *Config) Provide(ctx context.Context, serviceName string) error {
	if cfg.Enabled {
		var traceOptions []otlptracehttp.Option
		if cfg.HttpInsecure {
			traceOptions = append(traceOptions, otlptracehttp.WithInsecure())
		}
		traceOptions = append(traceOptions,
			otlptracehttp.WithEndpoint(
				fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			),
			otlptracehttp.WithHeaders(map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}),
		)
		exporter, err := otlptrace.New(
			ctx,
			otlptracehttp.NewClient(traceOptions...),
		)
		if err != nil {
			return err
		}
		tracerProvider := trace.NewTracerProvider(
			trace.WithBatcher(
				exporter,
				trace.WithBatchTimeout(cfg.BatchTimeout),
				trace.WithMaxExportBatchSize(cfg.MaxExportBatchSize),
			),
			trace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(serviceName),
				),
			),
		)
		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
	} else {
		slog.InfoContext(ctx, "Tracing disabled")
	}
	return nil
}
