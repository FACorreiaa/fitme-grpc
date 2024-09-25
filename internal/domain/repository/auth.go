package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pb "github.com/FACorreiaa/fitme-protos/modules/user/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

// NewAuthService creates a new AuthService
func NewAuthService(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *AuthService {
	return &AuthService{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (a *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = a.pgpool.Exec(ctx, `INSERT INTO "users" (username, email, password) VALUES ($1, $2, $3)`, req.Username, req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Message: "Registration successful"}, nil
}

func (a *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var passwordHash string
	var email string
	err := a.pgpool.QueryRow(ctx, `SELECT password, email FROM "users" WHERE username=$1`, req.Username).Scan(&passwordHash, &email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	sessionID, err := a.sessionManager.GenerateSession(auth.UserSession{
		Email:    req.Email,
		Username: req.Username,
	})

	fmt.Printf("sessionIDsessionID %s", sessionID)

	//err = a.redis.Set(ctx, req.Username, sessionID, 0).Err()
	//if err != nil {
	//	return nil, err
	//}

	return &pb.LoginResponse{Token: sessionID, Message: "Login successful!"}, nil
}

func (a *AuthService) Logout(ctx context.Context, req *pb.NilReq) (*pb.NilRes, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("unable to retrieve metadata")
	}

	authHeader := md["authorization"]
	if len(authHeader) != 1 {
		return nil, errors.New("invalid authorization header")
	}

	token := authHeader[0]
	sessionData, err := a.redis.Get(ctx, token).Result()
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	var session auth.UserSession
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	//username := session.Username
	//
	//err := a.sessionManager.SignOut(username)
	//if err != nil {
	//	return nil, err
	//}

	err = a.redis.Del(ctx, token).Err()
	if err != nil {
		return nil, errors.New("delete item")
	}

	return &pb.NilRes{}, nil
}

func (a *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	var passwordHash string
	err := a.pgpool.QueryRow(ctx, `SELECT password FROM "users" WHERE username=$1`, req.Username).Scan(&passwordHash)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.OldPassword))
	if err != nil {
		return nil, errors.New("invalid old password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = a.pgpool.Exec(ctx, `UPDATE "users" SET password=$1, updated_at=now() WHERE username=$2`, newHashedPassword, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.ChangePasswordResponse{Message: "Password changed successfully"}, nil
}

func (a *AuthService) ChangeEmail(ctx context.Context, req *pb.ChangeEmailRequest) (*pb.ChangeEmailResponse, error) {
	var passwordHash string
	err := a.pgpool.QueryRow(ctx, `SELECT password FROM "users" WHERE username=$1`, req.Username).Scan(&passwordHash)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	_, err = a.pgpool.Exec(ctx, `UPDATE "users" SET email=$1, updated_at=now() WHERE username=$2`, req.NewEmail, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.ChangeEmailResponse{Message: "Email changed successfully"}, nil
}

func (a *AuthService) GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error) {
	rows, err := a.pgpool.Query(ctx, `SELECT id, username, email, role, created_at, updated_at FROM "users"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var id, username, email, roleStr string
		var createdAt time.Time
		var updatedAt *time.Time

		err := rows.Scan(&id, &username, &email, &roleStr, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		fmt.Printf("role %s", roleStr)

		var role pb.User_Role
		switch roleStr {
		case "ADMIN":
			role = pb.User_ADMIN
		case "MODERATOR":
			role = pb.User_MODERATOR
		case "COACH":
			role = pb.User_COACH
		default:
			role = pb.User_ROLE_UNSPECIFIED
		}

		var updatedAtStr string
		if updatedAt != nil {
			updatedAtStr = updatedAt.Format(time.RFC3339)
		}

		users = append(users, &pb.User{
			Id:        id,
			Username:  username,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAtStr,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllUsersResponse{Users: users}, nil
}

func (a *AuthService) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	var u pb.User
	var createdAt time.Time
	var updatedAt *time.Time

	err := a.pgpool.QueryRow(ctx, `
			SELECT u.id, u.username, u.email, u.created_at, u.updated_at
			FROM "users" u
			WHERE id = $1`, req.Id).Scan(
		&u.Id, &u.Username, &u.Email, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	if updatedAt != nil {
		u.UpdatedAt = updatedAt.Format(time.RFC3339)
	}

	u.CreatedAt = createdAt.Format(time.RFC3339)

	return &pb.GetUserByIDResponse{User: &u}, nil
}

func (a *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Execute the delete query
	commandTag, err := a.pgpool.Exec(ctx, `DELETE FROM "users" WHERE user = $1`, req.Id)
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %v", err)
	}

	// Check if any row was deleted
	if commandTag.RowsAffected() == 0 {
		return nil, errors.New("user not found")
	}

	return &pb.DeleteUserResponse{Message: "user deleted successfully"}, nil
}

func (a *AuthService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Execute the update query
	commandTag, err := a.pgpool.Exec(ctx, `
		UPDATE "users"
		SET username = $1, email = $2, updated_at = NOW()
		WHERE id = $3`,
		req.User.Username, req.User.Email, req.User.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	// Check if any row was updated
	if commandTag.RowsAffected() == 0 {
		return nil, errors.New("user not found")
	}

	return &pb.UpdateUserResponse{Message: "user updated successfully"}, nil
}

// InsertUser change this later

func (a *AuthService) InsertUser(ctx context.Context, req *pb.InsertUserRequest) (*pb.InsertUserResponse, error) {
	// Insert the new user
	_, err := a.pgpool.Exec(ctx, `
		INSERT INTO "users" (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		req.User.Id, req.User.Username, req.User.Email, req.User.PasswordHash, req.User.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	return &pb.InsertUserResponse{Message: "user inserted successfully"}, nil
}
