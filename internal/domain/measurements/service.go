package measurements

import (
	"context"

	pbm "github.com/FACorreiaa/fitme-protos/modules/measurement/generated"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/FACorreiaa/fitme-grpc/internal/domain"
	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware/grpcrequest"
)

type ServiceMeasurement struct {
	pbm.UnimplementedUserMeasurementsServer
	ctx  context.Context
	repo domain.RepositoryMeasurement
}

func NewMeasurementService(ctx context.Context, repo domain.RepositoryMeasurement) *ServiceMeasurement {
	return &ServiceMeasurement{
		ctx:  ctx,
		repo: repo,
	}
}

func (s ServiceMeasurement) CreateWeight(ctx context.Context, req *pbm.CreateWeightReq) (*pbm.CreateWeightRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "CreateWeight")
	defer span.End()

	requestID, ok := ctx.Value(grpcrequest.RequestIDKey{}).(string)
	if !ok {
		return nil, status.Error(codes.Internal, "request id not found in context")
	}

	if req.Request == nil {
		req.Request = &pbm.BaseRequest{}
	}

	req.Request.RequestId = requestID

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "userID is missing in metadata")
	}

	req.UserId = userID

	res, err := s.repo.CreateWeight(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.CreateWeightRes{
		Success: true,
		Message: "Weight inserted correctly",
		Weight:  res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
