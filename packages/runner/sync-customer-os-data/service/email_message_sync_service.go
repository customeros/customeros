package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

const (
	INBOUND  string = "INBOUND"
	OUTBOUND string = "OUTBOUND"
)

type emailMessageSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultEmailMessageSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &emailMessageSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *emailMessageSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	span, ctx := tracing.StartTracerSpan(ctx, "EmailMessageSyncService.Sync")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	completed, failed, skipped := 0, 0, 0
	for {
		messages := dataService.GetDataForSync(ctx, common.EMAIL_MESSAGES, batchSize, runId)
		if len(messages) == 0 {
			s.log.Debugf("no email messages found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		s.log.Infof("syncing %d email messages from %s for tenant %s", len(messages), dataService.SourceId(), tenant)

		for _, v := range messages {
			s.syncEmailMessage(ctx, v.(entity.EmailMessageData), dataService, syncDate, tenant, runId, &completed, &failed, &skipped)
		}
		if len(messages) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *emailMessageSyncService) syncEmailMessage(ctx context.Context, messageInput entity.EmailMessageData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactSyncService.syncContact")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	var interactionEventId string
	messageInput.Normalize()

	if messageInput.Skip {
		if err := dataService.MarkProcessed(ctx, messageInput.SyncId, runId, true, true, messageInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tenant, syncDate, messageInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merge interaction session with external reference %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	if !failedSync {
		interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tenant, messageInput.ExternalSystem, syncDate, messageInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge interaction event with external reference %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}

		if !failedSync {
			err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tenant, interactionEventId, sessionId)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed to associate interaction event to session %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
				s.log.Errorf(reason)
			}
		}
	}

	//from
	if messageInput.Direction == OUTBOUND && !failedSync {
		emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, messageInput.FromEmail)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed retrieving email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
			s.log.Errorf(reason)
		}

		if emailId == "" && !failedSync {
			emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateUserByEmail(ctx, tenant, messageInput.FromEmail, messageInput.FromFirstName, messageInput.FromLastName, messageInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed creating contact with email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
				s.log.Errorf(reason)
			}
		}

		if !failedSync {
			err := s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
				s.log.Errorf(reason)
			}
		}
	} else if messageInput.Direction == INBOUND && !failedSync {
		//1. find email ( contact/organization/user )
		//2. if not found, create contact with email

		emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, messageInput.FromEmail)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed retrieving email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
			s.log.Errorf(reason)
		}

		if emailId == "" && !failedSync {
			emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateContactByEmail(ctx, tenant, messageInput.FromEmail, messageInput.FromFirstName, messageInput.FromLastName, messageInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed creating contact with email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
				s.log.Errorf(reason)
			}
		}

		if !failedSync {
			err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
				s.log.Errorf(reason)
			}
		}
		if !failedSync {
			s.services.OrganizationService.UpdateLastTouchpointByContactEmailId(ctx, tenant, emailId)
		}
	}

	//to
	if len(messageInput.ToEmail) > 0 && !failedSync {
		err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "TO", messageInput.ToEmail)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed set TO users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	//cc
	if len(messageInput.CcEmail) > 0 && !failedSync {
		err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "CC", messageInput.CcEmail)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed set CC users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	//bcc
	if len(messageInput.BccEmail) > 0 && !failedSync {
		err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "BCC", messageInput.BccEmail)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed set BCC users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if !failedSync {
		for _, v := range messageInput.ContactsExternalIds {
			s.services.OrganizationService.UpdateLastTouchpointByContactIdExternalId(ctx, tenant, v, messageInput.ExternalSystem)
		}
	}

	s.log.Debugf("successfully merged email message with external id %v to interaction session %v for tenant %v from %v", messageInput.ExternalId, sessionId, tenant, dataService.SourceId())
	if err := dataService.MarkProcessed(ctx, messageInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
