package grpctracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"

	"github.com/FACorreiaa/fitme-grpc/logger"
)

//func newLocalTracerProvider() (trace.SpanExporter, error) {
//	exporter, _ := stdouttrace.New()
//	return exporter, nil
//}

func newTracerProvider(endpoint, apiKey, caCertPath string, insecure bool) (*trace.TracerProvider, error) {
	var opts []otlptracegrpc.Option

	// Set endpoint
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	//opts = append(opts, otlptracegrpc.WithGRPCConn(conn))
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

	exp, err := otlptracegrpc.New(context.Background(), opts...)
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
	tp.Tracer("DeezNuts")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func InitTracer() (*trace.TracerProvider, error) {
	log := logger.Log
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if otlpEndpoint == "" {
		log.Error("You MUST set OTEL_EXPORTER_OTLP_TRACES_ENDPOINT env variable!")
	}

	tp, err := newTracerProvider(otlpEndpoint, "", "", true)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}

	// Ensure TracerProvider shuts down properly on exit
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err = tp.Shutdown(ctx); err != nil {
			log.Error("failed to shut down trace provider")
		}
	}()

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
