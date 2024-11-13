package measurements

import (
	"context"
	"database/sql"
	"time"

	pbm "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type RepositoryMeasurement struct {
	pbm.UnimplementedUserMeasurementsServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewRepositoryMeasurement(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *RepositoryMeasurement {
	return &RepositoryMeasurement{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (r *RepositoryMeasurement) CreateWeight(ctx context.Context, req *pbm.CreateWeightReq) (*pbm.XWeight, error) {
	query := `
		INSERT INTO weight_measure
		    (user_id, weight_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	currentTime := time.Now()
	var weightID string
	var updatedAt sql.NullTime
	if req.Weight.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: req.Weight.UpdatedAt.AsTime(), Valid: true}
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	err := r.pgpool.QueryRow(ctx, query, req.UserId, req.Weight.WeightValue, currentTime, updatedAt).Scan(&weightID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert weight: %v", err)
	}

	weightProto := &pbm.XWeight{
		WeightId:    weightID,
		UserId:      req.UserId,
		WeightValue: req.Weight.WeightValue,
		CreatedAt:   timestamppb.New(currentTime),
		UpdatedAt:   timestamppb.New(currentTime),
	}

	return weightProto, nil
}
