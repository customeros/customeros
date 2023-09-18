package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type interactionEventSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultInteractionEventSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &interactionEventSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *interactionEventSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	interactionEventSyncMutex := &sync.Mutex{}

	for {
		events := dataService.GetDataForSync(ctx, common.INTERACTION_EVENTS, batchSize, runId)

		if len(events) == 0 {
			break
		}

		s.log.Infof("Syncing %d interaction events", len(events))

		var wg sync.WaitGroup
		wg.Add(len(events))

		results := make(chan result, len(events))
		done := make(chan struct{})

		for _, v := range events {
			v := v

			go func(event entity.InteractionEventData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncInteractionEvent(ctx, interactionEventSyncMutex, event, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.InteractionEventData))
		}
		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()
		go func() {
			<-done
			close(results)
		}()

		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}

		if len(events) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *interactionEventSyncService) syncInteractionEvent(ctx context.Context, interactionEventSyncMutex *sync.Mutex, interactionEventInput entity.InteractionEventData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventSyncService.syncInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	interactionEventInput.Normalize()

	if interactionEventInput.ExternalSystem == "" {
		_ = dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
		*failed++
		return
	}

	if interactionEventInput.Skip {
		if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, true, true, interactionEventInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	if interactionEventInput.ContactRequired {
		found := false
		_, label, _ := s.getParticipantId(ctx, tenant, interactionEventInput.ExternalSystem, interactionEventInput.SentBy)
		if label == "Contact" {
			found = true
		}
		if !found {
			for _, sentTo := range interactionEventInput.SentTo {
				_, label, _ = s.getParticipantId(ctx, tenant, interactionEventInput.ExternalSystem, sentTo)
				if label == "Contact" {
					found = true
					break
				}
			}
		}
		if !found {
			if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, true, true, "No contact found as interaction event participant"); err != nil {
				*failed++
				span.LogFields(log.Bool("failedSync", true))
				return
			}
			*skipped++
			span.LogFields(log.Bool("skippedSync", true))
			return
		}
	}
	if interactionEventInput.SessionRequired {
		found := false
		id, _, _ := s.getReferencedEntityIdAndLabel(ctx, tenant, interactionEventInput.ExternalSystem, &interactionEventInput.PartOfSession)
		if id != "" {
			found = true
		}
		if !found {
			if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, true, true, "Interaction session not found: "+interactionEventInput.PartOfSession.ExternalId); err != nil {
				*failed++
				span.LogFields(log.Bool("failedSync", true))
				return
			}
			*skipped++
			span.LogFields(log.Bool("skippedSync", true))
			return
		}
	}

	interactionEventSyncMutex.Lock()
	interactionEventId, err := s.repositories.InteractionEventRepository.GetMatchedInteractionEvent(ctx, tenant, interactionEventInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched interaction event with external reference id %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
		s.log.Error(reason)
	}

	// Create new interaction event id if not found
	if interactionEventId == "" {
		ieUuid, _ := uuid.NewRandom()
		interactionEventId = ieUuid.String()
	}
	interactionEventInput.Id = interactionEventId
	span.LogFields(log.String("interactionEventId", interactionEventId))

	if !failedSync {
		err = s.repositories.InteractionEventRepository.MergeInteractionEvent(ctx, tenant, syncDate, interactionEventInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge interaction event with external reference %v for tenant %v :%v", interactionEventInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}
	interactionEventSyncMutex.Unlock()

	if !failedSync && interactionEventInput.HasSession() {
		err = s.repositories.InteractionEventRepository.MergeInteractionSessionForEvent(ctx, tenant, interactionEventId, interactionEventInput.ExternalSystem, syncDate, interactionEventInput.SessionDetails)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge interaction session by external id %v for tenant %v :%v", interactionEventInput.SessionDetails.ExternalId, tenant, err)
			s.log.Error(reason)
		}
	}

	if !failedSync && interactionEventInput.IsPartOf() {
		var id, label string
		if interactionEventInput.PartOfIssue.Available() {
			id, label, _ = s.getReferencedEntityIdAndLabel(ctx, tenant, interactionEventInput.ExternalSystem, &interactionEventInput.PartOfIssue)
		} else if interactionEventInput.PartOfSession.Available() {
			id, label, _ = s.getReferencedEntityIdAndLabel(ctx, tenant, interactionEventInput.ExternalSystem, &interactionEventInput.PartOfSession)
		}
		if id != "" {
			err = s.repositories.InteractionEventRepository.LinkInteractionEventAsPartOf(ctx, tenant, interactionEventId, id, label)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link interaction event %s to be part of %s %s for tenant %s :%s", interactionEventId, label, id, tenant, err.Error())
				s.log.Error(reason)
			}
		}
	}

	if !failedSync && interactionEventInput.HasSender() {
		sender := interactionEventInput.SentBy

		id, label, err := s.getParticipantId(ctx, tenant, interactionEventInput.ExternalSystem, sender)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed finding participant %v for tenant %s :%s", sender, tenant, err.Error())
			s.log.Error(reason)
		}
		if id != "" {
			err = s.repositories.InteractionEventRepository.LinkInteractionEventWithSenderById(ctx, tenant, interactionEventId, id, label)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link interaction event with job role for tenant %v :%v", tenant, err)
				s.log.Error(reason)
			}
		}
	}

	if !failedSync {
		for _, recipient := range interactionEventInput.SentTo {
			id, label, err := s.getParticipantId(ctx, tenant, interactionEventInput.ExternalSystem, recipient)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed finding participant %v for tenant %s :%s", recipient, tenant, err.Error())
				s.log.Error(reason)
			}
			if id != "" {
				err = s.repositories.InteractionEventRepository.LinkInteractionEventWithRecipientById(ctx, tenant, interactionEventId, id, label, recipient.RelationType)
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed link interaction event with job role for tenant %v :%v", tenant, err)
					s.log.Error(reason)
				}
			}
		}
	}

	if failedSync == false {
		s.log.Debugf("successfully merged interaction event with id %v for tenant %v from %v", interactionEventId, tenant, dataService.SourceId())
	}
	if err := dataService.MarkProcessed(ctx, interactionEventInput.SyncId, runId, failedSync == false, false, reason); err != nil {
		*failed++
		span.LogFields(log.Bool("failedSync", true))
		return
	}
	if failedSync == true {
		*failed++
	} else {
		*completed++
	}
	span.LogFields(log.Bool("failedSync", failedSync))
}

func (s *interactionEventSyncService) getParticipantId(ctx context.Context, tenant, externalSystemId string, participant entity.InteractionEventParticipant) (id string, label string, err error) {
	id = ""
	label = ""
	err = nil

	if participant.ReferencedUser.Available() {
		id, label, err = s.getReferencedEntityIdAndLabel(ctx, tenant, externalSystemId, &participant.ReferencedUser)
	}
	if id == "" && participant.ReferencedContact.Available() {
		id, label, err = s.getReferencedEntityIdAndLabel(ctx, tenant, externalSystemId, &participant.ReferencedContact)
	}
	if id == "" && participant.ReferencedOrganization.Available() {
		id, label, err = s.getReferencedEntityIdAndLabel(ctx, tenant, externalSystemId, &participant.ReferencedOrganization)
	}
	if id == "" && participant.ReferencedParticipant.Available() {
		id, label, err = s.getReferencedEntityIdAndLabel(ctx, tenant, externalSystemId, &participant.ReferencedParticipant)
	}
	if id == "" && participant.ReferencedJobRole.Available() {
		id, label, err = s.getReferencedEntityIdAndLabel(ctx, tenant, externalSystemId, &participant.ReferencedJobRole)
	}
	if id == "" {
		label = ""
	}
	return
}

func (s *interactionEventSyncService) getReferencedEntityIdAndLabel(ctx context.Context, tenant, externalSystemId string, refEntity entity.ReferencedEntity) (id string, label string, err error) {
	id = ""
	label = ""
	err = nil
	if !refEntity.Available() {
		return "", "", nil
	}
	switch r := refEntity.(type) {
	case *entity.ReferencedInteractionSession:
		id, err = s.GetIdForReferencedInteractionSession(ctx, tenant, externalSystemId, *r)
		if id != "" {
			label = "InteractionSession"
		}
	case *entity.ReferencedIssue:
		id, err = s.services.IssueService.GetIdForReferencedIssue(ctx, tenant, externalSystemId, *r)
		if id != "" {
			label = "Issue"
		}
	case *entity.ReferencedUser:
		id, err = s.services.UserService.GetIdForReferencedUser(ctx, tenant, externalSystemId, *r)
		if id != "" {
			label = "User"
		}
	case *entity.ReferencedContact:
		id, err = s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, *r)
		if id != "" {
			label = "Contact"
		}
	case *entity.ReferencedOrganization:
		id, err = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, *r)
		if id != "" {
			label = "Organization"
		}
	case *entity.ReferencedJobRole:
		contactId, _ := s.services.ContactService.GetIdForReferencedContact(ctx, tenant, externalSystemId, r.ReferencedContact)
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, externalSystemId, r.ReferencedOrganization)
		id, err = s.repositories.ContactRepository.GetJobRoleId(ctx, tenant, contactId, orgId)
		if id != "" {
			label = "JobRole"
		}
	case *entity.ReferencedParticipant:
		if id == "" {
			id, err = s.repositories.UserRepository.GetUserIdByExternalId(ctx, tenant, r.ExternalId, externalSystemId)
			if id != "" {
				label = "User"
			}
		}
		if id == "" {
			id, err = s.repositories.ContactRepository.GetContactIdByExternalId(ctx, tenant, r.ExternalId, externalSystemId)
			if id != "" {
				label = "Contact"
			}
		}
		if id == "" {
			id, err = s.repositories.OrganizationRepository.GetOrganizationIdByExternalId(ctx, tenant, r.ExternalId, externalSystemId)
			if id != "" {
				label = "Organization"
			}
		}
	}
	if id == "" {
		label = ""
	}
	return
}

func (s *interactionEventSyncService) GetIdForReferencedInteractionSession(ctx context.Context, tenant, externalSystemId string, interactionSession entity.ReferencedInteractionSession) (string, error) {
	if !interactionSession.Available() {
		return "", nil
	}

	if interactionSession.ReferencedByExternalId() {
		return s.repositories.InteractionEventRepository.GetInteractionSessionIdByExternalId(ctx, tenant, interactionSession.ExternalId, externalSystemId)
	}
	return "", nil
}
