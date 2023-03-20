package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type OrganizationSyncService interface {
	SyncOrganizations(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type organizationSyncService struct {
	repositories *repository.Repositories
}

func NewOrganizationSyncService(repositories *repository.Repositories) OrganizationSyncService {
	return &organizationSyncService{
		repositories: repositories,
	}
}

func (s *organizationSyncService) SyncOrganizations(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
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
			utils.LowercaseStrings(v.Domains)

			organizationId, err := s.repositories.OrganizationRepository.GetMatchedOrganizationId(ctx, tenant, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched organization with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			// Create new organization id if not found
			if len(organizationId) == 0 {
				orgUuid, _ := uuid.NewRandom()
				organizationId = orgUuid.String()
			}
			v.Id = organizationId

			if !failedSync {
				err = s.repositories.OrganizationRepository.MergeOrganization(ctx, tenant, syncDate, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge organization with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasDomains() && !failedSync {
				for _, domain := range v.Domains {
					err = s.repositories.OrganizationRepository.MergeOrganizationDomain(ctx, tenant, organizationId, domain, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed merge organization domain for organization %v, tenant %v :%v", organizationId, tenant, err)
					}
				}
			}

			if v.HasPhoneNumber() && !failedSync {
				if err = s.repositories.OrganizationRepository.MergePrimaryPhoneNumber(ctx, tenant, organizationId, v.PhoneNumber, v.ExternalSystem, v.CreatedAt); err != nil {
					failedSync = true
					logrus.Errorf("failed merge primary phone number for organization with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasLocation() && !failedSync {
				err = s.repositories.OrganizationRepository.MergeOrganizationDefaultPlace(ctx, tenant, organizationId, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge organization' place with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasNotes() && !failedSync {
				note := entity.NoteData{
					Html:           v.NoteContent,
					CreatedAt:      v.CreatedAt,
					ExternalId:     v.ExternalId,
					ExternalSystem: v.ExternalSystem,
				}
				noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, note)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed finding existing matched note with external reference id %v for tenant %v :%v", note.ExternalId, tenant, err)
				}
				// Create new note id if not found
				if len(noteId) == 0 {
					noteUuid, _ := uuid.NewRandom()
					noteId = noteUuid.String()
				}
				note.Id = noteId
				err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, note)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge organization note for organization %v, tenant %v :%v", organizationId, tenant, err)
				}
				err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(ctx, tenant, noteId, v.ExternalId, v.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link note with organization %v, tenant %v :%v", organizationId, tenant, err)
				}
			}

			if v.HasOrganizationType() && !failedSync {
				err = s.repositories.OrganizationRepository.MergeOrganizationType(ctx, tenant, organizationId, v.OrganizationTypeName)
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
