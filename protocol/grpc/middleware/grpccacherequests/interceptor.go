package grpccacherequests

import (
	"context"

	"google.golang.org/grpc"

	"github.com/FACorreiaa/fitme-grpc/protocol/grpc/middleware"
)

// CacheRequestInterceptor is a gRPC client-side interceptor that caches requests.
type CacheRequestInterceptor struct {
	cache map[string]interface{}
}

// NewCacheRequestInterceptor creates a new CacheRequestInterceptor.
func NewCacheRequestInterceptor() *CacheRequestInterceptor {
	return &CacheRequestInterceptor{
		cache: make(map[string]interface{}),
	}
}

// UnaryClientInterceptor is a gRPC client-side interceptor that caches requests.
func (c *CacheRequestInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// TODO: implement caching
	return invoker(ctx, method, req, reply, cc, opts...)
}

// StreamClientInterceptor is a gRPC client-side interceptor that caches requests.
func (c *CacheRequestInterceptor) StreamClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// TODO: implement caching
	return invoker(ctx, method, req, reply, cc, opts...)
}

// Interceptors returns the unary and stream client interceptors.
func Interceptors() (middleware.ClientInterceptor, middleware.ServerInterceptor) {
	//return NewCacheRequestInterceptor(), middleware.ServerInterceptor{}
	return middleware.ClientInterceptor{}, middleware.ServerInterceptor{}
}
