package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	interactioneventpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/interaction_event"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

type InteractionEventService interface {
	SyncInteractionEvents(ctx context.Context, contacts []model.InteractionEventData) (SyncResult, error)
}

type interactionEventService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewInteractionEventService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) InteractionEventService {
	return &interactionEventService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.InteractionEventSyncConcurrency,
	}
}

func (s *interactionEventService) SyncInteractionEvents(ctx context.Context, interactionEvents []model.InteractionEventData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.SyncInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("num of interaction events", len(interactionEvents)))

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate intraction events input before syncing
	for _, interactionEvent := range interactionEvents {
		if interactionEvent.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(interactionEvent.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", interactionEvent.ExternalSystem))
			return SyncResult{}, errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all interaction events concurrently
	for _, interactionEventData := range interactionEvents {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(interactionEventData model.InteractionEventData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncInteractionEvent(ctx, syncMutex, interactionEventData, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(interactionEventData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), interactionEvents[0].ExternalSystem,
		interactionEvents[0].AppSource, "interaction event", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *interactionEventService) syncInteractionEvent(ctx context.Context, syncMutex *sync.Mutex, interactionEventInput model.InteractionEventData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.syncInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, interactionEventInput.ExternalSystem)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "interactionEventInput", interactionEventInput)

	tenant := common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""
	var err error

	interactionEventInput.Normalize()
	err = s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, interactionEventInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", interactionEventInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if interaction event sync should be skipped
	if interactionEventInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(interactionEventInput.SkipReason)
	} else if interactionEventInput.ExternalId == "" {
		reason = fmt.Sprintf("missing external id for interaction event, tenant %s", tenant)
		s.log.Warnf("Skip interaction event sync: %v", reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason)
	}

	senderId, senderLabel := s.getSenderIdAndLabel(ctx, interactionEventInput, span)
	receiversIdAndRelationType := make(map[string]string)
	receiversIdAndLabel := make(map[string]string)
	s.getReceiversIdAndLabel(ctx, interactionEventInput, span, receiversIdAndLabel, receiversIdAndRelationType)
	syncStatus, done := s.checkRequiredContact(interactionEventInput, senderLabel, receiversIdAndLabel, tenant, span)
	if done {
		return syncStatus
	}

	// Lock interaction event creation
	syncMutex.Lock()
	defer syncMutex.Unlock()

	if interactionEventInput.HasSessionDetails() {
		interactionSessionId, err := s.services.InteractionSessionService.MergeInteractionSession(ctx, tenant, interactionEventInput.ExternalSystem, interactionEventInput.SessionDetails, syncDate)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed merging interaction session with external reference %s for tenant %s :%s", interactionEventInput.SessionDetails.ExternalId, tenant, err.Error())
			s.log.Error(reason)
			return NewFailedSyncStatus(reason)
		}
		if interactionSessionId != "" && !interactionEventInput.BelongsTo.Available() {
			interactionEventInput.BelongsTo = model.BelongsTo{
				Session: model.ReferencedInteractionSession{
					ExternalId: interactionEventInput.SessionDetails.ExternalId,
				},
			}
		}
	}

	parentId, parentLabel, syncStatus, done := s.getParentIdAndLabel(ctx, interactionEventInput, span)
	if done {
		return syncStatus
	}
	syncStatus, done = s.checkRequiredParent(interactionEventInput, parentId, tenant, span)
	if done {
		return syncStatus
	}

	// Check if interaction event already exists
	interactionEventId, err := s.repositories.InteractionEventRepository.GetMatchedInteractionEventId(ctx, tenant, interactionEventInput.ExternalId, interactionEventInput.ExternalSystem, interactionEventInput.ExternalSourceEntity)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched interaction event with external id %s for tenant %s :%s", interactionEventInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if !failedSync {
		matchingInteractionEventExists := interactionEventId != ""
		span.LogFields(log.Bool("found matching interaction event", matchingInteractionEventExists))

		// Create new interaction event id if not found
		interactionEventId = utils.NewUUIDIfEmpty(interactionEventId)
		interactionEventInput.Id = interactionEventId
		span.LogFields(log.String("interactionEventId", interactionEventId))

		// Create or update interaction event
		interactionEventGrpcRequest := interactioneventpb.UpsertInteractionEventGrpcRequest{
			Tenant:      tenant,
			Id:          interactionEventId,
			Content:     interactionEventInput.Content,
			ContentType: interactionEventInput.ContentType,
			Channel:     interactionEventInput.Channel,
			ChannelData: interactionEventInput.ChannelData,
			Identifier:  interactionEventInput.Identifier,
			EventType:   interactionEventInput.EventType,
			Hide:        interactionEventInput.Hide,
			CreatedAt:   utils.ConvertTimeToTimestampPtr(interactionEventInput.CreatedAt),
			UpdatedAt:   utils.ConvertTimeToTimestampPtr(interactionEventInput.UpdatedAt),
			SourceFields: &commonpb.SourceFields{
				Source:    interactionEventInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(interactionEventInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: interactionEventInput.ExternalSystem,
				ExternalId:       interactionEventInput.ExternalId,
				ExternalUrl:      interactionEventInput.ExternalUrl,
				ExternalIdSecond: interactionEventInput.ExternalIdSecond,
				ExternalSource:   interactionEventInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}
		if parentId != "" {
			switch parentLabel {
			case model2.NodeLabelIssue:
				interactionEventGrpcRequest.BelongsToIssueId = &parentId
			case model2.NodeLabelInteractionSession:
				interactionEventGrpcRequest.BelongsToSessionId = &parentId
			}
		}
		if !matchingInteractionEventExists {
			if senderId != "" {
				participant := interactioneventpb.Participant{
					Id: senderId,
				}
				s.setParticipantTypeForGrpcRequest(senderLabel, &participant)
				interactionEventGrpcRequest.Sender = &interactioneventpb.Sender{
					Participant:  &participant,
					RelationType: interactionEventInput.SentBy.RelationType,
				}
			}
			for receiverId, receiverLabel := range receiversIdAndLabel {
				participant := interactioneventpb.Participant{
					Id: receiverId,
				}
				s.setParticipantTypeForGrpcRequest(receiverLabel, &participant)
				interactionEventGrpcRequest.Receivers = append(interactionEventGrpcRequest.Receivers, &interactioneventpb.Receiver{
					Participant:  &participant,
					RelationType: receiversIdAndRelationType[receiverId],
				})
			}
		}
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = CallEventsPlatformGRPCWithRetry[*interactioneventpb.InteractionEventIdGrpcResponse](func() (*interactioneventpb.InteractionEventIdGrpcResponse, error) {
			return s.grpcClients.InteractionEventClient.UpsertInteractionEvent(ctx, &interactionEventGrpcRequest)
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertInteractionEvent"))
			reason = fmt.Sprintf("failed sending event to upsert interaction event with external reference %s for tenant %s :%s", interactionEventInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for interaction event to be created in neo4j
		if !failedSync && !matchingInteractionEventExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				issue, findErr := s.repositories.InteractionEventRepository.GetById(ctx, tenant, interactionEventId)
				if issue != nil && findErr == nil {
					break
				}
				time.Sleep(utils.BackOffExponentialDelay(i))
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}

func (s *interactionEventService) getParticipantIdAndLabel(ctx context.Context, externalSystemId string, participant model.InteractionEventParticipant) (id string, label string, err error) {
	id = ""
	label = ""
	err = nil

	if participant.ReferencedUser.Available() {
		id, label, err = s.services.FinderService.FindReferencedEntityId(ctx, externalSystemId, &participant.ReferencedUser)
	}
	if id == "" && participant.ReferencedContact.Available() {
		id, label, err = s.services.FinderService.FindReferencedEntityId(ctx, externalSystemId, &participant.ReferencedContact)
	}
	if id == "" && participant.ReferencedOrganization.Available() {
		id, label, err = s.services.FinderService.FindReferencedEntityId(ctx, externalSystemId, &participant.ReferencedOrganization)
	}
	if id == "" && participant.ReferencedParticipant.Available() {
		id, label, err = s.services.FinderService.FindReferencedEntityId(ctx, externalSystemId, &participant.ReferencedParticipant)
	}
	if id == "" && participant.ReferencedJobRole.Available() {
		id, label, err = s.services.FinderService.FindReferencedEntityId(ctx, externalSystemId, &participant.ReferencedJobRole)
	}
	if id == "" {
		label = ""
	}
	return
}

func (s *interactionEventService) getParentIdAndLabel(ctx context.Context, interactionEventInput model.InteractionEventData, span opentracing.Span) (string, string, SyncStatus, bool) {
	parentId, parentLabel := "", ""
	var err error
	if interactionEventInput.BelongsTo.Issue.Available() {
		parentId, parentLabel, err = s.services.FinderService.FindReferencedEntityId(ctx, interactionEventInput.ExternalSystem, &interactionEventInput.BelongsTo.Issue)
	} else if interactionEventInput.BelongsTo.Session.Available() {
		parentId, parentLabel, err = s.services.FinderService.FindReferencedEntityId(ctx, interactionEventInput.ExternalSystem, &interactionEventInput.BelongsTo.Session)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		reason := fmt.Sprintf("failed finding parent for interaction event %s for tenant %s :%s", interactionEventInput.ExternalId, common.GetTenantFromContext(ctx), err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return "", "", NewFailedSyncStatus(reason), true
	}
	return parentId, parentLabel, SyncStatus{}, false
}

func (s *interactionEventService) getSenderIdAndLabel(ctx context.Context, interactionEventInput model.InteractionEventData, span opentracing.Span) (string, string) {
	if interactionEventInput.SentBy.Available() {
		senderId, senderLabel, err := s.getParticipantIdAndLabel(ctx, interactionEventInput.ExternalSystem, interactionEventInput.SentBy)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Error(fmt.Sprintf("failed finding sender for interaction event %s for tenant %s :%s", interactionEventInput.ExternalId, common.GetTenantFromContext(ctx), err.Error()))
		}
		return senderId, senderLabel
	}
	return "", ""
}

func (s *interactionEventService) getReceiversIdAndLabel(ctx context.Context, interactionEventInput model.InteractionEventData, span opentracing.Span, idAndLabel map[string]string, idAndRelationType map[string]string) {
	for _, receiver := range interactionEventInput.SentTo {
		receiverId, receiverLabel, err := s.getParticipantIdAndLabel(ctx, interactionEventInput.ExternalSystem, receiver)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Error(fmt.Sprintf("failed finding receiver for interaction event %s for tenant %s :%s", interactionEventInput.ExternalId, common.GetTenantFromContext(ctx), err.Error()))
		}
		if receiverId != "" {
			idAndLabel[receiverId] = receiverLabel
			idAndRelationType[receiverId] = receiver.RelationType
		}
	}
}

func (s *interactionEventService) checkRequiredContact(interactionEventInput model.InteractionEventData, senderLabel string, receiversIdAndLabel map[string]string, tenant string, span opentracing.Span) (SyncStatus, bool) {
	if interactionEventInput.ContactRequired {
		found := false
		if senderLabel == model2.NodeLabelContact {
			found = true
		}
		for _, receiverLabel := range receiversIdAndLabel {
			if receiverLabel == model2.NodeLabelContact {
				found = true
				break
			}
		}
		if !found {
			reason := fmt.Sprintf("contact not found for interaction event %s for tenant %s", interactionEventInput.ExternalId, tenant)
			s.log.Warnf("Skip interaction event sync: %v", reason)
			span.LogFields(log.String("output", "skipped"))
			return NewSkippedSyncStatus(reason), true
		}
	}
	return SyncStatus{}, false
}

func (s *interactionEventService) checkRequiredParent(interactionEventInput model.InteractionEventData, parentId string, tenant string, span opentracing.Span) (SyncStatus, bool) {
	if interactionEventInput.ParentRequired && parentId == "" {
		reason := fmt.Sprintf("parent not found for interaction event %s for tenant %s", interactionEventInput.ExternalId, tenant)
		s.log.Warnf("Skip interaction event sync: %v", reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason), true
	}
	return SyncStatus{}, false
}

func (s *interactionEventService) setParticipantTypeForGrpcRequest(participantLabel string, participant *interactioneventpb.Participant) {
	switch participantLabel {
	case model2.NodeLabelContact:
		participant.ParticipantType = &interactioneventpb.Participant_Contact{
			Contact: &interactioneventpb.Contact{},
		}
	case model2.NodeLabelOrganization:
		participant.ParticipantType = &interactioneventpb.Participant_Organization{
			Organization: &interactioneventpb.Organization{},
		}
	case model2.NodeLabelUser:
		participant.ParticipantType = &interactioneventpb.Participant_User{
			User: &interactioneventpb.User{},
		}
	case model2.NodeLabelJobRole:
		participant.ParticipantType = &interactioneventpb.Participant_JobRole{
			JobRole: &interactioneventpb.JobRole{},
		}
	}

}
