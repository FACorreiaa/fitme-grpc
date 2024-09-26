package calculator

import (
	"context"
	"fmt"
	"time"

	pbc "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type CalculatorRepository struct {
	pbc.UnimplementedCalculatorServiceServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewCalculatorRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *CalculatorRepository {
	return &CalculatorRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func (c *CalculatorRepository) CreateUserMacro(ctx context.Context, req *pbc.CreateUserMacroRequest) (*pbc.CreateUserMacroResponse, error) {
	return nil, nil
}

func (c *CalculatorRepository) GetUsersMacros(ctx context.Context, req *pbc.GetAllUserMacrosRequest) (*pbc.GetAllUserMacrosResponse, error) {
	macroDistribution := make([]*pbc.UserMacroDistribution, 0)
	query := `SELECT user_id, age, height, weight,
                      gender, system, activity, activity_description, objective,
					  objective_description, calories_distribution, calories_distribution_description,
                      protein, fats, carbs, bmr, tdee, goal, created_at
				FROM user_macro_distribution
				WHERE id = $1
				ORDER BY created_at`

	rows, err := c.pgpool.Query(ctx, query, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert the timestamp to string for the proto message
		macro.CreatedAt = createdAt.Format(time.RFC3339)

		// Append the mapped macro to the slice
		macroDistribution = append(macroDistribution, &macro)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return &pbc.GetAllUserMacrosResponse{UserMacros: macroDistribution}, nil
}

func (c *CalculatorRepository) GetUserMacros(ctx context.Context, req *pbc.GetUserMacroRequest) (*pbc.GetUserMacroResponse, error) {
	var macroDistribution *pbc.UserMacroDistribution
	planID := req.PlanId
	query := `SELECT id, user_id, age, height, weight,
                      gender, system, activity, activity_description, objective,
					  objective_description, calories_distribution, calories_distribution_description,
                      protein, fats, carbs, bmr, tdee, goal, created_at
				FROM user_macro_distribution
				WHERE id = $1
				ORDER BY created_at`

	err := c.pgpool.QueryRow(ctx, query, query, planID)

	if err != nil {
		return nil, fmt.Errorf("macro not found: %w", err)
	}

	return &pbc.GetUserMacroResponse{UserMacro: macroDistribution}, nil
}

func (c *CalculatorRepository) InsertDietGoals(ctx context.Context, req *pbc.UserMacroDistribution) (*pbc.UserMacroDistribution, error) {
	return nil, nil
}
