package internal

import (
	"context"

	"github.com/FACorreiaa/fitme-protos/container"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/activity"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/calculator"
)

type ServiceContainer struct {
	Brokers           *container.Brokers
	AuthService       *auth.AuthService
	CustomerService   *domain.CustomerService
	CalculatorService *calculator.CalculatorService
	ServiceActivity   *activity.ServiceActivity
}

func NewServiceContainer(ctx context.Context, pgPool *pgxpool.Pool, redisClient *redis.Client, brokers *container.Brokers) *ServiceContainer {
	sessionManager := auth.NewSessionManager(pgPool, redisClient)
	authRepo := auth.NewAuthRepository(pgPool, redisClient, sessionManager)
	calculatorRepo := calculator.NewCalculatorRepository(pgPool, redisClient, sessionManager)
	activityRepo := activity.NewRepositoryActivity(pgPool, redisClient, sessionManager)
	authService := auth.NewAuthService(ctx, authRepo, pgPool, redisClient, sessionManager)
	customerService := domain.NewCustomerService(ctx, pgPool, redisClient)
	calculatorService := calculator.NewCalculatorService(ctx, calculatorRepo)
	ServiceActivity := activity.NewCalculatorService(ctx, activityRepo)
	return &ServiceContainer{
		Brokers:           brokers,
		AuthService:       authService,
		CustomerService:   customerService,
		CalculatorService: calculatorService,
		ServiceActivity:   ServiceActivity,
	}
}
