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

// Weights

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
func (s ServiceMeasurement) GetWeights(ctx context.Context, req *pbm.GetWeightsReq) (*pbm.GetWeightsRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWeights")
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

	res, err := s.repo.GetWeights(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.GetWeightsRes{
		Success: true,
		Message: "Weight fetched correctly",
		Weight:  res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) GetWeight(ctx context.Context, req *pbm.GetWeightReq) (*pbm.GetWeightRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWeight")
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

	res, err := s.repo.GetWeight(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	// TODO Change proto def
	return &pbm.GetWeightRes{
		//Success: true,
		//Message: "Weight fetched correctly",
		WeightId: res.WeightId,
		UserId:   userID,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) DeleteWeight(ctx context.Context, req *pbm.DeleteWeightReq) (*pbm.NilRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "DeleteWeight")
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

	res, err := s.repo.DeleteWeight(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}
func (s ServiceMeasurement) UpdateWeight(ctx context.Context, req *pbm.UpdateWeightReq) (*pbm.UpdateWeightRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "DeleteWeight")
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

	res, err := s.repo.UpdateWeight(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.UpdateWeightRes{
		Success: true,
		Message: "Weight updated correctly",
		Weight:  res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}

// Water Intake

func (s ServiceMeasurement) CreateWaterMeasurement(ctx context.Context, req *pbm.CreateWaterIntakeReq) (*pbm.CreateWaterIntakeRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "CreateWaterMeasurement")
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

	res, err := s.repo.CreateWaterMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.CreateWaterIntakeRes{
		Success: true,
		Message: "Water inserted correctly",
		Water:   res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) GetWaterMeasurements(ctx context.Context, req *pbm.GetWaterIntakesReq) (*pbm.GetWaterIntakesRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWaterMeasurements")
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

	res, err := s.repo.GetWaterMeasurements(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.GetWaterIntakesRes{
		Success: true,
		Message: "Weight fetched correctly",
		Water:   res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) GetWaterMeasurement(ctx context.Context, req *pbm.GetWaterIntakeReq) (*pbm.GetWaterIntakeRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWaterMeasurement")
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

	_, err := s.repo.GetWaterMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	// TODO Change proto def
	return &pbm.GetWaterIntakeRes{
		WaterIntakeId: req.WaterIntakeId,
		UserId:        userID,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) DeleteWaterMeasurement(ctx context.Context, req *pbm.DeleteWaterIntakeReq) (*pbm.NilRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "DeleteWaterMeasurement")
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

	res, err := s.repo.DeleteWaterMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}
func (s ServiceMeasurement) UpdateWaterMeasurement(ctx context.Context, req *pbm.UpdateWaterIntakeReq) (*pbm.UpdateWaterIntakeRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "UpdateWaterMeasurement")
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

	res, err := s.repo.UpdateWaterMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.UpdateWaterIntakeRes{
		Success: true,
		Message: "Weight updated correctly",
		Water:   res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}

// Waistline

func (s ServiceMeasurement) CreateWasteLineMeasurement(ctx context.Context, req *pbm.CreateWasteLineReq) (*pbm.CreateWasteLineRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "CreateWasteLineMeasurement")
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

	res, err := s.repo.CreateWasteLineMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.CreateWasteLineRes{
		Success:   true,
		Message:   "Water inserted correctly",
		WasteLine: res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) GetWasteLineMeasurements(ctx context.Context, req *pbm.GetWasteLinesReq) (*pbm.GetWasteLinesRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWasteLineMeasurements")
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

	res, err := s.repo.GetWasteLineMeasurements(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.GetWasteLinesRes{
		Success:   true,
		Message:   "Weight fetched correctly",
		WasteLine: res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) GetWasteLineMeasurement(ctx context.Context, req *pbm.GetWasteLineReq) (*pbm.GetWasteLineRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "GetWasteLineMeasurement")
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

	_, err := s.repo.GetWasteLineMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	// TODO Change proto def
	return &pbm.GetWasteLineRes{
		WasteLineId: req.WasteLineId,
		UserId:      userID,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
func (s ServiceMeasurement) DeleteWasteLineMeasurement(ctx context.Context, req *pbm.DeleteWasteLineReq) (*pbm.NilRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "DeleteWasteLineMeasurement")
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

	res, err := s.repo.DeleteWasteLineMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return res, nil
}
func (s ServiceMeasurement) UpdateWasteLineMeasurement(ctx context.Context, req *pbm.UpdateWasteLineReq) (*pbm.UpdateWasteLineRes, error) {
	tracer := otel.Tracer("UserMeasurements")
	ctx, span := tracer.Start(ctx, "UpdateWasteLineMeasurement")
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

	res, err := s.repo.UpdateWasteLineMeasurement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	span.SetAttributes(
		attribute.String("request.id", req.Request.RequestId),
		attribute.String("request.details", req.String()),
	)

	return &pbm.UpdateWasteLineRes{
		Success:   true,
		Message:   "Weight updated correctly",
		WasteLine: res,
		Response: &pbm.BaseResponse{
			RequestId: req.Request.RequestId,
			Upstream:  "workout-service",
		},
	}, nil
}
