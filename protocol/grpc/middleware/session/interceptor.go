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
		unauthenticatedMethods := map[string]bool{
			"/auth.Auth/Register":                  true,
			"/auth.Auth/Login":                     true,
			"/auth.Auth/GetAllUsers":               true,
			"calculator.Calculator/GetUsersMacros": true,
			"CalculatorService/GetUserMacros":      true,
			"CalculatorService/GetUserMacrosAll":   true,
		}
		if unauthenticatedMethods[info.FullMethod] {
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

		// TODO handle permissions later

		//requiredPermission, ok := MethodPermissions[info.FullMethod]
		//if !ok {
		//	return nil, status.Error(codes.PermissionDenied, "method not found")
		//}
		//
		//userPermissions := GetUserPermissions(userSession.Role)
		//fmt.Printf("userPermissions: %+v\n", userPermissions)
		//if !hasPermission(userPermissions, requiredPermission) {
		//	return nil, status.Error(codes.PermissionDenied, "user does not have permission to access this method")
		//}

		// Pass user session in context
		ctx = context.WithValue(ctx, "userSession", userSession)
		return handler(ctx, req)
	}
}

func hasPermission(userPermissions []string, requiredPermission string) bool {
	for _, perm := range userPermissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}
