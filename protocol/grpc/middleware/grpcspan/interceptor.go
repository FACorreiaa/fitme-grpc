package grpcspan

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware"
)

// Interceptors wraps OpenTelemetry gRPC interceptors.
// We use a wrapper for this as the OpenTelemetry ecosystem changes
// *very* frequently, so we want to contain this change to a single point
// rather than having to update the servers and clients individually.
func Interceptors() (middleware.ClientInterceptor, middleware.ServerInterceptor) {
	clientInterceptor := middleware.ClientInterceptor{
		// nolint:staticcheck
		Unary: otelgrpc.UnaryClientInterceptor(),
		// nolint:staticcheck
		Stream: otelgrpc.StreamClientInterceptor(),
	}

	serverInterceptor := middleware.ServerInterceptor{
		// nolint:staticcheck
		Unary: otelgrpc.UnaryServerInterceptor(),
		// nolint:staticcheck
		Stream: otelgrpc.StreamServerInterceptor(),
	}

	return clientInterceptor, serverInterceptor
}
