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
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type meetingSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultMeetingSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &meetingSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *meetingSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {

	completed, failed, skipped := 0, 0, 0

	for {

		meetings := dataService.GetDataForSync(ctx, common.MEETINGS, batchSize, runId)

		if len(meetings) == 0 {
			break
		}

		s.log.Infof("syncing %d meetings from %s for tenant %s", len(meetings), dataService.SourceId(), tenant)

		var wg sync.WaitGroup
		wg.Add(len(meetings))

		results := make(chan result, len(meetings))
		done := make(chan struct{})

		for _, v := range meetings {
			v := v

			go func(meeting entity.MeetingData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncMeeting(ctx, meeting, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.MeetingData))
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

		if len(meetings) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *meetingSyncService) syncMeeting(ctx context.Context, meetingInput entity.MeetingData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingSyncService.syncMeeting")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	meetingInput.Normalize()

	if meetingInput.Skip {
		if err := dataService.MarkProcessed(ctx, meetingInput.SyncId, runId, true, true, meetingInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	meetingId, err := s.repositories.MeetingRepository.GetMatchedMeetingId(ctx, tenant, meetingInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched meeting with external reference id %v for tenant %v :%v", meetingInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	// Create new meeting id if not found
	if len(meetingId) == 0 {
		meetingUuid, _ := uuid.NewRandom()
		meetingId = meetingUuid.String()
	}
	meetingInput.Id = meetingId
	span.LogFields(log.String("meetingId", meetingId))

	if !failedSync {
		err = s.repositories.MeetingRepository.MergeMeeting(ctx, tenant, syncDate, meetingInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge meeting with external reference %v for tenant %v :%v", meetingInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if meetingInput.HasLocation() && !failedSync {
		err = s.repositories.MeetingRepository.MergeMeetingLocation(ctx, tenant, meetingInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merge location for meeting %v for tenant %v :%v", meetingId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if meetingInput.HasUserCreator() && !failedSync {
		err = s.repositories.MeetingRepository.MeetingLinkWithCreatorUserByExternalId(ctx, tenant, meetingId, meetingInput.CreatorUserExternalId, meetingInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link meeting %v with user creator for tenant %v :%v", meetingId, tenant, err)
			s.log.Errorf(reason)
		}
		err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByUserByExternalId(ctx, tenant, meetingId, meetingInput.CreatorUserExternalId, meetingInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link meeting %v with user attended by for tenant %v :%v", meetingId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if meetingInput.HasContacts() && !failedSync {
		for _, contactExternalId := range meetingInput.ContactsExternalIds {
			err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByContactByExternalId(ctx, tenant, meetingId, contactExternalId, meetingInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link meeting %v with contact attended by for tenant %v :%v", meetingId, tenant, err)
				s.log.Errorf(reason)
				break
			}
			if !failedSync {
				s.services.OrganizationService.UpdateLastTouchpointByContactIdExternalId(ctx, tenant, contactExternalId, meetingInput.ExternalSystem)
			}
		}
	}
	if failedSync == false {
		s.log.Debugf("successfully merged meeting with id %v for tenant %v from %v", meetingId, tenant, dataService.SourceId())
	}
	if err = dataService.MarkProcessed(ctx, meetingInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
