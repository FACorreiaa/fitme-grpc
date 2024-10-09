package activity

import (
	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type ActivityRepository struct {
	pbc.UnimplementedCalculatorServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewActivityRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *ActivityRepository {
	return &ActivityRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}
