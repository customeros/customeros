package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/service"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
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

		syncRunDtls.CompletedContacts = completedContactCount
		syncRunDtls.FailedContacts = failedContactCount
		syncRunDtls.CompletedUsers = completedUserCount
		syncRunDtls.FailedUsers = failedUserCount
		syncRunDtls.CompletedOrganizations = completedOrganizationCount
		syncRunDtls.FailedOrganizations = failedOrganizationCount
		syncRunDtls.CompletedNotes = completedNoteCount
		syncRunDtls.FailedNotes = failedNoteCount

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

			err = s.repositories.ContactRepository.MergeContactAddress(tenant, contactId, v)
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

			err = s.repositories.OrganizationRepository.MergeOrganizationAddress(tenant, organizationId, v)
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

			for _, contactExternalId := range note.ContactsExternalIds {
				err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(tenant, noteId, contactExternalId, dataService.SourceId())
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
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

			logrus.Debugf("successfully merged note with id %v for tenant %v from %v", noteId, tenant, dataService.SourceId())
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

func (s *syncService) dataService(tenantToSync entity.TenantSyncSettings) (common.DataService, error) {
	// Use a map to store the different implementations of common.DataService as functions.
	dataServiceMap := map[entity.AirbyteSource]func() common.DataService{
		entity.HUBSPOT: func() common.DataService {
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
