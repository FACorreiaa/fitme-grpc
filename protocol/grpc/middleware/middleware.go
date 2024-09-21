package middleware

import (
	"log"
	"sync"
	"time"

	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	Log      *zap.Logger
	onceInit sync.Once
)

func initializeLogger() {
	// Initialize Zap logger
	Log, _ = zap.NewProduction()
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			zap.Error(err)
			log.Fatal(err)
		}
	}(Log) // Flushes buffer, if any
}

func KeepaliveEnforcementPolicy() keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second,
		PermitWithoutStream: true,
	}
}

func KeepAliveServerParams() keepalive.ServerParameters {
	return keepalive.ServerParameters{
		Time:    10 * time.Second,
		Timeout: 10 * time.Second,
	}
}

type ClientInterceptor struct {
	Unary  grpc.UnaryClientInterceptor
	Stream grpc.StreamClientInterceptor
}

func (cl *ClientInterceptor) UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	onceInit.Do(initializeLogger)

	// Logic before invoking the invoker
	start := time.Now()
	// Calls the invoker to execute RPC
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Logic after invoking the invoker
	Log.Info("Invoked RPC method",
		zap.String("method", method),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)
	return err
}

type ServerInterceptor struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

// UnaryServerInterceptor is a middleware for gRPC server requests.
func (si *ServerInterceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Server request method: %s", info.FullMethod)
	// TODO fill logic
	resp, err := handler(ctx, req)
	return resp, err
}
