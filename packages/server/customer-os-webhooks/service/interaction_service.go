package service

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonlogger "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	interactionsessionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_session"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	webhookmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
)

type InteractionService interface {
	ProcessInteractionSession(ctx context.Context, tenant string, data *webhookmodel.PostmarkEmailWebhookData) (string, error)
	GetIdForReferencedSession(ctx context.Context, tenant, externalSystem string, session webhookmodel.ReferencedInteractionSession) (string, error)
}

type interactionService struct {
	services    *Services
	logger      commonlogger.Logger
	grpcClients *grpc_client.Clients
}

func NewInteractionService(services *Services, grpcClients *grpc_client.Clients) InteractionService {
	return &interactionService{
		services:    services,
		logger:      services.Logger,
		grpcClients: grpcClients,
	}
}

func (s *interactionService) ProcessInteractionSession(ctx context.Context, tenant string, data *webhookmodel.PostmarkEmailWebhookData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionService.ProcessInteractionSession")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	sessionData := webhookmodel.InteractionSessionData{
		ExternalId:  data.MessageID,
		ExternalUrl: data.MessageURL,
		Channel:     "email",
		Status:      "completed",
		Type:        "email",
		Name:        data.Subject,
	}
	sessionData.Normalize()

	return s.mergeInteractionSession(ctx, tenant, "mailstack", sessionData, time.Now())
}

func (s *interactionService) mergeInteractionSession(ctx context.Context, tenant, externalSystem string, sessionData webhookmodel.InteractionSessionData, syncDate time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionService.mergeInteractionSession")
	defer span.Finish()

	sessionId, err := s.services.CommonServices.PostgresRepositories.InteractionSessionRepository.GetInteractionSessionIdByExternalId(
		ctx, tenant, sessionData.ExternalId, externalSystem,
	)
	if err != nil {
		return "", err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := s.callGRPCWithRetry(ctx, tenant, sessionId, externalSystem, sessionData, syncDate)
	if err != nil {
		return "", err
	}

	// Verify session was saved
	err = s.verifySessionSaved(ctx, tenant, response.Id, externalSystem)
	if err != nil {
		return "", err
	}

	span.LogFields(log.String("response.InteractionSessionId", response.Id))
	return response.Id, nil
}

func (s *interactionService) callGRPCWithRetry(ctx context.Context, tenant, sessionId, externalSystem string, sessionData model.InteractionSessionData, syncDate time.Time) (*interactionsessionpb.InteractionSessionIdGrpcResponse, error) {
	return CallEventsPlatformGRPCWithRetry[*interactionsessionpb.InteractionSessionIdGrpcResponse](func() (*interactionsessionpb.InteractionSessionIdGrpcResponse, error) {
		return s.grpcClients.InteractionSessionClient.UpsertInteractionSession(ctx, &interactionsessionpb.UpsertInteractionSessionGrpcRequest{
			Tenant: tenant,
			Id:     sessionId,
			SourceFields: &commonpb.SourceFields{
				Source:    externalSystem,
				AppSource: utils.StringFirstNonEmpty(sessionData.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: externalSystem,
				ExternalId:       sessionData.ExternalId,
				ExternalUrl:      sessionData.ExternalUrl,
				ExternalIdSecond: sessionData.ExternalIdSecond,
				ExternalSource:   sessionData.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
			Identifier:  sessionData.Identifier,
			Channel:     sessionData.Channel,
			ChannelData: sessionData.ChannelData,
			Status:      sessionData.Status,
			Type:        sessionData.Type,
			Name:        sessionData.Name,
		})
	})
}

func (s *interactionService) verifySessionSaved(ctx context.Context, tenant, sessionId, externalSystem string) error {
	for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
		found, err := s.services.CommonServices.Neo4jRepositories.CommonReadRepository.ExistsByIdLinkedTo(
			ctx,
			tenant,
			sessionId,
			commonmodel.NodeLabelInteractionSession,
			externalSystem,
			commonmodel.NodeLabelExternalSystem,
			"IS_LINKED_WITH",
		)
		if found && err == nil {
			return nil
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	return errors.New("failed to verify interaction session was saved")
}

func (s *interactionService) GetIdForReferencedSession(ctx context.Context, tenant, externalSystem string, session model.ReferencedInteractionSession) (string, error) {
	if !session.Available() {
		return "", nil
	}

	if session.ReferencedByExternalId() {
		return s.services.CommonServices.PostgresRepositories.InteractionSessionRepository.GetInteractionSessionIdByExternalId(
			ctx,
			tenant,
			session.ExternalId,
			externalSystem,
		)
	}
	return "", nil
}
