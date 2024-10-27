package workout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	pbw.UnimplementedWorkoutServer
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
	exercisesProtoList := make([]*pbw.XExercises, 0)
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
		exerciseProto := pbw.XExercises{}
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

		exerciseProto.ExerciseId = e.ID
		exerciseProto.Name = e.Name
		exerciseProto.ExerciseType = e.ExerciseType
		exerciseProto.MuscleGroup = e.MuscleGroup
		exerciseProto.Equipment = e.Equipment
		exerciseProto.Difficulty = e.Difficulty
		exerciseProto.Instruction = e.Instructions
		exerciseProto.Video = e.Video
		exerciseProto.CustomCreated = e.CustomCreated
		exerciseProto.CreatedAt = createdAt
		exerciseProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		newProtoExerciseList := &pbw.XExercises{
			ExerciseId:    exerciseProto.ExerciseId,
			Name:          exerciseProto.Name,
			ExerciseType:  exerciseProto.ExerciseType,
			MuscleGroup:   exerciseProto.MuscleGroup,
			Equipment:     exerciseProto.Equipment,
			Difficulty:    exerciseProto.Difficulty,
			Instruction:   exerciseProto.Instruction,
			Video:         exerciseProto.Video,
			CustomCreated: exerciseProto.CustomCreated,
			CreatedAt:     exerciseProto.CreatedAt,
			UpdatedAt:     exerciseProto.UpdatedAt,
		}

		exercisesProtoList = append(exercisesProtoList, newProtoExerciseList)
	}

	if len(exercisesProtoList) == 0 {
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
		Exercise: exercisesProtoList,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil

}

func (r *RepositoryWorkout) GetExerciseID(ctx context.Context, req *pbw.GetExerciseIDReq) (*pbw.GetExerciseIDRes, error) {
	exerciseProto := &pbw.XExercises{}
	exercise := &Exercises{}
	id := req.ExerciseId

	if id == "" {
		return &pbw.GetExerciseIDRes{}, status.Error(codes.InvalidArgument, "workout ID is required")
	}

	query := `SELECT 	id, name, type, muscle, equipment, difficulty,
						instructions, video, custom_created, created_at, updated_at
			   FROM exercise_list
			   WHERE id = $1`

	err := r.pgpool.QueryRow(ctx, query, id).Scan(
		&exercise.ID, &exercise.Name, &exercise.ExerciseType, &exercise.MuscleGroup, &exercise.Equipment,
		&exercise.Difficulty, &exercise.Instructions, &exercise.Video, &exercise.CustomCreated, &exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pbw.GetExerciseIDRes{
				Success: false,
				Message: "No exercises found",
				Response: &pbw.BaseResponse{
					Upstream:  "workout-service",
					RequestId: domain.GenerateRequestID(ctx),
				},
			}, status.Error(codes.NotFound, "workout ID not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to query exercise: %v", err)
	}

	createdAt := timestamppb.New(exercise.CreatedAt)
	var updatedAt sql.NullTime

	if exercise.UpdatedAt.Valid {
		updatedAt = exercise.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	exerciseProto.ExerciseId = exercise.ID
	exerciseProto.Name = exercise.Name
	exerciseProto.MuscleGroup = exercise.MuscleGroup
	exerciseProto.Equipment = exercise.Equipment
	exerciseProto.Difficulty = exercise.Difficulty
	exerciseProto.Instruction = exercise.Instructions
	exerciseProto.Video = exercise.Video
	exerciseProto.CustomCreated = exercise.CustomCreated
	exerciseProto.CreatedAt = createdAt
	exerciseProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

	return &pbw.GetExerciseIDRes{
		Success:  true,
		Message:  "Exercise retrieved successfully",
		Exercise: exerciseProto,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (r *RepositoryWorkout) CreateExercise(ctx context.Context, req *pbw.CreateExerciseReq) (*pbw.CreateExerciseRes, error) {
	fmt.Printf("User id: %s", req.UserId)
	tx, err := r.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	createdExerciseListQuery := `
				INSERT INTO exercise_list (name, type, muscle, equipment, difficulty,
                                   instructions, video,
                                   created_at, updated_at)
        		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                RETURNING id`

	currentTime := time.Now()

	var exerciseID string
	err = tx.QueryRow(ctx, createdExerciseListQuery,
		req.Exercise.Name,
		req.Exercise.ExerciseType,
		req.Exercise.MuscleGroup,
		req.Exercise.Equipment,
		req.Exercise.Difficulty,
		req.Exercise.Instruction,
		req.Exercise.Video,
		currentTime,
		currentTime,
	).Scan(&exerciseID)

	fmt.Println(exerciseID)

	setExerciseToUserQuery := `
				INSERT INTO user_exercises (user_id, exercise_id)
				VALUES ($1, $2)
				RETURNING user_id, exercise_id`

	req.Exercise.ExerciseId = exerciseID

	var userID, associatedExerciseID string

	err = tx.QueryRow(ctx, setExerciseToUserQuery, req.UserId, req.Exercise.ExerciseId).Scan(&userID, &associatedExerciseID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to associate exercise with user: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	exerciseProto := &pbw.XExercises{
		ExerciseId:    exerciseID,
		Name:          req.Exercise.Name,
		ExerciseType:  req.Exercise.ExerciseType,
		MuscleGroup:   req.Exercise.MuscleGroup,
		Equipment:     req.Exercise.Equipment,
		Difficulty:    req.Exercise.Difficulty,
		Instruction:   req.Exercise.Instruction,
		Video:         req.Exercise.Video,
		CustomCreated: true,
		CreatedAt:     timestamppb.New(currentTime), // Proto timestamp for created_at fix later
		UpdatedAt:     timestamppb.New(currentTime), // Proto timestamp for updated_at fix later
	}

	return &pbw.CreateExerciseRes{
		Success:  true,
		Message:  "Exercise created and associated with user successfully",
		Exercise: exerciseProto,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (r *RepositoryWorkout) DeleteExercise(ctx context.Context, req *pbw.DeleteExerciseReq) (*pbw.NilRes, error) {
	query := `DELETE FROM exercise_list WHERE id = $1
 			  AND exercise_list.custom_created = true`
	_, err := r.pgpool.Exec(ctx, query, req.ExerciseId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise: %w", err)
	}
	return &pbw.NilRes{}, nil
}
