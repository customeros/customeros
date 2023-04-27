package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type MeetingSyncService interface {
	SyncMeetings(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type meetingSyncService struct {
	repositories *repository.Repositories
}

func NewMeetingSyncService(repositories *repository.Repositories) MeetingSyncService {
	return &meetingSyncService{
		repositories: repositories,
	}
}

func (s *meetingSyncService) SyncMeetings(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		meetings := dataService.GetMeetingsForSync(batchSize, runId)
		if len(meetings) == 0 {
			logrus.Debugf("no meetings found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d meetings from %s for tenant %s", len(meetings), dataService.SourceId(), tenant)

		for _, meeting := range meetings {
			var failedSync = false

			meetingId, err := s.repositories.MeetingRepository.GetMatchedMeetingId(ctx, tenant, meeting)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched meeting with external reference id %v for tenant %v :%v", meeting.ExternalId, tenant, err)
			}

			// Create new meeting id if not found
			if len(meetingId) == 0 {
				meetingUuid, _ := uuid.NewRandom()
				meetingId = meetingUuid.String()
			}
			meeting.Id = meetingId

			if !failedSync {
				err = s.repositories.MeetingRepository.MergeMeeting(ctx, tenant, syncDate, meeting)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge meeting with external reference %v for tenant %v :%v", meeting.ExternalId, tenant, err)
				}
			}

			if meeting.HasLocation() && !failedSync {
				err = s.repositories.MeetingRepository.MergeMeetingLocation(ctx, tenant, meeting)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merge location for meeting %v for tenant %v :%v", meetingId, tenant, err)
				}
			}

			if meeting.HasUserCreator() && !failedSync {
				err = s.repositories.MeetingRepository.MeetingLinkWithCreatorUserByExternalId(ctx, tenant, meetingId, meeting.UserCreatorExternalId, meeting.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link meeting %v with user creator for tenant %v :%v", meetingId, tenant, err)
				}
				err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByUserByExternalId(ctx, tenant, meetingId, meeting.UserCreatorExternalId, meeting.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link meeting %v with user attended by for tenant %v :%v", meetingId, tenant, err)
				}
			}

			if meeting.HasContacts() && !failedSync {
				for _, contactExternalId := range meeting.ContactsExternalIds {
					err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByContactByExternalId(ctx, tenant, meetingId, contactExternalId, meeting.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link meeting %v with contact attended by for tenant %v :%v", meetingId, tenant, err)
					}
				}
			}
			if failedSync == false {
				logrus.Debugf("successfully merged meeting with id %v for tenant %v from %v", meetingId, tenant, dataService.SourceId())
			}
			if err := dataService.MarkMeetingProcessed(meeting.ExternalSyncId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(meetings) < batchSize {
			break
		}
	}
	return completed, failed
}
