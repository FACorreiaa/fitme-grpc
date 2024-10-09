package activity

import (
	"context"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"

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
	return nil, nil
}

func (a *ActivityService) GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error) {
	return nil, nil
}
func (a *ActivityService) GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error) {
	return nil, nil
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
