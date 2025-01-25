package calculator

import (
	"context"
	"fmt"
	"math"
	"time"

	"errors"

	pb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

type CalculatorService struct {
	pb.UnimplementedCalculatorServer
	ctx  context.Context
	repo domain.CalculatorRepository
}

func NewCalculatorService(ctx context.Context, repo domain.CalculatorRepository) *CalculatorService {
	return &CalculatorService{
		ctx:  ctx,
		repo: repo,
	}
}

func mapActivity(ctx context.Context, activity pb.Activity) (*ActivityList, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err() // Return if context is done
	default:
	}

	description, valid := activityDescriptionMap[activity]
	if !valid {
		return nil, errors.New("invalid activity")
	}
	return &ActivityList{
		Activity:    activity.String(),
		Description: description,
	}, nil
}

func mapActivityValues(ctx context.Context, activity pb.Activity) (ActivityValues, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0, ctx.Err() // Return if context is done
	default:
	}
	value, valid := activityValuesMap[activity]
	if !valid {
		return 0, errors.New("invalid activity value")
	}
	return value, nil
}

func mapObjective(ctx context.Context, objective pb.Objective) (*ObjectiveList, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err() // Return if context is done
	default:
	}

	description, valid := objectiveDescriptionMap[objective]
	if !valid {
		return nil, errors.New("invalid objective")
	}

	return &ObjectiveList{
		Objective:   objective.String(),
		Description: description,
	}, nil
}

func mapDistribution(ctx context.Context, distribution pb.CaloriesDistribution) (*CaloriesInfo, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err() // Return if context is done
	default:
	}

	description, valid := carbsDistribution[distribution]
	if !valid {
		return nil, errors.New("invalid distribution")
	}
	return &CaloriesInfo{
		CaloriesDistribution:            distribution,
		CaloriesDistributionDescription: description,
	}, nil
}

func ValidateValues(value, minValue, maxValue uint16, fieldName string) (uint16, error) {
	if value <= minValue || value >= maxValue {
		return 0, fmt.Errorf("invalid %s: %d (must be between %d and %d)", fieldName, value, minValue, maxValue)
	}
	return value, nil
}

func ValidateWeight(value, minValue, maxValue float64, fieldName string) (float64, error) {
	if value <= minValue || value >= maxValue {
		return 0, fmt.Errorf("invalid %s: %f (must be between %f and %f)", fieldName, value, minValue, maxValue)
	}
	return value, nil
}

//func validateWeight(weight uint16) (uint16, error) {
//	if weight <= minWeight || weight > maxWeight {
//		return 0, errors.New("invalid weight")
//	}
//
//	return weight, nil
//}
//
//func validateHeight(height uint8) (uint8, error) {
//	if height <= minHeight || height > maxHeight {
//		return 0, errors.New("invalid height")
//	}
//
//	return height, nil
//}

func convertWeight(weight float64, system System) float64 {
	if system == Metric {
		return weight
	}
	return float64(weight) / 0.453592 // 1 lb = 0.453592 kg
}
func convertHeight(height uint16, system System) float64 {
	if system == Metric {
		return float64(height)
	}
	return float64(height) / 2.54 // 1 in = 2.54 cm
}
func CalculateBMR(userData UserData, system System) (float64, error) {
	if userData.Weight <= 0 || userData.Height <= 0 || userData.Age <= 0 {
		return 0, errors.New("weight, height, and age must be positive values")
	}

	var ageFactor float64
	weight := convertWeight(userData.Weight, system)
	height := convertHeight(userData.Height, system)
	if userData.Gender == m {
		ageFactor = maleAgeFactor
	} else if userData.Gender == f {
		ageFactor = femaleAgeFactor
	} else {
		return 0, errors.New("gender must be 'male' or 'female'")
	}

	if system == Metric {
		return math.Round((10*weight + 6.25*height - 5.0*(float64(userData.Age))) + ageFactor), nil
	} else {
		return math.Round((4.536*weight + 15.88*height - 5.0*(float64(userData.Age))) + ageFactor), nil
	}
}

func calculateTDEE(bmr float64, activityValue ActivityValues) float64 {
	return math.Round(bmr * float64(activityValue))
}

func calculateGoals(tdee float64) Goals {
	var fatLoss = tdee - caloricDeficit
	var bulk = tdee + caloricPlus
	return Goals{
		Cutting:     uint16(fatLoss),
		Maintenance: uint16(tdee),
		Bulking:     uint16(bulk),
	}
}

func getGoal(tdeeValue float64, objective Objective) uint16 {
	goals := calculateGoals(tdeeValue)
	mapGoals := make(map[Objective]uint16)
	mapGoals[maintenance] = goals.Maintenance
	mapGoals[cutting] = goals.Cutting
	mapGoals[bulking] = goals.Bulking
	return mapGoals[objective]
}

func calculateMacroNutrients(calorieGoal float64, distribution CaloriesDistribution) Macros {
	if ratios, ok := macroRatios[distribution]; ok {
		protein := calculateMacroDistribution(ratios.ProteinRatio, calorieGoal, proteinGramValue)
		fats := calculateMacroDistribution(ratios.FatRatio, calorieGoal, fatGramValue)
		carbs := calculateMacroDistribution(ratios.CarbRatio, calorieGoal, carbGramValue)

		return Macros{
			Protein: uint16(protein),
			Fats:    uint16(fats),
			Carbs:   uint16(carbs),
		}
	}

	return Macros{}
}

func calculateMacroDistribution(calorieFactor float64, calorieGoal float64, caloriesPerGram int) float64 {
	return math.Round((calorieFactor * calorieGoal) / float64(caloriesPerGram))
}

func calculateUserPersonalMacros(ctx context.Context, params UserParams) (UserInfo, error) {
	userData, err := validateUserInput(ctx, params)
	if err != nil {
		return UserInfo{}, err
	}

	aEnum, err := StringToActivityEnum(params.Activity) // Convert to pb.Activity
	if err != nil {
		return UserInfo{}, fmt.Errorf("invalid activity: %s", params.Activity)
	}

	oEnum, err := StringToObjectiveEnum(params.Objective)
	if err != nil {
		return UserInfo{}, fmt.Errorf("invalid objective: %s", params.Objective)
	}

	bmr, err := CalculateBMR(userData, System(params.System))
	if err != nil {
		return UserInfo{}, err
	}
	a, err := mapActivity(ctx, aEnum)
	if err != nil {
		return UserInfo{}, err
	}

	o, err := mapObjective(ctx, oEnum)
	if err != nil {
		return UserInfo{}, err
	}

	v, err := mapActivityValues(ctx, aEnum)
	if err != nil {
		return UserInfo{}, err
	}

	d, err := mapDistribution(ctx, parseCaloriesDistribution(params.CaloriesDist))
	if err != nil {
		return UserInfo{}, err
	}

	tdee := calculateTDEE(bmr, v)
	goal := getGoal(tdee, Objective(params.Objective))

	macros := calculateMacroNutrients(tdee, CaloriesDistribution(params.CaloriesDist))
	return UserInfo{
		System: params.System,
		UserData: UserData{
			Age:    userData.Age,
			Height: userData.Height,
			Weight: userData.Weight,
			Gender: userData.Gender,
		},
		ActivityInfo: ActivityInfo{
			Activity:    a.Activity,
			Description: a.Description,
		},
		ObjectiveInfo: ObjectiveInfo{
			Objective:   o.Objective,
			Description: o.Description,
		},
		BMR:  uint16(bmr),
		TDEE: uint16(tdee),
		MacrosInfo: MacrosInfo{
			CaloriesInfo: CaloriesInfo{
				CaloriesDistribution:            d.CaloriesDistribution,
				CaloriesDistributionDescription: d.CaloriesDistributionDescription,
			},
			Macros: macros,
		},
		Goal: goal,
	}, nil
}

// CreateUserMacro implements the CreateUserMacro gRPC method
func validateUserInput(ctx context.Context, params UserParams) (UserData, error) {
	select {
	case <-ctx.Done():
		return UserData{}, ctx.Err() // Return if context is done
	default:
		validAge, err := ValidateValues(params.Age, minAge, maxAge, "age")
		if err != nil {
			return UserData{}, err
		}
		validHeight, err := ValidateValues(params.Height, minHeight, maxHeight, "height")
		if err != nil {
			return UserData{}, err
		}
		validWeight, err := ValidateWeight(float64(params.Weight), minHeight, maxHeight, "weight")
		if err != nil {
			return UserData{}, err
		}
		userInputData := UserData{
			Age:    validAge,
			Height: validHeight,
			Weight: validWeight,
			Gender: params.Gender,
		}
		return userInputData, nil
	}
}

func (s *CalculatorService) CreateUserMacro(ctx context.Context, req *pb.CreateUserMacroRequest) (*pb.CreateUserMacroResponse, error) {
	// Extracting request data
	tracer := otel.Tracer("FitSphere")

	ctx, span := tracer.Start(ctx, "CreateUserMacro")
	defer span.End()

	if req.UserMacro == nil {
		return nil, status.Error(codes.InvalidArgument, "user macro cannot be nil")
	}

	params := UserParams{
		Age:      uint16(req.UserMacro.Age),
		Height:   uint16(req.UserMacro.Height),
		Weight:   uint16(req.UserMacro.Weight),
		Gender:   req.UserMacro.Gender.String(),
		System:   req.UserMacro.System.String(),
		Activity: req.UserMacro.Activity.String(),
		//ActivityDesc:     req.UserMacro.ActivityDescription,
		Objective: req.UserMacro.Objective.String(),
		//ObjectiveDesc:    req.UserMacro.ObjectiveDescription,
		CaloriesDist: string(req.UserMacro.CaloriesDistribution),
		//CaloriesDistDesc: req.UserMacro.CaloriesDistributionDescription,
	}

	// Perform the offline calculations
	userInfo, err := calculateUserPersonalMacros(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("calculate user macro info: %w", err)
	}
	createdAt := timestamppb.New(time.Now())
	system, err := StringToSystemEnum(userInfo.System)
	if err != nil {
		return nil, fmt.Errorf("parse system enum: %w", err)
	}

	objective, err := StringToObjectiveEnum(userInfo.ObjectiveInfo.Objective)
	if err != nil {
		return nil, fmt.Errorf("parse objective enum: %w", err)
	}

	activity, err := StringToActivityEnum(userInfo.ActivityInfo.Activity)
	if err != nil {
		return nil, fmt.Errorf("parse activity enum: %w", err)
	}
	macroDistribution := &pb.UserMacroDistribution{
		Id:                              req.UserMacro.Id,
		UserId:                          req.UserMacro.UserId,
		Age:                             uint32(userInfo.UserData.Age),
		Height:                          uint32(userInfo.UserData.Height),
		Weight:                          userInfo.UserData.Weight,
		Gender:                          req.UserMacro.Gender,
		System:                          system,
		Activity:                        activity,
		ActivityDescription:             string(userInfo.ActivityInfo.Description),
		Objective:                       objective,
		ObjectiveDescription:            string(userInfo.ObjectiveInfo.Description),
		CaloriesDistribution:            userInfo.MacrosInfo.CaloriesInfo.CaloriesDistribution,
		CaloriesDistributionDescription: string(userInfo.MacrosInfo.CaloriesInfo.CaloriesDistributionDescription),
		Protein:                         uint32(userInfo.MacrosInfo.Macros.Protein),
		Fats:                            uint32(userInfo.MacrosInfo.Macros.Fats),
		Carbs:                           uint32(userInfo.MacrosInfo.Macros.Carbs),
		Bmr:                             uint32(userInfo.BMR),
		Tdee:                            uint32(userInfo.TDEE),
		Goal:                            uint32(userInfo.Goal),
		CreatedAt:                       createdAt,
	}

	req = &pb.CreateUserMacroRequest{
		UserMacro: macroDistribution,
	}

	savedMacro, err := s.repo.CreateUserMacro(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", fmt.Sprintf("%T", err)))
		return nil, fmt.Errorf("failed to insert diet goals: %w", err)
	}

	response := &pb.CreateUserMacroResponse{
		UserMacro: savedMacro,
	}

	span.SetAttributes(
		attribute.String("request.id", req.UserMacro.Id),
		attribute.String("request.details", req.UserMacro.Id),
	)

	return response, nil
}

// GetUsersMacros implements the GetAllUserMacros gRPC method
func (s *CalculatorService) GetUsersMacros(ctx context.Context, req *pb.GetAllUserMacrosRequest) (*pb.GetAllUserMacrosResponse, error) {
	userMacrosResponse, err := s.repo.GetUsersMacros(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pb.GetAllUserMacrosResponse{}, fmt.Errorf("Error ")
		}
		return nil, err
	}

	response := &pb.GetAllUserMacrosResponse{}
	createdAt := timestamppb.New(time.Now())
	for _, macro := range userMacrosResponse.UserMacros {
		response.UserMacros = append(response.UserMacros, &pb.UserMacroDistribution{
			Id:                              macro.Id,
			UserId:                          macro.UserId,
			Age:                             macro.Age,
			Height:                          macro.Height,
			Weight:                          macro.Weight,
			Gender:                          macro.Gender,
			System:                          macro.System,
			Activity:                        macro.Activity,
			ActivityDescription:             macro.ActivityDescription,
			Objective:                       macro.Objective,
			ObjectiveDescription:            macro.ObjectiveDescription,
			CaloriesDistribution:            macro.CaloriesDistribution,
			CaloriesDistributionDescription: macro.CaloriesDistributionDescription,
			Protein:                         macro.Protein,
			Fats:                            macro.Fats,
			Carbs:                           macro.Carbs,
			Bmr:                             macro.Bmr,
			Tdee:                            macro.Tdee,
			Goal:                            macro.Goal,
			CreatedAt:                       createdAt,
		})
	}

	return &pb.GetAllUserMacrosResponse{
		UserMacros: response.UserMacros,
	}, nil
}

// GetUserMacros TODO FIX createdAT
// GetUserMacros implements the GetUserMacro gRPC method
func (s *CalculatorService) GetUserMacros(ctx context.Context, req *pb.GetUserMacroRequest) (*pb.GetUserMacroResponse, error) {
	macro, err := s.repo.GetUserMacros(ctx, req)
	createdAt := timestamppb.New(time.Now())

	macro.UserMacro.CreatedAt = createdAt
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user macro: %v", err)
	}

	return macro, nil
}

func (s *CalculatorService) CreateOfflineUserMacro(ctx context.Context, req *pb.CreateOfflineUserMacroRequest) (*pb.CreateOfflineUserMacroResponse, error) {
	// Extracting request data
	if req.UserMacro == nil {
		return nil, status.Error(codes.InvalidArgument, "user macro cannot be nil")
	}

	params := UserParams{
		Age:          uint16(req.UserMacro.Age),
		Height:       uint16(req.UserMacro.Height),
		Weight:       uint16(req.UserMacro.Weight),
		Gender:       req.UserMacro.Gender,
		System:       req.UserMacro.System,
		Activity:     req.UserMacro.Activity,
		Objective:    req.UserMacro.Objective,
		CaloriesDist: string(req.UserMacro.CaloriesDistribution),
	}

	// Perform the offline calculations
	userInfo, err := calculateUserPersonalMacros(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate user info: %v", err)
	}

	// Creating response
	response := &pb.CreateOfflineUserMacroResponse{
		UserMacro: &pb.OfflineUserMacroDistribution{
			Age:                             uint32(userInfo.UserData.Age),
			Height:                          uint32(userInfo.UserData.Height),
			Weight:                          uint32(userInfo.UserData.Weight),
			Gender:                          req.UserMacro.Gender,
			System:                          userInfo.System,
			Activity:                        userInfo.ActivityInfo.Activity,
			ActivityDescription:             string(userInfo.ActivityInfo.Description),
			Objective:                       userInfo.ObjectiveInfo.Objective,
			ObjectiveDescription:            string(userInfo.ObjectiveInfo.Description),
			CaloriesDistribution:            userInfo.MacrosInfo.CaloriesInfo.CaloriesDistribution,
			CaloriesDistributionDescription: string(userInfo.MacrosInfo.CaloriesInfo.CaloriesDistributionDescription),
			Protein:                         uint32(userInfo.MacrosInfo.Macros.Protein),
			Fats:                            uint32(userInfo.MacrosInfo.Macros.Fats),
			Carbs:                           uint32(userInfo.MacrosInfo.Macros.Carbs),
			Bmr:                             uint32(userInfo.BMR),
			Tdee:                            uint32(userInfo.TDEE),
			Goal:                            uint32(userInfo.Goal),
			CreatedAt:                       time.Now().Format(time.RFC3339), // Timestamp for creation
		},
	}

	return response, nil
}

func (s *CalculatorService) DeleteUserMacro(ctx context.Context, req *pb.DeleteUserMacroRequest) (*pb.DeleteUserMacroResponse, error) {
	_, err := s.repo.DeleteUserMacro(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user macro: %v", err)
	}

	return &pb.DeleteUserMacroResponse{}, nil // Return an empty response on successful deletion
}

func StringToSystemEnum(s string) (pb.System, error) {
	switch s {
	case "METRIC":
		return pb.System_METRIC, nil
	case "IMPERIAL":
		return pb.System_IMPERIAL, nil
	default:
		return pb.System_SYSTEM_UNSPECIFIED, fmt.Errorf("invalid System: %s", s)
	}
}

func StringToObjectiveEnum(s string) (pb.Objective, error) {
	switch s {
	case "MAINTENANCE":
		return pb.Objective_MAINTENANCE, nil

	case "BULKING":
		return pb.Objective_BULKING, nil

	case "CUTTING":
		return pb.Objective_CUTTING, nil
	default:
		return pb.Objective_OBJECTIVE_UNSPECIFIED, fmt.Errorf("invalid Objective: %s", s)
	}
}

func StringToActivityEnum(s string) (pb.Activity, error) {
	switch s {
	case "ACTIVITY_UNSPECIFIED":
		return pb.Activity_ACTIVITY_UNSPECIFIED, nil

	case "SEDENTARY":
		return pb.Activity_SEDENTARY, nil
	case "LIGHT":
		return pb.Activity_LIGHT, nil
	case "MODERATE":
		return pb.Activity_MODERATE, nil
	case "HEAVY":
		return pb.Activity_HEAVY, nil
	case "EXTRA_HEAVY":
		return pb.Activity_EXTRA_HEAVY, nil
	default:
		return pb.Activity_ACTIVITY_UNSPECIFIED, fmt.Errorf("invalid Activity: %s", s)
	}
}

func (s *CalculatorService) SetActiveUserMacro(ctx context.Context, req *pb.SetActiveUserMacroRequest) (*pb.SetActiveUserMacroResponse, error) {
	return nil, nil
}

func parseCaloriesDistribution(s string) pb.CaloriesDistribution {
	switch s {
	case "CD_UNSPECIFIED":
		return pb.CaloriesDistribution_CD_UNSPECIFIED
	case "HIGH_CALORIE":
		return pb.CaloriesDistribution_HIGH_CALORIE
	case "MODERATE_CALORIE":
		return pb.CaloriesDistribution_MODERATE_CALORIE
	case "LOW_CALORIE":
		return pb.CaloriesDistribution_LOW_CALORIE
	default:
		// Return a default if the string is unknown
		return pb.CaloriesDistribution_CD_UNSPECIFIED
	}
}
