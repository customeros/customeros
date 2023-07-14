package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type contactSyncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewDefaultContactSyncService(repositories *repository.Repositories, services *Services) SyncService {
	return &contactSyncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *contactSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		contacts := dataService.GetDataForSync(common.CONTACTS, batchSize, runId)
		if len(contacts) == 0 {
			logrus.Debugf("no contacts found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d contacts from %s for tenant %s", len(contacts), dataService.SourceId(), tenant)

		for _, v := range contacts {
			var failedSync = false
			var reason string
			contactInput := v.(entity.ContactData)
			contactInput.Normalize()

			if contactInput.Skip {
				if err := dataService.MarkProcessed(contactInput.SyncId, runId, true, true, contactInput.SkipReason); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			contactInput.Email = strings.ToLower(contactInput.Email)
			utils.LowercaseStrings(contactInput.AdditionalEmails)

			contactId, err := s.repositories.ContactRepository.GetMatchedContactId(ctx, tenant, contactInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed finding existing matched contactInput with external reference %v for tenant %v :%v", contactInput.ExternalId, tenant, err)
				logrus.Errorf(reason)
			}

			// Create new contactInput id if not found
			if len(contactId) == 0 {
				contactUuid, _ := uuid.NewRandom()
				contactId = contactUuid.String()
			}
			contactInput.Id = contactId

			if !failedSync {
				err = s.repositories.ContactRepository.MergeContact(ctx, tenant, syncDate, contactInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge contactInput with external reference %v for tenant %v :%v", contactInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if len(contactInput.Email) > 0 && !failedSync {
				if err = s.repositories.ContactRepository.MergePrimaryEmail(ctx, tenant, contactId, contactInput.Email, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge primary email for contactInput with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if !failedSync {
				for _, additionalEmail := range contactInput.AdditionalEmails {
					if len(additionalEmail) > 0 {
						if err = s.repositories.ContactRepository.MergeAdditionalEmail(ctx, tenant, contactId, additionalEmail, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
							failedSync = true
							reason = fmt.Sprintf("failed merge additional email for contactInput with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
							logrus.Errorf(reason)
							break
						}
					}
				}
			}

			if contactInput.HasPhoneNumber() && !failedSync {
				if err = s.repositories.ContactRepository.MergePrimaryPhoneNumber(ctx, tenant, contactId, contactInput.PhoneNumber, contactInput.ExternalSystem, *contactInput.CreatedAt); err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge primary phone number for contactInput with external reference %v , tenant %v :%v", contactInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if contactInput.HasOrganizations() && !failedSync {
				for _, organizationExternalId := range contactInput.OrganizationsExternalIds {
					if organizationExternalId != "" {
						if err = s.repositories.ContactRepository.LinkContactWithOrganization(ctx, tenant, contactId, organizationExternalId, dataService.SourceId(), *contactInput.CreatedAt); err != nil {
							failedSync = true
							reason = fmt.Sprintf("failed link contactInput %v to organization with external id %v, tenant %v :%v", contactId, organizationExternalId, tenant, err)
							logrus.Errorf(reason)
							break
						}
					}
				}
			}

			if !failedSync {
				if err = s.repositories.RoleRepository.RemoveOutdatedJobRoles(ctx, tenant, contactId, dataService.SourceId(), contactInput.PrimaryOrganizationExternalId); err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed removing outdated roles for contactInput %v, tenant %v :%v", contactId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if len(contactInput.PrimaryOrganizationExternalId) > 0 && !failedSync {
				if err = s.repositories.RoleRepository.MergeJobRole(ctx, tenant, contactId, contactInput.JobTitle, contactInput.PrimaryOrganizationExternalId, dataService.SourceId(), *contactInput.CreatedAt); err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge primary role for contactInput %v, tenant %v :%v", contactId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if contactInput.HasOwner() && !failedSync {
				if err = s.repositories.ContactRepository.SetOwner(ctx, tenant, contactId, contactInput.UserExternalOwnerId, dataService.SourceId()); err != nil {
					// Do not mark sync as failed in case owner relationship is not set
					logrus.Errorf("failed set owner user for contactInput %s, tenant %s :%s", contactId, tenant, err)
				}
			}

			if contactInput.HasNotes() && !failedSync {
				for _, note := range contactInput.Notes {
					localNote := entity.NoteData{
						BaseData: entity.BaseData{
							CreatedAt:      contactInput.CreatedAt,
							ExternalId:     string(note.FieldSource) + "-" + contactInput.ExternalId,
							ExternalSystem: contactInput.ExternalSystem,
						},
						Html: note.Note,
					}
					noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, localNote)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed finding existing matched note with external reference id %s for tenant %s :%s", localNote.ExternalId, tenant, err)
						logrus.Errorf(reason)
						break
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
						reason = fmt.Sprintf("failed merge note for contactInput %v, tenant %v :%v", contactId, tenant, err)
						logrus.Errorf(reason)
						break
					}
					err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, contactInput.ExternalId, contactInput.ExternalSystem)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note with contactInput %v, tenant %v :%v", contactId, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			if contactInput.HasTextCustomFields() && !failedSync {
				for _, customField := range contactInput.TextCustomFields {
					if err = s.repositories.ContactRepository.MergeTextCustomField(ctx, tenant, contactId, customField); err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed merge custom field %v for contactInput %v, tenant %v :%v", customField.Name, contactId, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			if contactInput.HasLocation() && !failedSync {
				err = s.repositories.ContactRepository.MergeContactLocation(ctx, tenant, contactId, contactInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge location for contactInput %v, tenant %v :%v", contactId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if contactInput.HasTags() && !failedSync {
				for _, tag := range contactInput.Tags {
					err = s.repositories.ContactRepository.MergeTagForContact(ctx, tenant, contactId, tag, contactInput.ExternalSystem)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed to merge tag %v for contactInput %v, tenant %v :%v", tag, contactId, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, tenant, contactId)

			logrus.Debugf("successfully merged contactInput with id %v for tenant %v from %v", contactId, tenant, dataService.SourceId())
			if err = dataService.MarkProcessed(contactInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
	return completed, failed, skipped
}
