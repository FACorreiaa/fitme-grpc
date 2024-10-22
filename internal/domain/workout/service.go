package workout

import (
	"context"
	"fmt"

	"github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

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

func (s ServiceWorkout) GetExercises(ctx context.Context, req *generated.GetExercisesReq) (*generated.GetExercisesRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "Workout/GetExercises")
	defer span.End()

	exercisesResponse, err := s.repo.GetExercises(ctx, req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
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
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (s ServiceWorkout) GetExerciseID(ctx context.Context, req *generated.GetExerciseIDReq) (*generated.GetExerciseIDRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) CreateExercise(ctx context.Context, req *generated.CreateExerciseReq) (*generated.CreateExerciseRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) UpdateExercice(ctx context.Context, req *generated.UpdateExerciseReq) (*generated.UpdateExerciseRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) GetWorkoutPlanExercises(ctx context.Context, req *generated.GetWorkoutPlanExercisesReq) (*generated.GetWorkoutPlanExercisesRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceWorkout) GetExerciseByIdWorkoutPlan(ctx context.Context, req *generated.GetExerciseByIdWorkoutPlanReq) (*generated.GetExerciseByIdWorkoutPlanRes, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}
