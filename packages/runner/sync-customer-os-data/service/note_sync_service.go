package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type NoteSyncService interface {
	SyncNotes(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type noteSyncService struct {
	repositories *repository.Repositories
}

func NewNoteSyncService(repositories *repository.Repositories) NoteSyncService {
	return &noteSyncService{
		repositories: repositories,
	}
}

func (s *noteSyncService) SyncNotes(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
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

			if !failedSync {
				err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, note)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge note with external reference %v for tenant %v :%v", note.ExternalId, tenant, err)
				}
			}

			if note.HasNotedContacts() && !failedSync {
				for _, contactExternalId := range note.NotedContactsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, contactExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
					}
				}
			}

			if note.HasNotedOrganizations() && !failedSync {
				for _, organizationExternalId := range note.NotedOrganizationsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(ctx, tenant, noteId, organizationExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with organization for tenant %v :%v", noteId, tenant, err)
					}
				}
			}

			if note.HasNotedTickets() && !failedSync {
				for _, ticketExternalId := range note.NotedTicketsExternalIds {
					err = s.repositories.NoteRepository.NoteLinkWithTicketByExternalId(ctx, tenant, noteId, ticketExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
					}
				}
			}

			if note.HasCreatorUserOrContact() && !failedSync {
				err = s.repositories.NoteRepository.NoteLinkWithCreatorUserOrContactByExternalId(ctx, tenant, noteId, note.CreatorUserOrContactExternalId, dataService.SourceId())
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link note %v with user or contact for tenant %v :%v", noteId, tenant, err)
				}
			}

			if !failedSync {
				if note.HasCreatorUser() {
					err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalId(ctx, tenant, noteId, note.CreatorUserExternalId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
					}
				} else if note.HasCreatorUserOwner() {
					err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalOwnerId(ctx, tenant, noteId, note.CreatorUserExternalOwnerId, dataService.SourceId())
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
					}
				}
			}
			if failedSync == false {
				logrus.Debugf("successfully merged note with id %v for tenant %v from %v", noteId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkNoteProcessed(note.ExternalSyncId, runId, failedSync == false); err != nil {
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
