package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpclog"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcspan"
)

const (
	component      = "grpc-example"
	httpAddr       = ":8082"
	targetGRPCAddr = "localhost:8080"
)

func BootstrapClient(
	address string,
	log *zap.Logger,
	traceProvider trace.TracerProvider,
	promRegistry *prometheus.Registry,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	// -- OpenTelemetry interceptor setup
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	//clMetrics := grpcprom.NewClientMetrics(
	//	grpcprom.WithClientHandlingTimeHistogram(
	//		grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	//	),
	//)
	//promRegistry.MustRegister(clMetrics)

	//exemplarFromContext := func(ctx context.Context) prometheus.Labels {
	//	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
	//		return prometheus.Labels{"traceID": span.TraceID().String()}
	//	}
	//	return nil
	//}

	spanInterceptor, _ := grpcspan.Interceptors()

	// -- Zap logging interceptor setup
	logInterceptor, _ := grpclog.Interceptors(log)

	// EXPERIMENTAL
	//promCollectors := grpcprometheus.NewPrometheusMetricsCollectors()
	//_ = grpcprometheus.RegisterMetrics(promRegistry, promCollectors)
	//clientInterceptor, _, _ := grpcprometheus.Interceptors(promCollectors)

	// trace to grafana
	//ctx := context.Background()
	//if _, err := grpcprometheus.SetupTracing(ctx); err != nil {
	//	log.Error("Failed to set up trace exporter", zap.Error(err))
	//	return nil, errors.Wrap(err, "failed to setup tracing")
	//}

	// -- Prometheus exporter setup
	//prometheusCollectors := grpcprometheus.NewPrometheusMetricsCollectors()
	//if err := grpcprometheus.RegisterMetrics(prometheus, prometheusCollectors); err != nil {
	//	return nil, errors.Wrap(err, "failed to register grpc metrics")
	//}

	//clientMetrics := prometheusCollectors.Client

	connOptions := []grpc.DialOption{
		// We terminate TLS in the linkerd sidecar, so no need for TLS on the listen server
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// default config
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),

		//	"methodConfig": [{
		//	"name": [{"service": "your.package.ServiceName"}],
		//	"retryPolicy": {
		//	"maxAttempts": 5,
		//	"initialBackoff": "0.1s",
		//	"maxBackoff": "1s",
		//	"backoffMultiplier": 2,
		//	"retryableStatusCodes": ["UNAVAILABLE"]
		//}}]

		// Add the unary interceptors
		grpc.WithChainUnaryInterceptor(
			spanInterceptor.Unary,
			logInterceptor.Unary,
			//clMetrics.UnaryClientInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			//promCollectors.Client.UnaryClientInterceptor(),
			//clientInterceptor.Unary,
		),

		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		// Add the stream interceptors
		grpc.WithChainStreamInterceptor(
			spanInterceptor.Stream,
			logInterceptor.Stream,
			//clMetrics.StreamClientInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			//promCollectors.Client.StreamClientInterceptor(),
			//clientInterceptor.Stream,
		),
	}

	// Add any additional options
	connOptions = append(connOptions, opts...)

	return grpc.NewClient(address, connOptions...)
}
