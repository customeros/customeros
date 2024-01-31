package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func CallEventsPlatformGRPCWithRetry[T any](fn func() (T, error)) (T, error) {
	var err error
	var response T
	for attempt := 1; attempt <= constants.MaxRetryGrpcCallWhenUnavailable; attempt++ {
		response, err = fn()
		if err == nil {
			break
		}
		if grpcError, ok := status.FromError(err); ok && grpcError.Code() == codes.Unavailable {
			time.Sleep(utils.BackOffExponentialDelay(attempt))
		} else {
			break
		}
	}
	return response, err
}
