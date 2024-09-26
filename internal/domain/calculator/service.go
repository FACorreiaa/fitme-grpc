package calculator

import (
	"context"

	pb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"
	"github.com/jackc/pgx/v5"

	"errors"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

type CalculatorService struct {
	repo                                    domain.CalculatorRepository
	pb.UnimplementedCalculatorServiceServer // Required for forward compatibility
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
//
// // GetAllUserMacros implements the GetAllUserMacros gRPC method
//
//	func (s *CalculatorService) GetAllUserMacros(ctx context.Context, req *pb.GetAllUserMacrosRequest) (*pb.GetAllUserMacrosResponse, error) {
//		userMacros, err := s.repo(ctx, req.UserId)
//		if err != nil {
//			if errors.Is(err, db.ErrObjectNotFound{}) {
//				return &pb.GetAllUserMacrosResponse{}, nil
//			}
//			return nil, err
//		}
//
//		response := &pb.GetAllUserMacrosResponse{}
//		for _, macro := range userMacros {
//			response.UserMacros = append(response.UserMacros, &pb.UserMacroDistribution{
//				Id:                              macro.ID,
//				UserId:                          int32(macro.UserID),
//				Age:                             uint32(macro.Age),
//				Height:                          uint32(macro.Height),
//				Weight:                          uint32(macro.Weight),
//				Gender:                          macro.Gender,
//				System:                          macro.System,
//				Activity:                        macro.Activity,
//				ActivityDescription:             macro.ActivityDescription,
//				Objective:                       macro.Objective,
//				ObjectiveDescription:            macro.ObjectiveDescription,
//				CaloriesDistribution:            macro.CaloriesDistribution,
//				CaloriesDistributionDescription: macro.CaloriesDistributionDescription,
//				Protein:                         uint32(macro.Protein),
//				Fats:                            uint32(macro.Fats),
//				Carbs:                           uint32(macro.Carbs),
//				Bmr:                             uint32(macro.BMR),
//				Tdee:                            uint32(macro.TDEE),
//				Goal:                            uint32(macro.Goal),
//				CreatedAt:                       macro.CreatedAt.String(),
//			})
//		}
//
//		return response, nil
//	}
//

// GetUserMacro implements the GetUserMacro gRPC method
func (s *CalculatorService) GetUserMacros(ctx context.Context, req *pb.GetUserMacroRequest) (*pb.GetUserMacroResponse, error) {
	macro, err := s.repo.GetUserMacros(ctx, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	macros := macro.UserMacro

	response := &pb.GetUserMacroResponse{
		UserMacro: &pb.UserMacroDistribution{
			Id:                              macros.Id,
			UserId:                          int32(macros.UserId),
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
