package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type InteractionSessionService interface {
	GetIdForReferencedInteractionSession(ctx context.Context, tenant, externalSystem string, user model.ReferencedInteractionSession) (string, error)
	MergeInteractionSession(ctx context.Context, tenant, externalSystemId string, interactionSessionInput model.InteractionSessionData, syncDate time.Time) (string, error)
}

type interactionSessionService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewInteractionSessionService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) InteractionSessionService {
	return &interactionSessionService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *interactionSessionService) GetIdForReferencedInteractionSession(ctx context.Context, tenant, externalSystemId string, interactionSession model.ReferencedInteractionSession) (string, error) {
	if !interactionSession.Available() {
		return "", nil
	}

	if interactionSession.ReferencedByExternalId() {
		return s.repositories.InteractionSessionRepository.GetInteractionSessionIdByExternalId(ctx, tenant, interactionSession.ExternalId, externalSystemId)
	}
	return "", nil
}

func (s *interactionSessionService) MergeInteractionSession(ctx context.Context, tenant, externalSystem string, interactionSessionInput model.InteractionSessionData, syncDate time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionService.MergeInteractionSession")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	interactionSessionInput.Normalize()

	interactionSessionId, err := s.repositories.InteractionSessionRepository.GetInteractionSessionIdByExternalId(ctx, tenant, interactionSessionInput.ExternalId, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := CallEventsPlatformGRPCWithRetry[*interactionsessionpb.InteractionSessionIdGrpcResponse](func() (*interactionsessionpb.InteractionSessionIdGrpcResponse, error) {
		return s.grpcClients.InteractionSessionClient.UpsertInteractionSession(ctx, &interactionsessionpb.UpsertInteractionSessionGrpcRequest{
			Tenant: tenant,
			Id:     interactionSessionId,
			SourceFields: &commonpb.SourceFields{
				Source:    externalSystem,
				AppSource: utils.StringFirstNonEmpty(interactionSessionInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: externalSystem,
				ExternalId:       interactionSessionInput.ExternalId,
				ExternalUrl:      interactionSessionInput.ExternalUrl,
				ExternalIdSecond: interactionSessionInput.ExternalIdSecond,
				ExternalSource:   interactionSessionInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
			Identifier:  interactionSessionInput.Identifier,
			Channel:     interactionSessionInput.Channel,
			ChannelData: interactionSessionInput.ChannelData,
			Status:      interactionSessionInput.Status,
			Type:        interactionSessionInput.Type,
			Name:        interactionSessionInput.Name,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "UpsertInteractionSessionGrpcRequest")
	}
	for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
		// check if the interaction session is saved in db and is linked to external system
		found, findErr := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedTo(ctx, tenant, response.Id, model2.NodeLabelInteractionSession, externalSystem, model2.NodeLabelExternalSystem, "IS_LINKED_WITH")
		if found && findErr == nil {
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	span.LogFields(log.String("response.InteractionSessionId", response.Id))
	return response.Id, nil
}
