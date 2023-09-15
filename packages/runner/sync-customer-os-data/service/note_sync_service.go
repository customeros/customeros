package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type noteSyncService struct {
	repositories *repository.Repositories
	log          logger.Logger
}

func NewDefaultNoteSyncService(repositories *repository.Repositories, log logger.Logger) SyncService {
	return &noteSyncService{
		repositories: repositories,
		log:          log,
	}
}

func (s *noteSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	noteSyncMutex := &sync.Mutex{}

	for {

		notes := dataService.GetDataForSync(ctx, common.NOTES, batchSize, runId)

		if len(notes) == 0 {
			break
		}

		s.log.Infof("syncing %d notes from %s for tenant %s", len(notes), dataService.SourceId(), tenant)

		var wg sync.WaitGroup
		wg.Add(len(notes))

		results := make(chan result, len(notes))
		done := make(chan struct{})

		for _, v := range notes {
			v := v

			go func(note entity.NoteData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncNote(ctx, noteSyncMutex, note, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.NoteData))
		}
		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()
		go func() {
			<-done
			close(results)
		}()

		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}

		if len(notes) < batchSize {
			break
		}
	}

	return completed, failed, skipped
}

func (s *noteSyncService) syncNote(ctx context.Context, noteSyncMutex *sync.Mutex, noteInput entity.NoteData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteSyncService.syncNote")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	noteInput.Normalize()

	if noteInput.Skip {
		if err := dataService.MarkProcessed(ctx, noteInput.SyncId, runId, true, true, noteInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	noteSyncMutex.Lock()
	noteId, err := s.repositories.NoteRepository.GetMatchedNoteId(ctx, tenant, noteInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched note with external reference id %v for tenant %v :%v", noteInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	// Create new note id if not found
	if noteId == "" {
		noteUuid, _ := uuid.NewRandom()
		noteId = noteUuid.String()
	}
	noteInput.Id = noteId
	span.LogFields(log.String("noteId", noteId))

	if !failedSync {
		err = s.repositories.NoteRepository.MergeNote(ctx, tenant, syncDate, noteInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge note with external reference %v for tenant %v :%v", noteInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}
	noteSyncMutex.Unlock()

	if noteInput.HasNotedContacts() && !failedSync {
		for _, contactExternalId := range noteInput.NotedContactsExternalIds {
			err = s.repositories.NoteRepository.NoteLinkWithContactByExternalId(ctx, tenant, noteId, contactExternalId, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with contact for tenant %v :%v", noteId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if noteInput.HasNotedOrganizations() && !failedSync {
		for _, organizationExternalId := range noteInput.NotedOrganizationsExternalIds {
			err = s.repositories.NoteRepository.NoteLinkWithOrganizationByExternalId(ctx, tenant, noteId, organizationExternalId, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with organization for tenant %v :%v", noteId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if noteInput.HasMentionedTags() && !failedSync {
		for _, tagName := range noteInput.MentionedTags {
			err = s.repositories.NoteRepository.NoteMentionedTag(ctx, tenant, noteId, tagName, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with tag %s for tenant %v :%v", noteId, tagName, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if !failedSync && noteInput.HasMentionedIssue() {
		issue, err := s.repositories.IssueRepository.GetMatchedIssue(ctx, tenant, noteInput.ExternalSystem, noteInput.MentionedIssueExternalId)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed finding existing matched issue with external reference id %v for tenant %v :%v", noteInput.MentionedIssueExternalId, tenant, err)
			s.log.Errorf(reason)
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
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed link note %v with tag %s for tenant %v :%v", noteId, tagName, tenant, err)
					s.log.Errorf(reason)
				}
			}
			err := s.repositories.NoteRepository.NoteLinkWithIssueReporterContactOrOrganization(ctx, tenant, noteId, issueId, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with issue %s reporter for tenant %v : %v", noteId, issueId, tenant, err)
				s.log.Errorf(reason)
			}
		}
	}

	if noteInput.HasCreator() && !failedSync {
		err = s.repositories.NoteRepository.NoteLinkWithCreatorByExternalId(ctx, tenant, noteId, noteInput.CreatorExternalId, dataService.SourceId())
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link note %v with user or contact for tenant %v :%v", noteId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if !failedSync {
		if noteInput.HasCreatorUser() {
			err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalId(ctx, tenant, noteId, noteInput.CreatorUserExternalId, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
				s.log.Errorf(reason)
			}
		} else if noteInput.HasCreatorUserOwner() {
			err = s.repositories.NoteRepository.NoteLinkWithCreatorUserByExternalOwnerId(ctx, tenant, noteId, noteInput.CreatorUserExternalOwnerId, dataService.SourceId())
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link note %v with user for tenant %v :%v", noteId, tenant, err)
				s.log.Errorf(reason)
			}
		}
	}
	if failedSync == false {
		s.log.Debugf("successfully merged note with id %v for tenant %v from %v", noteId, tenant, dataService.SourceId())
	}
	if err := dataService.MarkProcessed(ctx, noteInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
