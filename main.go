package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/FACorreiaa/fitme-protos/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/FACorreiaa/fitme-grpc/config"
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

	zapLogger := logger.Log

	dbConfig, err := internal.NewDatabaseConfig()
	if err != nil {
		zapLogger.Error("failed to initialize database configuration", zap.Error(err))
		return nil, nil, err
	}

	pool, err := internal.Init(dbConfig.ConnectionURL)
	if err != nil {
		zapLogger.Error("failed to initialize database pool", zap.Error(err))
		return nil, nil, err
	}
	internal.WaitForDB(pool)
	zapLogger.Info("Connected to Postgres", zap.String("host", os.Getenv("DB_HOST")), zap.String("port", os.Getenv("DB_PORT")))

	redisClient, err := internal.NewRedisConfig()
	if err != nil {
		zapLogger.Error("failed to initialize Redis configuration", zap.Error(err))
		pool.Close()
		return nil, nil, err
	}

	zapLogger.Info("Connected to Redis", zap.String("host", cfg.Repositories.Redis.Host), zap.String("port", cfg.Repositories.Redis.Port))

	if err = internal.Migrate(pool); err != nil {
		zapLogger.Error("failed to migrate database", zap.Error(err))
		pool.Close()
		redisClient.Close()
		return nil, nil, err
	}

	return pool, redisClient, nil
}

func main() {
	f, perf := os.Create("cpu.pprof")
	if perf != nil {
		log.Fatal(perf)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := pprof.StartCPUProfile(f)
	if err != nil {
		return
	}
	defer pprof.StopCPUProfile()

	cfg, err := config.InitConfig()
	if err != nil {
		zap.L().Error("failed to initialize config", zap.Error(err))
		return
	}

	if err = logger.Init(
		zap.DebugLevel,
		zap.String("service", "example"),
		zap.String("version", "v42.0.69"),
		zap.Strings("maintainers", []string{"@fc", "@FACorreiaa"}),
	); err != nil || logger.Log == nil {
		panic("failed to initialize logging")
	}

	zapLogger := logger.Log

	pool, redisClient, err := run()
	if err != nil {
		zapLogger.Error("failed to run the application", zap.Error(err))
		return
	}
	defer pool.Close()
	defer redisClient.Close()

	tu := new(utils.TransportUtils)

	brokers := internal.ConfigureUpstreamClients(zapLogger, tu)
	if brokers == nil {
		zapLogger.Error("failed to configure brokers")
		return
	}

	metrics.InitPprof()

	container := internal.NewServiceContainer(ctx, pool, redisClient, brokers)

	//var wg sync.WaitGroup
	//wg.Add(2)

	go func() {
		//defer wg.Done()
		if err = internal.ServeGRPC(ctx, cfg.Server.GrpcPort, container); err != nil {
			zapLogger.Error("failed to serve grpc", zap.Error(err))
			return
		}
	}()

	//go func() {
	//	defer wg.Done()
	//	if err = internal.ServeHTTP(cfg.Server.HTTPPort); err != nil {
	//		zapLogger.Error("failed to serve http", zap.Error(err))
	//		return
	//	}
	//}()
	//
	//wg.Wait()
	if err = internal.ServeHTTP(cfg.Server.HTTPPort); err != nil {
		zapLogger.Error("failed to serve http", zap.Error(err))
		return
	}
}
