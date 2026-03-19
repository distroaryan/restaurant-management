package observability

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitTelemetry() func(context.Context) error {
	exporter, err := prometheus.New()
	if err != nil {
		slog.Error("Failed to initialize prometheus exporter", slog.String("error", err.Error()))
		panic(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	slog.Info("OpenTelemetry Prometheus exporter initialized successfully")

	return func(ctx context.Context) error {
		// Flush remaining telemetry data before shutdown
		return provider.Shutdown(ctx)
	}
}
