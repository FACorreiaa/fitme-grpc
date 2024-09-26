package calculator

import (
	"context"
	"fmt"

	pb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"errors"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

type CalculatorService struct {
	pb.UnimplementedCalculatorServer // Required for forward compatibilit
	repo                             domain.CalculatorRepository
}

func NewCalculatorService(repo domain.CalculatorRepository) *CalculatorService {
	return &CalculatorService{
		repo: repo,
	}
}

// CreateUserMacro implements the CreateUserMacro gRPC method
//
//	func (s *CalculatorService) CreateUserMacro(ctx context.Context, req *pb.CreateUserMacroRequest) (*pb.CreateUserMacroResponse, error) {
//		userMacro := UserMacroDistribution{
//			ID:                              req.UserMacro.Id,
//			UserID:                          int(req.UserMacro.UserId),
//			Age:                             uint8(req.UserMacro.Age),
//			Height:                          uint8(req.UserMacro.Height),
//			Weight:                          uint16(req.UserMacro.Weight),
//			Gender:                          req.UserMacro.Gender,
//			System:                          req.UserMacro.System,
//			Activity:                        req.UserMacro.Activity,
//			ActivityDescription:             req.UserMacro.ActivityDescription,
//			Objective:                       req.UserMacro.Objective,
//			ObjectiveDescription:            req.UserMacro.ObjectiveDescription,
//			CaloriesDistribution:            req.UserMacro.CaloriesDistribution,
//			CaloriesDistributionDescription: req.UserMacro.CaloriesDistributionDescription,
//			Protein:                         uint16(req.UserMacro.Protein),
//			Fats:                            uint16(req.UserMacro.Fats),
//			Carbs:                           uint16(req.UserMacro.Carbs),
//			BMR:                             uint16(req.UserMacro.Bmr),
//			TDEE:                            uint16(req.UserMacro.Tdee),
//			Goal:                            uint16(req.UserMacro.Goal),
//			CreatedAt:                       time.Now(),
//		}
//
//		diet, err := s.repo.InsertDietGoals(userMacro)
//		if err != nil {
//			return nil, err
//		}
//
//		response := &pb.CreateUserMacroResponse{
//			UserMacro: &pb.UserMacroDistribution{
//				Id:                              diet.ID,
//				UserId:                          int32(diet.UserID),
//				Age:                             uint32(diet.Age),
//				Height:                          uint32(diet.Height),
//				Weight:                          uint32(diet.Weight),
//				Gender:                          diet.Gender,
//				System:                          diet.System,
//				Activity:                        diet.Activity,
//				ActivityDescription:             diet.ActivityDescription,
//				Objective:                       diet.Objective,
//				ObjectiveDescription:            diet.ObjectiveDescription,
//				CaloriesDistribution:            diet.CaloriesDistribution,
//				CaloriesDistributionDescription: diet.CaloriesDistributionDescription,
//				Protein:                         uint32(diet.Protein),
//				Fats:                            uint32(diet.Fats),
//				Carbs:                           uint32(diet.Carbs),
//				Bmr:                             uint32(diet.BMR),
//				Tdee:                            uint32(diet.TDEE),
//				Goal:                            uint32(diet.Goal),
//				CreatedAt:                       diet.CreatedAt.String(),
//			},
//		}
//
//		return response, nil
//	}

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
	for _, macro := range userMacrosResponse.UserMacros {
		//id, err := uuid.Parse(macro.Id)
		//if err != nil {
		//	return nil, status.Errorf(codes.InvalidArgument,
		//		"invalid UUID format for ID: %v",
		//		err.Error())
		//}
		//
		//userID, err := uuid.Parse(macro.UserId)
		//if err != nil {
		//	return nil, status.Errorf(codes.InvalidArgument,
		//		"invalid UUID format for user ID: %v",
		//		err.Error())
		//}
		response.UserMacros = append(response.UserMacros, &pb.UserMacroDistribution{
			Id:                              macro.Id,
			UserId:                          macro.UserId,
			Age:                             uint32(macro.Age),
			Height:                          uint32(macro.Height),
			Weight:                          uint32(macro.Weight),
			Gender:                          macro.Gender,
			System:                          macro.System,
			Activity:                        macro.Activity,
			ActivityDescription:             macro.ActivityDescription,
			Objective:                       macro.Objective,
			ObjectiveDescription:            macro.ObjectiveDescription,
			CaloriesDistribution:            macro.CaloriesDistribution,
			CaloriesDistributionDescription: macro.CaloriesDistributionDescription,
			Protein:                         uint32(macro.Protein),
			Fats:                            uint32(macro.Fats),
			Carbs:                           uint32(macro.Carbs),
			Bmr:                             uint32(macro.Bmr),
			Tdee:                            uint32(macro.Tdee),
			Goal:                            uint32(macro.Goal),
			CreatedAt:                       macro.CreatedAt,
		})
	}

	return &pb.GetAllUserMacrosResponse{
		UserMacros: response.UserMacros,
	}, nil
}

// GetUserMacros implements the GetUserMacro gRPC method
func (s *CalculatorService) GetUserMacros(ctx context.Context, req *pb.GetUserMacroRequest) (*pb.GetUserMacroResponse, error) {
	macro, err := s.repo.GetUserMacros(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "macro not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user macro: %v", err)
	}

	macros := macro.UserMacro

	response := &pb.GetUserMacroResponse{
		UserMacro: &pb.UserMacroDistribution{
			Id:                              macros.Id,
			UserId:                          macros.UserId,
			Age:                             uint32(macros.Age),
			Height:                          uint32(macros.Height),
			Weight:                          uint32(macros.Weight),
			Gender:                          macros.Gender,
			System:                          macros.System,
			Activity:                        macros.Activity,
			ActivityDescription:             macros.ActivityDescription,
			Objective:                       macros.Objective,
			ObjectiveDescription:            macros.ObjectiveDescription,
			CaloriesDistribution:            macros.CaloriesDistribution,
			CaloriesDistributionDescription: macros.CaloriesDistributionDescription,
			Protein:                         uint32(macros.Protein),
			Fats:                            uint32(macros.Fats),
			Carbs:                           uint32(macros.Carbs),
			Bmr:                             uint32(macros.Bmr),
			Tdee:                            uint32(macros.Tdee),
			Goal:                            uint32(macros.Goal),
			CreatedAt:                       macros.CreatedAt,
		},
	}

	return &pb.GetUserMacroResponse{
		UserMacro: response.UserMacro,
	}, nil
}

func validateUserMacro(macro *pb.UserMacroDistribution) error {
	if macro.Age < minAge || macro.Age > maxAge {
		return errors.New("invalid age")
	}
	// Add other validation checks
	return nil
}
