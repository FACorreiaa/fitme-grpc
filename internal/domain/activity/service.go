package activity

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

var mu sync.Mutex

type ServiceActivity struct {
	pba.UnimplementedActivityServer
	ctx              context.Context
	repo             domain.RepositoryActivity
	exerciseSessions map[string]*pba.XExerciseSession
	pausedTimers     map[string]time.Time
}

func NewCalculatorService(ctx context.Context, repo domain.RepositoryActivity) *ServiceActivity {
	return &ServiceActivity{
		ctx:              ctx,
		repo:             repo,
		exerciseSessions: make(map[string]*pba.XExerciseSession),
		pausedTimers:     make(map[string]time.Time),
	}
}

func (a *ServiceActivity) GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error) {
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

func (a *ServiceActivity) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
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

func (a *ServiceActivity) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
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

// we need to review the logic of all the sevices in the end
// an user can only see its own sessions
// so the userID comes from a session
// but a PT can search and select several userID on its network
func (a *ServiceActivity) GetUserExerciseSession(ctx context.Context, req *pba.GetUserExerciseSessionReq) (*pba.GetUserExerciseSessionRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivitiesByName")
	defer span.End()

	if req.PublicId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req name is required")
	}

	exerciseSession, err := a.repo.GetUserExerciseSession(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "exercise session not found")
		}
		return &pba.GetUserExerciseSessionRes{
			Success: false,
			Message: "Failed to retrieve exercise session",
			Session: exerciseSession.Session,
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

	return &pba.GetUserExerciseSessionRes{
		Success: true,
		Message: "Exercise session retrieved successfully",
		Session: exerciseSession.Session,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil

}

func (a *ServiceActivity) GetUserExerciseTotalData(ctx context.Context, req *pba.GetUserExerciseTotalDataReq) (*pba.GetUserExerciseTotalDataRes, error) {
	//userSession, ok := ctx.Value(auth.SessionManagerKey{}).(*auth.UserSession)
	//if !ok || userSession == nil {
	//	return nil, status.Error(codes.Unauthenticated, "failed to retrieve user session")
	//}

	sessionStats, err := a.repo.GetUserExerciseTotalData(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "exercise session not found")
		}
		return &pba.GetUserExerciseTotalDataRes{
			Success: false,
			Message: "Failed to retrieve exercise session",
			Session: sessionStats.Session,
			Response: &pba.BaseResponse{
				Upstream:  "activity-service",
				RequestId: domain.GenerateRequestID(ctx),
			},
		}, nil
	}

	return &pba.GetUserExerciseTotalDataRes{
		Success: true,
		Message: "Total session retrieved successfully",
		Session: sessionStats.Session,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

func (a *ServiceActivity) GetUserExerciseSessionStats(ctx context.Context, req *pba.GetUserExerciseSessionStatsReq) (*pba.GetUserExerciseSessionStatsRes, error) {
	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivitiesByName")
	defer span.End()

	userID := req.PublicId
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	stats, err := a.repo.GetUserExerciseSessionStats(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "exercise session not found")
		}
		return &pba.GetUserExerciseSessionStatsRes{
			Success:       false,
			Message:       "Failed to retrieve exercise stats",
			ExerciseCount: stats.ExerciseCount,
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

	return &pba.GetUserExerciseSessionStatsRes{
		Success:       true,
		Message:       "Exercise session stats retrieved successfully",
		ExerciseCount: stats.ExerciseCount,
		Response: &pba.BaseResponse{
			Upstream:  "activity-service",
			RequestId: domain.GenerateRequestID(ctx),
		},
	}, nil
}

// GetExerciseSessionStats maybe delete later
func (a *ServiceActivity) GetExerciseSessionStats(ctx context.Context, req *pba.GetExerciseSessionStatsOccurrenceReq) (*pba.GetExerciseSessionStatsOccurrenceRes, error) {
	return nil, nil
}

func (a *ServiceActivity) StartActivityTracker(ctx context.Context, req *pba.StartActivityTrackerReq) (*pba.StartActivityTrackerRes, error) {
	activityID := req.ActivityId
	userID := req.UserId

	tracer := otel.Tracer("FITDEV")
	ctx, span := tracer.Start(ctx, "GetActivity")
	defer span.End()

	if activityID == "" {
		return nil, status.Error(codes.InvalidArgument, "Activity ID is required")
	}
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	mu.Lock()
	_, found := a.exerciseSessions[userID]
	mu.Unlock()
	if found {
		return nil, status.Error(codes.FailedPrecondition, "activity tracker already started")
	}

	activityRes, err := a.repo.GetActivitiesByID(ctx, &pba.GetActivityIDReq{PublicId: activityID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve activity: %v", err)
	}

	currentTime := time.Now()
	exerciseSession := &pba.XExerciseSession{
		ExerciseSessionId: uuid.NewString(),
		UserId:            userID,
		ActivityId:        activityRes.Activity.ActivityId,
		SessionName:       activityRes.Activity.Name,
		StartTime:         timestamppb.New(time.Now()),
		CreatedAt:         timestamppb.New(currentTime),
	}

	a.exerciseSessions[exerciseSession.ExerciseSessionId] = exerciseSession

	span.SetAttributes(
		attribute.String("request.id", req.ActivityId),
		attribute.String("request.details", req.String()),
	)

	return &pba.StartActivityTrackerRes{
		Success:         true,
		Message:         "Activity tracker started",
		ExerciseSession: exerciseSession,
	}, nil
}

func (a *ServiceActivity) PauseActivityTracker(ctx context.Context, req *pba.PauseActivityTrackerReq) (*pba.PauseActivityTrackerRes, error) {
	// this is the user session! change after
	sessionID := req.SessionId

	if sessionID == "" {
		return nil, status.Error(codes.InvalidArgument, "Session ID is required")
	}

	mu.Lock()
	a.pausedTimers[sessionID] = time.Now()
	mu.Unlock()
	return &pba.PauseActivityTrackerRes{
		Success: true,
		Message: "Activity tracker paused",
	}, nil
}

func (a *ServiceActivity) ResumeActivityTracker(ctx context.Context, req *pba.ResumeActivityTrackerReq) (*pba.ResumeActivityTrackerRes, error) {
	sessionID := req.SessionId
	if sessionID == "" {
		return nil, status.Error(codes.InvalidArgument, "Session ID is required")
	}

	//for id, session := range a.exerciseSessions {
	//	fmt.Printf("Exercise Session ID: %s\n", id)
	//	fmt.Printf("Session Name: %s\n", session.SessionName)
	//	fmt.Printf("ExerciseSessionId: %s\n", session.ExerciseSessionId)
	//}

	mu.Lock()
	session, found := a.exerciseSessions[sessionID]
	mu.Unlock()

	if !found {
		return nil, status.Error(codes.FailedPrecondition, "activity tracker session not found")
	}

	if sessionID != session.ExerciseSessionId {
		return nil, status.Error(codes.FailedPrecondition, "activity tracker session not found")
	}

	//session, _ := a.exerciseSessions[sessionID]
	//if !found {
	//	return nil, status.Error(codes.FailedPrecondition, "activity tracker session not found")
	//}

	pausedTime, found := a.pausedTimers[sessionID]
	if !found {
		return nil, status.Error(codes.FailedPrecondition, "activity tracker was not paused")
	}

	startTime := session.StartTime.AsTime()

	pausedDuration := time.Since(pausedTime)
	adjustedStartTime := startTime.Add(pausedDuration)

	session.StartTime = timestamppb.New(adjustedStartTime)
	delete(a.pausedTimers, sessionID)

	return &pba.ResumeActivityTrackerRes{
		Success:         true,
		Message:         "Activity tracker resumed successfully",
		ExerciseSession: session, // Return the updated session with adjusted start time
	}, nil
}

func (a *ServiceActivity) StopActivityTracker(ctx context.Context, req *pba.StopActivityTrackerReq) (*pba.StopActivityTrackerRes, error) {
	// review protos.
	sessionID := req.SessionId
	if sessionID == "" {
		return nil, status.Error(codes.InvalidArgument, "Session ID is required")
	}
	mu.Lock()
	session, found := a.exerciseSessions[sessionID]
	mu.Unlock()
	if !found {
		return nil, status.Error(codes.FailedPrecondition, "activity tracker session not found")
	}

	activityRes, err := a.repo.GetActivitiesByID(ctx, &pba.GetActivityIDReq{PublicId: session.ActivityId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting activity: %v", err)
	}

	startUpTime := session.StartTime.AsTime()

	totalDurationSeconds := int(time.Since(startUpTime).Seconds())
	session.DurationHours = uint32(totalDurationSeconds / 3600)
	session.DurationMinutes = uint32((totalDurationSeconds % 3600) / 60)
	session.DurationSeconds = uint32(totalDurationSeconds % 60)

	caloriesPerSecond := activityRes.Activity.CaloriesPerHour / 3600
	session.CaloriesBurned = uint32(caloriesPerSecond * float32(totalDurationSeconds))
	//session.EndTime = timestamppb.New(time.Now()) // change protobuf generation
	session.EndTime = timestamppb.New(time.Now())
	err = a.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error saving exercise session to DB: %v", err)
	}
	delete(a.exerciseSessions, sessionID)

	return &pba.StopActivityTrackerRes{
		Success:         true,
		Message:         "Activity tracker stopped successfully",
		ExerciseSession: session,
	}, nil

}

func (a *ServiceActivity) DeleteExerciseSession(ctx context.Context, req *pba.DeleteExerciseSessionReq) (*pba.NilRes, error) {
	sessionID := req.PublicId

	if sessionID == "" {
		return nil, status.Error(codes.InvalidArgument, "Session ID is required")
	}

	_, err := a.repo.DeleteExerciseSession(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error deleting exercise session: %v", err)
	}

	return &pba.NilRes{}, nil
}

func (a *ServiceActivity) DeleteAllExercisesSession(ctx context.Context, req *pba.DeleteAllExercisesSessionReq) (*pba.NilRes, error) {
	_, err := a.repo.DeleteAllExercisesSession(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error deleting exercise session: %v", err)
	}
	return &pba.NilRes{}, nil
}
