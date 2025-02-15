package meals

import (
	"context"
	"fmt"
	"strconv"

	pbml "github.com/FACorreiaa/fitme-protos/modules/meal/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcrequest"
)

type MealServices interface {
	GetMealPlanService() MealPlanService
	GetDietPreferenceService() DietPreferenceService
	GetFoodLogService() FoodLogService
	GetIngredientService() IngredientService
	GetTrackMealProgressService() TrackMealProgressService
	GetGoalRecommendationService() GoalRecommendationService
	GetMealReminderService() MealReminderService
}

type MealPlanService struct {
	pbml.UnimplementedMealPlanServer
	repo domain.MealPlanRepository
	db   *pgxpool.Pool
}

//func (m MealPlanService) mustEmbedUnimplementedMealServer() {
//	//TODO implement me
//	panic("implement me")
//}

type DietPreferenceService struct {
	pbml.UnimplementedDietPreferenceServiceServer
	repo domain.DietPreferenceRepository
}

func (d DietPreferenceService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.AddIngredientRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DietPreferenceService) mustEmbedUnimplementedMealServer() {
	//TODO implement me
	panic("implement me")
}

type FoodLogService struct {
	pbml.UnimplementedFoodLogServiceServer
	repo domain.FoodLogRepository
}

func (f FoodLogService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (f FoodLogService) mustEmbedUnimplementedMealServer() {
	//TODO implement me
	panic("implement me")
}

type IngredientService struct {
	pbml.UnimplementedIngredientsServer
	repo domain.IngredientsRepository
}

//
//func (i IngredientService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (i IngredientService) mustEmbedUnimplementedMealServer() {
//	//TODO implement me
//	panic("implement me")
//}

type TrackMealProgressService struct {
	pbml.UnimplementedTrackMealProgressServer
	repo domain.TrackMealProgressRepository
}

func (t TrackMealProgressService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (t TrackMealProgressService) mustEmbedUnimplementedMealServer() {
	//TODO implement me
	panic("implement me")
}

type GoalRecommendationService struct {
	pbml.UnimplementedGoalRecommendationServer
	repo domain.GoalRecommendationRepository
}

func (g GoalRecommendationService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoalRecommendationService) mustEmbedUnimplementedMealServer() {
	//TODO implement me
	panic("implement me")
}

type MealReminderService struct {
	pbml.UnimplementedMealReminderServer
	repo domain.MealReminderRepository
}

func (m MealReminderService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	//TODO implement me
	panic("implement me")
}

func (m MealReminderService) mustEmbedUnimplementedMealServer() {
	//TODO implement me
	panic("implement me")
}

func NewMealPlanService(ctx context.Context, repo domain.MealPlanRepository) *MealPlanService {
	return &MealPlanService{repo: repo}
}

func NewDietPreferenceService(ctx context.Context, repo domain.DietPreferenceRepository) *DietPreferenceService {
	return &DietPreferenceService{repo: repo}
}

func NewFoodLogService(ctx context.Context, repo domain.FoodLogRepository) *FoodLogService {
	return &FoodLogService{repo: repo}
}

func NewIngredientService(ctx context.Context, repo domain.IngredientsRepository) *IngredientService {
	return &IngredientService{repo: repo}
}

func NewTrackMealProgressService(ctx context.Context, repo domain.TrackMealProgressRepository) *TrackMealProgressService {
	return &TrackMealProgressService{repo: repo}
}

func NewGoalRecommendationService(ctx context.Context, repo domain.GoalRecommendationRepository) *GoalRecommendationService {
	return &GoalRecommendationService{repo: repo}
}

func NewMealReminderService(ctx context.Context, repo domain.MealReminderRepository) *MealReminderService {
	return &MealReminderService{repo: repo}
}

type mealServices struct {
	mealPlanService           MealPlanService
	dietPreferenceService     DietPreferenceService
	foodLogService            FoodLogService
	ingredientService         IngredientService
	trackMealProgressService  TrackMealProgressService
	goalRecommendationService GoalRecommendationService
	mealReminderService       MealReminderService
}

func (m mealServices) GetMealPlanService() MealPlanService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetDietPreferenceService() DietPreferenceService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetFoodLogService() FoodLogService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetIngredientService() IngredientService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetTrackMealProgressService() TrackMealProgressService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetGoalRecommendationService() GoalRecommendationService {
	//TODO implement me
	panic("implement me")
}

func (m mealServices) GetMealReminderService() MealReminderService {
	//TODO implement me
	panic("implement me")
}

func NewMealServices(
	mealPlanService MealPlanService,
	dietPreferenceService DietPreferenceService,
	foodLogService FoodLogService,
	ingredientService IngredientService,
	trackMealProgressService TrackMealProgressService,
	goalRecommendationService GoalRecommendationService,
	mealReminderService MealReminderService,
) MealServices {
	return &mealServices{
		mealPlanService:           mealPlanService,
		dietPreferenceService:     dietPreferenceService,
		foodLogService:            foodLogService,
		ingredientService:         ingredientService,
		trackMealProgressService:  trackMealProgressService,
		goalRecommendationService: goalRecommendationService,
		mealReminderService:       mealReminderService,
	}
}

func (m *MealPlanService) GetMealPlanService() MealPlanService {
	return m.GetMealPlanService()
}

func (m *MealPlanService) GetDietPreferenceService() DietPreferenceService {
	return m.GetDietPreferenceService()
}

func (m *MealPlanService) GetFoodLogService() FoodLogService {
	return m.GetFoodLogService()
}

func (m *MealPlanService) GetIngredientService() IngredientService {
	return m.GetIngredientService()
}

func (m *MealPlanService) GetTrackMealProgressService() TrackMealProgressService {
	return m.GetTrackMealProgressService()
}

func (m *MealPlanService) GetGoalRecommendationService() GoalRecommendationService {
	return m.GetGoalRecommendationService()
}

func (m *MealPlanService) GetMealReminderService() MealReminderService {
	return m.GetMealReminderService()
}

func (m *MealPlanService) GetMealReminder() MealReminderService {
	return m.GetMealReminder()
}

func (m *MealPlanService) GetDietPreference() DietPreferenceService {
	return m.GetDietPreference()
}

func (m *MealPlanService) GetFoodLog() FoodLogService {
	return m.GetFoodLog()
}

func (m *MealPlanService) GetIngredient() IngredientService {
	return m.GetIngredient()
}

func (m *MealPlanService) GetTrackMealProgress() TrackMealProgressService {
	return m.GetTrackMealProgress()
}

func (m *MealPlanService) GetGoalRecommendation() GoalRecommendationService {
	return m.GetGoalRecommendation()
}

func (i *IngredientService) GetIngredient(ctx context.Context, req *pbml.GetIngredientReq) (*pbml.GetIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Ingrediet/GetIngredient")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	ingredient, err := i.repo.GetIngredient(ctx, req)
	if err != nil {
		return &pbml.GetIngredientRes{
			Success: false,
			Message: "Ingredient fetch failed",
			Response: &pbml.BaseResponse{
				Upstream:  "ingredient-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get ingredient: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return ingredient, nil
}

func (i *IngredientService) GetIngredients(ctx context.Context, req *pbml.GetIngredientsReq) (*pbml.GetIngredientsRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Ingrediet/GetIngredients")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID
	if userID := ctx.Value("userID"); userID != nil {
		req.UserId = userID.(string)
	}

	ingredients, err := i.repo.GetIngredients(ctx, req)
	if err != nil {
		return &pbml.GetIngredientsRes{
			Success: false,
			Message: "Ingredients fetch failed",
			Response: &pbml.BaseResponse{
				Upstream:  "ingredient-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get ingredients: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return ingredients, nil
}

func (i *IngredientService) CreateIngredient(ctx context.Context, req *pbml.CreateIngredientReq) (*pbml.CreateIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Ingrediet/CreateIngredient")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	ingredient, err := i.repo.CreateIngredient(ctx, req)
	if err != nil {
		return &pbml.CreateIngredientRes{
			Success: false,
			Message: "Ingredient creation failed",
			Response: &pbml.BaseResponse{
				Upstream:  "ingredient-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to create ingredient: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.CreateIngredientRes{
		Success:    true,
		Message:    "Ingredient created successfully",
		Ingredient: ingredient,
		Response: &pbml.BaseResponse{
			Upstream:  "ingredient-service",
			RequestId: requestID,
		},
	}, nil
}

func (i *IngredientService) DeleteIngredient(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Ingrediet/DeleteIngredient")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	_, err := i.repo.DeleteIngredient(ctx, req)
	if err != nil {
		return &pbml.NilRes{}, status.Errorf(codes.Internal, "failed to delete ingredient: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.NilRes{}, nil
}

func (i *IngredientService) UpdateIngredient(ctx context.Context, req *pbml.UpdateIngredientReq) (*pbml.UpdateIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Ingrediet/UpdateIngredient")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	ingredient, err := i.repo.UpdateIngredient(ctx, req)
	if err != nil {
		return &pbml.UpdateIngredientRes{
			Success: false,
			Message: "Ingredient update failed",
			Response: &pbml.BaseResponse{
				Upstream:  "ingredient-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to update ingredient: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.UpdateIngredientRes{
		Success:    true,
		Message:    "Ingredient updated successfully",
		Ingredient: ingredient,
		Response: &pbml.BaseResponse{
			Upstream:  "ingredient-service",
			RequestId: requestID,
		},
	}, nil
}

// TODO GetIngredientsByName
// TODO GetMealsByDate

func (m *MealPlanService) CreateMeal(ctx context.Context, req *pbml.CreateMealReq) (*pbml.CreateMealRes, error) {

	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/CreateMeal")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	//tx, err := m.db.BeginTx(ctx, pgx.TxOptions{})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, "failed to start transaction")
	//}
	//defer func() {
	//	if err != nil {
	//		_ = tx.Rollback(ctx)
	//	}
	//}()

	meal, err := m.repo.CreateMeal(ctx, req)
	if err != nil {
		return &pbml.CreateMealRes{
			Success: false,
			Message: "Meal creation failed",
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to create meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.CreateMealRes{
		Success: true,
		Message: "Meal created successfully",
		Meal:    meal,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) GetMeal(ctx context.Context, req *pbml.GetMealReq) (*pbml.GetMealRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMeal")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	meal, err := m.repo.GetMeal(ctx, req)
	if err != nil {
		return &pbml.GetMealRes{
			Success: false,
			Message: "Meal creation failed",
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.GetMealRes{
		Success: true,
		Message: "Meal fetched successfully",
		Meal:    meal,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) GetMeals(ctx context.Context, req *pbml.GetMealsReq) (*pbml.GetMealsRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMeals")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	meals, err := m.repo.GetMeals(ctx, req)
	if err != nil {
		return &pbml.GetMealsRes{
			Success: false,
			Message: "Meal creation failed",
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get meals: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.GetMealsRes{
		Success: true,
		Message: "Meals fetched successfully",
		Meals:   meals,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) DeleteMeal(ctx context.Context, req *pbml.DeleteMealReq) (*pbml.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/DeleteMeal")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	_, err := m.repo.DeleteMeal(ctx, req)
	if err != nil {
		return &pbml.NilRes{}, status.Errorf(codes.Internal, "failed to delete meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.NilRes{}, nil
}

func (m *MealPlanService) UpdateMeal(ctx context.Context, req *pbml.UpdateMealReq) (*pbml.UpdateMealRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/UpdateMeal")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	meal, err := m.repo.UpdateMeal(ctx, req)
	if err != nil {
		return &pbml.UpdateMealRes{
			Success: false,
			Message: "Meal update failed",
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to update meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.UpdateMealRes{
		Success: true,
		Message: "Meal updated successfully",
		Meal:    meal,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) AddIngredientToMeal(ctx context.Context, req *pbml.AddIngredientReq) (*pbml.AddIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/AddIngredientToMeal")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	ingredient, err := m.repo.AddIngredientToMeal(ctx, req)
	if err != nil {
		return &pbml.AddIngredientRes{
			Success: false,
			Message: "Failed to add ingredient to meal",
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to add ingredient to meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbml.AddIngredientRes{
		Success:       true,
		Message:       "Ingredient added to meal successfully",
		NewIngredient: ingredient,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) RemoveIngredientFromMeal(ctx context.Context, req *pbml.DeleteIngredientReq) (*pbml.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/RemoveIngredientFromMeal")
	defer span.End()
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	req.Request.RequestId = requestID
	req.UserId = userID
	_, err := m.repo.RemoveIngredientFromMeal(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove ingredient from meal: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.NilRes{}, nil
}

func (m *MealPlanService) GetMealIngredients(ctx context.Context, req *pbml.GetMealIngredientsReq) (*pbml.GetMealIngredientsRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMealIngredients")
	defer span.End()
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	req.Request.RequestId = requestID
	req.UserId = userID
	ingredients, err := m.repo.GetMealIngredients(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get meal ingredients: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.GetMealIngredientsRes{
		Success:         true,
		Message:         "Meal ingredients fetched successfully",
		MealIngredients: ingredients,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) GetMealIngredient(ctx context.Context, req *pbml.GetMealIngredientReq) (*pbml.GetMealIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMealPlan")
	defer span.End()
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	req.Request.RequestId = requestID
	req.UserId = userID
	mealPlan, err := m.repo.GetMealIngredient(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get meal plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.GetMealIngredientRes{
		Success:         true,
		Message:         "Meal plan fetched successfully",
		MealIngredients: mealPlan,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) UpdateIngredientInMeal(ctx context.Context, req *pbml.UpdateMealIngredientReq) (*pbml.UpdateMealIngredientRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/UpdateIngredientInMeal")
	defer span.End()
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	req.Request.RequestId = requestID
	req.UserId = userID
	ing, err := m.repo.UpdateIngredientInMeal(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update meal plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.UpdateMealIngredientRes{
		Success:        true,
		Message:        "Meal plan updated successfully",
		MealIngredient: ing,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) CreateMealPlan(ctx context.Context, req *pbml.CreateMealPlanReq) (*pbml.CreateMealPlanRes, error) {
	var warningMsg string
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/CreateMealPlan")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	activeMacro, err := m.repo.GetUserCalorieLimit(ctx, req.MealPlan.UserId)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to get active macro: %v", err)
	}

	req.Request.RequestId = requestID
	req.MealPlan.UserId = userID

	mp, err := m.repo.CreateMealPlan(ctx, req)
	if err != nil {
		return &pbml.CreateMealPlanRes{
			Success: false,
			Message: "Failed to add ingredient to meal",
			Status:  strconv.Itoa(int(codes.NotFound)),
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to create new meal plan: %v", err)
	}

	if mp.TotalMealNutrients != nil && mp.TotalMealNutrients.Calories > activeMacro {
		warningMsg = fmt.Sprintf("Warning: The total calories of your meal plan (%.0f) exceed your active calorie objective (%.0f).", mp.TotalMealNutrients.Calories, activeMacro)
		span.SetAttributes(attribute.String("mealplan.warning", warningMsg))
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.CreateMealPlanRes{
		Success:  true,
		Message:  "Meal plan inserted successfully",
		MealPlan: mp,
		Status:   strconv.Itoa(int(codes.OK)),
		Warning:  warningMsg,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) GetMealPlans(ctx context.Context, req *pbml.GetMealPlansReq) (*pbml.GetMealPlansRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMealPlans")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	mps, err := m.repo.GetMealPlans(ctx, req)
	if err != nil {
		return &pbml.GetMealPlansRes{
			Success: false,
			Message: "Failed to fetch meal plans",
			//Status:  strconv.Itoa(int(codes.NotFound)),
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get meal plans: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.GetMealPlansRes{
		Success:                true,
		Message:                "Meal plan fetched successfully",
		MealPlan:               mps.MealPlan,
		CreatedAt:              mps.CreatedAt,
		UpdatedAt:              mps.UpdatedAt,
		TotalMealPlanNutrients: mps.TotalMealPlanNutrients,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) GetMealPlan(ctx context.Context, req *pbml.GetMealPlanReq) (*pbml.GetMealPlanRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/GetMealPlan")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	mps, err := m.repo.GetMealPlan(ctx, req)
	if err != nil {
		return &pbml.GetMealPlanRes{
			Success: false,
			Message: "Failed to fetch meal plan",
			//Status:  strconv.Itoa(int(codes.NotFound)),
			Response: &pbml.BaseResponse{
				Upstream:  "meal-service",
				RequestId: requestID,
			},
		}, status.Errorf(codes.Internal, "failed to get meal plans: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.GetMealPlanRes{
		Success:                true,
		Message:                "Meal plan fetched successfully",
		MealPlan:               mps.MealPlan,
		CreatedAt:              mps.CreatedAt,
		UpdatedAt:              mps.UpdatedAt,
		TotalMealPlanNutrients: mps.TotalMealPlanNutrients,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}

func (m *MealPlanService) DeleteMealPlan(ctx context.Context, req *pbml.DeleteMealPlanReq) (*pbml.NilRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/DeleteMealPlan")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	d, err := m.repo.DeleteMealPlan(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete meal plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return d, nil
}

func (m *MealPlanService) UpdateMealPlan(ctx context.Context, req *pbml.UpdateMealPlanReq) (*pbml.UpdateMealPlanRes, error) {
	tracer := otel.Tracer("FitSphere")
	ctx, span := tracer.Start(ctx, "Meal/UpdateMealPlan")
	defer span.End()
	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in context")
	}

	if req.Request == nil {
		req.Request = &pbml.BaseRequest{}
	}

	req.Request.RequestId = requestID
	req.UserId = userID

	mp, err := m.repo.UpdateMealPlan(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete meal plan: %v", err)
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()))

	return &pbml.UpdateMealPlanRes{
		Success:  true,
		Message:  "meal plan updated successfully",
		MealPlan: mp,
		Response: &pbml.BaseResponse{
			Upstream:  "meal-service",
			RequestId: requestID,
		},
	}, nil
}
