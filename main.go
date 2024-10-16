package main

import (
	"context"
	"os"

	"github.com/FACorreiaa/fitme-protos/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	config "github.com/FACorreiaa/fitme-grpc/config"
	"github.com/FACorreiaa/fitme-grpc/internal"
	"github.com/FACorreiaa/fitme-grpc/internal/metrics"
	"github.com/FACorreiaa/fitme-grpc/logger"
)

func run() (*pgxpool.Pool, *redis.Client, error) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Log.Error("failed to initialize config", zap.Error(err))
		return nil, nil, err
	}

	log := logger.Log

	dbConfig, err := internal.NewDatabaseConfig()
	if err != nil {
		log.Error("failed to initialize database configuration", zap.Error(err))
		return nil, nil, err
	}

	pool, err := internal.Init(dbConfig.ConnectionURL)
	if err != nil {
		log.Error("failed to initialize database pool", zap.Error(err))
		return nil, nil, err
	}
	internal.WaitForDB(pool)
	log.Info("Connected to Postgres", zap.String("host", os.Getenv("DB_HOST")), zap.String("port", os.Getenv("DB_PORT")))

	redisClient, err := internal.NewRedisConfig()
	if err != nil {
		log.Error("failed to initialize Redis configuration", zap.Error(err))
		pool.Close()
		return nil, nil, err
	}

	log.Info("Connected to Redis", zap.String("host", cfg.Repositories.Redis.Host), zap.String("port", cfg.Repositories.Redis.Port))

	if err = internal.Migrate(pool); err != nil {
		log.Error("failed to migrate database", zap.Error(err))
		pool.Close()
		redisClient.Close()
		return nil, nil, err
	}

	return pool, redisClient, nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.InitConfig()
	if err != nil {
		zap.L().Error("failed to initialize config", zap.Error(err))
		return
	}

	if err := logger.Init(
		zap.DebugLevel,
		zap.String("service", "example"),
		zap.String("version", "v42.0.69"),
		zap.Strings("maintainers", []string{"@fc", "@FACorreiaa"}),
	); err != nil || logger.Log == nil {
		panic("failed to initialize logging")
	}

	log := logger.Log

	pool, redisClient, err := run()
	if err != nil {
		log.Error("failed to run the application", zap.Error(err))
		return
	}
	defer pool.Close()
	defer redisClient.Close()

	tu := new(utils.TransportUtils)

	brokers := internal.ConfigureUpstreamClients(log, tu)
	if brokers == nil {
		log.Error("failed to configure brokers")
		return
	}

	metrics.InitPprof()

	container := internal.NewServiceContainer(ctx, pool, redisClient, brokers)

	go func() {
		if err := internal.ServeGRPC(ctx, cfg.Server.GrpcPort, container); err != nil {
			log.Error("failed to serve grpc", zap.Error(err))
			return
		}
	}()

	if err := internal.ServeHTTP(cfg.Server.HTTPPort); err != nil {
		log.Error("failed to serve http", zap.Error(err))
		return
	}
}
