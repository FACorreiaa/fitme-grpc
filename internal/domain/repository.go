package domain

import (
	"context"

	pba "github.com/FACorreiaa/fitme-protos/modules/activity/generated"
	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	pbm "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	pb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
	pbw "github.com/FACorreiaa/fitme-protos/modules/workout/generated"
)

type AuthRepository interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Logout(ctx context.Context, req *pb.NilReq) (*pb.NilRes, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error)

	// GetAllUsers Users Methods
	GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error)
	GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	InsertUser(ctx context.Context, req *pb.InsertUserRequest) (*pb.InsertUserResponse, error)
}

type CalculatorRepository interface {
	CreateUserMacro(ctx context.Context, req *pbc.UserMacroDistribution) (*pbc.UserMacroDistribution, error)
	GetUsersMacros(ctx context.Context, req *pbc.GetAllUserMacrosRequest) (*pbc.GetAllUserMacrosResponse, error)
	GetUserMacros(ctx context.Context, req *pbc.GetUserMacroRequest) (*pbc.GetUserMacroResponse, error)
	DeleteUserMacro(ctx context.Context, req *pbc.DeleteUserMacroRequest) (*pbc.DeleteUserMacroResponse, error)
}

type RepositoryActivity interface {
	GetActivity(ctx context.Context, req *pba.GetActivityReq) (*pba.GetActivityRes, error)                                                         // done
	GetActivitiesByID(ctx context.Context, req *pba.GetActivityIDReq) (*pba.GetActivityIDRes, error)                                               // done
	GetActivitiesByName(ctx context.Context, req *pba.GetActivityNameReq) (*pba.GetActivityNameRes, error)                                         // done
	GetUserExerciseSession(ctx context.Context, req *pba.GetUserExerciseSessionReq) (*pba.GetUserExerciseSessionRes, error)                        // done
	GetUserExerciseTotalData(ctx context.Context, req *pba.GetUserExerciseTotalDataReq) (*pba.GetUserExerciseTotalDataRes, error)                  // done
	GetUserExerciseSessionStats(ctx context.Context, req *pba.GetUserExerciseSessionStatsReq) (*pba.GetUserExerciseSessionStatsRes, error)         // done
	GetExerciseSessionStats(ctx context.Context, req *pba.GetExerciseSessionStatsOccurrenceReq) (*pba.GetExerciseSessionStatsOccurrenceRes, error) //done
	StartActivityTracker(ctx context.Context, req *pba.StartActivityTrackerReq) (*pba.StartActivityTrackerRes, error)
	PauseActivityTracker(ctx context.Context, req *pba.PauseActivityTrackerReq) (*pba.PauseActivityTrackerRes, error)
	ResumeActivityTracker(ctx context.Context, req *pba.ResumeActivityTrackerReq) (*pba.ResumeActivityTrackerRes, error)
	StopActivityTracker(ctx context.Context, req *pba.StopActivityTrackerReq) (*pba.StopActivityTrackerRes, error)
	DeleteExerciseSession(ctx context.Context, req *pba.DeleteExerciseSessionReq) (*pba.NilRes, error)
	DeleteAllExercisesSession(ctx context.Context, req *pba.DeleteAllExercisesSessionReq) (*pba.NilRes, error)
	SaveSession(ctx context.Context, req *pba.XExerciseSession) error
}

type RepositoryWorkout interface {
	GetExercises(ctx context.Context, req *pbw.GetExercisesReq) (*pbw.GetExercisesRes, error)
	GetExerciseID(ctx context.Context, req *pbw.GetExerciseIDReq) (*pbw.GetExerciseIDRes, error)
	CreateExercise(ctx context.Context, req *pbw.CreateExerciseReq) (*pbw.CreateExerciseRes, error)
	DeleteExercise(ctx context.Context, req *pbw.DeleteExerciseReq) (*pbw.NilRes, error)
	UpdateExercise(ctx context.Context, req *pbw.UpdateExerciseReq) (*pbw.UpdateExerciseRes, error)
	CreateWorkoutPlan(ctx context.Context, req *pbw.InsertWorkoutPlanReq) (*pbw.InsertWorkoutPlanRes, error)
	GetWorkoutPlans(ctx context.Context, req *pbw.GetWorkoutPlansReq) (*pbw.GetWorkoutPlansRes, error)
	GetWorkoutPlan(ctx context.Context, req *pbw.GetWorkoutPlanReq) (*pbw.GetWorkoutPlanRes, error)
	DeleteWorkoutPlan(ctx context.Context, req *pbw.DeleteWorkoutPlanReq) (*pbw.NilRes, error)
	UpdateWorkoutPlan(ctx context.Context, req *pbw.UpdateWorkoutPlanReq) (*pbw.UpdateWorkoutPlanRes, error)
	GetWorkoutPlanExercises(ctx context.Context, req *pbw.GetWorkoutPlanExercisesReq) (*pbw.GetWorkoutPlanExercisesRes, error)
	GetWorkoutPlanExercisesByID(ctx context.Context, req *pbw.GetExerciseByIdWorkoutPlanReq) (*pbw.GetExerciseByIdWorkoutPlanRes, error)
	InsertExerciseWorkoutPlan(ctx context.Context, req *pbw.InsertExerciseWorkoutPlanReq) (*pbw.NilRes, error)
	DeleteExerciseWorkoutPlan(ctx context.Context, req *pbw.DeleteExerciseByIdWorkoutPlanReq) (*pbw.NilRes, error)
	UpdateExerciseWorkoutPLan(ctx context.Context, req *pbw.UpdateExerciseByIdWorkoutPlanReq) (*pbw.UpdateExerciseByIdWorkoutPlanRes, error)
}

type RepositoryMeasurement interface {

	// Weight
	CreateWeight(ctx context.Context, req *pbm.CreateWeightReq) (*pbm.XWeight, error)
	GetWeights(ctx context.Context) ([]*pbm.XWeight, error)
	GetWeight(ctx context.Context, req *pbm.GetWeightReq) (*pbm.XWeight, error)
	DeleteWeight(ctx context.Context, req *pbm.DeleteWeightReq) (*pbm.NilRes, error)
	UpdateWeight(ctx context.Context, req *pbm.UpdateWeightReq) (*pbm.XWeight, error)

	// waterIntake
	CreateWaterMeasurement(ctx context.Context, req *pbm.CreateWaterIntakeReq) (*pbm.XWaterIntake, error)
	GetWaterMeasurements(ctx context.Context) ([]*pbm.XWaterIntake, error)
	GetWaterMeasurement(ctx context.Context, req *pbm.GetWaterIntakeReq) (*pbm.XWaterIntake, error)
	DeleteWaterMeasurement(ctx context.Context, req *pbm.DeleteWaterIntakeReq) (*pbm.NilRes, error)
	UpdateWaterMeasurement(ctx context.Context, req *pbm.UpdateWaterIntakeReq) (*pbm.XWaterIntake, error)

	// wasteline
	CreateWasteLineMeasurement(ctx context.Context, req *pbm.CreateWasteLineReq) (*pbm.XWasteLine, error)
	GetWasteLineMeasurements(ctx context.Context) ([]*pbm.XWasteLine, error)
	GetWasteLineMeasurement(ctx context.Context, req *pbm.GetWasteLineReq) (*pbm.XWasteLine, error)
	DeleteWasteLineMeasurement(ctx context.Context, req *pbm.DeleteWasteLineReq) (*pbm.NilRes, error)
	UpdateWasteLineMeasurement(ctx context.Context, req *pbm.UpdateWasteLineReq) (*pbm.XWasteLine, error)
}
