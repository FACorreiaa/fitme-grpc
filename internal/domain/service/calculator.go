package service

//type CalculatorServiceServer struct {
//	repo                                    *Repository
//	pb.UnimplementedCalculatorServiceServer // Required for forward compatibility
//}
//
//func NewCalculatorService(repo *Repository) *CalculatorServiceServer {
//	return &CalculatorServiceServer{
//		repo: repo,
//	}
//}
//
//// CreateUserMacro implements the CreateUserMacro gRPC method
//func (s *CalculatorServiceServer) CreateUserMacro(ctx context.Context, req *pb.CreateUserMacroRequest) (*pb.CreateUserMacroResponse, error) {
//	userMacro := UserMacroDistribution{
//		ID:                              req.UserMacro.Id,
//		UserID:                          int(req.UserMacro.UserId),
//		Age:                             uint8(req.UserMacro.Age),
//		Height:                          uint8(req.UserMacro.Height),
//		Weight:                          uint16(req.UserMacro.Weight),
//		Gender:                          req.UserMacro.Gender,
//		System:                          req.UserMacro.System,
//		Activity:                        req.UserMacro.Activity,
//		ActivityDescription:             req.UserMacro.ActivityDescription,
//		Objective:                       req.UserMacro.Objective,
//		ObjectiveDescription:            req.UserMacro.ObjectiveDescription,
//		CaloriesDistribution:            req.UserMacro.CaloriesDistribution,
//		CaloriesDistributionDescription: req.UserMacro.CaloriesDistributionDescription,
//		Protein:                         uint16(req.UserMacro.Protein),
//		Fats:                            uint16(req.UserMacro.Fats),
//		Carbs:                           uint16(req.UserMacro.Carbs),
//		BMR:                             uint16(req.UserMacro.Bmr),
//		TDEE:                            uint16(req.UserMacro.Tdee),
//		Goal:                            uint16(req.UserMacro.Goal),
//		CreatedAt:                       time.Now(),
//	}
//
//	diet, err := s.repo.InsertDietGoals(userMacro)
//	if err != nil {
//		return nil, err
//	}
//
//	response := &pb.CreateUserMacroResponse{
//		UserMacro: &pb.UserMacroDistribution{
//			Id:                              diet.ID,
//			UserId:                          int32(diet.UserID),
//			Age:                             uint32(diet.Age),
//			Height:                          uint32(diet.Height),
//			Weight:                          uint32(diet.Weight),
//			Gender:                          diet.Gender,
//			System:                          diet.System,
//			Activity:                        diet.Activity,
//			ActivityDescription:             diet.ActivityDescription,
//			Objective:                       diet.Objective,
//			ObjectiveDescription:            diet.ObjectiveDescription,
//			CaloriesDistribution:            diet.CaloriesDistribution,
//			CaloriesDistributionDescription: diet.CaloriesDistributionDescription,
//			Protein:                         uint32(diet.Protein),
//			Fats:                            uint32(diet.Fats),
//			Carbs:                           uint32(diet.Carbs),
//			Bmr:                             uint32(diet.BMR),
//			Tdee:                            uint32(diet.TDEE),
//			Goal:                            uint32(diet.Goal),
//			CreatedAt:                       diet.CreatedAt.String(),
//		},
//	}
//
//	return response, nil
//}
//
//// GetAllUserMacros implements the GetAllUserMacros gRPC method
//func (s *CalculatorServiceServer) GetAllUserMacros(ctx context.Context, req *pb.GetAllUserMacrosRequest) (*pb.GetAllUserMacrosResponse, error) {
//	userMacros, err := s.repo.GetUserDietGoals(ctx, int(req.UserId))
//	if err != nil {
//		if errors.Is(err, db.ErrObjectNotFound{}) {
//			return &pb.GetAllUserMacrosResponse{}, nil
//		}
//		return nil, err
//	}
//
//	response := &pb.GetAllUserMacrosResponse{}
//	for _, macro := range userMacros {
//		response.UserMacros = append(response.UserMacros, &pb.UserMacroDistribution{
//			Id:                              macro.ID,
//			UserId:                          int32(macro.UserID),
//			Age:                             uint32(macro.Age),
//			Height:                          uint32(macro.Height),
//			Weight:                          uint32(macro.Weight),
//			Gender:                          macro.Gender,
//			System:                          macro.System,
//			Activity:                        macro.Activity,
//			ActivityDescription:             macro.ActivityDescription,
//			Objective:                       macro.Objective,
//			ObjectiveDescription:            macro.ObjectiveDescription,
//			CaloriesDistribution:            macro.CaloriesDistribution,
//			CaloriesDistributionDescription: macro.CaloriesDistributionDescription,
//			Protein:                         uint32(macro.Protein),
//			Fats:                            uint32(macro.Fats),
//			Carbs:                           uint32(macro.Carbs),
//			Bmr:                             uint32(macro.BMR),
//			Tdee:                            uint32(macro.TDEE),
//			Goal:                            uint32(macro.Goal),
//			CreatedAt:                       macro.CreatedAt.String(),
//		})
//	}
//
//	return response, nil
//}
//
//// GetUserMacro implements the GetUserMacro gRPC method
//func (s *CalculatorServiceServer) GetUserMacro(ctx context.Context, req *pb.GetUserMacroRequest) (*pb.GetUserMacroResponse, error) {
//	macro, err := s.repo.GetUserDietGoal(ctx, req.PlanId)
//	if err != nil {
//		if errors.Is(err, db.ErrObjectNotFound{}) {
//			return nil, err
//		}
//		return nil, err
//	}
//
//	response := &pb.GetUserMacroResponse{
//		UserMacro: &pb.UserMacroDistribution{
//			Id:                              macro.ID,
//			UserId:                          int32(macro.UserID),
//			Age:                             uint32(macro.Age),
//			Height:                          uint32(macro.Height),
//			Weight:                          uint32(macro.Weight),
//			Gender:                          macro.Gender,
//			System:                          macro.System,
//			Activity:                        macro.Activity,
//			ActivityDescription:             macro.ActivityDescription,
//			Objective:                       macro.Objective,
//			ObjectiveDescription:            macro.ObjectiveDescription,
//			CaloriesDistribution:            macro.CaloriesDistribution,
//			CaloriesDistributionDescription: macro.CaloriesDistributionDescription,
//			Protein:                         uint32(macro.Protein),
//			Fats:                            uint32(macro.Fats),
//			Carbs:                           uint32(macro.Carbs),
//			Bmr:                             uint32(macro.BMR),
//			Tdee:                            uint32(macro.TDEE),
//			Goal:                            uint32(macro.Goal),
//			CreatedAt:                       macro.CreatedAt.String(),
//		},
//	}
//
//	return response, nil
//}
