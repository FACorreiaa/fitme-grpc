package calculator

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type CalculatorRepository struct {
	pbc.UnimplementedCalculatorServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewCalculatorRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *CalculatorRepository {
	return &CalculatorRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (c *CalculatorRepository) GetUsersMacros(ctx context.Context, req *pbc.GetAllUserMacrosRequest) (*pbc.GetAllUserMacrosResponse, error) {
	macroDistribution := make([]*pbc.UserMacroDistribution, 0)
	query := `SELECT id, user_id, age, height, weight,
                      gender, system, activity, activity_description, objective,
					  objective_description, calories_distribution, calories_distribution_description,
                      protein, fats, carbs, bmr, tdee, goal, created_at
				FROM user_macro_distribution
				ORDER BY created_at`

	//if req.UserId != "" {
	//	query += ` WHERE user_id = $1 ORDER BY created_at`
	//} else {
	//	query += ` ORDER BY created_at`
	//}

	rows, err := c.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(rows.Err(), pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "No user macros found")
		}
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "operation cancelled: %v", ctx.Err())
		default:
		}

		var macro pbc.UserMacroDistribution
		var createdAt time.Time

		err := rows.Scan(
			&macro.Id, &macro.UserId, &macro.Age, &macro.Height, &macro.Weight,
			&macro.Gender, &macro.System, &macro.Activity, &macro.ActivityDescription,
			&macro.Objective, &macro.ObjectiveDescription, &macro.CaloriesDistribution,
			&macro.CaloriesDistributionDescription, &macro.Protein, &macro.Fats,
			&macro.Carbs, &macro.Bmr, &macro.Tdee, &macro.Goal, &createdAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, status.Error(codes.NotFound, "User macro not found")
			}
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		macroDistribution = append(macroDistribution, &macro)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return &pbc.GetAllUserMacrosResponse{UserMacros: macroDistribution}, nil
}

func (c *CalculatorRepository) GetUserMacros(ctx context.Context, req *pbc.GetUserMacroRequest) (*pbc.GetUserMacroResponse, error) {
	var macroDistribution pbc.UserMacroDistribution
	planID := req.PlanId
	var createdAt time.Time

	if planID == "" {
		return nil, status.Error(codes.InvalidArgument, "planID is required")
	}

	query := `SELECT id, user_id, age, height, weight,
                      gender, system, activity, activity_description, objective,
					  objective_description, calories_distribution, calories_distribution_description,
                      protein, fats, carbs, bmr, tdee, goal, created_at
				FROM user_macro_distribution
				WHERE id = $1
				ORDER BY created_at`

	err := c.pgpool.QueryRow(ctx, query, planID).Scan(
		&macroDistribution.Id, &macroDistribution.UserId, &macroDistribution.Age, &macroDistribution.Height,
		&macroDistribution.Weight, &macroDistribution.Gender, &macroDistribution.System, &macroDistribution.Activity,
		&macroDistribution.ActivityDescription, &macroDistribution.Objective, &macroDistribution.ObjectiveDescription,
		&macroDistribution.CaloriesDistribution, &macroDistribution.CaloriesDistributionDescription, &macroDistribution.Protein,
		&macroDistribution.Fats, &macroDistribution.Carbs, &macroDistribution.Bmr, &macroDistribution.Tdee,
		&macroDistribution.Goal, &createdAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return nil, fmt.Errorf("failed to retrieve user macro: %w", err)
	}

	macroDistribution.CreatedAt = timestamppb.New(createdAt)

	return &pbc.GetUserMacroResponse{UserMacro: &macroDistribution}, nil
}

func (c *CalculatorRepository) CreateUserMacro(ctx context.Context, req *pbc.CreateUserMacroRequest) (*pbc.UserMacroDistribution, error) {
	// Start a transaction if you want to ensure atomic update
	tx, err := c.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// If the incoming macro should be "current," set all other macros for this user to false
	if req.IsCurrent {
		_, err = tx.Exec(ctx, `
	 UPDATE user_macro_distribution
	 SET is_current = false
	 WHERE user_id = $1
	`, req.UserMacro.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to reset is_current: %w", err)
		}
	}

	// Then insert the new macro
	query := `
    INSERT INTO user_macro_distribution (
      user_id, age, height, weight, gender, system, activity, activity_description,
      objective, objective_description, calories_distribution, calories_distribution_description,
      protein, fats, carbs, bmr, tdee, goal, created_at, is_current
    )
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,now(),$19)
    RETURNING
      id, user_id, age, height, weight, gender, system, activity, activity_description,
      objective, objective_description, calories_distribution, calories_distribution_description,
      protein, fats, carbs, bmr, tdee, goal, created_at, is_current`

	var macro pbc.UserMacroDistribution
	var createdAt time.Time
	var isCurrent bool
	userMacro := req.UserMacro
	row := tx.QueryRow(ctx, query,
		userMacro.UserId, userMacro.Age, userMacro.Height, userMacro.Weight, userMacro.Gender, userMacro.System, userMacro.Activity,
		userMacro.ActivityDescription, userMacro.Objective, userMacro.ObjectiveDescription,
		userMacro.CaloriesDistribution, userMacro.CaloriesDistributionDescription,
		userMacro.Protein, userMacro.Fats, userMacro.Carbs, userMacro.Bmr, userMacro.Tdee, userMacro.Goal, req.IsCurrent,
	)

	err = row.Scan(
		&macro.Id, &macro.UserId, &macro.Age, &macro.Height, &macro.Weight, &macro.Gender,
		&macro.System, &macro.Activity, &macro.ActivityDescription, &macro.Objective,
		&macro.ObjectiveDescription, &macro.CaloriesDistribution,
		&macro.CaloriesDistributionDescription, &macro.Protein,
		&macro.Fats, &macro.Carbs, &macro.Bmr, &macro.Tdee,
		&macro.Goal, &createdAt, &isCurrent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	macro.CreatedAt = timestamppb.New(createdAt)
	req.IsCurrent = isCurrent

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &macro, nil
}

func (c *CalculatorRepository) DeleteUserMacro(ctx context.Context, req *pbc.DeleteUserMacroRequest) (*pbc.DeleteUserMacroResponse, error) {
	macroID := req.MacroId

	if macroID == "" {
		return nil, status.Error(codes.InvalidArgument, "macroID is required")
	}

	query := `DELETE FROM user_macro_distribution WHERE id = $1`

	cmdTag, err := c.pgpool.Exec(ctx, query, macroID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user macro: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "macro not found")
	}

	return &pbc.DeleteUserMacroResponse{}, nil
}

func (c *CalculatorRepository) SetActiveUserMacro(ctx context.Context, userID, macroID string) (*pbc.UserMacroDistribution, error) {
	tx, err := c.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
        UPDATE user_macro_distribution
        SET is_current = false
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to reset active macros: %w", err)
	}

	query := `
        UPDATE user_macro_distribution
        SET is_current = true
        WHERE id = $1 AND user_id = $2
        RETURNING
            id, user_id, age, height, weight, gender, system, activity, activity_description,
            objective, objective_description, calories_distribution, calories_distribution_description,
            protein, fats, carbs, bmr, tdee, goal, created_at, is_current
    `
	row := tx.QueryRow(ctx, query, macroID, userID)

	var macro pbc.UserMacroDistribution
	var createdAt time.Time
	var isCurrent bool

	err = row.Scan(
		&macro.Id, &macro.UserId, &macro.Age, &macro.Height, &macro.Weight, &macro.Gender,
		&macro.System, &macro.Activity, &macro.ActivityDescription, &macro.Objective,
		&macro.ObjectiveDescription, &macro.CaloriesDistribution, &macro.CaloriesDistributionDescription,
		&macro.Protein, &macro.Fats, &macro.Carbs, &macro.Bmr, &macro.Tdee, &macro.Goal,
		&createdAt, &isCurrent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan updated macro: %w", err)
	}
	macro.CreatedAt = timestamppb.New(createdAt)

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &macro, nil
}
