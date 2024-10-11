package grpclog

import (
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware"
)

func Interceptors(instance *zap.Logger) (middleware.ClientInterceptor, middleware.ServerInterceptor) {
	opt := []grpcZap.Option{
		grpcZap.WithLevels(codeToLevel),
	}

	grpcZap.ReplaceGrpcLoggerV2WithVerbosity(instance, int(zap.WarnLevel))

	//serverUnaryInterceptor := grpc.UnaryServerInterceptor(func(
	//	ctx context.Context,
	//	req interface{},
	//	info *grpc.UnaryServerInfo,
	//	handler grpc.UnaryHandler,
	//) (interface{}, error) {
	//	// Log request size
	//	requestSize := proto.Size(req.(proto.Message))
	//	instance.Info("Received request",
	//		zap.String("method", info.FullMethod),
	//		zap.Int("request_size_bytes", requestSize),
	//	)
	//
	//	// Handle the request
	//	res, err := handler(ctx, req)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	// Log response size
	//	responseSize := proto.Size(res.(proto.Message))
	//	instance.Info("Sending response",
	//		zap.String("method", info.FullMethod),
	//		zap.Int("response_size_bytes", responseSize),
	//	)
	//
	//	return res, nil
	//})

	clientInterceptor := middleware.ClientInterceptor{
		Unary:  grpcZap.UnaryClientInterceptor(instance, opt...),
		Stream: grpcZap.StreamClientInterceptor(instance, opt...),
	}

	serverInterceptor := middleware.ServerInterceptor{
		Unary:  grpcZap.UnaryServerInterceptor(instance, opt...),
		Stream: grpcZap.StreamServerInterceptor(instance, opt...),
	}

	return clientInterceptor, serverInterceptor
}

// codeToLevel translates a GRPC status code to a zap logging level
func codeToLevel(code codes.Code) zapcore.Level {
	// override OK to DebugLevel
	if code == codes.OK {
		return zap.DebugLevel
	}

	return grpcZap.DefaultCodeToLevel(code)
}
