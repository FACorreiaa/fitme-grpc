package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpclog"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcspan"
)

// BootstrapClient creates a gRPC client connection with basic OTel + logging interceptors.
func BootstrapClient(
	address string,
	log *zap.Logger,
	traceProvider trace.TracerProvider, // [currently not used directly, but available if needed]
	promRegistry *prometheus.Registry, // [not used, but can be integrated if you'd like client metrics]
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {

	// OTel and logging interceptors.
	spanInterceptor, _ := grpcspan.Interceptors()
	logInterceptor, _ := grpclog.Interceptors(log)

	// Base gRPC dial options.
	connOptions := []grpc.DialOption{
		// For local dev/demo, we use insecure. For production, use TLS.
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Basic load balancing config (round_robin).
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),

		// Chain the unary interceptors.
		grpc.WithChainUnaryInterceptor(
			spanInterceptor.Unary,
			logInterceptor.Unary,
		),

		// OTel gRPC stats handler for more telemetry coverage.
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),

		// Chain the stream interceptors.
		grpc.WithChainStreamInterceptor(
			spanInterceptor.Stream,
			logInterceptor.Stream,
		),
	}

	// Append any additional dial options.
	connOptions = append(connOptions, opts...)

	return grpc.Dial(address, connOptions...)
}
