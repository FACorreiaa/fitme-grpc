package internal

import (
	"github.com/FACorreiaa/fitme-protos/container"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type ServiceContainer struct {
	PgPool          *pgxpool.Pool
	RedisClient     *redis.Client
	Brokers         *container.Brokers
	AuthService     *auth.AuthService
	CustomerService *domain.CustomerService
}

func NewServiceContainer(pgPool *pgxpool.Pool, redisClient *redis.Client, brokers *container.Brokers) *ServiceContainer {
	sessionManager := auth.NewSessionManager(pgPool, redisClient)
	authRepo := auth.NewAuthRepository(pgPool, redisClient, sessionManager)
	authService := auth.NewAuthService(authRepo, pgPool, redisClient, sessionManager)
	customerService := domain.NewCustomerService(pgPool, redisClient)

	return &ServiceContainer{
		PgPool:          pgPool,
		RedisClient:     redisClient,
		Brokers:         brokers,
		AuthService:     authService,
		CustomerService: customerService,
	}
}
