package trace

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	oteltrace "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const ginTraceName = "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

func Tracer() oteltrace.Tracer {
	return otel.GetTracerProvider().Tracer(
		ginTraceName,
		oteltrace.WithInstrumentationVersion(otelgin.Version()),
	)
}

type OTelConfig struct {
	ServiceName    string
	ServiceVersion string
	TraceProvider  struct {
		Enable           bool
		OtlpHttpEndpoint string
		OtlpHttpInsecure bool
	}
	MetricProvider struct {
		Enable bool
	}
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, cfg OTelConfig) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up resource.
	res, err := newResource(cfg.ServiceName, cfg.ServiceVersion)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	var tracerProvider *trace.TracerProvider
	if cfg.TraceProvider.Enable {
		tracerProvider, err = newTraceProvider(ctx, cfg, res)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	// Set up meter provider.
	var meterProvider *metric.MeterProvider
	if cfg.MetricProvider.Enable {
		meterProvider, err = newMeterProvider(ctx, cfg, res)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}

	return
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newStdoutTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(3*time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newTraceProvider(ctx context.Context, cfg OTelConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	var options []otlptracehttp.Option
	options = append(options, otlptracehttp.WithEndpoint(cfg.TraceProvider.OtlpHttpEndpoint))
	if cfg.TraceProvider.OtlpHttpInsecure {
		options = append(options, otlptracehttp.WithInsecure())
	}
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cfg.TraceProvider.OtlpHttpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, cfg OTelConfig, res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}
