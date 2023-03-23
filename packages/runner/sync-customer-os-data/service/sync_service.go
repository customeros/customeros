package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	hubspot_service "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot/service"
	zendesk_support_service "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

const batchSize = 100

type SyncService interface {
	Sync(ctx context.Context, runId string)
}

type syncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewSyncService(repositories *repository.Repositories, services *Services) SyncService {
	return &syncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *syncService) Sync(ctx context.Context, runId string) {
	tenantsToSync, err := s.repositories.TenantSyncSettingsRepository.GetTenantsForSync()
	if err != nil {
		logrus.Error("failed to get tenants for sync")
		return
	}

	for _, v := range tenantsToSync {
		syncRunDtls := entity.SyncRun{
			StartAt:              time.Now().UTC(),
			RunId:                runId,
			TenantSyncSettingsId: v.ID,
		}

		dataService, err := s.sourceDataService(v)
		if err != nil {
			logrus.Errorf("failed to get data service for tenant %v: %v", v.Tenant, err)
			continue
		}

		defer func() {
			dataService.Close()
		}()

		syncDate := time.Now().UTC()

		s.syncExternalSystem(ctx, dataService, v.Tenant)

		userSyncService, err := s.userSyncService(v)
		completedUserCount, failedUserCount := userSyncService.SyncUsers(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedUsers = completedUserCount
		syncRunDtls.FailedUsers = failedUserCount

		organizationSyncService, err := s.organizationSyncService(v)
		completedOrganizationCount, failedOrganizationCount := organizationSyncService.SyncOrganizations(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedOrganizations = completedOrganizationCount
		syncRunDtls.FailedOrganizations = failedOrganizationCount

		contactSyncService, err := s.contactSyncService(v)
		completedContactCount, failedContactCount := contactSyncService.SyncContacts(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedContacts = completedContactCount
		syncRunDtls.FailedContacts = failedContactCount

		ticketSyncService, err := s.ticketSyncService(v)
		completedTicketCount, failedTicketCount := ticketSyncService.SyncTickets(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedTickets = completedTicketCount
		syncRunDtls.FailedTickets = failedTicketCount

		noteSyncService, err := s.noteSyncService(v)
		completedNoteCount, failedNoteCount := noteSyncService.SyncNotes(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedNotes = completedNoteCount
		syncRunDtls.FailedNotes = failedNoteCount

		completedEmailMessageCount, failedEmailMessageCount := s.syncEmailMessages(ctx, dataService, syncDate, v.Tenant, runId)
		syncRunDtls.CompletedEmailMessages = completedEmailMessageCount
		syncRunDtls.FailedEmailMessages = failedEmailMessageCount

		syncRunDtls.TotalFailedEntities = failedUserCount + failedOrganizationCount + failedContactCount + failedTicketCount + failedNoteCount + failedEmailMessageCount
		syncRunDtls.TotalCompletedEntities = completedUserCount + completedOrganizationCount + completedContactCount + completedTicketCount + completedNoteCount + completedEmailMessageCount
		syncRunDtls.EndAt = time.Now().UTC()

		s.repositories.SyncRunRepository.Save(syncRunDtls)
	}
}

func (s *syncService) syncExternalSystem(ctx context.Context, dataService common.SourceDataService, tenant string) {
	_ = s.repositories.ExternalSystemRepository.Merge(ctx, tenant, dataService.SourceId())
}

func (s *syncService) syncEmailMessages(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		messages := dataService.GetEmailMessagesForSync(batchSize, runId)
		if len(messages) == 0 {
			logrus.Debugf("no email messages found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d email messages from %s for tenant %s", len(messages), dataService.SourceId(), tenant)

		for _, message := range messages {
			var failedSync = false
			var interactionEventId string

			sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tenant, syncDate, message)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge interaction session with external reference %v for tenant %v :%v", message.ExternalId, tenant, err)
			}

			if !failedSync {
				interactionEventId, err = s.repositories.InteractionEventRepository.MergeInteractionEvent(ctx, tenant, message.ExternalSystem, syncDate, message)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge interaction event with external reference %v for tenant %v :%v", message.ExternalId, tenant, err)
				}

				if !failedSync {
					err = s.repositories.InteractionEventRepository.MergeInteractionEventToSession(ctx, tenant, interactionEventId, sessionId)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed to associate interaction event to session %v for tenant %v :%v", message.ExternalId, tenant, err)
					}
				}
			}

			//from
			if message.Direction == entity.OUTBOUND && !failedSync {
				emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, message.FromEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed retrieving email %v for tenant %v :%v", message.FromEmail, tenant, err)
				}

				if emailId == "" && !failedSync {
					emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateUserByEmail(ctx, tenant, message.FromEmail, message.FromFirstName, message.FromLastName, message.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed creating contact with email %v for tenant %v :%v", message.FromEmail, tenant, err)
					}
				}

				if !failedSync {
					err := s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
					}
				}
			} else if message.Direction == entity.INBOUND && !failedSync {
				//1. find email ( contact/organization/user )
				//2. if not found, create contact with email

				emailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, message.FromEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed retrieving email %v for tenant %v :%v", message.FromEmail, tenant, err)
				}

				if emailId == "" && !failedSync {
					emailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateContactByEmail(ctx, tenant, message.FromEmail, message.FromFirstName, message.FromLastName, message.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed creating contact with email %v for tenant %v :%v", message.FromEmail, tenant, err)
					}
				}

				if !failedSync {
					err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, emailId)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed set sender for interaction event %v in tenant %v :%v", interactionEventId, tenant, err)
					}
				}
			}

			//to
			if len(message.ToEmail) > 0 && !failedSync {
				err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "TO", message.ToEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set TO users %v for interaction event %v in tenant %v :%v", message.ContactsExternalIds, sessionId, tenant, err)
				}
			}

			//cc
			if len(message.CcEmail) > 0 && !failedSync {
				err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "CC", message.CcEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set CC users %v for interaction event %v in tenant %v :%v", message.ContactsExternalIds, sessionId, tenant, err)
				}
			}

			//bcc
			if len(message.BccEmail) > 0 && !failedSync {
				err := s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "BCC", message.BccEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set BCC users %v for interaction event %v in tenant %v :%v", message.ContactsExternalIds, sessionId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged email message with external id %v to interaction session %v for tenant %v from %v", message.ExternalId, sessionId, tenant, dataService.SourceId())
			if err := dataService.MarkEmailMessageProcessed(message.ExternalId, runId, failedSync == false); err != nil {
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
	return completed, failed
}

func (s *syncService) sourceDataService(tenantToSync entity.TenantSyncSettings) (common.SourceDataService, error) {
	// Use a map to store the different implementations of common.SourceDataService as functions.
	dataServiceMap := map[string]func() common.SourceDataService{
		string(entity.AirbyteSourceHubspot): func() common.SourceDataService {
			return hubspot_service.NewHubspotDataService(s.repositories.Dbs.AirbyteStoreDB, tenantToSync.Tenant)
		},
		string(entity.AirbyteSourceZendeskSupport): func() common.SourceDataService {
			return zendesk_support_service.NewZendeskSupportDataService(s.repositories.Dbs.AirbyteStoreDB, tenantToSync.Tenant)
		},
		// Add additional implementations here.
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createDataService, ok := dataServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createDataService function to create a new instance of common.SourceDataService.
	dataService := createDataService()

	// Call the Refresh method on the sourceDataService instance.
	dataService.Refresh()

	return dataService, nil
}

func (s *syncService) userSyncService(tenantToSync entity.TenantSyncSettings) (UserSyncService, error) {
	userSyncServiceMap := map[string]func() UserSyncService{
		string(entity.AirbyteSourceHubspot): func() UserSyncService {
			return s.services.UserSyncService
		},
		string(entity.AirbyteSourceZendeskSupport): func() UserSyncService {
			return s.services.UserSyncService
		},
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createUserSyncService, ok := userSyncServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createUserSyncService function to create a new instance of common.SourceDataService.
	userSyncService := createUserSyncService()

	return userSyncService, nil
}

func (s *syncService) organizationSyncService(tenantToSync entity.TenantSyncSettings) (OrganizationSyncService, error) {
	organizationSyncServiceMap := map[string]func() OrganizationSyncService{
		string(entity.AirbyteSourceHubspot): func() OrganizationSyncService {
			return s.services.OrganizationSyncService
		},
		string(entity.AirbyteSourceZendeskSupport): func() OrganizationSyncService {
			return s.services.OrganizationSyncService
		},
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createOrganizationSyncService, ok := organizationSyncServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createOrganizationSyncService function to create a new instance of common.SourceDataService.
	organizationSyncService := createOrganizationSyncService()

	return organizationSyncService, nil
}

func (s *syncService) contactSyncService(tenantToSync entity.TenantSyncSettings) (ContactSyncService, error) {
	contactSyncServiceMap := map[string]func() ContactSyncService{
		string(entity.AirbyteSourceHubspot): func() ContactSyncService {
			return s.services.ContactSyncService
		},
		string(entity.AirbyteSourceZendeskSupport): func() ContactSyncService {
			return s.services.ContactSyncService
		},
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createContactSyncService, ok := contactSyncServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createContactSyncService function to create a new instance of common.SourceDataService.
	contactSyncService := createContactSyncService()

	return contactSyncService, nil
}

func (s *syncService) ticketSyncService(tenantToSync entity.TenantSyncSettings) (TicketSyncService, error) {
	ticketSyncServiceMap := map[string]func() TicketSyncService{
		string(entity.AirbyteSourceHubspot): func() TicketSyncService {
			return s.services.TicketSyncService
		},
		string(entity.AirbyteSourceZendeskSupport): func() TicketSyncService {
			return s.services.TicketSyncService
		},
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createTicketSyncService, ok := ticketSyncServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	ticketSyncService := createTicketSyncService()

	return ticketSyncService, nil
}

func (s *syncService) noteSyncService(tenantToSync entity.TenantSyncSettings) (NoteSyncService, error) {
	noteSyncServiceMap := map[string]func() NoteSyncService{
		string(entity.AirbyteSourceHubspot): func() NoteSyncService {
			return s.services.NoteSyncService
		},
		string(entity.AirbyteSourceZendeskSupport): func() NoteSyncService {
			return s.services.NoteSyncService
		},
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createNoteSyncService, ok := noteSyncServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	noteSyncService := createNoteSyncService()

	return noteSyncService, nil
}
