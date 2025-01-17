package internal

import (
	"context"

	"github.com/FACorreiaa/fitme-protos/container"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/activity"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/calculator"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/meals"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/measurements"
	"github.com/FACorreiaa/fitme-grpc/internal/domain/workout"
)

type MealServiceContainer struct {
	MealPlanService           *meals.MealPlanService
	DietPreferenceService     *meals.DietPreferenceService
	FoodLogService            *meals.FoodLogService
	IngredientService         *meals.IngredientService
	TrackMealProgressService  *meals.TrackMealProgressService
	GoalRecommendationService *meals.GoalRecommendationService
	MealReminderService       *meals.MealReminderService
}

type ServiceContainer struct {
	Brokers     *container.Brokers
	AuthService *auth.Service
	//CustomerService    *domain.CustomerService
	CalculatorService  *calculator.CalculatorService
	ServiceActivity    *activity.ServiceActivity
	WorkoutService     *workout.ServiceWorkout
	MeasurementService *measurements.ServiceMeasurement
	MealServices       *MealServiceContainer
}

func NewServiceContainer(ctx context.Context, pgPool *pgxpool.Pool, redisClient *redis.Client, brokers *container.Brokers) *ServiceContainer {
	sessionManager := auth.NewSessionManager(pgPool, redisClient)
	authRepo := auth.NewRepository(pgPool, redisClient, sessionManager)
	calculatorRepo := calculator.NewCalculatorRepository(pgPool, redisClient, sessionManager)
	activityRepo := activity.NewRepositoryActivity(pgPool, redisClient, sessionManager)
	workoutRepo := workout.NewRepositoryWorkout(pgPool, redisClient, sessionManager)
	measurementRepo := measurements.NewRepositoryMeasurement(pgPool, redisClient, sessionManager)
	authService := auth.NewService(ctx, authRepo, pgPool, redisClient, sessionManager)
	//customerService := domain.NewCustomerService(ctx, pgPool, redisClient)
	calculatorService := calculator.NewCalculatorService(ctx, calculatorRepo)
	activityService := activity.NewCalculatorService(ctx, activityRepo)
	workoutService := workout.NewServiceWorkout(ctx, workoutRepo)
	measurementService := measurements.NewMeasurementService(ctx, measurementRepo)

	// meals
	mealPlanRepo := meals.NewMealPlanRepository(pgPool, redisClient, sessionManager)
	dietPreferenceRepo := meals.NewDietPreferenceRepository(pgPool, redisClient, sessionManager)
	foodLogRepo := meals.NewFoodLogRepository(pgPool, redisClient, sessionManager)
	ingredientRepo := meals.NewIngredientRepository(pgPool, redisClient, sessionManager)
	trackMealProgressRepo := meals.NewTrackMealProgressRepository(pgPool, redisClient, sessionManager)
	goalRecommendationRepo := meals.NewGoalRecommendationRepository(pgPool, redisClient, sessionManager)
	mealReminderRepo := meals.NewMealReminderRepository(pgPool, redisClient, sessionManager)

	mealServices := &MealServiceContainer{
		MealPlanService:           meals.NewMealPlanService(ctx, mealPlanRepo),
		DietPreferenceService:     meals.NewDietPreferenceService(ctx, dietPreferenceRepo),
		FoodLogService:            meals.NewFoodLogService(ctx, foodLogRepo),
		IngredientService:         meals.NewIngredientService(ctx, ingredientRepo),
		TrackMealProgressService:  meals.NewTrackMealProgressService(ctx, trackMealProgressRepo),
		GoalRecommendationService: meals.NewGoalRecommendationService(ctx, goalRecommendationRepo),
		MealReminderService:       meals.NewMealReminderService(ctx, mealReminderRepo),
	}

	return &ServiceContainer{
		Brokers:     brokers,
		AuthService: authService,
		//CustomerService:    customerService,
		CalculatorService:  calculatorService,
		ServiceActivity:    activityService,
		WorkoutService:     workoutService,
		MeasurementService: measurementService,
		MealServices:       mealServices,
	}
}
