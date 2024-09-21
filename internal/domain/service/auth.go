package service

import (
	"context"

	"github.com/FACorreiaa/ink-app-backend-grpc/internal/domain/repository"
	pb "github.com/FACorreiaa/ink-app-backend-protos/modules/user/generated"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.repo.Register(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.repo.Login(ctx, req)
}

func (s *AuthService) Logout(ctx context.Context, req *pb.NilReq) (*pb.NilRes, error) {
	return s.repo.Logout(ctx, req)
}

func (s *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return s.repo.ChangePassword(ctx, req)
}

func (s *AuthService) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error) {
	return s.repo.ChangeEmail(ctx, req)
}

func (s *AuthService) GetAllUsers(ctx context.Context, _ *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := s.repo.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return s.repo.DeleteUser(ctx, req)
}

func (s *AuthService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return s.repo.UpdateUser(ctx, req)
}

func (s *AuthService) InsertUser(ctx context.Context, req *pb.InsertUserRequest) (*pb.InsertUserResponse, error) {
	return s.repo.InsertUser(ctx, req)
}
