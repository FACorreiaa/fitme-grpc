package internal

import (
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
	ActivityService   *activity.ActivityService
}

func NewServiceContainer(pgPool *pgxpool.Pool, redisClient *redis.Client, brokers *container.Brokers) *ServiceContainer {
	sessionManager := auth.NewSessionManager(pgPool, redisClient)
	authRepo := auth.NewAuthRepository(pgPool, redisClient, sessionManager)
	calculatorRepo := calculator.NewCalculatorRepository(pgPool, redisClient, sessionManager)
	activityRepo := activity.NewActivityRepository(pgPool, redisClient, sessionManager)
	authService := auth.NewAuthService(authRepo, pgPool, redisClient, sessionManager)
	customerService := domain.NewCustomerService(pgPool, redisClient)
	calculatorService := calculator.NewCalculatorService(calculatorRepo)
	activityService := activity.NewCalculatorService(activityRepo)
	return &ServiceContainer{
		Brokers:           brokers,
		AuthService:       authService,
		CustomerService:   customerService,
		CalculatorService: calculatorService,
		ActivityService:   activityService,
	}
}
