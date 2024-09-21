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
