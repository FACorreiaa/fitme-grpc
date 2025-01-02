package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/FACorreiaa/fitme-protos/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/internal"
	"github.com/FACorreiaa/fitme-grpc/internal/metrics"
	"github.com/FACorreiaa/fitme-grpc/logger"
)

func initializeLogger() error {
	return logger.Init(
		zap.DebugLevel,
		zap.String("service", "example"),
		zap.String("version", "v42.0.69"),
		zap.Strings("maintainers", []string{"@fc", "@FACorreiaa"}),
	)
}

func setupDatabases(cfg *config.Config) (*pgxpool.Pool, *redis.Client, error) {
	dbConfig, err := internal.NewDatabaseConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize database configuration: %w", err)
	}

	pool, err := internal.Init(dbConfig.ConnectionURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize database pool: %w", err)
	}

	internal.WaitForDB(pool)
	logger.Log.Info("Connected to Postgres",
		zap.String("host", os.Getenv("POSTGRES_HOST")),
		zap.String("port", os.Getenv("POSTGRES_PORT")))

	redisClient, err := internal.NewRedisConfig()
	if err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to initialize Redis configuration: %w", err)
	}

	logger.Log.Info("Connected to Redis",
		zap.String("host", cfg.Repositories.Redis.Host),
		zap.String("port", cfg.Repositories.Redis.Port))

	if err = internal.Migrate(pool); err != nil {
		pool.Close()
		redisClient.Close()
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return pool, redisClient, nil
}

func startServices(ctx context.Context, cfg *config.Config, container *internal.ServiceContainer, reg *prometheus.Registry) error {
	errChan := make(chan error, 2)

	// Start gRPC server
	go func() {
		if err := internal.ServeGRPC(ctx, cfg.Server.GrpcPort, container, reg); err != nil {
			logger.Log.Error("gRPC server error", zap.Error(err))
			errChan <- err
		}
	}()
	logger.Log.Info("Serving gRPC", zap.String("port", cfg.Server.GrpcPort))

	// Start HTTP server
	go func() {
		if err := internal.ServeHTTP(cfg.Server.HTTPPort, reg); err != nil {
			logger.Log.Error("HTTP server error", zap.Error(err))
			errChan <- err
		}
	}()

	logger.Log.Info("Serving HTTP", zap.String("port", cfg.Server.HTTPPort))

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func run(cfg *config.Config) (*pgxpool.Pool, *redis.Client, error) {
	pool, redisClient, err := setupDatabases(cfg)
	if err != nil {
		return nil, nil, err
	}
	return pool, redisClient, nil
}

func main() {
	println("Fitme dev app starting...")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	reg := prometheus.NewRegistry()
	println("Loaded prometheus registry")

	if err := initializeLogger(); err != nil {
		panic("failed to initialize logging")
	}

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Log.Error("failed to initialize config", zap.Error(err))
		return
	}

	pool, redisClient, err := run(&cfg)
	if err != nil {
		logger.Log.Error("failed to run the application", zap.Error(err))
		return
	}
	defer pool.Close()
	defer redisClient.Close()

	tu := new(utils.TransportUtils)
	brokers := internal.ConfigureUpstreamClients(logger.Log, tu)
	if brokers == nil {
		logger.Log.Error("failed to configure brokers")
		return
	}

	metrics.InitPprof()

	container := internal.NewServiceContainer(ctx, pool, redisClient, brokers)

	if err := startServices(ctx, &cfg, container, reg); err != nil {
		logger.Log.Error("service error", zap.Error(err))
	}
}
