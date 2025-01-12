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

func (m *MealRepository) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.XMeal, error) {
	tx, err := m.pgpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

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

func (m *MealRepository) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.XMeal, error) {
	mealProto := &pbml.XMeal{}
	id := req.MealId
	meal := &Meal{}
	var nutritionData map[string]float64

	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "meal ID cannot be empty")
	}

	//query := `SELECT * FROM meals WHERE id = $1`
	query := `
		SELECT 
			m.id,
			m.user_id,
			m.meal_number,
			m.meal_description,
			(m.total_macros->>'calories')::DOUBLE PRECISION as total_calories,
			(m.total_macros->>'protein')::DOUBLE PRECISION as total_protein,
			(m.total_macros->>'carbohydrates_total')::DOUBLE PRECISION as total_carbs,
			(m.total_macros->>'fat_total')::DOUBLE PRECISION as total_fat,
			(m.total_macros->>'fiber')::DOUBLE PRECISION as total_fiber,
			(m.total_macros->>'sugar')::DOUBLE PRECISION as total_sugar,
			(m.total_macros->>'sodium')::DOUBLE PRECISION as total_sodium,
			(m.total_macros->>'potassium')::DOUBLE PRECISION as total_potassium,
			(m.total_macros->>'cholesterol')::DOUBLE PRECISION as total_cholesterol,
			(m.total_macros->>'fat_saturated')::DOUBLE PRECISION as total_fat_saturated,
			
			m.created_at,
			m.updated_at
		FROM meals m
		LEFT JOIN meal_ingredients mi ON m.id = mi.meal_id
		WHERE m.user_id = $1 AND m.id = $2
		GROUP BY m.id, m.user_id, m.meal_number, m.meal_description, m.total_macros
	`
	// ARRAY_AGG(
	// 	jsonb_build_object(
	// 		'ingredient_id', mi.ingredient_id,
	// 		'quantity', mi.quantity,
	// 		'calories', mi.calories
	// 	)
	// ) as ingredients,

	calories := nutritionData["calories"]
	protein := nutritionData["protein"]
	carbohydratesTotal := nutritionData["carbohydrates_total"]
	fatTotal := nutritionData["fat_total"]
	fatSaturated := nutritionData["fat_saturated"]
	fiber := nutritionData["fiber"]
	sugar := nutritionData["sugar"]
	sodium := nutritionData["sodium"]
	potassium := nutritionData["potassium"]
	cholesterol := nutritionData["cholesterol"]

	if err := m.pgpool.QueryRow(ctx, query, req.UserId, id).Scan(&meal.ID,
		&meal.UserID,
		&meal.MealNumber,
		&meal.MealDescription,
		&calories,
		&protein,
		&carbohydratesTotal,
		&fatTotal,
		&fiber,
		&sugar,
		&sodium,
		&potassium,
		&cholesterol,
		&fatSaturated,
		//&nutritionData,
		&meal.CreatedAt,
		&meal.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "meal with id %s not found", id)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch meal: %v", err)
	}

	createdAt := timestamppb.New(meal.CreatedAt)
	var updatedAt sql.NullTime
	if meal.UpdatedAt.Valid {
		updatedAt = meal.UpdatedAt
	} else {
		updatedAt = sql.NullTime{Valid: false}
	}

	mealProto.MealId = id
	mealProto.UserId = meal.UserID.String()
	mealProto.MealNumber = int32(meal.MealNumber)
	mealProto.MealDescription = meal.MealDescription
	mealProto.CreatedAt = createdAt
	mealProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)
	for m := range mealProto.MealIngredients {
		m.MealId = id
		m.Calories = float32(calories)
		m.protein = float32(protein)
		m.carbohydratesTotal = float32(carbohydratesTotal)
		m.fatTotal = float32(fatTotal)
		m.fatSaturated = float32(fatSaturated)
		m.fiber = float32(fiber)
		m.sugar = float32(sugar)
		m.sodium = float32(sodium)
		m.potassium = float32(potassium)
		m.cholesterol = float32(cholesterol)
		mealProto.MealIngredients = append(mealProto.MealIngredients, m)

	}
	return mealProto, nil
}

func (m *MealRepository) GetMeals(ctx context.Context, req *pbml.GetMealsReq) ([]*pbml.XMeal, error) {
	mealsProto := make([]*pbml.XMeal, 0)
	query := `
		SELECT * FROM meals WHERE user_id = $1`

	rows, err := m.pgpool.Query(ctx, query, req.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no meals found: %w", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch meals: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		mealProto := pbml.XMeal{}
		meal := &Meal{}
		var nutritionData map[string]float64

		if err := rows.Scan(&meal.ID, &meal.UserID, &meal.MealNumber, &meal.MealDescription, &meal.CreatedAt, &meal.UpdatedAt, &nutritionData); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("no meals found: %w", err)
			}
			return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}

		createdAt := timestamppb.New(meal.CreatedAt)
		var updatedAt sql.NullTime
		if meal.UpdatedAt.Valid {
			updatedAt = meal.UpdatedAt
		} else {
			updatedAt = sql.NullTime{Valid: false}
		}

		mealProto.MealId = meal.ID.String()
		mealProto.UserId = meal.UserID.String()
		mealProto.MealNumber = int32(meal.MealNumber)
		mealProto.MealDescription = meal.MealDescription
		mealProto.CreatedAt = createdAt
		mealProto.UpdatedAt = nullTimeToTimestamppb(updatedAt)

		mealsProto = append(mealsProto, &mealProto)
	}

	return mealsProto, nil
}

func (m *MealRepository) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.XMeal, error) {
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

func (m *MealRepository) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
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

func (m *MealRepository) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
	return nil, nil
}

func (m *MealRepository) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	return nil, nil
}

func (m *MealRepository) UpdateIngredientInMeal(ctx context.Context, req *pbml.UpdateIngredientReq) (*pbml.NilRes, error) {
	return nil, nil
}

func (m *MealRepository) GetMealIngredients(ctx context.Context, req *pbml.GetMealIngredientsReq) (*pbml.GetMealIngredientsRes, error) {
	return nil, nil
}

func (m *MealRepository) GetMealIngredient(ctx context.Context, req *pbml.GetMealIngredientReq) (*pbml.GetMealIngredientRes, error) {
	return nil, nil
}
