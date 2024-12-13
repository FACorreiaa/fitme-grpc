package grpctracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/FACorreiaa/fitme-grpc/logger"
)

//func newTracerProvider(endpoint, apiKey, caCertPath string, insecure bool) (*trace.TracerProvider, error) {
//	var opts []otlptracegrpc.Option
//
//	// Set endpoint
//	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
//	//opts = append(opts, otlptracegrpc.WithGRPCConn(conn))
//	// Handle insecure or TLS configuration
//	if insecure {
//		opts = append(opts, otlptracegrpc.WithInsecure())
//	} else {
//		c, err := credentials.NewClientTLSFromFile(caCertPath, "")
//		if err != nil {
//			return nil, fmt.Errorf("failed to create TLS credentials: %w", err)
//		}
//		opts = append(opts, otlptracegrpc.WithTLSCredentials(c))
//	}
//
//	// Add authorization header if needed (uncomment if using API keys)
//	// opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
//	// 	"Authorization": "Bearer " + apiKey,
//	// }))
//
//	exp, err := otlptracegrpc.New(context.Background(), opts...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
//	}
//
//	res := resource.NewWithAttributes(
//		semconv.SchemaURL,
//		semconv.ServiceNameKey.String("fitme-app-dev"),
//		semconv.ServiceName("fitme-app-dev"),
//		semconv.ServiceVersionKey.String("0.1"),
//	)
//
//	tp := trace.NewTracerProvider(
//		trace.WithBatcher(exp),
//		trace.WithResource(res),
//	)
//
//	otel.SetTracerProvider(tp)
//	tp.Tracer("DeezNuts")
//	otel.SetTextMapPropagator(propagation.TraceContext{})
//
//	return tp, nil
//}
//
//func InitTracer() (*trace.TracerProvider, error) {
//	log := logger.Log
//	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
//	if otlpEndpoint == "" {
//		log.Error("You MUST set OTEL_EXPORTER_OTLP_TRACES_ENDPOINT env variable!")
//	}
//
//	tp, err := newTracerProvider(otlpEndpoint, "", "", true)
//
//	if err != nil {
//		return nil, fmt.Errorf("failed to create trace provider: %w", err)
//	}
//
//	// Ensure TracerProvider shuts down properly on exit
//	go func() {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//		defer cancel()
//		if err = tp.Shutdown(ctx); err != nil {
//			log.Error("failed to shut down trace provider")
//		}
//	}()
//
//	return tp, nil
//}

func NewOTLPExporter(ctx context.Context) (trace.SpanExporter, error) {
	// Change default HTTPS -> HTTP
	log := logger.Log
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	fmt.Printf("otlp endpoint %s\n", otlpEndpoint)
	if otlpEndpoint == "" {
		log.Error("You MUST set OTEL_EXPORTER_OTLP_TRACES_ENDPOINT env variable!")
	}

	insecureOpt := otlptracehttp.WithInsecure()

	// Update default OTLP reciver endpoint
	endpointOpt := otlptracehttp.WithEndpoint(otlpEndpoint)

	//timeout := otlptracehttp.WithTimeout(30 * time.Second)

	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

func NewTraceProvider(exp trace.SpanExporter) *trace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fitme-app-dev")))

	if err != nil {
		panic(err)
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r))
}
