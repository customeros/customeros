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
	"time"
)

type noteSyncService struct {
	repositories *repository.Repositories
}

func NewDefaultNoteSyncService(repositories *repository.Repositories) SyncService {
	return &noteSyncService{
		repositories: repositories,
	}
}

func (s *noteSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		notes := dataService.GetDataForSync(common.NOTES, batchSize, runId)
		if len(notes) == 0 {
			logrus.Debugf("no notes found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d notes from %s for tenant %s", len(notes), dataService.SourceId(), tenant)

		for _, v := range notes {
			var failedSync = false
			var reason string

			noteInput := v.(entity.NoteData)
			noteInput.Normalize()

			if noteInput.Skip {
				if err := dataService.MarkProcessed(noteInput.SyncId, runId, true, true, noteInput.SkipReason); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, noteInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed finding existing matched note with external reference id %v for tenant %v :%v", noteInput.ExternalId, tenant, err)
				logrus.Errorf(reason)
			}

			// Create new note id if not found
			if noteId == "" {
				noteUuid, _ := uuid.NewRandom()
				noteId = noteUuid.String()
			}
			noteInput.Id = noteId

			if !failedSync {
				err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, noteInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge note with external reference %v for tenant %v :%v", noteInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if noteInput.HasNotedContacts() && !failedSync {
				for _, contactExternalId := range noteInput.NotedContactsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, contactExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			if noteInput.HasNotedOrganizations() && !failedSync {
				for _, organizationExternalId := range noteInput.NotedOrganizationsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(ctx, tenant, noteId, organizationExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with organization for tenant %v :%v", noteId, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			if noteInput.HasMentionedTags() && !failedSync {
				for _, tagName := range noteInput.MentionedTags {
					err = s.repositories.NoteRepository.NoteMentionedTag(ctx, tenant, noteId, tagName, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with tag %s for tenant %v :%v", noteId, tagName, tenant, err)
						logrus.Errorf(reason)
						break
					}
				}
			}

			if !failedSync && noteInput.HasMentionedIssue() {
				issue, err := s.repositories.IssueRepository.GetMatchedIssue(ctx, tenant, noteInput.ExternalSystem, noteInput.MentionedIssueExternalId)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed finding existing matched issue with external reference id %v for tenant %v :%v", noteInput.MentionedIssueExternalId, tenant, err)
					logrus.Errorf(reason)
				}
				if issue != nil {
					props := issue.Props
					subject := utils.GetStringPropOrEmpty(props, "subject")
					issueId := utils.GetStringPropOrEmpty(props, "id")
					if subject != "" {
						tagName := fmt.Sprintf("%s - %s", subject, noteInput.MentionedIssueExternalId)
						err = s.repositories.NoteRepository.NoteMentionedTag(ctx, tenant, noteId, tagName, dataService.SourceId())
						if err != nil {
							failedSync = true
							reason = fmt.Sprintf("failed link note %v with tag %s for tenant %v :%v", noteId, tagName, tenant, err)
							logrus.Errorf(reason)
							break
						}
					}
					err := s.repositories.NoteRepository.NoteLinkWithIssueReporterContactOrOrganization(ctx, tenant, noteId, issueId, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with issue %s reporter for tenant %v : %v", noteId, issueId, tenant, err)
						logrus.Errorf(reason)
					}
				}
			}

			if noteInput.HasCreator() && !failedSync {
				err = s.repositories.NoteRepository.NoteLinkWithCreatorByExternalId(ctx, tenant, noteId, noteInput.CreatorExternalId, dataService.SourceId())
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed link note %v with user or contact for tenant %v :%v", noteId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if !failedSync {
				if noteInput.HasCreatorUser() {
					err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalId(ctx, tenant, noteId, noteInput.CreatorUserExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
						logrus.Errorf(reason)
					}
				} else if noteInput.HasCreatorUserOwner() {
					err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalOwnerId(ctx, tenant, noteId, noteInput.CreatorUserExternalOwnerId, dataService.SourceId())
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
						logrus.Errorf(reason)
					}
				}
			}
			if failedSync == false {
				logrus.Debugf("successfully merged note with id %v for tenant %v from %v", noteId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkProcessed(noteInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
	return completed, failed, skipped
}
