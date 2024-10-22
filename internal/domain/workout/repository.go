package workout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type RepositoryWorkout struct {
	pba.UnimplementedActivityServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewRepositoryWorkout(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *RepositoryWorkout {
	return &RepositoryWorkout{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func nullTimeToTimestamppb(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

func (r *RepositoryWorkout) GetExercises(ctx context.Context, req *pbw.GetExercisesReq) (*pbw.GetExercisesRes, error) {
	exercises := make([]*pbw.XExercises, 0)
	query := `SELECT DISTINCT
    			id, name, type, muscle, equipment, difficulty,
				instructions, video, custom_created, created_at, updated_at
				FROM exercise_list`

	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pbw.GetExercisesRes{
				Success: false,
				Message: "No exercises found",
				Response: &pbw.BaseResponse{
					Upstream:  "workout-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, fmt.Errorf("exercises not found %w", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to query exercises: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		exercise := pbw.XExercises{}
		e := &Exercises{}

		err := rows.Scan(
			&e.ID, &e.Name, &e.ExerciseType, &e.MuscleGroup, &e.Equipment, &e.Difficulty,
			&e.Instructions, &e.Video, &e.CustomCreated, &e.CreatedAt, &e.UpdatedAt,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &pbw.GetExercisesRes{
					Success: false,
					Message: "No exercises found",
					Response: &pbw.BaseResponse{
						Upstream:  "workout-service",
						RequestId: domain.GenerateRequestID(ctx),
					},
				}, fmt.Errorf("exercises not found: %w", err)
			}
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		createdAt := timestamppb.New(e.CreatedAt)
		var updatedAt sql.NullTime
		if e.UpdatedAt.Valid {
			updatedAt = e.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		exercise.ExerciseId = e.ID
		exercise.Name = e.Name
		exercise.ExerciseType = e.ExerciseType
		exercise.MuscleGroup = e.MuscleGroup
		exercise.Equipment = e.Equipment
		exercise.Difficulty = e.Difficulty
		exercise.Instruction = e.Instructions
		exercise.Video = e.Video
		exercise.CustomCreated = e.CustomCreated
		exercise.CreatedAt = createdAt
		exercise.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		exerciseProto := &pbw.XExercises{
			ExerciseId:    exercise.ExerciseId,
			Name:          exercise.Name,
			ExerciseType:  exercise.ExerciseType,
			MuscleGroup:   exercise.MuscleGroup,
			Equipment:     exercise.Equipment,
			Difficulty:    exercise.Difficulty,
			Instruction:   exercise.Instruction,
			Video:         exercise.Video,
			CustomCreated: exercise.CustomCreated,
			CreatedAt:     exercise.CreatedAt,
			UpdatedAt:     exercise.UpdatedAt,
		}

		exercises = append(exercises, exerciseProto)
	}

	if len(exercises) == 0 {
		return &pbw.GetExercisesRes{
			Success: false,
			Message: "No exercises found",
			Response: &pbw.BaseResponse{
				Upstream:  "workout-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}

	return &pbw.GetExercisesRes{
		Success:  true,
		Message:  "Exercises retrieved successfully",
		Exercise: exercises,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil

}
