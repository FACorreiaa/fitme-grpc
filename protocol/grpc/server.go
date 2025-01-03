package grpc

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpclog"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcprometheus"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcratelimit"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcrecovery"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcrequest"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcspan"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/session"
)

// BootstrapServer creates a new gRPC server with the base middleware included.
// We ignore the first return value from each of the interceptors as we only
// need the server interceptors.
func BootstrapServer(
	port string,
	log *zap.Logger,
	registry *prometheus.Registry,
	traceProvider trace.TracerProvider,
	opts ...grpc.ServerOption,
) (*grpc.Server, net.Listener, error) {
	// setup basic metrics

	// initiate the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create listener")
	}

	// set redis cache
	//newCache, _ := grpccacherequests.NewCache()

	// -- Prometheus interceptor Setup
	promCollectors := grpcprometheus.NewPrometheusMetricsCollectors()
	// Must be called before using Prometheus interceptors
	err = grpcprometheus.RegisterMetrics(registry, promCollectors)
	_, promInterceptor, err := grpcprometheus.Interceptors(promCollectors)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Prometheus interceptors")
	}

	// -- OpenTelemetry interceptor setup
	//otel.SetTracerProvider(traceProvider)
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
	//	propagation.TraceContext{},
	//	propagation.Baggage{},
	//))

	_, spanInterceptor := grpcspan.Interceptors()
	_, serverInterceptor, _ := grpcprometheus.Interceptors(promCollectors)
	// -- Zap logging interceptor setup
	_, logInterceptor := grpclog.Interceptors(log)
	// -- Recovery interceptor setup
	_, recoveryInterceptor := grpcrecovery.Interceptors(grpcrecovery.RegisterMetrics(registry))
	// session
	sessionInterceptor := session.InterceptorSession()
	requestGeneratorSession := grpcrequest.RequestIDMiddleware()
	//globalRateLimiter := grpcratelimit.Interceptors()

	rateLimiter := grpcratelimit.NewRateLimiter(10, 20)
	//cache := grpccacherequests.UnaryCachingInterceptor(newCache)
	// Configure server options from our base configuration
	serverOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(middleware.KeepaliveEnforcementPolicy()),
		grpc.KeepaliveParams(middleware.KeepAliveServerParams()),

		// Note: Order of interceptors for the connection matters.
		// Trace interceptor must be the first interceptor in the chain for cross-service spans
		// to work correctly as it needs the original inbound context.

		// Add the unary interceptors
		grpc.ChainUnaryInterceptor(
			spanInterceptor.Unary,
			promInterceptor.Unary,
			serverInterceptor.Unary,
			logInterceptor.Unary,
			sessionInterceptor,
			requestGeneratorSession,
			recoveryInterceptor.Unary,
			rateLimiter.UnaryServerInterceptor(),
		),

		//grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// Add the stream interceptors
		grpc.ChainStreamInterceptor(
			spanInterceptor.Stream,
			promInterceptor.Stream,
			serverInterceptor.Stream,
			logInterceptor.Stream,
			recoveryInterceptor.Stream,
			recoveryInterceptor.Stream,
		),
	}

	// Append any additional options passed in while connecting
	serverOptions = append(serverOptions, opts...)

	server := grpc.NewServer(serverOptions...)

	return server, listener, nil
}
