package domain

import (
	"context"

	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	pb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
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
	GetUsersMacros(ctx context.Context) (*pbc.GetAllUserMacrosResponse, error)
	GetUserMacros(ctx context.Context, req *pbc.GetUserMacroRequest) (*pbc.GetUserMacroResponse, error)
}
