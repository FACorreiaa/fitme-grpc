package meals

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func nullTimeToTimestamppb(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

func calculateMacro(macro float64, quantity float64) float64 {
	return macro * quantity / 100
}

type MealPlanRepository struct {
	pbml.UnimplementedMealPlanServer
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

func NewMealPlanRepository(db *pgxpool.Pool, redis *redis.Client, sessionManager *auth.SessionManager) *MealPlanRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}
	return &MealPlanRepository{pgpool: db, redis: redis, sessionManager: sessionManager}
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

func (m *MealPlanRepository) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.XMeal, error) {
	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var exists bool
	err = m.pgpool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, req.UserId).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "user_id does not exist")
	}

	query := `
		INSERT INTO meals (user_id, meal_number, meal_description, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var mealID string
	err = tx.QueryRow(ctx, query,
		req.UserId,
		req.Meal.MealNumber,
		req.Meal.MealDescription,
		time.Now()).Scan(&mealID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert meal: %v", err)
	}

	totalMacros := &pbml.XMealIngredient{
		Calories:           0,
		Protein:            0,
		CarbohydratesTotal: 0,
		FatTotal:           0,
		FatSaturated:       0,
		Fiber:              0,
		Sugar:              0,
		Sodium:             0,
		Potassium:          0,
		Cholesterol:        0,
	}

	for _, ingredient := range req.Meal.MealIngredients {
		if len(ingredient.IngredientId) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "ingredient ID cannot be empty")
		}
		for _, ingredientID := range ingredient.IngredientId {
			ingredientUUID, err := uuid.Parse(ingredientID)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient ID: %v", err)
			}

			ingredientRow := tx.QueryRow(ctx, `
			SELECT calories, protein, carbohydrates_total, fat_total, fat_saturated, fiber, sugar, sodium, potassium, cholesterol
			FROM ingredients
			WHERE id = $1
		`, ingredientUUID)

			var calories, protein, carbohydratesTotal, fatTotal, fatSaturated, fiber, sugar, sodium, potassium, cholesterol float64
			if err = ingredientRow.Scan(&calories, &protein, &carbohydratesTotal,
				&fatTotal, &fatSaturated, &fiber,
				&sugar, &sodium, &potassium, &cholesterol); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, status.Errorf(codes.NotFound, "ingredient with id %s not found", ingredient.IngredientId)
				}
				return nil, status.Errorf(codes.Internal, "failed to fetch ingredient macros: %v", err)
			}

			calories = calculateMacro(calories, ingredient.Quantity)
			protein = calculateMacro(protein, ingredient.Quantity)
			carbohydratesTotal = calculateMacro(carbohydratesTotal, ingredient.Quantity)
			fatTotal = calculateMacro(fatTotal, ingredient.Quantity)
			fatSaturated = calculateMacro(fatSaturated, ingredient.Quantity)
			fiber = calculateMacro(fiber, ingredient.Quantity)
			sugar = calculateMacro(sugar, ingredient.Quantity)
			sodium = calculateMacro(sodium, ingredient.Quantity)
			potassium = calculateMacro(potassium, ingredient.Quantity)
			cholesterol = calculateMacro(cholesterol, ingredient.Quantity)

			_, err = tx.Exec(ctx, `
			INSERT INTO meal_ingredients
			(meal_id, ingredient_id, quantity, calories, protein, carbohydrates_total, fat_total, fat_saturated, fiber, sugar, sodium, potassium, cholesterol, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
		`, mealID, ingredientUUID, ingredient.Quantity, calories, protein,
				carbohydratesTotal, fatTotal, fatSaturated, fiber, sugar, sodium, potassium,
				cholesterol)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to insert ingredient: %v", err)
			}

			totalMacros.Calories += calories
			totalMacros.Protein += protein
			totalMacros.CarbohydratesTotal += carbohydratesTotal
			totalMacros.FatTotal += fatTotal
			totalMacros.FatSaturated += fatSaturated
			totalMacros.Fiber += fiber
			totalMacros.Sugar += sugar
			totalMacros.Sodium += sodium
			totalMacros.Potassium += potassium
			totalMacros.Cholesterol += cholesterol
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE meals
		SET total_macros = $1::jsonb
		WHERE id = $2
	`, totalMacros, mealID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update total macros: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// TODO return full macros with quantity
	return &pbml.XMeal{
		MealId:          mealID,
		UserId:          req.UserId,
		MealNumber:      req.Meal.MealNumber,
		MealDescription: req.Meal.MealDescription,
		MealIngredients: req.Meal.MealIngredients,
		CreatedAt:       timestamppb.New(time.Now()),
	}, nil
}

func (m *MealPlanRepository) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.XMeal, error) {
	mealProto := &pbml.XMeal{}
	id := req.MealId

	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "meal ID cannot be empty")
	}

	query := `
		SELECT
			m.id,
			m.user_id,
			m.meal_number,
			m.meal_description,
			(m.total_macros->>'calories')::DOUBLE PRECISION as total_calories,
			(m.total_macros->>'protein')::DOUBLE PRECISION as total_protein,
			(m.total_macros->>'carbohydrates_total')::DOUBLE PRECISION as total_carbohydrates_total,
			(m.total_macros->>'fat_total')::DOUBLE PRECISION as total_fat_total,
			(m.total_macros->>'fat_saturated')::DOUBLE PRECISION as total_fat_saturated,
			(m.total_macros->>'fiber')::DOUBLE PRECISION as total_fiber,
			(m.total_macros->>'sugar')::DOUBLE PRECISION as total_sugar,
			(m.total_macros->>'sodium')::DOUBLE PRECISION as total_sodium,
			(m.total_macros->>'potassium')::DOUBLE PRECISION as total_potassium,
			(m.total_macros->>'cholesterol')::DOUBLE PRECISION as total_cholesterol,
			m.created_at,
			m.updated_at,
			COALESCE(
				jsonb_agg(jsonb_build_object(
					'ingredient_id', mi.ingredient_id,
					'name', i.name,
					'quantity', mi.quantity,
					'calories', mi.calories,
					'protein', mi.protein,
					'carbohydrates_total', mi.carbohydrates_total,
					'fat_total', mi.fat_total,
					'fat_saturated', mi.fat_saturated,
					'fiber', mi.fiber,
					'sugar', mi.sugar,
					'sodium', mi.sodium,
					'potassium', mi.potassium,
					'cholesterol', mi.cholesterol
				)), '[]'::jsonb
			) AS ingredients
		FROM meals m
		LEFT JOIN meal_ingredients mi ON m.id = mi.meal_id
		LEFT JOIN ingredients i ON mi.ingredient_id = i.id
		WHERE m.id = $1
		GROUP BY m.id, m.user_id, m.meal_number, m.meal_description, m.total_macros
	`

	var rawIngredients []byte
	meal := &Meal{
		TotalMacros: &TotalNutrients{},
	}

	if err := m.pgpool.QueryRow(ctx, query, id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.MealNumber,
		&meal.MealDescription,
		&meal.TotalMacros.Calories,
		&meal.TotalMacros.Protein,
		&meal.TotalMacros.CarbohydratesTotal,
		&meal.TotalMacros.FatTotal,
		&meal.TotalMacros.FatSaturated,
		&meal.TotalMacros.Fiber,
		&meal.TotalMacros.Sugar,
		&meal.TotalMacros.Sodium,
		&meal.TotalMacros.Potassium,
		&meal.TotalMacros.Cholesterol,
		&meal.CreatedAt,
		&meal.UpdatedAt,
		&rawIngredients,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "meal with id %s not found", id)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch meal: %v", err)
	}

	mealProto.TotalMealNutrients = &pbml.XTotalMealNutrients{
		Calories:           convertFloat(meal.TotalMacros.Calories),
		Protein:            convertFloat(meal.TotalMacros.Protein),
		CarbohydratesTotal: convertFloat(meal.TotalMacros.CarbohydratesTotal),
		FatTotal:           convertFloat(meal.TotalMacros.FatTotal),
		FatSaturated:       convertFloat(meal.TotalMacros.FatSaturated),
		Fiber:              convertFloat(meal.TotalMacros.Fiber),
		Sugar:              convertFloat(meal.TotalMacros.Sugar),
		Sodium:             convertFloat(meal.TotalMacros.Sodium),
		Potassium:          convertFloat(meal.TotalMacros.Potassium),
		Cholesterol:        convertFloat(meal.TotalMacros.Cholesterol),
	}

	// Parse rawIngredients into Ingredients slice
	var ingredients []map[string]interface{}
	if err := json.Unmarshal(rawIngredients, &ingredients); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse ingredients: %v", err)
	}

	for _, ing := range ingredients {
		ingredient := Ingredient{
			ID:                 uuid.MustParse(ing["ingredient_id"].(string)),
			Name:               ing["name"].(string),
			Calories:           ing["calories"].(float64),
			Protein:            ing["protein"].(float64),
			CarbohydratesTotal: ing["carbohydrates_total"].(float64),
			FatTotal:           ing["fat_total"].(float64),
			FatSaturated:       ing["fat_saturated"].(float64),
			Fiber:              ing["fiber"].(float64),
			Sugar:              ing["sugar"].(float64),
			Sodium:             ing["sodium"].(float64),
			Potassium:          ing["potassium"].(float64),
			Cholesterol:        ing["cholesterol"].(float64),
		}

		meal.Ingredients = append(meal.Ingredients, ingredient)
	}

	// Map Meal to gRPC Response
	mealProto.MealId = meal.ID.String()
	mealProto.UserId = meal.UserID.String()
	mealProto.MealNumber = int32(meal.MealNumber)
	mealProto.MealDescription = meal.MealDescription
	mealProto.CreatedAt = timestamppb.New(meal.CreatedAt)
	mealProto.UpdatedAt = nullTimeToTimestamppb(meal.UpdatedAt)

	for _, ing := range meal.Ingredients {
		mealProto.MealIngredients = append(mealProto.MealIngredients, &pbml.XMealIngredient{
			IngredientId:       []string{ing.ID.String()},
			MealId:             mealProto.MealId,
			Name:               ing.Name,
			Quantity:           0,
			Calories:           ing.Calories,
			Protein:            ing.Protein,
			CarbohydratesTotal: ing.CarbohydratesTotal,
			FatTotal:           ing.FatTotal,
			FatSaturated:       ing.FatSaturated,
			Fiber:              ing.Fiber,
			Sugar:              ing.Sugar,
			Sodium:             ing.Sodium,
			Potassium:          ing.Potassium,
			Cholesterol:        ing.Cholesterol,
		})
	}

	return mealProto, nil
}

// Helper function to calculate nutrient totals
func calculateTotals(ingredients []*pbml.XMealIngredient) *pbml.XTotalMealNutrients {
	totals := &pbml.XTotalMealNutrients{}
	for _, ing := range ingredients {
		totals.Calories += ing.Calories
		totals.Protein += ing.Protein
		totals.CarbohydratesTotal += ing.CarbohydratesTotal
		totals.FatTotal += ing.FatTotal
		totals.FatSaturated += ing.FatSaturated
		totals.Fiber += ing.Fiber
		totals.Sugar += ing.Sugar
		totals.Sodium += ing.Sodium
		totals.Potassium += ing.Potassium
		totals.Cholesterol += ing.Cholesterol
	}
	return totals
}

func convertFloat(n sql.NullFloat64) float64 {
	if n.Valid {
		return n.Float64
	}
	return 0
}

func (m *MealPlanRepository) GetMeals(ctx context.Context, req *pbml.GetMealsReq) ([]*pbml.XMeal, error) {
	// Check if UserID is valid
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	query := `
		SELECT
			m.id,
			m.user_id,
			m.meal_number,
			m.meal_description,
			(m.total_macros->>'calories')::DOUBLE PRECISION as total_calories,
			(m.total_macros->>'protein')::DOUBLE PRECISION as total_protein,
			(m.total_macros->>'carbohydrates_total')::DOUBLE PRECISION as total_carbohydrates_total,
			(m.total_macros->>'fat_total')::DOUBLE PRECISION as total_fat_total,
			(m.total_macros->>'fat_saturated')::DOUBLE PRECISION as total_fat_saturated,
			(m.total_macros->>'fiber')::DOUBLE PRECISION as total_fiber,
			(m.total_macros->>'sugar')::DOUBLE PRECISION as total_sugar,
			(m.total_macros->>'sodium')::DOUBLE PRECISION as total_sodium,
			(m.total_macros->>'potassium')::DOUBLE PRECISION as total_potassium,
			(m.total_macros->>'cholesterol')::DOUBLE PRECISION as total_cholesterol,
			m.created_at,
			m.updated_at,
			COALESCE(
				jsonb_agg(jsonb_build_object(
					'ingredient_id', mi.ingredient_id::TEXT,
					'name', i.name,
					'quantity', mi.quantity,
					'calories', mi.calories,
					'protein', mi.protein,
					'carbohydrates_total', mi.carbohydrates_total,
					'fat_total', mi.fat_total,
					'fat_saturated', mi.fat_saturated,
					'fiber', mi.fiber,
					'sugar', mi.sugar,
					'sodium', mi.sodium,
					'potassium', mi.potassium,
					'cholesterol', mi.cholesterol
				)), '[]'::jsonb
			) AS ingredients
		FROM meals m
		LEFT JOIN meal_ingredients mi ON m.id = mi.meal_id
		LEFT JOIN ingredients i ON mi.ingredient_id = i.id
		WHERE m.user_id = $1
		GROUP BY m.id, m.user_id, m.meal_number, m.meal_description, m.total_macros
	`

	rows, err := m.pgpool.Query(ctx, query, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch meals: %v", err)
	}
	defer rows.Close()

	var mealsProto []*pbml.XMeal

	for rows.Next() {
		var rawIngredients []byte
		meal := &Meal{
			TotalMacros: &TotalNutrients{},
		}

		if err := rows.Scan(
			&meal.ID,
			&meal.UserID,
			&meal.MealNumber,
			&meal.MealDescription,
			&meal.TotalMacros.Calories,
			&meal.TotalMacros.Protein,
			&meal.TotalMacros.CarbohydratesTotal,
			&meal.TotalMacros.FatTotal,
			&meal.TotalMacros.FatSaturated,
			&meal.TotalMacros.Fiber,
			&meal.TotalMacros.Sugar,
			&meal.TotalMacros.Sodium,
			&meal.TotalMacros.Potassium,
			&meal.TotalMacros.Cholesterol,
			&meal.CreatedAt,
			&meal.UpdatedAt,
			&rawIngredients,
		); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse row: %v", err)
		}

		mealProto := &pbml.XMeal{
			MealId:          meal.ID.String(),
			UserId:          meal.UserID.String(),
			MealNumber:      int32(meal.MealNumber),
			MealDescription: meal.MealDescription,
			CreatedAt:       timestamppb.New(meal.CreatedAt),
			UpdatedAt:       nullTimeToTimestamppb(meal.UpdatedAt),
			TotalMealNutrients: &pbml.XTotalMealNutrients{
				Calories:           convertFloat(meal.TotalMacros.Calories),
				Protein:            convertFloat(meal.TotalMacros.Protein),
				CarbohydratesTotal: convertFloat(meal.TotalMacros.CarbohydratesTotal),
				FatTotal:           convertFloat(meal.TotalMacros.FatTotal),
				FatSaturated:       convertFloat(meal.TotalMacros.FatSaturated),
				Fiber:              convertFloat(meal.TotalMacros.Fiber),
				Sugar:              convertFloat(meal.TotalMacros.Sugar),
				Sodium:             convertFloat(meal.TotalMacros.Sodium),
				Potassium:          convertFloat(meal.TotalMacros.Potassium),
				Cholesterol:        convertFloat(meal.TotalMacros.Cholesterol),
			},
		}

		var ingredients []map[string]interface{}
		if err := json.Unmarshal(rawIngredients, &ingredients); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse ingredients: %v", err)
		}

		for _, ing := range ingredients {
			mealProto.MealIngredients = append(mealProto.MealIngredients, &pbml.XMealIngredient{
				IngredientId:       []string{ing["ingredient_id"].(string)},
				MealId:             mealProto.MealId,
				Name:               ing["name"].(string),
				Quantity:           ing["quantity"].(float64),
				Calories:           ing["calories"].(float64),
				Protein:            ing["protein"].(float64),
				CarbohydratesTotal: ing["carbohydrates_total"].(float64),
				FatTotal:           ing["fat_total"].(float64),
				FatSaturated:       ing["fat_saturated"].(float64),
				Fiber:              ing["fiber"].(float64),
				Sugar:              ing["sugar"].(float64),
				Sodium:             ing["sodium"].(float64),
				Potassium:          ing["potassium"].(float64),
				Cholesterol:        ing["cholesterol"].(float64),
			})
		}

		mealsProto = append(mealsProto, mealProto)
	}

	return mealsProto, nil
}

func (m *MealPlanRepository) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.XMeal, error) {
	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE meals SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1

	updatedMeal := &pbml.XMeal{}

	// Add updated_at timestamp
	setClauses = append(setClauses, "updated_at = NOW()")

	for _, update := range req.Updates {
		switch update.Field {
		case "meal_number":
			newValue, err := strconv.ParseInt(update.NewValue, 10, 32)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid meal number: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("meal_number = $%d", argIndex))
			args = append(args, int32(newValue))
			updatedMeal.MealNumber = int32(newValue)
			argIndex++
		case "meal_description":
			setClauses = append(setClauses, fmt.Sprintf("meal_description = $%d", argIndex))
			args = append(args, update.NewValue)
			updatedMeal.MealDescription = update.NewValue
			argIndex++
		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid field %s", update.Field)
		}
	}

	if len(setClauses) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no updates provided")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d RETURNING id, meal_number, meal_description, created_at, updated_at, user_id",
		argIndex, argIndex+1)
	args = append(args, req.MealId, req.UserId)

	var createdAt time.Time
	var updatedAt sql.NullTime
	var userID uuid.UUID

	err = tx.QueryRow(ctx, query, args...).Scan(
		&updatedMeal.MealId,
		&updatedMeal.MealNumber,
		&updatedMeal.MealDescription,
		&createdAt,
		&updatedAt,
		&userID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "meal not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update meal: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	updatedMeal.CreatedAt = timestamppb.New(createdAt)
	if updatedAt.Valid {
		updatedMeal.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	updatedMeal.UserId = userID.String()

	return updatedMeal, nil
}

func (m *MealPlanRepository) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	if req.MealId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "meal ID cannot be empty")
	}

	query := `
		DELETE FROM meals
		WHERE id = $1 AND user_id = $2`

	_, err := m.pgpool.Exec(ctx, query, req.MealId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete meal: %w", err)
	}
	return &pbml.NilRes{}, nil
}

func (m *MealPlanRepository) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NewIngredient, error) {
	if req.MealId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "meal plan ID cannot be empty")
	}
	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity must be greater than zero")
	}
	if req.NewIngredient == nil || req.NewIngredient.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "new ingredient details are required")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	newIngredient := req.NewIngredient
	var ingredientID uuid.UUID
	var ingredientMacros struct {
		Calories           float64
		Protein            float64
		CarbohydratesTotal float64
		FatTotal           float64
		FatSaturated       float64
		Fiber              float64
		Sugar              float64
		Sodium             float64
		Potassium          float64
		Cholesterol        float64
	}

	checkQuery := `
        SELECT id, calories, protein, carbohydrates_total, fat_total,
               fat_saturated, fiber, sugar, sodium, potassium, cholesterol
        FROM ingredients
        WHERE name = $1 AND user_id = $2
    `
	err = tx.QueryRow(ctx, checkQuery, newIngredient.Name, req.UserId).Scan(
		&ingredientID,
		&ingredientMacros.Calories,
		&ingredientMacros.Protein,
		&ingredientMacros.CarbohydratesTotal,
		&ingredientMacros.FatTotal,
		&ingredientMacros.FatSaturated,
		&ingredientMacros.Fiber,
		&ingredientMacros.Sugar,
		&ingredientMacros.Sodium,
		&ingredientMacros.Potassium,
		&ingredientMacros.Cholesterol,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		createIngQuery := `
            INSERT INTO ingredients (
                name, calories, protein, carbohydrates_total,
                fat_total, fat_saturated, fiber, sugar,
                sodium, potassium, cholesterol, user_id
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            RETURNING id
        `
		err = tx.QueryRow(ctx, createIngQuery,
			newIngredient.Name,
			newIngredient.Calories,
			newIngredient.Protein,
			newIngredient.CarbohydratesTotal,
			newIngredient.FatTotal,
			newIngredient.FatSaturated,
			newIngredient.Fiber,
			newIngredient.Sugar,
			newIngredient.Sodium,
			newIngredient.Potassium,
			newIngredient.Cholesterol,
			req.UserId,
		).Scan(&ingredientID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create ingredient: %v", err)
		}
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check ingredient existence: %v", err)
	} else {
		newIngredient.Calories = ingredientMacros.Calories
		newIngredient.Protein = ingredientMacros.Protein
		newIngredient.CarbohydratesTotal = ingredientMacros.CarbohydratesTotal
		newIngredient.FatTotal = ingredientMacros.FatTotal
		newIngredient.FatSaturated = ingredientMacros.FatSaturated
		newIngredient.Fiber = ingredientMacros.Fiber
		newIngredient.Sugar = ingredientMacros.Sugar
		newIngredient.Sodium = ingredientMacros.Sodium
		newIngredient.Potassium = ingredientMacros.Potassium
		newIngredient.Cholesterol = ingredientMacros.Cholesterol
	}

	addQuery := `
        INSERT INTO meal_ingredients (
            meal_id, ingredient_id, quantity, calories,
            protein, carbohydrates_total, fat_total,
            fat_saturated, fiber, sugar, sodium,
            potassium, cholesterol
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
	_, err = tx.Exec(ctx, addQuery,
		req.MealId,
		ingredientID,
		req.Quantity,
		calculateMacro(newIngredient.Calories, req.Quantity),
		calculateMacro(newIngredient.Protein, req.Quantity),
		calculateMacro(newIngredient.CarbohydratesTotal, req.Quantity),
		calculateMacro(newIngredient.FatTotal, req.Quantity),
		calculateMacro(newIngredient.FatSaturated, req.Quantity),
		calculateMacro(newIngredient.Fiber, req.Quantity),
		calculateMacro(newIngredient.Sugar, req.Quantity),
		calculateMacro(newIngredient.Sodium, req.Quantity),
		calculateMacro(newIngredient.Potassium, req.Quantity),
		calculateMacro(newIngredient.Cholesterol, req.Quantity),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add ingredient to meal: %v", err)
	}

	updateMealMacros := `
        UPDATE meals
        SET total_macros = (
            SELECT jsonb_build_object(
                'calories', SUM(calories),
                'protein', SUM(protein),
                'carbohydrates_total', SUM(carbohydrates_total),
                'fat_total', SUM(fat_total),
                'fat_saturated', SUM(fat_saturated),
                'fiber', SUM(fiber),
                'sugar', SUM(sugar),
                'sodium', SUM(sodium),
                'potassium', SUM(potassium),
                'cholesterol', SUM(cholesterol)
            )
            FROM meal_ingredients
            WHERE meal_id = $1
        )
        WHERE id = $1
    `
	_, err = tx.Exec(ctx, updateMealMacros, req.MealId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update meal macros: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &pbml.NewIngredient{
		IngredientId:       ingredientID.String(),
		Name:               newIngredient.Name,
		Quantity:           req.Quantity,
		Calories:           calculateMacro(newIngredient.Calories, req.Quantity),
		Protein:            calculateMacro(newIngredient.Protein, req.Quantity),
		CarbohydratesTotal: calculateMacro(newIngredient.CarbohydratesTotal, req.Quantity),
		FatTotal:           calculateMacro(newIngredient.FatTotal, req.Quantity),
		FatSaturated:       calculateMacro(newIngredient.FatSaturated, req.Quantity),
		Fiber:              calculateMacro(newIngredient.Fiber, req.Quantity),
		Sugar:              calculateMacro(newIngredient.Sugar, req.Quantity),
		Sodium:             calculateMacro(newIngredient.Sodium, req.Quantity),
		Potassium:          calculateMacro(newIngredient.Potassium, req.Quantity),
		Cholesterol:        calculateMacro(newIngredient.Cholesterol, req.Quantity),
	}, nil
}

func (m *MealPlanRepository) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	if req.MealPlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "meal plan id is required")
	}

	if req.IngredientId == "" {
		return nil, status.Error(codes.InvalidArgument, "ingredient id is required")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}

	defer tx.Rollback(ctx)
	deleteQuery := `
        DELETE FROM meal_ingredients
        WHERE meal_id = $1 AND ingredient_id = $2
    `
	res, err := tx.Exec(ctx, deleteQuery, req.MealPlanId, req.IngredientId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete ingredient from meal: %v", err)
	}

	if res.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "ingredient not found in meal")
	}

	updateMealMacros := `
        UPDATE meals
        SET total_macros = (
            SELECT jsonb_build_object(
                'calories', COALESCE(SUM(calories), 0),
                'protein', COALESCE(SUM(protein), 0),
                'carbohydrates_total', COALESCE(SUM(carbohydrates_total), 0),
                'fat_total', COALESCE(SUM(fat_total), 0),
                'fat_saturated', COALESCE(SUM(fat_saturated), 0),
                'fiber', COALESCE(SUM(fiber), 0),
                'sugar', COALESCE(SUM(sugar), 0),
                'sodium', COALESCE(SUM(sodium), 0),
                'potassium', COALESCE(SUM(potassium), 0),
                'cholesterol', COALESCE(SUM(cholesterol), 0)
            )
            FROM meal_ingredients
            WHERE meal_id = $1
        )
        WHERE id = $1
    `

	_, err = tx.Exec(ctx, updateMealMacros, req.MealPlanId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update meal macros: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &pbml.NilRes{}, nil

}

func (m *MealPlanRepository) GetMealIngredients(ctx context.Context, req *pbml.GetMealIngredientsReq) ([]*pbml.XMealIngredient, error) {
	if req.MealId == "" {
		return nil, status.Error(codes.InvalidArgument, "meal id is required")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	query := `
		SELECT m.id, i.name, m.quantity, m.calories, m.protein, m.carbohydrates_total,
			   m.fat_total, m.fat_saturated, m.fiber, m.sugar, m.sodium,
			   m.potassium, m.cholesterol
		FROM meal_ingredients m
		JOIN ingredients i ON m.ingredient_id = i.id
		WHERE meal_id = $1
	`

	rows, err := tx.Query(ctx, query, req.MealId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get meal ingredients: %v", err)
	}

	defer rows.Close()

	ingredients := make([]*pbml.XMealIngredient, 0)

	for rows.Next() {
		ingredientProto := &pbml.XMealIngredient{}
		ingredient := &Ingredient{}
		err := rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.ServingSize,
			&ingredient.Calories,
			&ingredient.Protein,
			&ingredient.CarbohydratesTotal,
			&ingredient.FatTotal,
			&ingredient.FatSaturated,
			&ingredient.Fiber,
			&ingredient.Sugar,
			&ingredient.Sodium,
			&ingredient.Potassium,
			&ingredient.Cholesterol,
		)

		// createdAt := timestamppb.New(ingredient.CreatedAt)
		// var updatedAt sql.NullTime

		// if ingredient.UpdatedAt.Valid {
		// 	updatedAt = ingredient.UpdatedAt
		// } else {
		// 	updatedAt = sql.NullTime{Valid: false}
		// }

		ingredientProto.IngredientId = []string{ingredient.ID.String()}
		ingredientProto.Name = ingredient.Name
		ingredientProto.Quantity = ingredient.ServingSize
		ingredientProto.Calories = ingredient.Calories
		ingredientProto.Protein = ingredient.Protein
		ingredientProto.CarbohydratesTotal = ingredient.CarbohydratesTotal
		ingredientProto.FatTotal = ingredient.FatTotal
		ingredientProto.FatSaturated = ingredient.FatSaturated
		ingredientProto.Fiber = ingredient.Fiber
		ingredientProto.Sugar = ingredient.Sugar
		ingredientProto.Sodium = ingredient.Sodium
		ingredientProto.Potassium = ingredient.Potassium
		ingredientProto.Cholesterol = ingredient.Cholesterol

		// TODO add later
		// ingredientProto.CreatedAt = createdAt
		// ingredientProto.UpdatedAt = updatedAt

		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}

		ingredients = append(ingredients, ingredientProto)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to iterate rows: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return ingredients, nil

}

func (m *MealPlanRepository) GetMealIngredient(ctx context.Context, req *pbml.GetMealIngredientReq) (*pbml.XMealIngredient, error) {
	if req.MealId == "" {
		return nil, status.Error(codes.InvalidArgument, "meal id is required")
	}

	if req.IngredientId == "" {
		return nil, status.Error(codes.InvalidArgument, "ingredient id is required")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	query := `
		SELECT m.id, i.name, m.quantity, m.calories, m.protein, m.carbohydrates_total,
			   m.fat_total, m.fat_saturated, m.fiber, m.sugar, m.sodium,
			   m.potassium, m.cholesterol
		FROM meal_ingredients m
		JOIN ingredients i ON m.ingredient_id = i.id
		WHERE meal_id = $1 AND m.ingredient_id = $2
	`

	row := tx.QueryRow(ctx, query, req.MealId, req.IngredientId)

	ingredientProto := &pbml.XMealIngredient{}
	ingredient := &Ingredient{}
	err = row.Scan(
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.ServingSize,
		&ingredient.Calories,
		&ingredient.Protein,
		&ingredient.CarbohydratesTotal,
		&ingredient.FatTotal,
		&ingredient.FatSaturated,
		&ingredient.Fiber,
		&ingredient.Sugar,
		&ingredient.Sodium,
		&ingredient.Potassium,
		&ingredient.Cholesterol,
	)

	// createdAt := timestamppb.New(ingredient.CreatedAt)
	// var updatedAt sql.NullTime

	// if ingredient.UpdatedAt.Valid {
	// 	updatedAt = ingredient.UpdatedAt
	// } else {
	// 	updatedAt = sql.NullTime{Valid: false}
	// }

	ingredientProto.IngredientId = []string{ingredient.ID.String()}
	ingredientProto.Name = ingredient.Name
	ingredientProto.Quantity = ingredient.ServingSize
	ingredientProto.Calories = ingredient.Calories
	ingredientProto.Protein = ingredient.Protein
	ingredientProto.CarbohydratesTotal = ingredient.CarbohydratesTotal
	ingredientProto.FatTotal = ingredient.FatTotal
	ingredientProto.FatSaturated = ingredient.FatSaturated
	ingredientProto.Fiber = ingredient.Fiber
	ingredientProto.Sugar = ingredient.Sugar
	ingredientProto.Sodium = ingredient.Sodium
	ingredientProto.Potassium = ingredient.Potassium
	ingredientProto.Cholesterol = ingredient.Cholesterol

	// TODO add later
	// ingredientProto.CreatedAt = createdAt
	// ingredientProto.UpdatedAt = updatedAt

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return ingredientProto, nil
}

func (m *MealPlanRepository) UpdateIngredientInMeal(ctx context.Context, req *pbml.UpdateMealIngredientReq) (*pbml.XMealIngredient, error) {
	query := `UPDATE meal_ingredients SET `
	var setClauses []string
	var args []interface{}
	argIndex := 1
	updatedIngredient := &pbml.XMealIngredient{}

	if req.MealId == "" {
		return nil, status.Error(codes.InvalidArgument, "meal id is required")
	}

	if req.IngredientId == "" {
		return nil, status.Error(codes.InvalidArgument, "ingredient id is required")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	// query := `
	// 	UPDATE meal_ingredients
	// 	SET quantity = $1, calories = $2, protein = $3, carbohydrates_total = $4, fat_total = $5, fat_saturated = $6, fiber = $7, sugar = $8, sodium = $9, potassium = $10, cholesterol = $11
	// 	WHERE meal_id = $12 AND ingredient_id = $13
	// `

	for _, update := range req.Updates {
		switch update.Field {
		case "name":
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Name = update.NewValue
		case "quantity":
			quantity, _ := parseStringToFloat(update.NewValue, "invalid quantity value")
			setClauses = append(setClauses, fmt.Sprintf("quantity = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Quantity = quantity
		case "calories":
			calories, _ := parseStringToFloat(update.NewValue, "invalid calories value")
			setClauses = append(setClauses, fmt.Sprintf("calories = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Calories = calories
		case "protein":
			protein, _ := parseStringToFloat(update.NewValue, "invalid protein value")
			setClauses = append(setClauses, fmt.Sprintf("protein = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Protein = protein
		case "carbohydrates_total":
			carbohydrates, _ := parseStringToFloat(update.NewValue, "invalid carbohydrates value")
			setClauses = append(setClauses, fmt.Sprintf("carbohydrates_total = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.CarbohydratesTotal = carbohydrates
		case "fat_total":
			fat, _ := parseStringToFloat(update.NewValue, "invalid fat value")
			setClauses = append(setClauses, fmt.Sprintf("fat_total = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.FatTotal = fat
		case "fat_saturated":
			fatSaturated, _ := parseStringToFloat(update.NewValue, "invalid fat saturated value")
			setClauses = append(setClauses, fmt.Sprintf("fat_saturated = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.FatSaturated = fatSaturated
		case "fiber":
			fiber, _ := parseStringToFloat(update.NewValue, "invalid fiber value")
			setClauses = append(setClauses, fmt.Sprintf("fiber = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Fiber = fiber
		case "sugar":
			sugar, _ := parseStringToFloat(update.NewValue, "invalid sugar value")
			setClauses = append(setClauses, fmt.Sprintf("sugar = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Sugar = sugar
		case "sodium":
			sodium, _ := parseStringToFloat(update.NewValue, "invalid sodium value")
			setClauses = append(setClauses, fmt.Sprintf("sodium = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Sodium = sodium
		case "potassium":
			potassium, _ := parseStringToFloat(update.NewValue, "invalid potassium value")
			setClauses = append(setClauses, fmt.Sprintf("potassium = $%d", argIndex))
			args = append(args, update.NewValue)
			argIndex++
			updatedIngredient.Potassium = potassium
		case "cholesterol":
			cholesterol, _ := parseStringToFloat(update.NewValue, "invalid cholesterol value")
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
	query += fmt.Sprintf(" WHERE meal_id = $%d AND ingredient_id = $%d", argIndex, argIndex+1)
	args = append(args, req.MealId, req.IngredientId)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update ingredient in meal: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return updatedIngredient, nil
}

func parseStringToFloat(s, errorMessage string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "%s: %v", errorMessage, err)
	}
	return val, nil
}

func (m *MealPlanRepository) GetMealPlan(ctx context.Context, req *pbml.GetMealPlanReq) (*pbml.XMealPlan, error) {
	return nil, nil
}

func (m *MealPlanRepository) GetMealPlans(ctx context.Context, req *pbml.GetMealPlansReq) (*pbml.GetMealPlansRes, error) {
	// Check if UserID is valid
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	//query := `
	//			SELECT
	//				mp.id AS id,
	//				mp.user_id AS user_id,
	//				mp.name AS name,
	//				mp.description AS description,
	//				mp.notes AS notes,
	//				mp.rating AS rating,
	//				mp.created_at AS created_at,
	//				mp.updated_at AS updated_at,
	//				COALESCE(
	//					jsonb_agg(jsonb_build_object(
	//						'meal_id', m.id::TEXT,
	//						'meal_number', m.meal_number,
	//						'meal_description', m.meal_description,
	//						'total_macros', m.total_macros
	//					)), '[]'::jsonb
	//				) AS meals
	//			FROM meal_plans mp
	//			LEFT JOIN meals m ON mp.id = m.meal_plan_id
	//			WHERE mp.user_id = $1
	//			GROUP BY mp.id, mp.user_id, mp.name, mp.description, mp.notes, mp.rating, mp.created_at, mp.updated_at
	//`

	query := `
				SELECT
					mp.id AS id,
					mp.user_id AS user_id,
					mp.name AS name,
					mp.description AS description,
					mp.notes AS notes,
					mp.rating AS rating,
					mp.created_at AS created_at,
					mp.updated_at AS updated_at,
					COALESCE(
						jsonb_agg(
							CASE
								WHEN m.id IS NOT NULL THEN jsonb_build_object(
									'meal_id', m.id::TEXT,
									'meal_number', m.meal_number,
									'meal_description', m.meal_description,
									'total_macros', m.total_macros
								)
								ELSE NULL
							END
						) FILTER (WHERE m.id IS NOT NULL), '[]'::jsonb
					) AS meals
				FROM meal_plans mp
				LEFT JOIN meals m ON mp.id = m.meal_plan_id
				WHERE mp.user_id = $1
				GROUP BY mp.id, mp.user_id, mp.name, mp.description, mp.notes, mp.rating, mp.created_at, mp.updated_at
`

	rows, err := m.pgpool.Query(ctx, query, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch meal plans: %v", err)
	}
	defer rows.Close()

	mealPlansProto := make([]*pbml.XMealPlan, 0)

	for rows.Next() {
		var rawMeals []byte
		mealPlan := &MealPlan{
			TotalMacros: &TotalNutrients{},
		}

		//createdAt := timestamppb.New(mealPlan.CreatedAt)
		//var updatedAt *timestamppb.Timestamp
		//var name, description, notes sql.NullString
		//var rating sql.NullFloat64
		//if mealPlan.UpdatedAt.Valid {
		//	updatedAt = timestamppb.New(mealPlan.UpdatedAt.Time)
		//} else {
		//	updatedAt = nil
		//}

		// Scan the database row into the struct
		if err = rows.Scan(
			&mealPlan.ID,
			&mealPlan.UserID,
			&mealPlan.Name,
			&mealPlan.Description,
			&mealPlan.Notes,
			&mealPlan.Rating,
			&mealPlan.CreatedAt,
			&mealPlan.UpdatedAt,
			&rawMeals,
		); err != nil {
			log.Printf("Scan Error: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to parse row: %v", err)
		}

		// Parse meals for the meal plan
		var meals []map[string]interface{}
		if err = json.Unmarshal(rawMeals, &meals); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse meals: %v", err)
		}

		fmt.Println("Raw Meals:", string(rawMeals))

		mealProtos := make([]*pbml.XMeal, 0)
		totalNutrients := &pbml.XTotalMealNutrients{}

		for _, meal := range meals {

			// Convert meal data to `XMeal` format
			mealProto := &pbml.XMeal{
				MealId:          safeString(meal["meal_id"].(string)),
				MealNumber:      safeFloat64ToInt32(int32(meal["meal_number"].(float64))),
				MealDescription: safeString(meal["meal_description"].(string)),
			}

			if mealId, ok := meal["meal_id"].(string); ok {
				mealProto.MealId = mealId
			} else {
				mealProto.MealId = "" // Provide a default value
			}

			if mealNumber, ok := meal["meal_number"].(float64); ok {
				mealProto.MealNumber = int32(mealNumber)
			} else {
				mealProto.MealNumber = 0 // Provide a default value
			}

			if mealDescription, ok := meal["meal_description"].(string); ok {
				mealProto.MealDescription = mealDescription
			} else {
				mealProto.MealDescription = "" // Provide a default value
			}

			// Parse and sum up nutrients
			if totalMacros, ok := meal["total_macros"].(map[string]interface{}); ok {
				mealProto.TotalMealNutrients = &pbml.XTotalMealNutrients{
					Calories:           safeFloat64(totalMacros["calories"]),
					Protein:            safeFloat64(totalMacros["protein"]),
					CarbohydratesTotal: safeFloat64(totalMacros["carbohydrates_total"]),
					FatTotal:           safeFloat64(totalMacros["fat_total"]),
					FatSaturated:       safeFloat64(totalMacros["fat_saturated"]),
					Fiber:              safeFloat64(totalMacros["fiber"]),
					Sugar:              safeFloat64(totalMacros["sugar"]),
					Sodium:             safeFloat64(totalMacros["sodium"]),
					Potassium:          safeFloat64(totalMacros["potassium"]),
					Cholesterol:        safeFloat64(totalMacros["cholesterol"]),
				}

				totalNutrients.Calories += mealProto.TotalMealNutrients.Calories
				totalNutrients.Protein += mealProto.TotalMealNutrients.Protein
				totalNutrients.CarbohydratesTotal += mealProto.TotalMealNutrients.CarbohydratesTotal
				totalNutrients.FatTotal += mealProto.TotalMealNutrients.FatTotal
				totalNutrients.FatSaturated += mealProto.TotalMealNutrients.FatSaturated
				totalNutrients.Fiber += mealProto.TotalMealNutrients.Fiber
				totalNutrients.Sugar += mealProto.TotalMealNutrients.Sugar
				totalNutrients.Sodium += mealProto.TotalMealNutrients.Sodium
				totalNutrients.Potassium += mealProto.TotalMealNutrients.Potassium
				totalNutrients.Cholesterol += mealProto.TotalMealNutrients.Cholesterol
			} else {
				mealProto.TotalMealNutrients = &pbml.XTotalMealNutrients{}
			}

			mealProtos = append(mealProtos, mealProto)
		}

		createdAt := timestamppb.New(mealPlan.CreatedAt)
		var updatedAt sql.NullTime

		if mealPlan.UpdatedAt.Valid {
			updatedAt = mealPlan.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		// Build `XMealPlan`
		mealPlanProto := &pbml.XMealPlan{
			MealPlanId:         mealPlan.ID.String(),
			UserId:             mealPlan.UserID.String(),
			Name:               mealPlan.Name.String,
			CreatedAt:          createdAt,
			UpdatedAt:          nullTimeToTimestamppb(updatedAt),
			Meal:               mealProtos,
			TotalMealNutrients: totalNutrients,
		}

		mealPlansProto = append(mealPlansProto, mealPlanProto)
	}

	// Construct final response
	return &pbml.GetMealPlansRes{
		Success:   true,
		Message:   "Meal plans fetched successfully",
		MealPlan:  mealPlansProto,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (m *MealPlanRepository) CreateMealPlan(ctx context.Context, req *pbml.CreateMealPlanReq) (*pbml.XMealPlan, error) {
	mealPlanProto := &pbml.XMealPlan{}
	// totalMealNutrients := &pbml.XTotalMealNutrients{}

	if m.pgpool == nil {
		return nil, status.Errorf(codes.Internal, "pgpool is nil")
	}

	if req == nil || req.UserId == "" || req.Name == "" || req.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: required fields are missing")
	}
	if req.Meal == nil {
		return nil, status.Errorf(codes.InvalidArgument, "meal list cannot be nil")
	}

	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var exists bool
	err = m.pgpool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, req.UserId).Scan(&exists)
	if err != nil || !exists {
		return nil, status.Errorf(codes.InvalidArgument, "user_id does not exist")
	}

	var mealPlanID string
	err = tx.QueryRow(ctx, `
       INSERT INTO meal_plans (user_id, name, description, created_at)
       VALUES ($1, $2, $3, $4)
       RETURNING id
   `, req.UserId, req.Name, req.Description, time.Now()).Scan(&mealPlanID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert meal plan: %v", err)
	}

	for i, meal := range req.Meal {
		mealOrder := i + 1
		if meal == nil || meal.MealId == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid meal at index %d", i)
		}

		// Query or create meal
		var mealID string
		err := m.pgpool.QueryRow(ctx, `
           SELECT id FROM meals WHERE user_id = $1 AND id = $2
       `, req.UserId, meal.MealId).Scan(&mealID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if m.CreateMeal == nil {
					return nil, status.Errorf(codes.Internal, "CreateMeal function is nil")
				}

				mealReq := &pbml.CreateMealReq{
					UserId: req.UserId,
					Meal:   meal,
				}
				createdMeal, err := m.CreateMeal(ctx, mealReq)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "failed to create meal: %v", err)
				}
				mealID = createdMeal.MealId
			} else {
				return nil, status.Errorf(codes.Internal, "failed to query meal: %v", err)
			}
		}

		_, err = tx.Exec(ctx, `
           INSERT INTO meal_plan_meals (meal_plan_id, meal_id, meal_order)
           VALUES ($1, $2, $3)
       `, mealPlanID, mealID, mealOrder)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to associate meal with meal plan: %v", err)
		}
	}

	// Add each macro value to the meal plan's total macros
	//if totalMealNutrients != nil {
	//	totalMealNutrients.Calories += totalMealNutrients.Calories
	//	totalMealNutrients.Protein += totalMealNutrients.Protein
	//	totalMealNutrients.CarbohydratesTotal += totalMealNutrients.CarbohydratesTotal
	//	totalMealNutrients.FatTotal += totalMealNutrients.FatTotal
	//	totalMealNutrients.FatSaturated += totalMealNutrients.FatSaturated
	//	totalMealNutrients.Fiber += totalMealNutrients.Fiber
	//	totalMealNutrients.Sugar += totalMealNutrients.Sugar
	//	totalMealNutrients.Sodium += totalMealNutrients.Sodium
	//	totalMealNutrients.Potassium += totalMealNutrients.Potassium
	//	totalMealNutrients.Cholesterol += totalMealNutrients.Cholesterol
	//}

	mealPlanProto.Meal = req.Meal
	mealPlanProto.MealPlanId = mealPlanID
	mealPlanProto.Description = req.Description
	mealPlanProto.Name = req.Name
	//mealPlanProto.TotalMealNutrients = totalMealNutrients

	if err = tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return mealPlanProto, nil
}

func (m *MealPlanRepository) UpdateMealPlan(ctx context.Context, req *pbml.UpdateMealPlanReq) (*pbml.UpdateMealPlanRes, error) {
	return nil, nil
}

func (m *MealPlanRepository) DeleteMealPlan(ctx context.Context, req *pbml.DeleteMealPlanReq) (*pbml.NilRes, error) {
	return nil, nil
}

func safeString(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	}
	return ""
}

func safeFloat64ToInt32(value interface{}) int32 {
	if v, ok := value.(float64); ok {
		return int32(v)
	}
	return 0
}

func safeFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return 0
}
