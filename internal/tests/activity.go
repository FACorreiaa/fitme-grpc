package tests

import (
	"context"
	"testing"

	"github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/activity"
)

//func TestSaveExerciseSession(t *testing.T) {
//	repo := &MockRepositoryActivity{} // Assuming you've implemented a mock repository
//	service := NewServiceActivity(repo)
//
//	session := &ExerciseSession{UserID: 1, ActivityID: 1, SessionName: "Morning Run"}
//	err := service.SaveExerciseSession(context.Background(), session)
//	if err != nil {
//		t.Fatalf("expected no error, got %v", err)
//	}
//}

type MockRepositoryActivity struct {
	mock.Mock
}

type ServiceActivityTestSuite struct {
	suite.Suite
	ctx      context.Context
	mockRepo *MockRepositoryActivity
	service  *activity.ServiceActivity
}

func (suite *ServiceActivityTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRepo = new(MockRepositoryActivity)
	suite.service = activity.NewCalculatorService(suite.ctx, suite.mockRepo)
}

func (m *MockRepositoryActivity) GetActivity(ctx context.Context, req *generated.GetActivityReq) (*generated.GetActivityRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetActivityRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetActivitiesByID(ctx context.Context, req *generated.GetActivityIDReq) (*generated.GetActivityIDRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetActivityIDRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetActivitiesByName(ctx context.Context, req *generated.GetActivityNameReq) (*generated.GetActivityNameRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetActivityNameRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetUserExerciseSession(ctx context.Context, req *generated.GetUserExerciseSessionReq) (*generated.GetUserExerciseSessionRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetUserExerciseSessionRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetUserExerciseTotalData(ctx context.Context, req *generated.GetUserExerciseTotalDataReq) (*generated.GetUserExerciseTotalDataRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetUserExerciseTotalDataRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetUserExerciseSessionStats(ctx context.Context, req *generated.GetUserExerciseSessionStatsReq) (*generated.GetUserExerciseSessionStatsRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetUserExerciseSessionStatsRes), args.Error(1)
}

func (m *MockRepositoryActivity) GetExerciseSessionStats(ctx context.Context, req *generated.GetExerciseSessionStatsOccurrenceReq) (*generated.GetExerciseSessionStatsOccurrenceRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.GetExerciseSessionStatsOccurrenceRes), args.Error(1)
}

func (m *MockRepositoryActivity) StartActivityTracker(ctx context.Context, req *generated.StartActivityTrackerReq) (*generated.StartActivityTrackerRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.StartActivityTrackerRes), args.Error(1)
}

func (m *MockRepositoryActivity) PauseActivityTracker(ctx context.Context, req *generated.PauseActivityTrackerReq) (*generated.PauseActivityTrackerRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.PauseActivityTrackerRes), args.Error(1)
}

func (m *MockRepositoryActivity) ResumeActivityTracker(ctx context.Context, req *generated.ResumeActivityTrackerReq) (*generated.ResumeActivityTrackerRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.ResumeActivityTrackerRes), args.Error(1)
}

func (m *MockRepositoryActivity) StopActivityTracker(ctx context.Context, req *generated.StopActivityTrackerReq) (*generated.StopActivityTrackerRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.StopActivityTrackerRes), args.Error(1)
}

func (m *MockRepositoryActivity) DeleteExerciseSession(ctx context.Context, req *generated.DeleteExerciseSessionReq) (*generated.NilRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.NilRes), args.Error(1)
}

func (m *MockRepositoryActivity) DeleteAllExercisesSession(ctx context.Context, req *generated.DeleteAllExercisesSessionReq) (*generated.NilRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*generated.NilRes), args.Error(1)
}

func (m *MockRepositoryActivity) SaveSession(ctx context.Context, req *generated.XExerciseSession) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// setup to complete after the workout rewrite

// TestServiceActivity_GetActivity placeholder
func TestServiceActivity_GetActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := new(MockRepositoryActivity)
	ctx := context.Background()

	service := NewCalculatorService(ctx, mockRepo)

	// Test case: success scenario
	req := &generated.GetActivityReq{PublicId: "123"}
	mockRepo.On("GetActivity", ctx, req).Return(&domain.ActivityResponse{
		Activity: []*domain.Activity{
			{
				ActivityId:        "123",
				UserId:            "user123",
				Name:              "Running",
				CaloriesPerHour:   500,
				DurationInMinutes: 60,
				TotalCalories:     500,
				CreatedAt:         timestamppb.Now(),
				UpdatedAt:         timestamppb.Now(),
			},
		},
	}, nil)

	res, err := service.GetActivity(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Activities retrieved successfully", res.Message)
	assert.Equal(t, "Running", res.Activity[0].Name)

	reqEmpty := &generated.GetActivityReq{PublicId: ""}
	_, err = service.GetActivity(ctx, reqEmpty)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	mockRepo.On("GetActivity", ctx, req).Return(nil, errors.New("database error"))
	_, err = service.GetActivity(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

// TestServiceActivity_GetActivitiesByID placeholder
func TestServiceActivity_GetActivitiesByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := new(MockRepositoryActivity)
	ctx := context.Background()

	service := NewCalculatorService(ctx, mockRepo)

	req := &generated.GetActivityIDReq{PublicId: "123"}
	mockRepo.On("GetActivitiesByID", ctx, req).Return(&domain.ActivityResponse{
		Activity: &domain.Activity{
			ActivityId:        "123",
			UserId:            "user123",
			Name:              "Running",
			CaloriesPerHour:   500,
			DurationInMinutes: 60,
			TotalCalories:     500,
			CreatedAt:         timestamppb.Now(),
			UpdatedAt:         timestamppb.Now(),
		},
	}, nil)

	res, err := service.GetActivitiesByID(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Activity retrieved successfully", res.Message)

	mockRepo.On("GetActivitiesByID", ctx, req).Return(nil, pgx.ErrNoRows)
	_, err = service.GetActivitiesByID(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

// TestServiceActivity_GetUserExerciseSession placeholder
func TestServiceActivity_GetUserExerciseSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := new(MockRepositoryActivity)
	ctx := context.Background()

	service := NewCalculatorService(ctx, mockRepo)

	req := &generated.GetUserExerciseSessionReq{PublicId: "session123"}
	mockRepo.On("GetUserExerciseSession", ctx, req).Return(&domain.SessionStats{
		Session: &generated.XExerciseSession{
			ExerciseSessionId: "session123",
			UserId:            "user123",
			ActivityId:        "activity123",
			SessionName:       "Running",
		},
	}, nil)

	res, err := service.GetUserExerciseSession(ctx, req)
	assert.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Running", res.Session.SessionName)

	mockRepo.On("GetUserExerciseSession", ctx, req).Return(nil, pgx.ErrNoRows)
	_, err = service.GetUserExerciseSession(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
