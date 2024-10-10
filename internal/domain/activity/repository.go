package activity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type ActivityRepository struct {
	pba.UnimplementedActivityServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewActivityRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *ActivityRepository {
	return &ActivityRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (a *ActivityRepository) GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error) {
	activities := make([]*pba.XActivity, 0)
	query := `SELECT id, user_id, name,
					duration_minutes, total_calories, calories_per_hour,
					created_at, updated_at
			FROM activity`

	rows, err := a.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pba.GetActivityRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found %w", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to query activities: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a pba.XActivity
		var userId, name sql.NullString
		var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

		var createdAt time.Time
		var updatedAt sql.NullTime

		err := rows.Scan(
			&a.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
			&createdAt, &updatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &pba.GetActivityRes{
					Success: false,
					Message: "No activities found",
					Response: &pba.BaseResponse{
						Upstream:  "activity-service",
						RequestId: domain.GenerateRequestID(ctx),
					},
				}, fmt.Errorf("activities not found: %w", err)
			}
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		// Populate activity fields
		a.UserId = userId.String
		a.Name = name.String

		if durationInMinutes.Valid {
			a.DurationInMinutes = float32(durationInMinutes.Float64)
		}
		if totalCalories.Valid {
			a.TotalCalories = float32(totalCalories.Float64)
		}
		if caloriesPerHour.Valid {
			a.CaloriesPerHour = float32(caloriesPerHour.Float64)
		}

		// Use the time values directly to create protobuf timestamps
		a.CreatedAt = timestamppb.New(createdAt)

		if updatedAt.Valid {
			a.UpdatedAt = timestamppb.New(updatedAt.Time)
		} else {
			a.UpdatedAt = timestamppb.New(time.Now()) // default value if updated_at is NULL
		}

		activities = append(activities, &a)
	}

	if len(activities) == 0 {
		return &pba.GetActivityRes{
			Success: false,
			Message: "No activities found",
			Response: &pba.BaseResponse{
				Upstream:  "activity-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}

	return &pba.GetActivityRes{
		Success:  true,
		Message:  "Activities retrieved successfully",
		Activity: activities,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityRepository) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
	activity := &pba.XActivity{}
	activityID := req.PublicId
	var createdAt time.Time
	var updatedAt sql.NullTime

	var userId, name sql.NullString
	var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

	if activityID == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE id = $1`

	err := a.pgpool.QueryRow(ctx, query, activityID).Scan(
		&activity.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pba.GetActivityIDRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found: %w", err)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	activity.UserId = userId.String
	activity.Name = name.String
	activity.DurationInMinutes = float32(durationInMinutes.Float64)
	activity.TotalCalories = float32(totalCalories.Float64)
	activity.CaloriesPerHour = float32(caloriesPerHour.Float64)

	activity.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		activity.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &pba.GetActivityIDRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityRepository) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
	activity := &pba.XActivity{}
	nameReq := req.PublicId
	var createdAt time.Time
	var updatedAt sql.NullTime

	var userId, name sql.NullString
	var durationInMinutes, totalCalories, caloriesPerHour sql.NullFloat64

	log.Printf("Searching for activity with name: '%s'", nameReq)

	if nameReq == "" {
		return nil, status.Error(codes.InvalidArgument, "activity ID is required")
	}

	query := `SELECT 	id, user_id, name, duration_minutes,
       					total_calories, calories_per_hour, created_at,
       					updated_at
			   FROM activity
			   WHERE name LIKE '%' || $1 || '%'`

	err := a.pgpool.QueryRow(ctx, query, nameReq).Scan(
		&activity.ActivityId, &userId, &name, &durationInMinutes, &totalCalories, &caloriesPerHour,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pba.GetActivityNameRes{
				Success: false,
				Message: "No activities found",
				Response: &pba.BaseResponse{
					Upstream:  "activity-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("activities not found: %w", err)
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if userId.Valid {
		activity.UserId = userId.String
	}
	if name.Valid {
		activity.Name = name.String
	}
	if durationInMinutes.Valid {
		activity.DurationInMinutes = float32(durationInMinutes.Float64)
	}
	if totalCalories.Valid {
		activity.TotalCalories = float32(totalCalories.Float64)
	}
	if caloriesPerHour.Valid {
		activity.CaloriesPerHour = float32(caloriesPerHour.Float64)
	}

	activity.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		activity.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &pba.GetActivityNameRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}
