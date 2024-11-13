package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	apb "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	ccpb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	cpb "github.com/FACorreiaa/fitme-protos/modules/customer/generated"
	mpb "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	upb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
	wpb "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"

	config "github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/logger"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc"
)

// --- Server components

// isReady is used for kube liveness probes, it's only latched to true once
// the gRPC server is ready to handle requests
var isReady atomic.Value

func ServeGRPC(ctx context.Context, port string, container *ServiceContainer) error {
	log := logger.Log

	// Configure Prometheus registry and trace provider
	promRegistry, err := setupPrometheusRegistry(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure Prometheus registry")
	}

	// Initialize OpenTelemetry trace provider with options as needed
	traceProvider, err := SetupTracerProvider(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure OpenTelemetry trace provider")
	}

	// Ensure TracerProvider shuts down properly on exit
	defer func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			log.Error("failed to shut down trace provider", zap.Error(err))
		}
	}()

	// Bootstrap the gRPC server
	server, listener, err := grpc.BootstrapServer(port, log, promRegistry, traceProvider)
	if err != nil {
		return errors.Wrap(err, "failed to configure gRPC server")
	}

	// Register your services
	cpb.RegisterCustomerServer(server, container.CustomerService)
	upb.RegisterAuthServer(server, container.AuthService)
	ccpb.RegisterCalculatorServer(server, container.CalculatorService)
	apb.RegisterActivityServer(server, container.ServiceActivity)
	wpb.RegisterWorkoutServer(server, container.WorkoutService)
	mpb.RegisterUserMeasurementsServer(server, container.MeasurementService)
	// Enable gRPC reflection for easier debugging

	reflection.Register(server)

	// Start serving
	log.Info("gRPC server starting", zap.String("port", port))
	if err := server.Serve(listener); err != nil {
		return errors.Wrap(err, "gRPC server failed to serve")
	}

	return nil
}

// ServeHTTP creates a simple server to serve Prometheus metrics for
// the collector, and (not included) healthcheck endpoints for K8S to
// query readiness. By default, these should serve on "/healthz" and "/readyz"
func ServeHTTP(port string) error {
	log := logger.Log
	log.Info("running http server", zap.String("port", port))

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("failed to initialize config", zap.Error(err))
		return err
	}
	//ctx := context.Background()
	//collectors := grpcprometheus.NewPrometheusMetricsCollectors()
	//
	//// Set up Prometheus registry and register collectors
	//registry := prometheus.NewRegistry()
	//if err = grpcprometheus.RegisterMetrics(registry, collectors); err != nil {
	//	log.Fatal("failed to register Prometheus metrics")
	//}
	//
	//if err = grpcprometheus.SetupTracing(ctx); err != nil {
	//	log.Fatal("failed to set up tracing")
	//}
	//
	//if err != nil {
	//	log.Error("failed to initialize config", zap.Error(err))
	//	return err
	//}
	promRegistry, err := setupPrometheusRegistry(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to configure Prometheus registry")
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(promRegistry)

	server := http.NewServeMux()
	// Add healthcheck endpoints
	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Perform health check logic here
		// For example, check if the server is healthy
		// Respond with appropriate status code
		w.WriteHeader(http.StatusOK)
	})

	server.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Perform readiness check logic here
		// For example, check if the server is ready to receive traffic
		// Respond with appropriate status code
		w.WriteHeader(http.StatusOK)
	})

	//server.HandleFunc("/metrics", promhttp.Handler().ServeHTTP) // This should use the correct registry.
	server.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	listener := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: cfg.Server.Timeout,
		Handler:           server,
	}

	if err := listener.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "failed to create telemetry server")
	}

	return nil
}
