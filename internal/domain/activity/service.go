package activity

import (
	"context"
	"errors"
	"fmt"
	"time"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

type ActivityService struct {
	pba.UnimplementedActivityServer // Required for forward compatibilit
	repo                            domain.ActivityRepository
}

func NewCalculatorService(repo domain.ActivityRepository) *ActivityService {
	return &ActivityService{
		repo: repo,
	}
}

func (a *ActivityService) GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivity")
	defer span.End()

	if req.PublicId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req PublicId is required")
	}

	activityResponse, err := a.repo.GetActivity(ctx, req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &pba.GetActivityRes{}

	for _, a := range activityResponse.Activity {
		//createdAtFormatted := a.CreatedAt.AsTime().Format("2006-01-02 15:04:05.999999")
		//updatedAtFormatted := a.UpdatedAt.AsTime().Format("2006-01-02 15:04:05.999999")

		fmt.Printf("Created at: %#v", a.CreatedAt.AsTime())
		response.Activity = append(response.Activity, &pba.XActivity{
			ActivityId:        a.ActivityId,
			UserId:            a.UserId,
			Name:              a.Name,
			CaloriesPerHour:   a.CaloriesPerHour,
			DurationInMinutes: a.DurationInMinutes,
			TotalCalories:     a.TotalCalories,
			CreatedAt:         a.CreatedAt,
			UpdatedAt:         a.UpdatedAt,
		})
	}

	span.SetAttributes(
		attribute.String("request.id", req.PublicId),
		attribute.String("request.details", req.String()),
	)

	return &pba.GetActivityRes{
		Success:  true,
		Message:  "Activities retrieved successfully",
		Activity: response.Activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ActivityService) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivitiesByID")
	defer span.End()

	if req.PublicId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req PublicId is required")
	}

	activity, err := a.repo.GetActivitiesByID(ctx, req)

	createdAt := timestamppb.New(time.Now())

	activity.Activity.CreatedAt = createdAt
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return &pba.GetActivityIDRes{
			Success:  false,
			Message:  "Failed to retrieve activity by ID",
			Activity: activity.Activity,
			Response: &pba.BaseResponse{
				Upstream:  "activity-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}
	span.SetAttributes(
		attribute.String("request.id", req.PublicId),
		attribute.String("request.details", req.String()),
	)

	return &pba.GetActivityIDRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity.Activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil

}

func (a *ActivityService) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivitiesByName")
	defer span.End()

	if req.PublicId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req name is required")
	}

	activity, err := a.repo.GetActivitiesByName(ctx, req)

	createdAt := timestamppb.New(time.Now())

	activity.Activity.CreatedAt = createdAt
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return &pba.GetActivityNameRes{
			Success:  false,
			Message:  "Failed to retrieve activity by name",
			Activity: activity.Activity,
			Response: &pba.BaseResponse{
				Upstream:  "activity-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}
	span.SetAttributes(
		attribute.String("request.id", req.PublicId),
		attribute.String("request.details", req.String()),
	)

	return &pba.GetActivityNameRes{
		Success:  true,
		Message:  "Activity retrieved successfully",
		Activity: activity.Activity,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil

}

func (a *ActivityService) GetUserExerciseSession(ctx context.Context, req *pba.GetUserExerciseSessionReq) (*pba.GetUserExerciseSessionRes, error) {
	return nil, nil
}

func (a *ActivityService) GetUserExerciseTotalData(ctx context.Context, req *pba.GetUserExerciseTotalDataReq) (*pba.GetUserExerciseTotalDataRes, error) {
	return nil, nil
}

func (a *ActivityService) GetUserExerciseSessionStats(ctx context.Context, req *pba.GetUserExerciseSessionStatsReq) (*pba.GetUserExerciseSessionStatsRes, error) {
	return nil, nil
}

func (a *ActivityService) GetExerciseSessionStats(ctx context.Context, req *pba.GetExerciseSessionStatsOccurrenceReq) (*pba.GetExerciseSessionStatsOccurrenceRes, error) {
	return nil, nil
}

func (a *ActivityService) StartActivityTracker(ctx context.Context, req *pba.StartActivityTrackerReq) (*pba.StartActivityTrackerRes, error) {
	return nil, nil
}

func (a *ActivityService) PauseActivityTracker(ctx context.Context, req *pba.PauseActivityTrackerReq) (*pba.PauseActivityTrackerRes, error) {
	return nil, nil
}

func (a *ActivityService) ResumeActivityTracker(ctx context.Context, req *pba.ResumeActivityTrackerReq) (*pba.ResumeActivityTrackerRes, error) {
	return nil, nil
}

func (a *ActivityService) StopActivityTracker(ctx context.Context, req *pba.StopActivityTrackerReq) (*pba.StopActivityTrackerRes, error) {
	return nil, nil
}

func (a *ActivityService) DeleteExerciseSession(ctx context.Context, req *pba.DeleteExerciseSessionReq) (*pba.NilRes, error) {
	return nil, nil
}

func (a ActivityService) DeleteAllExercisesSession(ctx context.Context, req *pba.DeleteAllExercisesSessionReq) (*pba.NilRes, error) {
	return nil, nil
}
