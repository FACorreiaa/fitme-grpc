package calculator

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (c *CalculatorRepository) GetUsersMacros(ctx context.Context) (*pbc.GetAllUserMacrosResponse, error) {
	macroDistribution := make([]*pbc.UserMacroDistribution, 0)
	query := `SELECT id, user_id, age, height, weight,
                      gender, system, activity, activity_description, objective,
					  objective_description, calories_distribution, calories_distribution_description,
                      protein, fats, carbs, bmr, tdee, goal, created_at
				FROM user_macro_distribution
				ORDER BY created_at`

	rows, err := c.pgpool.Query(ctx, query)
	if err != nil {
		if errors.Is(rows.Err(), pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "No user macros found")
		}
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}
	defer rows.Close()

	for rows.Next() {
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

func (c *CalculatorRepository) CreateUserMacro(ctx context.Context, req *pbc.UserMacroDistribution) (*pbc.UserMacroDistribution, error) {
	query := `INSERT INTO user_macro_distribution (user_id, age, height, weight,
                                    gender, system, activity, activity_description, objective,
                                    objective_description, calories_distribution, calories_distribution_description,
                                    protein, fats, carbs, bmr, tdee, goal, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
              RETURNING *`

	var macro pbc.UserMacroDistribution
	var createdAt time.Time

	rows, err := c.pgpool.Query(ctx, query,
		req.UserId, req.Age, req.Height, req.Weight, req.Gender, req.System, req.Activity,
		req.ActivityDescription, req.Objective, req.ObjectiveDescription, req.CaloriesDistribution,
		req.CaloriesDistributionDescription, req.Protein, req.Fats, req.Carbs, req.Bmr, req.Tdee,
		req.Goal, createdAt,
	)

	if err != nil {
		log.Printf("Query execution error: %v", err) // Log detailed error
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	//accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[pbc.UserMacroDistribution])
	//if err != nil {
	//	return nil, fmt.Errorf("failed collecting rows: %w", err)
	//}
	//
	//if len(accounts) == 0 {
	//	return nil, fmt.Errorf("no rows returned")
	//}
	if rows.Next() {
		err = rows.Scan(
			&macro.Id, &macro.UserId, &macro.Age, &macro.Height, &macro.Weight, &macro.Gender, &macro.System,
			&macro.Activity, &macro.ActivityDescription, &macro.Objective, &macro.ObjectiveDescription,
			&macro.CaloriesDistribution, &macro.CaloriesDistributionDescription, &macro.Protein,
			&macro.Fats, &macro.Carbs, &macro.Bmr, &macro.Tdee, &macro.Goal, &createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert `createdAt` (Go time.Time) to Protobuf Timestamp
		macro.CreatedAt = timestamppb.New(createdAt)
	}

	return &macro, nil
}

func (c *CalculatorRepository) DeleteUserMacro(ctx context.Context, macroID string) error {
	if macroID == "" {
		return status.Error(codes.InvalidArgument, "macroID is required")
	}

	query := `DELETE FROM user_macro_distribution WHERE id = $1`

	cmdTag, err := c.pgpool.Exec(ctx, query, macroID)
	if err != nil {
		return fmt.Errorf("failed to delete user macro: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return status.Error(codes.NotFound, "macro not found")
	}

	return nil
}
