package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"

	apb "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	ccpb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	cpb "github.com/FACorreiaa/fitme-protos/modules/customer/generated"
	upb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
	wpb "github.com/FACorreiaa/fitme-protos/modules/workout/generated"

	config "github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/logger"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

// --- Server components

// isReady is used for kube liveness probes, it's only latched to true once
// the gRPC server is ready to handle requests
var isReady atomic.Value

func ServeGRPC(ctx context.Context, port string, container *ServiceContainer) error {
	log := logger.Log
	// dependencies

	//customerService := domain.NewCustomerService(pgPool, redisClient)
	//
	//// implement brokers
	//
	//sessionManager := auth.NewSessionManager(pgPool, redisClient)
	//
	//authRepo := repository.NewAuthService(pgPool, redisClient, sessionManager)
	//authService := service.NewAuthService(authRepo, pgPool, redisClient, sessionManager)

	// When you have a configured prometheus registry and OTEL trace provider,
	// pass in as param 3 & 4

	// configure prometheus registry
	registry, err := setupPrometheusRegistry(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure prometheus registry")
	}
	tp, err := otelTraceProvider(ctx, true, "", "", "", "localhost:7077/metrics")
	if err != nil {
		return errors.Wrap(err, "failed to configure jaeger trace provider")
	}
	server, listener, err := grpc.BootstrapServer(port, logger.Log, registry, tp, container.AuthService.SessionManager)
	if err != nil {
		return errors.Wrap(err, "failed to configure grpc server")
	}

	// Replace with your actual generated registration method
	//generated.RegisterDummyServer(server, implementation)
	//client := generated.NewCustomerClient(brokers.Customer)

	//customerService and any implementation is a dependency that is injected to dest and delete

	cpb.RegisterCustomerServer(server, container.CustomerService)
	upb.RegisterAuthServer(server, container.AuthService)
	ccpb.RegisterCalculatorServer(server, container.CalculatorService)
	apb.RegisterActivityServer(server, container.ServiceActivity)
	wpb.RegisterWorkoutServer(server, container.WorkoutService)
	// Enable reflection to be able to use grpcui or insomnia without
	// having to manually maintain .proto files

	reflection.Register(server)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Warn("shutting down grpc server")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	isReady.Store(true)

	log.Info("running grpc server", zap.String("port", port))
	return server.Serve(listener)
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
	//reg := prometheus.NewRegistry()

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
	server.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	//server.Handle("/prometheus/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

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
