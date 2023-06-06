package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CheckApiKeyInterceptor(apiKey string) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "metadata not provided")
		}
		receivedApiKey := md.Get("X-Openline-API-KEY")
		if len(receivedApiKey) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "api key not provided")
		}
		if receivedApiKey[0] != apiKey {
			return nil, status.Errorf(codes.InvalidArgument, "invalid api key")
		}

		return handler(ctx, req)
	}
}
