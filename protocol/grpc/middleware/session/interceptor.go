package session

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/FACorreiaa/fitme-grpc/internal/domain/auth"
)

func InterceptorSession(sessionManager *auth.SessionManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/auth.Auth/Register" {
			return handler(ctx, req)
		}

		fmt.Printf("info.FullMethod: %s\n", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing context metadata")
		}

		// Extract token
		token := md["authorization"]
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing auth token")
		}

		// Validate session
		userSession, err := sessionManager.GetSession(token[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid session token")
		}

		// Pass user session in context
		ctx = context.WithValue(ctx, "userSession", userSession)
		return handler(ctx, req)
	}
}

//// Add a custom type for context key
//type contextKey string
//
//const userKey contextKey = "userID"
//
//// In your UnaryInterceptor, you could do something like this:
//func (i *StaticInterceptor) UnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
//	err := i.engines.IsAuthorized(ctx)
//	if err != nil {
//		if status.Code(err) == codes.PermissionDenied {
//			logger.Infof("unauthorized RPC request rejected for method %s: %v", info.FullMethod, err)
//			return nil, status.Errorf(codes.PermissionDenied, "unauthorized request to %s", info.FullMethod)
//		}
//		return nil, err
//	}
//
//	// Assuming user information is set in the context after authorization
//	userID, ok := ctx.Value(userKey).(string)
//	if ok {
//		logger.Infof("User %s authorized for method %s", userID, info.FullMethod)
//	}
//	return handler(ctx, req)
//}
