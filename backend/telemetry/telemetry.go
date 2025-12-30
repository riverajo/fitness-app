package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitOTel initializes the OpenTelemetry tracer and logger providers.
func InitOTel(ctx context.Context, appEnv string, enableAlloy bool) (func(context.Context) error, error) {
	// 1. Create the Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("fitness-app"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 2. Metrics (OTLP)
	// 2. Metrics
	var metricReader metric.Reader

	if appEnv == "production" && enableAlloy {
		metricsExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to create otlp metric exporter: %w", err)
		}
		metricReader = metric.NewPeriodicReader(metricsExporter)
	} else {
		// Dev/Test: No metrics or stdout
		// For now we can just use a manual reader that does nothing effectively for this test case
		// or just skip setting the global provider if we want.
		// Let's use a manual reader to satisfy the interface without network calls.
		metricReader = metric.NewManualReader()
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return nil, fmt.Errorf("failed to start runtime metrics: %w", err)
	}

	var traceExporter sdktrace.SpanExporter
	loggerProviderOptions := []log.LoggerProviderOption{
		log.WithResource(res),
	}

	if appEnv == "production" && enableAlloy {
		// Production: Use OTLP Exporters (Alloy)
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}

		// 1. OTLP Log Exporter
		otlpLogExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
		}
		loggerProviderOptions = append(loggerProviderOptions, log.WithProcessor(log.NewBatchProcessor(otlpLogExporter)))

		// 2. Stdout Log Exporter (for Docker logs)
		stdoutLogExporter, err := stdoutlog.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
		}
		loggerProviderOptions = append(loggerProviderOptions, log.WithProcessor(log.NewBatchProcessor(stdoutLogExporter)))

	} else {
		// Development: Use Stdout Exporters (Console)
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}

		stdoutLogExporter, err := stdoutlog.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
		}
		loggerProviderOptions = append(loggerProviderOptions, log.WithProcessor(log.NewBatchProcessor(stdoutLogExporter)))
	}

	// 2. Create the TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// 3. Create the LoggerProvider
	lp := log.NewLoggerProvider(loggerProviderOptions...)
	global.SetLoggerProvider(lp)

	// 4. Set the global Propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return a shutdown function that flushes both traces and logs.
	return func(ctx context.Context) error {
		var errs []error
		if err := tp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown TracerProvider: %w", err))
		}
		if err := lp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown LoggerProvider: %w", err))
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown MeterProvider: %w", err))
		}
		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	}, nil
}
