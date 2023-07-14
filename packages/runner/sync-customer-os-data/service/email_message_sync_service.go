package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	INBOUND  string = "INBOUND"
	OUTBOUND string = "OUTBOUND"
)

type emailMessageSyncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewDefaultEmailMessageSyncService(repositories *repository.Repositories, services *Services) SyncService {
	return &emailMessageSyncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *emailMessageSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		messages := dataService.GetDataForSync(common.EMAIL_MESSAGES, batchSize, runId)
		if len(messages) == 0 {
			logrus.Debugf("no email messages found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d email messages from %s for tenant %s", len(messages), dataService.SourceId(), tenant)

		for _, v := range messages {
			var failedSync = false
			var reason string
			var interactionEventId string
			messageInput := v.(entity.EmailMessageData)
			messageInput.Normalize()

			if messageInput.Skip {
				if err := dataService.MarkProcessed(messageInput.SyncId, runId, true, true, messageInput.SkipReason); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tenant, syncDate, messageInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed merge interaction session with external reference %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
				logrus.Errorf(reason)
			}

			if !failedSync {
				interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tenant, messageInput.ExternalSystem, syncDate, messageInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge interaction event with external reference %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}

				if !failedSync {
					err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tenant, interactionEventId, sessionId)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed to associate interaction event to session %v for tenant %v :%v", messageInput.ExternalId, tenant, err)
						logrus.Errorf(reason)
					}
				}
			}

			//from
			if messageInput.Direction == OUTBOUND && !failedSync {
				emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, messageInput.FromEmail)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed retrieving email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
					logrus.Errorf(reason)
				}

				if emailId == "" && !failedSync {
					emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateUserByEmail(ctx, tenant, messageInput.FromEmail, messageInput.FromFirstName, messageInput.FromLastName, messageInput.ExternalSystem)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed creating contact with email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
						logrus.Errorf(reason)
					}
				}

				if !failedSync {
					err := s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
						logrus.Errorf(reason)
					}
				}
			} else if messageInput.Direction == INBOUND && !failedSync {
				//1. find email ( contact/organization/user )
				//2. if not found, create contact with email

				emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, messageInput.FromEmail)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed retrieving email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
					logrus.Errorf(reason)
				}

				if emailId == "" && !failedSync {
					emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateContactByEmail(ctx, tenant, messageInput.FromEmail, messageInput.FromFirstName, messageInput.FromLastName, messageInput.ExternalSystem)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed creating contact with email %v for tenant %v :%v", messageInput.FromEmail, tenant, err)
						logrus.Errorf(reason)
					}
				}

				if !failedSync {
					err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
						logrus.Errorf(reason)
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
					reason = fmt.Sprintf("failed set TO users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			//cc
			if len(messageInput.CcEmail) > 0 && !failedSync {
				err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "CC", messageInput.CcEmail)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed set CC users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			//bcc
			if len(messageInput.BccEmail) > 0 && !failedSync {
				err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "BCC", messageInput.BccEmail)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed set BCC users %v for interaction event %v in tenant %v :%v", messageInput.ContactsExternalIds, sessionId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if !failedSync {
				for _, v := range messageInput.ContactsExternalIds {
					s.services.OrganizationService.UpdateLastTouchpointByContactIdExternalId(ctx, tenant, v, messageInput.ExternalSystem)
				}
			}

			logrus.Debugf("successfully merged email message with external id %v to interaction session %v for tenant %v from %v", messageInput.ExternalId, sessionId, tenant, dataService.SourceId())
			if err := dataService.MarkProcessed(messageInput.SyncId, runId, failedSync == false, false, reason); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(messages) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}
