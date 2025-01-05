package meals

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	pbml "github.com/FACorreiaa/fitme-protos/modules/meal/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

type MealPlanRepository struct {
	pbml.UnimplementedMealPlanServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type MealRepository struct {
	pbml.UnimplementedMealServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type DietPreferenceRepository struct {
	pbml.UnimplementedDietPreferenceServiceServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type FoodLogRepository struct {
	pbml.UnimplementedFoodLogServiceServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type IngredientRepository struct {
	pbml.UnimplementedIngredientsServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type TrackMealProgressRepository struct {
	pbml.UnimplementedTrackMealProgressServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type GoalRecommendationRepository struct {
	pbml.UnimplementedGoalRecommendationServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

type MealReminderRepository struct {
	pbml.UnimplementedMealReminderServer
	pgpool         *pgxpool.Pool
	redis          *redis.Client
	sessionManager *auth.SessionManager
}

func NewMealPlanRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *MealPlanRepository {
	return &MealPlanRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewMealRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *MealRepository {
	return &MealRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewDietPreferenceRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *DietPreferenceRepository {
	return &DietPreferenceRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewFoodLogRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *FoodLogRepository {
	return &FoodLogRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewIngredientRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *IngredientRepository {
	return &IngredientRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewTrackMealProgressRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *TrackMealProgressRepository {
	return &TrackMealProgressRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewGoalRecommendationRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *GoalRecommendationRepository {
	return &GoalRecommendationRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func NewMealReminderRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *MealReminderRepository {
	return &MealReminderRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
}

func nullTimeToTimestamppb(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

func (i *IngredientRepository) GetIngredient(ctx context.Context, req *pbml.GetIngredientReq) (*pbml.GetIngredientRes, error) {
	ingredient := &Ingredient{}
	ingredientProto := &pbml.XIngredient{}

	tx, err := i.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		SELECT id, name, calories, protein, carbohydrates_total, fat_total
		FROM ingredients
		WHERE id = $1 AND (user_id = $2 OR user_id IS NULL) -- Restrict to user's or global ingredients
	`, req.IngredientId, req.UserId).Scan(
		&ingredient.ID, &ingredient.Name, &ingredient.Calories,
		&ingredient.Protein, &ingredient.CarbohydratesTotal, &ingredient.FatTotal,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "ingredient not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get ingredient: %v", err)
	}

	createdAt := timestamppb.New(ingredient.CreatedAt)
	var updatedAt sql.NullTime

	if ingredient.UpdatedAt.Valid {
		updatedAt = ingredient.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	ingredientProto.IngredientId = ingredient.ID.String()
	ingredientProto.Name = ingredient.Name
	ingredientProto.Calories = ingredient.Calories
	ingredientProto.Protein = ingredient.Protein
	ingredientProto.CarbohydratesTotal = ingredient.CarbohydratesTotal
	ingredientProto.FatTotal = ingredient.FatTotal
	ingredientProto.CreatedAt = createdAt
	ingredientProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &pbml.GetIngredientRes{
		Success:    true,
		Message:    "Ingredients retrieved successfully",
		Ingredient: ingredientProto,
		Response: &pbml.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil
}

func (i *IngredientRepository) GetIngredients(ctx context.Context, req *pbml.GetIngredientsReq) (*pbml.GetIngredientsRes, error) {
	ingredients := make([]*pbml.XIngredient, 0)
	query := `
		SELECT id, name, calories, protein, carbohydrates_total, fat_total
		FROM ingredients
		WHERE (user_id = $1 OR user_id IS NULL) -- Restrict to user's or global ingredients
	`

	rows, err := i.pgpool.Query(ctx, query, req.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no ingredients found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ingredientProto := pbml.XIngredient{}
		ingredient := &Ingredient{}

		err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Calories, &ingredient.Protein, &ingredient.CarbohydratesTotal, &ingredient.FatTotal)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("no ingredients found: %w", err)
			}
			return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}

		createdAt := timestamppb.New(ingredient.CreatedAt)
		var updatedAt sql.NullTime

		if ingredient.UpdatedAt.Valid {
			updatedAt = ingredient.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		ingredientProto.IngredientId = ingredient.ID.String()
		ingredientProto.Name = ingredient.Name
		ingredientProto.Calories = ingredient.Calories
		ingredientProto.Protein = ingredient.Protein
		ingredientProto.CarbohydratesTotal = ingredient.CarbohydratesTotal
		ingredientProto.FatTotal = ingredient.FatTotal
		ingredientProto.CreatedAt = createdAt
		ingredientProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		ingredients = append(ingredients, &ingredientProto)
	}

	return &pbml.GetIngredientsRes{
		Success:     true,
		Message:     "Ingredients retrieved successfully",
		Ingredients: ingredients,
		Response:    &pbml.BaseResponse{Upstream: "workout-service"},
	}, nil

}

func (i *IngredientRepository) CreateIngredient(ctx context.Context, req *pbml.CreateIngredientReq) (*pbml.XIngredient, error) {
	ingredientProto := &pbml.XIngredient{}
	ingredient := &Ingredient{}

	tx, err := i.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO ingredients (name, calories, serving_size,
			protein, fat_total, fat_saturated, carbohydrates_total, fiber, sugar, sodium, potassium, cholesterol,
			created_at, user_id
		)
		Values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	currentTime := time.Now()
	var ingredientID uuid.UUID

	err = tx.QueryRow(ctx, query,
		req.Name,
		req.Calories,
		req.ServingSize,
		req.Protein,
		req.FatTotal,
		req.FatSaturated,
		req.CarbohydratesTotal,
		req.Fiber,
		req.Sugar,
		req.Sodium,
		req.Potassium,
		req.Cholesterol,
		currentTime,
		req.UserId).Scan(&ingredientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ingredient: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	createdAt := timestamppb.New(ingredient.CreatedAt)
	var updatedAt sql.NullTime

	if ingredient.UpdatedAt.Valid {
		updatedAt = ingredient.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	ingredientProto.IngredientId = ingredientID.String()
	ingredientProto.Name = ingredient.Name
	ingredientProto.Calories = ingredient.Calories
	ingredientProto.Protein = ingredient.Protein
	ingredientProto.CarbohydratesTotal = ingredient.CarbohydratesTotal
	ingredientProto.FatTotal = ingredient.FatTotal
	ingredientProto.CreatedAt = createdAt
	ingredientProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)
	ingredientProto.UserId = req.UserId
	return ingredientProto, nil
}

func (i *IngredientRepository) UpdateIngredient(ctx context.Context, req *pbml.UpdateIngredientReq) (*pbml.XIngredient, error) {
	query := `UPDATE ingredients SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedIngredient := &pbml.XIngredient{}

	for _, update := range req.Updates {
		switch update.Field {
		case "name":
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Name = update.NewValue
		case "calories":
			calories, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid calories value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("calories = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Calories = calories
		case "serving_size":
			servingSize, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid serving size value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("serving_size = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.ServingSize = servingSize
		case "protein":
			protein, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid protein value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("protein = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Protein = protein
		case "fat_total":
			fatTotal, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid fat total value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("fat_total = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.FatTotal = fatTotal
		case "fat_saturated":
			fatSaturated, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid fat saturated value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("fat_saturated = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.FatSaturated = fatSaturated
		case "carbohydrates_total":
			carbohydratesTotal, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid carbohydrates total value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("carbohydrates_total = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.CarbohydratesTotal = carbohydratesTotal
		case "fiber":
			fiber, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid fiber value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("fiber = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Fiber = fiber
		case "sugar":
			sugar, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid sugar value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("sugar = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Sugar = sugar
		case "sodium":
			sodium, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid sodium value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("sodium = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Sodium = sodium
		case "potassium":
			potassium, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid potassium value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("potassium = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Potassium = potassium
		case "cholesterol":
			cholesterol, err := strconv.ParseFloat(update.NewValue, 64)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid cholesterol value: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("cholesterol = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Cholesterol = cholesterol
		default:
			return nil, fmt.Errorf("unsupported update field: %s", update.Field)
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, req.IngredientId, req.UserId)

	_, err := i.pgpool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	return updatedIngredient, nil
}

func (i *IngredientRepository) DeleteIngredient(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	query := `
		DELETE FROM ingredients
		WHERE id = $1 AND user_id = $2
	`

	_, err := i.pgpool.Exec(ctx, query, req.IngredientId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete exercise: %w", err)
	}
	return &pbml.NilRes{}, nil
}

//func (m *MealRepository) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
//	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
//	if err != nil {
//		return nil, status.Error(codes.Internal, "failed to start transaction")
//	}
//	defer tx.Rollback(ctx)
//
//	// Create meal
//	var mealID uuid.UUID
//	err = tx.QueryRow(ctx, `
//        INSERT INTO meals (user_id, meal_number, meal_description)
//        VALUES ($1, $2, $3)
//        RETURNING id
//    `, req.UserId, req.MealNumber, req.Description).Scan(&mealID)
//
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to create meal: %v", err)
//	}
//
//	// Process ingredients
//	for _, ingredient := range req.Ingredients {
//		var ingredientID uuid.UUID
//
//		// Check if ingredient exists
//		err = tx.QueryRow(ctx, `
//            SELECT id FROM ingredients
//            WHERE name = $1
//        `, ingredient.Name).Scan(&ingredientID)
//
//		if errors.Is(err, pgx.ErrNoRows) {
//			// Create new ingredient
//			err = tx.QueryRow(ctx, `
//                INSERT INTO ingredients (name, calories_per_100g, protein_per_100g, carbs_per_100g, fats_per_100g)
//                VALUES ($1, $2, $3, $4, $5)
//                RETURNING id
//            `, ingredient.Name, ingredient.CaloriesPer100g, ingredient.ProteinPer100g,
//				ingredient.CarbsPer100g, ingredient.FatsPer100g).Scan(&ingredientID)
//
//			if err != nil {
//				return nil, status.Errorf(codes.Internal, "failed to create ingredient: %v", err)
//			}
//		} else if err != nil {
//			return nil, status.Errorf(codes.Internal, "failed to check ingredient: %v", err)
//		}
//
//		// Add ingredient to meal
//		_, err = tx.Exec(ctx, `
//            INSERT INTO meal_ingredients (meal_id, ingredient_id, quantity)
//            VALUES ($1, $2, $3)
//        `, mealID, ingredientID, ingredient.Quantity)
//
//		if err != nil {
//			return nil, status.Errorf(codes.Internal, "failed to add ingredient to meal: %v", err)
//		}
//	}
//
//	if err = tx.Commit(ctx); err != nil {
//		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
//	}
//
//	return &pbml.CreateMealRes{
//		Id: mealID.String(),
//	}, nil
//}
