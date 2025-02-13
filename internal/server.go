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
	mlpb "github.com/FACorreiaa/fitme-protos/modules/meal/generated"
	mpb "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	upb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
	wpb "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"

	config "github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/logger"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpctracing"
)

// --- Server components

// isReady is used for kube liveness probes, it's only latched to true once
// the gRPC server is ready to handle requests
var isReady atomic.Value

func ServeGRPC(ctx context.Context, port string, container *ServiceContainer, reg *prometheus.Registry) error {
	log := logger.Log

	// Initialize OpenTelemetry trace provider with options as needed
	err := grpctracing.InitOTELToCollector(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to configure OpenTelemetry trace provider")
	}

	tp := otel.GetTracerProvider()

	// Bootstrap the gRPC server
	server, listener, err := grpc.BootstrapServer(port, log, reg, tp)
	if err != nil {
		return errors.Wrap(err, "failed to configure gRPC server")
	}

	// Register your services
	//cpb.RegisterCustomerServer(server, container.CustomerService)
	upb.RegisterAuthServer(server, container.AuthService)
	ccpb.RegisterCalculatorServer(server, container.CalculatorService)
	apb.RegisterActivityServer(server, container.ServiceActivity)
	wpb.RegisterWorkoutServer(server, container.WorkoutService)
	mpb.RegisterUserMeasurementsServer(server, container.MeasurementService)

	mlpb.RegisterMealPlanServer(server, container.MealServices.MealPlanService)
	mlpb.RegisterDietPreferenceServiceServer(server, container.MealServices.DietPreferenceService)
	mlpb.RegisterFoodLogServiceServer(server, container.MealServices.FoodLogService)
	mlpb.RegisterIngredientsServer(server, container.MealServices.IngredientService)
	mlpb.RegisterTrackMealProgressServer(server, container.MealServices.TrackMealProgressService)
	mlpb.RegisterGoalRecommendationServer(server, container.MealServices.GoalRecommendationService)
	mlpb.RegisterMealReminderServer(server, container.MealServices.MealReminderService)

	// meal services
	//mealServices := []mlpb.MealServer{
	//	container.MealServices.MealPlanService,
	//	container.MealServices.MealService,
	//	container.MealServices.DietPreferenceService,
	//	container.MealServices.FoodLogService,
	//	container.MealServices.IngredientService,
	//	container.MealServices.TrackMealProgressService,
	//	container.MealServices.GoalRecommendationService,
	//	container.MealServices.MealReminderService,
	//}
	//
	//for _, service := range mealServices {
	//	mlpb.RegisterMealServer(server, service)
	//}

	//mlpb.RegisterMealServer(server, container.MealServices.MealService)
	// Enable gRPC reflection for easier debugging
	reflection.Register(server)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Warn("shutting down grpc server")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// Start serving
	log.Info("gRPC server starting", zap.String("port", port))
	if err = server.Serve(listener); err != nil {
		return errors.Wrap(err, "gRPC server failed to serve")
	}

	isReady.Store(true)
	logger.Log.Info("running grpc server", zap.String("port", port))

	return server.Serve(listener)
}

// ServeHTTP creates a simple server to serve Prometheus metrics for
// the collector, and (not included) healthcheck endpoints for K8S to
// query readiness. By default, these should serve on "/healthz" and "/readyz"
func ServeHTTP(port string, reg *prometheus.Registry) error {
	log := logger.Log
	log.Info("running http server", zap.String("port", port))

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("failed to initialize config", zap.Error(err))
		return err
	}

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
	server.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

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
