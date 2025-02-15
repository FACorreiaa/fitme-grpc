package workout

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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

// GetExercises Exercises
func (s ServiceWorkout) GetExercises(ctx context.Context, req *pbw.GetExercisesReq) (*pbw.GetExercisesRes, error) {
	traceContext, span := otel.Tracer("fitme-dev").Start(ctx, "GetExercises")
	defer span.End()
	traceID := span.SpanContext().TraceID().String()
	println(traceID)
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	exercisesResponse, err := s.repo.GetExercises(traceContext, req)

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

	span.SetAttributes(
		semconv.ServiceNameKey.String("Workout"),
		attribute.String("grpc.method", "GetExercises"),
		attribute.String("request.id", req.GetRequest().RequestId),
	)

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
	tracer := otel.Tracer("FitSphere")
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
	tracer := otel.Tracer("FitSphere")
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
	tracer := otel.Tracer("FitSphere")
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
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExercise")
	defer span.End()

	if err != nil {
		logger.Error("failed to update exercise", zap.Error(err))
		return &pbw.UpdateExerciseRes{
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

// GetWorkoutPlanExercises Workout plan
func (s ServiceWorkout) GetWorkoutPlanExercises(ctx context.Context, req *pbw.GetWorkoutPlanExercisesReq) (*pbw.GetWorkoutPlanExercisesRes, error) {
	tracer := otel.Tracer("FitSphere")
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

func (s ServiceWorkout) GetExerciseByIdWorkoutPlan(ctx context.Context, req *pbw.GetExerciseByIdWorkoutPlanReq) (*pbw.GetExerciseByIdWorkoutPlanRes, error) {
	exerciseID := req.ExerciseWorkoutPlan
	tracer := otel.Tracer("FitSphere")
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

func (s ServiceWorkout) DeleteExerciseByIdWorkoutPlan(ctx context.Context, req *pbw.DeleteExerciseByIdWorkoutPlanReq) (*pbw.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
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

	//userID := ctx.Value("userID").(string)
	//if userID == "" {
	//	return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	//}
	//req.UserId = userID

	_, err := s.repo.DeleteExerciseWorkoutPlan(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout plan")
	}

	span.SetAttributes(
		attribute.String("request.id", requestID),
		attribute.String("request.details", req.String()),
	)

	return &pbw.NilRes{}, nil
}

func (s ServiceWorkout) UpdateExerciseByIdWorkoutPlan(ctx context.Context, req *pbw.UpdateExerciseByIdWorkoutPlanReq) (*pbw.UpdateExerciseByIdWorkoutPlanRes, error) {
	if req.WorkoutPlanId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}
	res, err := s.repo.UpdateExerciseWorkoutPLan(ctx, req)
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Workout/UpdateExerciseByIdWorkoutPlan")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	if err != nil {
		logger.Error("failed to update exercise", zap.Error(err))
		return &pbw.UpdateExerciseByIdWorkoutPlanRes{
			Success: false,
			Message: "failed to update workout: " + err.Error(),
		}, nil
	}

	// change later
	span.SetAttributes(
		attribute.String("request.id", requestID),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}

func (s ServiceWorkout) InsertExerciseWorkoutPlan(ctx context.Context, req *pbw.InsertExerciseWorkoutPlanReq) (*pbw.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
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

//func (s ServiceWorkout) DeleteExerciseWorkoutPlan(ctx context.Context, req *pbw.DeleteExerciseByIdWorkoutPlanReq) *pbw.NilRes {
//	return nil
//}

// InsertWorkoutPlan Insert workout plan
func (s ServiceWorkout) InsertWorkoutPlan(
	ctx context.Context,
	req *pbw.InsertWorkoutPlanReq,
) (*pbw.InsertWorkoutPlanRes, error) {

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	// We'll build a new array of “resolved” days, each with a slice of *pbw.XExercises.
	resolvedDays := make([]*pbw.XWorkoutPlanDay, len(req.PlanDay))

	for i, planDay := range req.PlanDay {
		dayExercises := make([]*pbw.XExercises, len(planDay.Exercises))
		dayName := planDay.Day

		// Loop over each XExerciseInput
		for j, exInput := range planDay.Exercises {
			// If the request has an exercise_id, try to fetch from DB
			if exInput.ExerciseId != "" {
				getReq := &pbw.GetExerciseIDReq{
					ExerciseId: exInput.ExerciseId,
					// Optionally set the .request field if you want:
					// Request: &pbw.BaseRequest{ request_id: "..."},
				}

				existing, err := s.repo.GetExerciseID(ctx, getReq)
				if err != nil {
					// If the repo returns NOT_FOUND, create a new exercise.
					_, ok := status.FromError(err)
					if !ok {
						// Construct a CreateExerciseReq
						createReq := &pbw.CreateExerciseReq{
							Exercise: &pbw.XExercises{
								ExerciseId:   exInput.ExerciseId,
								Name:         exInput.Name,
								ExerciseType: exInput.ExerciseType,
								MuscleGroup:  exInput.MuscleGroup,
								Equipment:    exInput.Equipment,
								Difficulty:   exInput.Difficulty,
								Instruction:  exInput.Instruction,
								Video:        exInput.Video,
								CreatedAt:    exInput.CreatedAt, // or timestamppb.Now()
							},
							UserId: userID,
							// Optionally fill .request if needed
						}
						newExerciseRes, err2 := s.repo.CreateExercise(ctx, createReq)
						if err2 != nil {
							return nil, status.Errorf(codes.Internal, "failed to create new exercise: %v", err2)
						}
						// Now we have a *pbw.CreateExerciseRes, which has .Exercise
						dayExercises[j] = newExerciseRes.Exercise
					} else {
						// Some other database error
						return nil, status.Errorf(codes.Internal, "failed to get exercise details: %v", err)
					}
				} else {
					// existing is *pbw.GetExerciseIDRes; use .Exercise
					dayExercises[j] = existing.Exercise
				}
			} else {
				// No exercise_id => definitely create a new exercise
				newID := uuid.NewString()
				createReq := &pbw.CreateExerciseReq{
					Exercise: &pbw.XExercises{
						ExerciseId:   newID,
						Name:         exInput.Name,
						ExerciseType: exInput.ExerciseType,
						MuscleGroup:  exInput.MuscleGroup,
						Equipment:    exInput.Equipment,
						Difficulty:   exInput.Difficulty,
						Instruction:  exInput.Instruction,
						Video:        exInput.Video,
						CreatedAt:    exInput.CreatedAt,
					},
					UserId: userID,
				}
				newExerciseRes, err2 := s.repo.CreateExercise(ctx, createReq)
				if err2 != nil {
					return nil, status.Errorf(codes.Internal, "failed to auto-create exercise: %v", err2)
				}
				dayExercises[j] = newExerciseRes.Exercise
			}
		} // end for exInput

		// Build final day object with fully resolved exercises
		resolvedDays[i] = &pbw.XWorkoutPlanDay{
			Day:       dayName,
			Exercises: dayExercises,
		}
	}

	// Overwrite the “workout_plan_day” in the request so the repo can insert them
	req.Workout.WorkoutPlanDay = resolvedDays
	req.Workout.UserId = userID
	// If the client didn’t provide a workout_id, generate one
	if req.Workout.WorkoutId == "" {
		req.Workout.WorkoutId = uuid.NewString()
	}
	req.Request.RequestId = requestID

	// Now call the repo’s CreateWorkoutPlan with the entire InsertWorkoutPlanReq
	// that includes the final resolved exercises. The repo presumably returns
	// (*pbw.InsertWorkoutPlanRes, error).
	response, err := s.repo.CreateWorkoutPlan(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert workout plan: %v", err)
	}

	// The repo is expected to fill 'response.Workout' with the final data
	return response, nil
}

func (s ServiceWorkout) GetWorkoutPlans(ctx context.Context, req *pbw.GetWorkoutPlansReq) (*pbw.GetWorkoutPlansRes, error) {
	tracer := otel.Tracer("FitSphere")
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

	result, err := s.repo.GetWorkoutPlans(ctx, req)
	if err != nil {
		return &pbw.GetWorkoutPlansRes{
			Success: false,
			Message: "Workout fetch failed",
			Response: &pbw.BaseResponse{
				Upstream:  "workout-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get workout plans: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.GetWorkoutPlansRes{
		Success:     true,
		Message:     "Workout fetch succeeded",
		WorkoutPlan: result.WorkoutPlan,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: requestID,
		},
	}, nil
}

func (s ServiceWorkout) GetWorkoutPlan(ctx context.Context, req *pbw.GetWorkoutPlanReq) (*pbw.GetWorkoutPlanRes, error) {
	tracer := otel.Tracer("FitSphere")
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

	workout, err := s.repo.GetWorkoutPlan(ctx, req)
	if err != nil {
		return &pbw.GetWorkoutPlanRes{
			Success: false,
			Message: "Workout plan fetch failed",
			Response: &pbw.BaseResponse{
				Upstream:  "workout-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get workout plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbw.GetWorkoutPlanRes{
		Success:     true,
		Message:     "Workout plan fetch succeed",
		WorkoutPlan: workout.WorkoutPlan,
		Response: &pbw.BaseResponse{
			Upstream:  "workout-service",
			RequestId: requestID,
		},
	}, nil
}

func (s ServiceWorkout) DeleteWorkoutPlan(ctx context.Context, req *pbw.DeleteWorkoutPlanReq) (*pbw.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
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
	req.UserId = userID

	_, err := s.repo.DeleteWorkoutPlan(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete workout plan")
	}

	span.SetAttributes(
		attribute.String("request.id", requestID),
		attribute.String("request.details", req.String()),
	)

	return &pbw.NilRes{}, nil
}

func (s ServiceWorkout) UpdateWorkoutPlan(ctx context.Context, req *pbw.UpdateWorkoutPlanReq) (*pbw.UpdateWorkoutPlanRes, error) {
	if req.WorkoutId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req id is required")
	}
	res, err := s.repo.UpdateWorkoutPlan(ctx, req)
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Workout/UpdateWorkoutPlan")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbw.BaseRequest{}
	}

	req.Request.RequestId = requestID

	if err != nil {
		logger.Error("failed to update exercise", zap.Error(err))
		return &pbw.UpdateWorkoutPlanRes{
			Success: false,
			Message: "failed to update workout: " + err.Error(),
		}, nil
	}

	// change later
	span.SetAttributes(
		attribute.String("request.id", requestID),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}

func (s ServiceWorkout) DownloadWorkoutPlan(req *pbw.DownloadWorkoutPlanRequest, stream pbw.Workout_DownloadWorkoutPlanServer) (err error) {
	if req.WorkoutPlanId == "" {
		return status.Errorf(codes.InvalidArgument, "req id is required")
	}
	ctx := stream.Context()

	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Workout/UpdateWorkoutPlan")
	defer span.End()

	var fileData []byte
	var fileName, contentType string
	//requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	//if !ok {
	//	return status.Error(codes.Internal, "request id not found in context")
	//}

	baseReq := &pbw.BaseRequest{
		Downstream: "todo",
		//RequestId:  requestID,
	}

	workoutPlanReq := &pbw.GetWorkoutPlanReq{
		WorkoutPlanId: req.WorkoutPlanId,
		Request:       baseReq,
	}

	workoutPlan, err := s.repo.GetWorkoutPlan(ctx, workoutPlanReq)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to download workout plan request: %v", err)
	}

	switch req.Format {
	case pbw.FileFormat_CSV:
		fileData, fileName, contentType, err = generateCSV(ctx, workoutPlan)
	case pbw.FileFormat_PDF:
		fileData, fileName, contentType, err = generatePDF(ctx, workoutPlan)
	case pbw.FileFormat_EXCEL:
		fileData, fileName, contentType, err = generateExcel(ctx, workoutPlan)
	default:
		return status.Errorf(codes.InvalidArgument, "unknown format")
	}

	if err != nil {
		return status.Errorf(codes.Internal, "failed to download workout plan request: %v", err)
	}

	const chunkSize = 64 * 1024
	for current := 0; current < len(fileData); current += chunkSize {
		end := current + chunkSize
		if end > len(fileData) {
			end = len(fileData)
		}

		chunk := &pbw.FileChunk{
			Content: fileData[current:end],
		}

		if current == 0 {
			// add some metadata on the first chunk
			chunk.IsFirstChunk = true
			chunk.FileName = fileName
			chunk.ContentType = contentType
		}

		if err = stream.Send(chunk); err != nil {
			return status.Errorf(codes.Internal, "failed to download workout plan request: %v", err)
		}
	}

	span.SetAttributes(
		attribute.String("request.id", workoutPlanReq.Request.RequestId),
		attribute.String("request.details", workoutPlanReq.String()))

	return nil

}

func generateCSV(ctx context.Context, workoutPlan *pbw.GetWorkoutPlanRes) ([]byte, string, string, error) {
	if workoutPlan == nil || workoutPlan.WorkoutPlan == nil {
		return nil, "", "", errors.New("no workout plan provided")
	}

	plan := workoutPlan.WorkoutPlan

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{
		"workout_plan_id",
		"user_id",
		"description",
		"notes",
		"rating",
		"day",
		"exercises",
	}

	if err := writer.Write(header); err != nil {
		return nil, "", "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// If the proto eventually uses repeated WorkoutDayResponse, handle each day.
	if len(plan.WorkoutDay) > 0 {
		for _, wd := range plan.WorkoutDay {
			record := []string{
				plan.WorkoutPlanId,
				plan.UserId,
				plan.Description,
				plan.Notes,
				strconv.Itoa(int(plan.Rating)),
				wd.Day,
				strings.Join(wd.Exercises, ", "),
			}
			if err := writer.Write(record); err != nil {
				return nil, "", "", fmt.Errorf("failed to write CSV record: %w", err)
			}
		}
	} else {
		// Fallback: use plan.Day and plan.Exercises
		record := []string{
			plan.WorkoutPlanId,
			plan.UserId,
			plan.Description,
			plan.Notes,
			strconv.Itoa(int(plan.Rating)),
			plan.Day,
			strings.Join(plan.Exercises, ", "),
		}
		if err := writer.Write(record); err != nil {
			return nil, "", "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", "", fmt.Errorf("error flushing CSV writer: %w", err)
	}

	fileName := "workout_plan.csv"
	contentType := "text/csv"

	return buf.Bytes(), fileName, contentType, nil
}

func generateExcel(ctx context.Context, workoutPlan *pbw.GetWorkoutPlanRes) ([]byte, string, string, error) {
	return nil, "", "", nil
}

func generatePDF(ctx context.Context, workoutPlan *pbw.GetWorkoutPlanRes) ([]byte, string, string, error) {
	return nil, "", "", nil
}
