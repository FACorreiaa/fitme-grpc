package grpcprometheus

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Collectors struct {
	Client *grpcprom.ClientMetrics
	Server *grpcprom.ServerMetrics
}

func NewPrometheusMetricsCollectors() *Collectors {
	return &Collectors{
		Client: clientMetrics(),
		Server: serverMetrics(),
	}
}

//func otelTraceProvider(ctx context.Context, insecure bool, caFile, certFile, keyFile, endpoint string) (*trace.TracerProvider, error) {
//	var opts []otlptracegrpc.Option
//
//	// auth header
//	opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
//		"Authorization": "Bearer " + os.Getenv("OTEL_EXPORTER_API_KEY"),
//	}))
//	opts = append(opts, otlptracegrpc.WithEndpoint("1056716.https://otlp-gateway-prod-eu-west-2.grafana.net/otlp"))
//
//	// If insecure is set to true, use an insecure gRPC connection
//	if insecure {
//		opts = append(opts, otlptracegrpc.WithInsecure())
//	} else {
//		// Load TLS credentials if provided
//		creds, err := credentials.NewClientTLSFromFile(caFile, "")
//		if err != nil {
//			return nil, err
//		}
//		opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
//	}
//
//	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
//
//	exp, err := otlptracegrpc.New(ctx, opts...)
//	if err != nil {
//		return nil, err
//	}
//
//	res := resource.NewWithAttributes(
//		semconv.SchemaURL,
//		semconv.ServiceNameKey.String("FITDEV"),
//		semconv.DeploymentEnvironmentKey.String("production"),
//	)
//
//	tp := trace.NewTracerProvider(
//		trace.WithBatcher(exp),
//		trace.WithResource(res),
//	)
//
//	otel.SetTracerProvider(tp)
//
//	return tp, nil
//}

func otelTraceProvider(ctx context.Context, endpoint, apiKey string) (*trace.TracerProvider, error) {
	var opts []otlptracegrpc.Option

	// auth header
	opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
		"Authorization": "Bearer " + apiKey,
	}))

	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))

	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("FITDEV"),
		semconv.DeploymentEnvironmentKey.String("production"),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

// RegisterMetrics must be called before the Prometheus interceptor is used.
func RegisterMetrics(registry *prometheus.Registry, collectors *Collectors) error {
	if registry == nil {
		return errors.New("must provide a Prometheus registry")
	}

	if collectors == nil {
		return errors.New("must provide Prometheus collectors")
	}

	if collectors.Client != nil {
		registry.MustRegister(collectors.Client)
	}

	if collectors.Server != nil {
		registry.MustRegister(collectors.Server)
	}

	return nil
}

func SetupTracing(ctx context.Context) (*trace.TracerProvider, error) {
	apiKey := os.Getenv("OTEL_EXPORTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing OTEL_EXPORTER_API_KEY environment variable")
	}

	tp, err := otelTraceProvider(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("OTEL_EXPORTER_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}

	// Ensure the TracerProvider is properly closed when the application shuts down
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = tp.Shutdown(shutdownCtx); err != nil {
			log.Printf("failed to shut down trace provider: %v", err)
		}
	}()

	return tp, nil
}

// grpcHandlingTimeHistogramBuckets is the default set of buckets used by both
// server and client histograms.
var grpcHandlingTimeHistogramBuckets = []float64{
	0.001, 0.01, 0.1, 0.3,
	0.6, 1, 3, 6, 9, 20,
	30, 60, 90, 120,
}

// clientMetrics attaches prometheus metrics to the grpc client
func clientMetrics() *grpcprom.ClientMetrics {
	return grpcprom.NewClientMetrics(
		grpcprom.WithClientHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(grpcHandlingTimeHistogramBuckets),
		),
	)
}

// clientMetrics attaches prometheus metrics to the grpc server
func serverMetrics() *grpcprom.ServerMetrics {
	return grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(grpcHandlingTimeHistogramBuckets),
		),
	)
}
