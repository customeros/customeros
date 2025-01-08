package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ApiKeyHeader = "X-Openline-API-KEY"

func ApiKeyEnricher(apiKey string) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctxWithAPIKey := metadata.AppendToOutgoingContext(ctx, ApiKeyHeader, apiKey)
		return invoker(ctxWithAPIKey, method, req, reply, cc, opts...)
	}
}
