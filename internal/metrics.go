package internal

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	expo "go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

const meterName = "github.com/open-telemetry/opentelemetry-go/example/prometheus"

func setupPrometheusRegistry(ctx context.Context) (*prometheus.Registry, error) {
	// Initialize Prometheus registry
	reg := prometheus.NewRegistry()
	//#nosec
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	exporter, err := expo.New(expo.WithRegisterer(reg))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OpenTelemetry Prometheus exporter")
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(meterName)
	opt := api.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)
	// Register the promhttp handler for serving metrics
	http.Handle("/prometheus/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
	if err != nil {
		zap.Error(err)
	}
	counter.Add(ctx, 5, opt)

	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.Float64Histogram(
		"baz",
		api.WithDescription("a histogram with custom buckets and rename"),
		api.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	)
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 136, opt)
	histogram.Record(ctx, 64, opt)
	histogram.Record(ctx, 701, opt)
	histogram.Record(ctx, 830, opt)

	return reg, nil
}

func otelTraceProvider(ctx context.Context, insecure bool, caFile, certFile, keyFile, endpoint string) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx)

	if err != nil {
		zap.Error(err)
		return nil, err
	}

	var opts []otlptracegrpc.Option

	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		// Load TLS credentials if provided
		creds, err := credentials.NewClientTLSFromFile(caFile, "")
		if err != nil {
			return nil, err
		}
		opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("FITDEV"),
			semconv.DeploymentEnvironmentKey.String("production"),
		)),
	)

	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exp))
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			zap.Error(err)
			panic(err)
		}
	}()
	otel.SetTracerProvider(tracerProvider)
	return tp, nil
}
