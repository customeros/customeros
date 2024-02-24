package service

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func WaitForObjectCreationAndLogSpan(ctx context.Context, s *repository.Repositories, id, nodeLabel string, span opentracing.Span) {
	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		found, findErr := s.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), id, nodeLabel)
		if found && findErr == nil {
			span.LogFields(log.Bool(fmt.Sprintf("response - %s saved in db", nodeLabel), true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

	span.LogFields(log.String(fmt.Sprintf("response - created %s with id", nodeLabel), id))
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

	response, err := backoff.RetryWithData(operationWithData, utils.BackOffForInvokingEventsPlatformGrpcClient())
	return response, err
}
