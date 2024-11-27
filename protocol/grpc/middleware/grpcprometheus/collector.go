package grpcprometheus

import (
	"context"
	"fmt"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc/credentials"
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

func otelTraceProvider(ctx context.Context, endpoint, apiKey, caCertPath string, insecure bool) (*trace.TracerProvider, error) {
	var opts []otlptracegrpc.Option

	// Set endpoint
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	opts = append(opts, otlptracegrpc.WithInsecure())

	// Handle insecure or TLS configuration
	if insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		c, err := credentials.NewClientTLSFromFile(caCertPath, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS credentials: %w", err)
		}
		opts = append(opts, otlptracegrpc.WithTLSCredentials(c))
	}

	// Add authorization header if needed (uncomment if using API keys)
	// opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
	// 	"Authorization": "Bearer " + apiKey,
	// }))

	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("FITDEV"),
		semconv.ServiceNamespaceKey.String("FitME"),
		semconv.ServiceVersionKey.String("0.1"),
		semconv.DeploymentEnvironmentKey.String("production"),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

//func otelTraceProvider(ctx context.Context, endpoint string) (*trace.TracerProvider, error) {
//	var opts []otlptracegrpc.Option
//
//	// Set endpoint
//	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
//	opts = append(opts, otlptracegrpc.WithInsecure())
//
//	// Add authorization header if needed (uncomment if using API keys)
//	// opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
//	// 	"Authorization": "Bearer " + apiKey,
//	// }))
//
//	exp, err := otlptracegrpc.New(ctx, opts...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
//	}
//
//	res := resource.NewWithAttributes(
//		semconv.SchemaURL,
//		semconv.ServiceNameKey.String("fitme-dev"),
//		semconv.ServiceNamespaceKey.String("FitME"),
//		semconv.ServiceVersionKey.String("0.1"),
//		semconv.DeploymentEnvironmentKey.String("production"),
//	)
//
//	tp := trace.NewTracerProvider(
//		trace.WithBatcher(exp),
//		trace.WithResource(res),
//		trace.WithSampler(trace.AlwaysSample()),
//	)
//
//	otel.SetTracerProvider(tp)
//	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
//	defer func() { _ = exp.Shutdown(context.Background()) }()
//
//	return tp, nil
//}

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

//func SetupTracing(ctx context.Context) (*trace.TracerProvider, error) {
//	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
//	if endpoint == "" {
//		return nil, errors.New("missing OTEL_EXPORTER_OTLP_ENDPOINT environment variable")
//	}
//
//	insecure := os.Getenv("OTEL_EXPORTER_INSECURE") == "true"
//	caCertPath := os.Getenv("OTEL_EXPORTER_CA_CERT_PATH")
//	apiKey := os.Getenv("OTEL_EXPORTER_API_KEY")
//
//	tp, err := otelTraceProvider(ctx, endpoint, apiKey, caCertPath, insecure)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create trace provider: %w", err)
//	}
//
//	return tp, nil
//}

func SetupTracing(ctx context.Context) (*trace.TracerProvider, error) {
	//endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	//if endpoint == "" {
	//	return nil, errors.New("missing OTEL_EXPORTER_OTLP_ENDPOINT environment variable")
	//}
	//
	//insecure := os.Getenv("OTEL_EXPORTER_INSECURE") == "true"
	//caCertPath := os.Getenv("OTEL_EXPORTER_CA_CERT_PATH")
	//apiKey := os.Getenv("OTEL_EXPORTER_API_KEY")

	tp, err := otelTraceProvider(ctx, "http://0.0.0.0:16686/", "", "", true)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}

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
