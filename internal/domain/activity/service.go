package activity

import (
	pb "github.com/FACorreiaa/fitme-protos/modules/calculator/generated"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
)

type ActivityService struct {
	pb.UnimplementedCalculatorServer // Required for forward compatibilit
	repo                             domain.ActivityRepository
}

func NewCalculatorService(repo domain.ActivityRepository) *ActivityService {
	return &ActivityService{
		repo: repo,
	}
}
