package subscriptions

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func WaitCheckNodeExistsInNeo4j(ctx context.Context, neo4jRepository *neo4jrepository.Repositories, tenant, id, nodeLabel string) bool {
	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4j; i++ {
		found, findErr := neo4jRepository.CommonReadRepository.ExistsById(ctx, tenant, id, nodeLabel)
		if found && findErr == nil {
			return true
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	return false
}

func CallEventsPlatformGRPCWithRetry[T any](operation func() (T, error)) (T, error) {
	operationWithData := func() (T, error) {
		result, opErr := operation()
		if opErr != nil {
			grpcError, ok := status.FromError(opErr)
			if ok && (grpcError.Code() == codes.Unavailable || grpcError.Code() == codes.DeadlineExceeded) {
				return result, opErr
			}
			return result, backoff.Permanent(opErr)
		}
		return result, nil
	}

	response, err := backoff.RetryWithData(operationWithData, utils.BackOffForInvoking())
	return response, err
}
