package meals

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
	pbml "github.com/FACorreiaa/fitme-protos/modules/meal/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		//Success:  true,
		// Message:  "Ingredients retrieved successfully",
		Ingredient: ingredientProto,
		Response: &pbml.BaseResponse{
			Upstream: "workout-service",
		},
	}, nil

	//return &pbml.GetIngredientsRes{
	//	Ingredients: []*pbml.XIngredient{ingredient},
	//	Response:    {},
	//}, nil
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
		Ingredients: ingredients,
		Response:    &pbml.BaseResponse{Upstream: "workout-service"},
	}, nil

}

func (i *IngredientRepository) CreateIngredient(ctx context.Context, req *pbml.CreateIngredientReq) (*pbml.CreateIngredientRes, error) {
	return nil, nil
}

func (i *IngredientRepository) UpdateIngredient(ctx context.Context, req *pbml.UpdateIngredientReq) (*pbml.UpdateIngredientRes, error) {
	return nil, nil
}

func (i *IngredientRepository) DeleteIngredient(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	return nil, nil
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
