package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ContactSyncService interface {
	SyncContacts(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int)
}

type contactSyncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewContactSyncService(repositories *repository.Repositories, services *Services) ContactSyncService {
	return &contactSyncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *contactSyncService) SyncContacts(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int) {
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
			v.PrimaryEmail = strings.ToLower(v.PrimaryEmail)
			utils.LowercaseStrings(v.AdditionalEmails)

			contactId, err := s.repositories.ContactRepository.GetMatchedContactId(ctx, tenant, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched contact with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			// Create new contact id if not found
			if len(contactId) == 0 {
				contactUuid, _ := uuid.NewRandom()
				contactId = contactUuid.String()
			}
			v.Id = contactId

			if !failedSync {
				err = s.repositories.ContactRepository.MergeContact(ctx, tenant, syncDate, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge contact with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if len(v.PrimaryEmail) > 0 && !failedSync {
				if err = s.repositories.ContactRepository.MergePrimaryEmail(ctx, tenant, contactId, v.PrimaryEmail, v.ExternalSystem, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if !failedSync {
				for _, additionalEmail := range v.AdditionalEmails {
					if len(additionalEmail) > 0 {
						if err = s.repositories.ContactRepository.MergeAdditionalEmail(ctx, tenant, contactId, additionalEmail, v.ExternalSystem, v.CreatedAt); err != nil {
							failedSync = true
							logrus.Errorf("failed merge additional email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
						}
					}
				}
			}

			if v.HasPhoneNumber() && !failedSync {
				if err = s.repositories.ContactRepository.MergePrimaryPhoneNumber(ctx, tenant, contactId, v.PhoneNumber, v.ExternalSystem, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary phone number for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasOrganizations() && !failedSync {
				for _, organizationExternalId := range v.OrganizationsExternalIds {
					if organizationExternalId != "" {
						if err = s.repositories.ContactRepository.LinkContactWithOrganization(ctx, tenant, contactId, organizationExternalId, dataService.SourceId(), v.CreatedAt); err != nil {
							failedSync = true
							logrus.Errorf("failed link contact %v to organization with external id %v, tenant %v :%v", contactId, organizationExternalId, tenant, err)
						}
					}
				}
			}

			if !failedSync {
				if err = s.repositories.RoleRepository.RemoveOutdatedJobRoles(ctx, tenant, contactId, dataService.SourceId(), v.PrimaryOrganizationExternalId); err != nil {
					failedSync = true
					logrus.Errorf("failed removing outdated roles for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if len(v.PrimaryOrganizationExternalId) > 0 && !failedSync {
				if err = s.repositories.RoleRepository.MergeJobRole(ctx, tenant, contactId, v.JobTitle, v.PrimaryOrganizationExternalId, dataService.SourceId(), v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary role for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if v.HasOwner() && !failedSync {
				if err = s.repositories.ContactRepository.SetOwner(ctx, tenant, contactId, v.UserExternalOwnerId, dataService.SourceId()); err != nil {
					// Do not mark sync as failed in case owner relationship is not set
					logrus.Errorf("failed set owner user for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if v.HasNotes() && !failedSync {
				for _, note := range v.Notes {
					localNote := entity.NoteData{
						Html:           note.Note,
						CreatedAt:      common_utils.TimePtr(v.CreatedAt),
						ExternalId:     string(note.FieldSource) + "-" + v.ExternalId,
						ExternalSystem: v.ExternalSystem,
					}
					noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, localNote)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed finding existing matched note with external reference id %v for tenant %v :%v", localNote.ExternalId, tenant, err)
					}
					// Create new note id if not found
					if len(noteId) == 0 {
						noteUuid, _ := uuid.NewRandom()
						noteId = noteUuid.String()
					}
					localNote.Id = noteId
					err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, localNote)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed merge note for contact %v, tenant %v :%v", contactId, tenant, err)
					}
					err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, v.ExternalId, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note with contact %v, tenant %v :%v", contactId, tenant, err)
					}
				}
			}

			if v.HasTextCustomFields() && !failedSync {
				for _, customField := range v.TextCustomFields {
					if err = s.repositories.ContactRepository.MergeTextCustomField(ctx, tenant, contactId, customField); err != nil {
						failedSync = true
						logrus.Errorf("failed merge custom field %v for contact %v, tenant %v :%v", customField.Name, contactId, tenant, err)
					}
				}
			}

			if v.HasLocation() && !failedSync {
				err = s.repositories.ContactRepository.MergeContactLocation(ctx, tenant, contactId, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge location for contact %v, tenant %v :%v", contactId, tenant, err)
				}
			}

			if v.HasTags() && !failedSync {
				for _, tag := range v.Tags {
					err = s.repositories.ContactRepository.MergeTagForContact(ctx, tenant, contactId, tag, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed to merge tag %v for contact %v, tenant %v :%v", tag, contactId, tenant, err)
					}
				}
			}

			s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, tenant, contactId)

			logrus.Debugf("successfully merged contact with id %v for tenant %v from %v", contactId, tenant, dataService.SourceId())
			if err = dataService.MarkContactProcessed(v.ExternalSyncId, runId, failedSync == false); err != nil {
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
