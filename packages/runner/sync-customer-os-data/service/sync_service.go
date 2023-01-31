package service

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot/service"
	"github.com/sirupsen/logrus"
	"time"
)

const batchSize = 100

type SyncService interface {
	Sync(runId string)
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

func (s *syncService) Sync(runId string) {
	tenantsToSync, err := s.repositories.TenantSyncSettingsRepository.GetTenantsForSync()
	if err != nil {
		logrus.Error("failed to get tenants for sync")
		return
	}

	for _, v := range tenantsToSync {

		syncRunDtls := entity.SyncRun{
			StarAt:               time.Now().UTC(),
			RunId:                runId,
			TenantSyncSettingsId: v.ID,
		}

		dataService, err := s.dataService(v)
		if err != nil {
			logrus.Errorf("failed to get data service for tenant %v: %v", v.Tenant, err)
			continue
		}

		defer func() {
			dataService.Close()
		}()

		syncDate := time.Now().UTC()

		s.syncExternalSystem(dataService, v.Tenant)
		completedUserCount, failedUserCount := s.syncUsers(dataService, syncDate, v.Tenant, runId)
		completedOrganizationCount, failedOrganizationCount := s.syncOrganizations(dataService, syncDate, v.Tenant, runId)
		completedContactCount, failedContactCount := s.syncContacts(dataService, syncDate, v.Tenant, runId)
		completedNoteCount, failedNoteCount := s.syncNotes(dataService, syncDate, v.Tenant, runId)
		completedEmailMessageCount, failedEmailMessageCount := s.syncEmailMessages(dataService, syncDate, v.Tenant, runId)

		syncRunDtls.CompletedContacts = completedContactCount
		syncRunDtls.FailedContacts = failedContactCount
		syncRunDtls.CompletedUsers = completedUserCount
		syncRunDtls.FailedUsers = failedUserCount
		syncRunDtls.CompletedOrganizations = completedOrganizationCount
		syncRunDtls.FailedOrganizations = failedOrganizationCount
		syncRunDtls.CompletedNotes = completedNoteCount
		syncRunDtls.FailedNotes = failedNoteCount
		syncRunDtls.CompletedEmailMessages = completedEmailMessageCount
		syncRunDtls.FailedEmailMessages = failedEmailMessageCount

		syncRunDtls.EndAt = time.Now().UTC()

		s.repositories.SyncRunRepository.Save(syncRunDtls)
	}
}

func (s *syncService) syncExternalSystem(dataService common.DataService, tenant string) {
	_ = s.repositories.ExternalSystemRepository.Merge(tenant, dataService.SourceId())
}

func (s *syncService) syncContacts(dataService common.DataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		contacts := dataService.GetContactsForSync(batchSize, runId)
		if len(contacts) == 0 {
			logrus.Debugf("no contacts found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d contacts from %s for tenant %s", len(contacts), dataService.SourceId(), tenant)

		for _, v := range contacts {
			var failedSync = false

			contactId, err := s.repositories.ContactRepository.MergeContact(tenant, syncDate, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge contact with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			if len(v.PrimaryEmail) > 0 {
				if err = s.repositories.ContactRepository.MergePrimaryEmail(tenant, contactId, v.PrimaryEmail, v.ExternalSystem, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			for _, additionalEmail := range v.AdditionalEmails {
				if len(additionalEmail) > 0 {
					if err = s.repositories.ContactRepository.MergeAdditionalEmail(tenant, contactId, additionalEmail, v.ExternalSystem, v.CreatedAt); err != nil {
						failedSync = true
						logrus.Errorf("failed merge additional email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
					}
				}
			}

			if len(v.PrimaryE164) > 0 {
				if err = s.repositories.ContactRepository.MergePrimaryPhoneNumber(tenant, contactId, v.PrimaryE164, v.ExternalSystem, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary phone number for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			for _, organizationExternalId := range v.OrganizationsExternalIds {
				if err = s.repositories.RoleRepository.MergeRole(tenant, contactId, organizationExternalId, dataService.SourceId()); err != nil {
					failedSync = true
					logrus.Errorf("failed merge role for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if err = s.repositories.RoleRepository.RemoveOutdatedRoles(tenant, contactId, dataService.SourceId(), v.OrganizationsExternalIds); err != nil {
				failedSync = true
				logrus.Errorf("failed removing outdated roles for contact %v, tenant %v :%v", contactId, tenant, err)
			}

			if len(v.PrimaryOrganizationExternalId) > 0 {
				if err = s.repositories.RoleRepository.MergePrimaryRole(tenant, contactId, v.JobTitle, v.PrimaryOrganizationExternalId, dataService.SourceId()); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary role for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if len(v.UserExternalOwnerId) > 0 {
				if err = s.repositories.ContactRepository.SetOwnerRelationship(tenant, contactId, v.UserExternalOwnerId, dataService.SourceId()); err != nil {
					failedSync = true
					logrus.Errorf("failed set owner user for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			for _, f := range v.TextCustomFields {
				if err = s.repositories.ContactRepository.MergeTextCustomField(tenant, contactId, f, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge custom field %v for contact %v, tenant %v :%v", f.Name, contactId, tenant, err)
				}
			}

			err = s.repositories.ContactRepository.MergeContactDefaultPlace(tenant, contactId, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge address for contact %v, tenant %v :%v", contactId, tenant, err)
			}

			if len(v.ContactTypeName) > 0 {
				err = s.repositories.ContactRepository.MergeContactType(tenant, contactId, v.ContactTypeName)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge contact type for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged contact with id %v for tenant %v from %v", contactId, tenant, dataService.SourceId())
			if err := dataService.MarkContactProcessed(v.ExternalId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(contacts) < batchSize {
			break
		}
	}
	return completed, failed
}

func (s *syncService) syncOrganizations(dataService common.DataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		organizations := dataService.GetOrganizationsForSync(batchSize, runId)
		if len(organizations) == 0 {
			logrus.Debugf("no organizations found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d organizations from %s for tenant %s", len(organizations), dataService.SourceId(), tenant)

		for _, v := range organizations {
			var failedSync = false

			organizationId, err := s.repositories.OrganizationRepository.MergeOrganization(tenant, syncDate, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge organization with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			err = s.repositories.OrganizationRepository.MergeOrganizationDefaultPlace(tenant, organizationId, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge organization' address with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			if len(v.OrganizationTypeName) > 0 {
				err = s.repositories.OrganizationRepository.MergeOrganizationType(tenant, organizationId, v.OrganizationTypeName)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge organization type for organization %v, tenant %v :%v", organizationId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged organization with id %v for tenant %v from %v", organizationId, tenant, dataService.SourceId())
			if err := dataService.MarkOrganizationProcessed(v.ExternalId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(organizations) < batchSize {
			break
		}
	}
	return completed, failed
}

func (s *syncService) syncUsers(dataService common.DataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		users := dataService.GetUsersForSync(batchSize, runId)
		if len(users) == 0 {
			logrus.Debugf("no users found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d users from %s for tenant %s", len(users), dataService.SourceId(), tenant)

		for _, v := range users {
			var failedSync = false

			userId, err := s.repositories.UserRepository.MergeUser(tenant, syncDate, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merging user with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			logrus.Debugf("successfully merged user with id %v for tenant %v from %v", userId, tenant, dataService.SourceId())
			if err := dataService.MarkUserProcessed(v.ExternalOwnerId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(users) < batchSize {
			break
		}
	}
	return completed, failed
}

func (s *syncService) syncNotes(dataService common.DataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		notes := dataService.GetNotesForSync(batchSize, runId)
		if len(notes) == 0 {
			logrus.Debugf("no notes found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d notes from %s for tenant %s", len(notes), dataService.SourceId(), tenant)

		for _, note := range notes {
			var failedSync = false

			noteId, err := s.repositories.NoteRepository.MergeNote(tenant, syncDate, note)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge note with external reference %v for tenant %v :%v", note.ExternalId, tenant, err)
			}

			if len(noteId) > 0 {
				for _, contactExternalId := range note.ContactsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(tenant, noteId, contactExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
					}
				}

				for _, organizationExternalId := range note.OrganizationsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(tenant, noteId, organizationExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with organization for tenant %v :%v", noteId, tenant, err)
					}
				}

				if len(note.UserExternalId) > 0 {
					err = s.repositories.NoteRepository.NoteLinkWithUserByExternalId(tenant, noteId, note.UserExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
					}
				} else if len(note.UserExternalOwnerId) > 0 {
					err = s.repositories.NoteRepository.NoteLinkWithUserByExternalOwnerId(tenant, noteId, note.UserExternalOwnerId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
					}
				}
			}
			if failedSync == false {
				logrus.Debugf("successfully merged note with id %v for tenant %v from %v", noteId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkNoteProcessed(note.ExternalId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(notes) < batchSize {
			break
		}
	}
	return completed, failed
}

func (s *syncService) syncEmailMessages(dataService common.DataService, syncDate time.Time, tenant, runId string) (int, int) {
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

			var initiatorUsername = ""
			conversationId, messageCount, initiatorUsername, err := s.repositories.ConversationRepository.MergeEmailConversation(tenant, syncDate, message)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merge email message with external reference %v for tenant %v :%v", message.ExternalId, tenant, err)
			}

			var fromContactId string

			if message.Direction == entity.INBOUND {
				fromContactId, err = s.repositories.ContactRepository.GetOrCreateContactId(
					tenant, message.FromEmail, message.FromFirstName, message.FromLastName, message.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed creating contact with email %v for tenant %v :%v", message.FromEmail, tenant, err)
				}
			}

			// set initiator for new conversation
			if messageCount == 0 {
				initiatorUsername = message.FromEmail

				if message.Direction == entity.OUTBOUND {
					initiator := entity.ConversationInitiator{
						ExternalSystem: message.ExternalSystem,
						ExternalId:     message.UserExternalId,
						FirstName:      message.FromFirstName,
						LastName:       message.FromLastName,
						Email:          message.FromEmail,
						InitiatorType:  entity.USER,
					}
					err := s.repositories.ConversationRepository.UserInitiateConversation(tenant, conversationId, initiator)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed set user initiator for conversation %v in tenant %v :%v", conversationId, tenant, err)
					}
				} else if message.Direction == entity.INBOUND {
					initiator := entity.ConversationInitiator{
						Id:            fromContactId,
						FirstName:     message.FromFirstName,
						LastName:      message.FromLastName,
						Email:         message.FromEmail,
						InitiatorType: entity.CONTACT,
					}
					err := s.repositories.ConversationRepository.ContactInitiateConversation(tenant, conversationId, initiator)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed set contact initiator for conversation %v in tenant %v :%v", conversationId, tenant, err)
					}
				}
			}

			// set contact participants
			if len(fromContactId) > 0 {
				err := s.repositories.ConversationRepository.ContactByIdParticipateInConversation(tenant, conversationId, fromContactId)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set contact participat %v for conversation %v in tenant %v :%v", fromContactId, conversationId, tenant, err)
				}
			}
			if len(message.ContactsExternalIds) > 0 {
				err := s.repositories.ConversationRepository.ContactsByExternalIdParticipateInConversation(tenant, conversationId, message.ExternalSystem, message.ContactsExternalIds)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set contact participants by external id %v for conversation %v in tenant %v :%v", message.ContactsExternalIds, conversationId, tenant, err)
				}
			}

			// set user participants
			if len(message.UserExternalId) > 0 {
				err := s.repositories.ConversationRepository.UserByExternalIdParticipateInConversation(tenant, conversationId, message.ExternalSystem, message.UserExternalId)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set user participant by external id %v for conversation %v in tenant %v :%v", message.UserExternalId, conversationId, tenant, err)
				}
			}
			if len(message.ToEmail) > 0 && message.Direction == entity.INBOUND {
				err := s.repositories.ConversationRepository.UsersByEmailParticipateInConversation(tenant, conversationId, message.ToEmail)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set contact participants by external id %v for conversation %v in tenant %v :%v", message.ContactsExternalIds, conversationId, tenant, err)
				}
			}

			// increment message count
			if failedSync == false {
				err = s.repositories.ConversationRepository.IncrementMessageCount(tenant, conversationId, message.CreatedAt)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed set contact participants by external id %v for conversation %v in tenant %v :%v", message.ContactsExternalIds, conversationId, tenant, err)
				}
			}

			if failedSync == false {
				conversationEvent := entity.ConversationEvent{
					TenantName:        tenant,
					ConversationId:    conversationId,
					Type:              entity.EMAIL,
					Subtype:           message.EmailThreadId,
					Source:            message.ExternalSystem,
					ExternalId:        message.ExternalId,
					CreateDate:        message.CreatedAt,
					SenderUsername:    message.FromEmail,
					InitiatorUsername: initiatorUsername,
				}
				if message.Direction == entity.INBOUND {
					conversationEvent.Direction = entity.INBOUND
					conversationEvent.SenderType = entity.CONTACT
					conversationEvent.SenderId = fromContactId
				} else {
					conversationEvent.Direction = entity.OUTBOUND
					conversationEvent.SenderType = entity.USER
					userId, err := s.repositories.UserRepository.GetUserIdForExternalId(tenant, message.UserExternalId, message.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed to get user id for external id %v for tenant %v :%v", message.UserExternalId, tenant, err)
					}
					conversationEvent.SenderId = userId
				}
				emailContent := entity.EmailContent{
					MessageId: message.EmailMessageId,
					Subject:   message.Subject,
					Html:      message.Html,
					From:      message.FromEmail,
					To:        message.ToEmail,
					Cc:        message.CcEmail,
					Bcc:       message.BccEmail,
				}
				jsonContent, err := json.Marshal(emailContent)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed to marshal email content with external id %v for conversation %v in tenant %v :%v", message.ExternalId, conversationId, tenant, err)
				}
				if failedSync == false {
					conversationEvent.Content = string(jsonContent)
					err = s.repositories.ConversationEventRepository.Save(conversationEvent)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed save message with external id %v in message store for conversation %v in tenant %v :%v", message.ExternalId, conversationId, tenant, err)
					}
				}
			}

			logrus.Debugf("successfully merged email message with external id %v to conversation %v for tenant %v from %v", message.ExternalId, conversationId, tenant, dataService.SourceId())
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

func (s *syncService) dataService(tenantToSync entity.TenantSyncSettings) (common.DataService, error) {
	// Use a map to store the different implementations of common.DataService as functions.
	dataServiceMap := map[string]func() common.DataService{
		string(entity.AirbyteSourceHubspot): func() common.DataService {
			return service.NewHubspotDataService(s.repositories.Dbs.AirbyteStoreDB, tenantToSync.Tenant)
		},
		// Add additional implementations here.
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createDataService, ok := dataServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createDataService function to create a new instance of common.DataService.
	dataService := createDataService()

	// Call the Refresh method on the dataService instance.
	dataService.Refresh()

	return dataService, nil
}
