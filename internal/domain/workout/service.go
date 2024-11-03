package workout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcrequest"
)

var logger *zap.Logger

type ServiceWorkout struct {
	pbw.UnimplementedWorkoutServer
	ctx  context.Context
	repo domain.RepositoryWorkout
}

func NewServiceWorkout(ctx context.Context, repo domain.RepositoryWorkout) *ServiceWorkout {
	return &ServiceWorkout{
		ctx:  ctx,
		repo: repo,
	}
}

func (s ServiceWorkout) GetExercises(ctx context.Context, req *pbw.GetExercisesReq) (*pbw.GetExercisesRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/GetExercises")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	exercisesResponse, err := s.repo.GetExercises(ctx, req)

	if err != nil {
		return &pbw.GetExercisesRes{
			Success: false,
			Message: "No exercises found",
			Response: &pbw.BaseResponse{
				Upstream:  "workout-service",
				RequestId: requestID,
			},
		}, fmt.Errorf("exercises not found %w", err)
	}

	response := &pbw.GetExercisesRes{}

	for _, a := range exercisesResponse.Exercise {
		//createdAtFormatted := a.CreatedAt.AsTime().Format("2006-01-02 15:04:05.999999")
		//updatedAtFormatted := a.UpdatedAt.AsTime().Format("2006-01-02 15:04:05.999999")

		fmt.Printf("Created at: %#v", a.CreatedAt.AsTime())
		response.Exercise = append(response.Exercise, &pbw.XExercises{
			ExerciseId:    a.ExerciseId,
			Name:          a.Name,
			ExerciseType:  a.ExerciseType,
			MuscleGroup:   a.MuscleGroup,
			Equipment:     a.Equipment,
			Difficulty:    a.Difficulty,
			Instruction:   a.Instruction,
			Video:         a.Video,
			CustomCreated: a.CustomCreated,
			CreatedAt:     a.CreatedAt,
			UpdatedAt:     a.UpdatedAt,
		})
	}

	if req.GetRequest() != nil {
		span.SetAttributes(
			attribute.String("request.id", req.GetRequest().RequestId),
			attribute.String("request.details", req.GetRequest().String()),
		)
	} else {
		span.SetAttributes(
			attribute.String("request.id", "unknown"),
			attribute.String("request.details", "no details available"),
		)
	}

	return &pbw.GetExercisesRes{
		Success:  true,
		Message:  "Activities retrieved successfully",
		Exercise: response.Exercise,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: requestID,
		},
	}, nil
}

func (s ServiceWorkout) GetExerciseID(ctx context.Context, req *pbw.GetExerciseIDReq) (*pbw.GetExerciseIDRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/GetExerciseID")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	if req.ExerciseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}

	exercise, err := s.repo.GetExerciseID(ctx, req)

	createdAt := timestamppb.New(time.Now())

	exercise.Exercise.CreatedAt = createdAt
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return &pbw.GetExerciseIDRes{
			Success:  false,
			Message:  "Failed to retrieve activity by name",
			Exercise: exercise.Exercise,
			Response: &pbw.BaseResponse{
				Upstream:  "activity-service",
				RequestId: requestID,
			},
		}, nil
	}
	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.GetExerciseIDRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Exercise: exercise.Exercise,
		Response: &pbw.BaseResponse{
			Upstream:  "activity-service",
			RequestId: requestID,
		},
	}, nil
}

func (s ServiceWorkout) CreateExercise(ctx context.Context, req *pbw.CreateExerciseReq) (*pbw.CreateExerciseRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/GetExercises")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	request := &pbw.CreateExerciseReq{
		Exercise: req.Exercise,
		UserId:   userID,
		Request: &pbw.BaseRequest{
			Downstream: "workout-service",
			RequestId:  requestID,
		},
	}

	response, err := s.repo.CreateExercise(ctx, request)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pbw.CreateExerciseRes{
				Success:  false,
				Message:  "Exercises creation failed",
				Exercise: response.Exercise,
				Response: &pbw.BaseResponse{
					Upstream:  "workout-service",
					RequestId: requestID,
				},
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to create exercise: %v", err)
	}
	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.CreateExerciseRes{
		Success:  true,
		Message:  "Exercises created successfully",
		Exercise: response.Exercise,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: requestID,
		},
	}, nil
}

func (s ServiceWorkout) DeleteExercise(ctx context.Context, req *pbw.DeleteExerciseReq) (*pbw.NilRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/DeleteExercise")
	defer span.End()
	if req.ExerciseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}

	_, err := s.repo.DeleteExercise(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error deleting exercise session: %v", err)
	}

	// change later
	span.SetAttributes(
		attribute.String("request.id", req.ExerciseId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.NilRes{}, nil
}

func (s ServiceWorkout) UpdateExercise(ctx context.Context, req *pbw.UpdateExerciseReq) (*pbw.UpdateExerciseRes, error) {
	if req.ExerciseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}
	res, err := s.repo.UpdateExercise(ctx, req)
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExercise")
	defer span.End()

	if err != nil {
		logger.Error("failed to update exercise", zap.Error(err))
		return &generated.UpdateExerciseRes{
			Success: false,
			Message: "failed to update exercise: " + err.Error(),
		}, nil
	}

	// change later
	span.SetAttributes(
		attribute.String("request.id", req.ExerciseId),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}

func (s ServiceWorkout) GetWorkoutPlanExercises(ctx context.Context, req *generated.GetWorkoutPlanExercisesReq) (*generated.GetWorkoutPlanExercisesRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExercise")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	res, err := s.repo.GetWorkoutPlanExercises(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pbw.GetWorkoutPlanExercisesRes{
				Success: false,
				Message: "No exercises found",
				Response: &pbw.BaseResponse{
					Upstream:  "workout-service",
					RequestId: requestID,
				},
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get exercises: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)
	return res, nil
}

func (s ServiceWorkout) GetExerciseByIdWorkoutPlan(ctx context.Context, req *generated.GetExerciseByIdWorkoutPlanReq) (*generated.GetExerciseByIdWorkoutPlanRes, error) {
	exerciseID := req.ExerciseWorkoutPlan
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExercise")
	defer span.End()
	if exerciseID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}

	workout, err := s.repo.GetWorkoutPlanExercisesByID(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no exercise with id %s found", exerciseID)
		}
		return nil, status.Errorf(codes.Internal, "failed to get exercises: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return workout, nil
}

func (s ServiceWorkout) DeleteExerciseByIdWorkoutPlan(ctx context.Context, req *generated.DeleteExerciseByIdWorkoutPlanReq) (*generated.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) UpdateExerciseByIdWorkoutPlan(ctx context.Context, req *generated.UpdateExerciseByIdWorkoutPlanReq) (*generated.UpdateExerciseByIdWorkoutPlanRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) InsertExerciseWorkoutPlan(ctx context.Context, req *generated.InsertExerciseWorkoutPlanReq) (*generated.NilRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExercise")
	defer span.End()

	_, err := s.repo.InsertExerciseWorkoutPlan(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert exercise_workout_plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.NilRes{}, nil
}

func (s ServiceWorkout) GetWorkoutPlans(ctx context.Context, req *generated.GetWorkoutPlansReq) (*generated.GetWorkoutPlansRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) GetWorkoutPlan(ctx context.Context, req *generated.GetWorkoutPlanReq) (*generated.GetWorkoutPlanRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) DeleteWorkoutPlan(ctx context.Context, req *generated.DeleteWorkoutPlanReq) (*generated.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) UpdateWorkoutPlan(ctx context.Context, req *generated.UpdateWorkoutPlanReq) (*generated.UpdateWorkoutPlanRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) InsertWorkoutPlan(ctx context.Context, req *generated.InsertWorkoutPlanReq) (*generated.InsertWorkoutPlanRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/GetExercises")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	//workoutPlan := req.Workout
	//workoutPlan.UserId = userID

	// create a number of workday days for a plan
	workoutDays := make([]*pbw.XWorkoutPlanDay, len(req.PlanDay))

	for i, planDay := range req.PlanDay {
		exercises := make([]*pbw.XExercises, len(planDay.ExerciseId))
		// var wg sync.WaitGroup
		// var mu sync.Mutex
		// var fetchError error

		for j, exercise := range planDay.ExerciseId {
			exerciseReq := &pbw.GetExerciseIDReq{
				ExerciseId: exercise,
			}
			exerciseDetails, err := s.repo.GetExerciseID(ctx, exerciseReq)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get exercise details: %v", err)
			}
			exercises[j] = exerciseDetails.Exercise
		}
		workoutDays[i] = &pbw.XWorkoutPlanDay{
			Day:       planDay.Day,
			Exercises: exercises,
		}
	}

	workoutPlan := &pbw.XWorkoutPlan{
		WorkoutId:      uuid.NewString(),
		UserId:         userID,
		Description:    req.Workout.Description,
		Notes:          req.Workout.Notes,
		Rating:         req.Workout.Rating,
		WorkoutPlanDay: workoutDays,
		CreatedAt:      timestamppb.New(time.Now()),
	}

	req.Workout = workoutPlan

	response, err := s.repo.CreateWorkoutPlan(ctx, req)
	if err != nil {
		return &pbw.InsertWorkoutPlanRes{
			Success: false,
			Message: "Workout creation failed",
			Response: &pbw.BaseResponse{
				Upstream:  "workout-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to insert workout: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)
	return &pbw.InsertWorkoutPlanRes{
		Success: false,
		Message: "Workout creation failed",
		Workout: response.Workout,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: requestID,
		},
	}, nil
}
