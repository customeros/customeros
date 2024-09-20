package subscriptions

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

	response, err := backoff.RetryWithData(operationWithData, utils.BackOffForInvokingEventsPlatformGrpcClient())
	return response, err
}

func EventCompleted(ctx context.Context, tenant, entity, entityId, entityType string, grpcClients *grpc_client.Clients) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventCompleted")
	defer span.Finish()
	span.LogKV("tenant", tenant, "entity", entity, "entityID", entityId, "entityType", entityType)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
		//return grpcClients.EventCompletionClient.NotifyEventProcessed(ctx, &eventcompletionpb.NotifyEventProcessedRequest{
		//	Tenant:    tenant,
		//	EventType: entityType,
		//	Entity:    entity,
		//	EntityId:  entityId,
		//})
		return &emptypb.Empty{}, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
}
