package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/FACorreiaa/fitme-protos/container"

	"github.com/FACorreiaa/fitme-grpc/logger"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

// --- Server components

// isReady is used for kube liveness probes, it's only latched to true once
// the gRPC server is ready to handle requests
var isReady atomic.Value

func ServeGRPC(ctx context.Context, port string, _ *container.Brokers, pgPool *pgxpool.Pool, redisClient *redis.Client) error {
	log := logger.Log

	// When you have a configured prometheus registry and OTEL trace provider,
	// pass in as param 3 & 4

	// configure prometheus registry
	registry, err := setupPrometheusRegistry(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure prometheus registry")
	}
	tp, err := otelTraceProvider(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure jaeger trace provider")
	}
	server, listener, err := grpc.BootstrapServer(port, logger.Log, registry, tp)
	if err != nil {
		return errors.Wrap(err, "failed to configure grpc server")
	}

	// Replace with your actual generated registration method
	// generated.RegisterDummyServer(server, implementation)
	// client := generated.NewCustomerClient(brokers.Customer)

	// customerService and any implementation is a dependency that is injected to dest and delete
	//customerService := domain.NewCustomerService(pgPool, redisClient)

	// implement brokers

	//authRepo := repository.NewAuthService(pgPool, redisClient)
	//authService := service.NewAuthService(authRepo)

	//cpb.RegisterCustomerServer(server, customerService)
	//upb.RegisterAuthServer(server, authService)

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

	log.Info("running http server", zap.String("port", port))

	//cfg, err := config.InitConfig()
	//
	//if err != nil {
	//	log.Error("failed to initialize config", zap.Error(err))
	//	return err
	//}

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

	//listener := &http.Server{
	//	Addr:              fmt.Sprintf(":%s", port),
	//	ReadHeaderTimeout: cfg.Server.Timeout,
	//	Handler:           server,
	//}
	//
	//if err := listener.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//	return errors.Wrap(err, "failed to create telemetry server")
	//}

	return nil
}
