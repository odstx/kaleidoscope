package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"

	"kaleidoscope/config"
)

type Telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	logger         *zap.Logger
}

func InitTelemetry(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Telemetry, error) {
	if !cfg.OTEL.Enabled {
		logger.Info("OpenTelemetry is disabled")
		return &Telemetry{logger: logger}, nil
	}

	exporterEndpoint := fmt.Sprintf("%s:%s", cfg.OTEL.ExporterHost, cfg.OTEL.ExporterPort)
	logger.Info("Initializing OpenTelemetry",
		zap.String("service_name", cfg.OTEL.ServiceName),
		zap.String("endpoint", exporterEndpoint))

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(exporterEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.OTEL.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("OpenTelemetry initialized successfully")

	return &Telemetry{
		tracerProvider: tracerProvider,
		logger:         logger,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.tracerProvider == nil {
		return nil
	}

	t.logger.Info("Shutting down OpenTelemetry")
	if err := t.tracerProvider.Shutdown(ctx); err != nil {
		t.logger.Error("Failed to shutdown OpenTelemetry", zap.Error(err))
		return err
	}

	t.logger.Info("OpenTelemetry shutdown completed")
	return nil
}
