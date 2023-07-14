package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/sirupsen/logrus"
	"time"
)

type meetingSyncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewDefaultMeetingSyncService(repositories *repository.Repositories, services *Services) SyncService {
	return &meetingSyncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *meetingSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		meetings := dataService.GetDataForSync(common.MEETINGS, batchSize, runId)
		if len(meetings) == 0 {
			logrus.Debugf("no meetings found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d meetings from %s for tenant %s", len(meetings), dataService.SourceId(), tenant)

		for _, v := range meetings {
			var failedSync = false
			var reason string

			meetingInput := v.(entity.MeetingData)
			meetingInput.Normalize()

			if meetingInput.Skip {
				if err := dataService.MarkProcessed(meetingInput.SyncId, runId, true, true, meetingInput.SkipReason); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			meetingId, err := s.repositories.MeetingRepository.GetMatchedMeetingId(ctx, tenant, meetingInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed finding existing matched meeting with external reference id %v for tenant %v :%v", meetingInput.ExternalId, tenant, err)
				logrus.Errorf(reason)
			}

			// Create new meeting id if not found
			if len(meetingId) == 0 {
				meetingUuid, _ := uuid.NewRandom()
				meetingId = meetingUuid.String()
			}
			meetingInput.Id = meetingId

			if !failedSync {
				err = s.repositories.MeetingRepository.MergeMeeting(ctx, tenant, syncDate, meetingInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge meeting with external reference %v for tenant %v :%v", meetingInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if meetingInput.HasLocation() && !failedSync {
				err = s.repositories.MeetingRepository.MergeMeetingLocation(ctx, tenant, meetingInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merge location for meeting %v for tenant %v :%v", meetingId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if meetingInput.HasUserCreator() && !failedSync {
				err = s.repositories.MeetingRepository.MeetingLinkWithCreatorUserByExternalId(ctx, tenant, meetingId, meetingInput.CreatorUserExternalId, meetingInput.ExternalSystem)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed link meeting %v with user creator for tenant %v :%v", meetingId, tenant, err)
					logrus.Errorf(reason)
				}
				err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByUserByExternalId(ctx, tenant, meetingId, meetingInput.CreatorUserExternalId, meetingInput.ExternalSystem)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed link meeting %v with user attended by for tenant %v :%v", meetingId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if meetingInput.HasContacts() && !failedSync {
				for _, contactExternalId := range meetingInput.ContactsExternalIds {
					err = s.repositories.MeetingRepository.MeetingLinkWithAttendedByContactByExternalId(ctx, tenant, meetingId, contactExternalId, meetingInput.ExternalSystem)
					if err != nil {
						failedSync = true
						reason = fmt.Sprintf("failed link meeting %v with contact attended by for tenant %v :%v", meetingId, tenant, err)
						logrus.Errorf(reason)
						break
					}
					if !failedSync {
						s.services.OrganizationService.UpdateLastTouchpointByContactIdExternalId(ctx, tenant, contactExternalId, meetingInput.ExternalSystem)
					}
				}
			}
			if failedSync == false {
				logrus.Debugf("successfully merged meeting with id %v for tenant %v from %v", meetingId, tenant, dataService.SourceId())
			}
			if err = dataService.MarkProcessed(meetingInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
	return completed, failed, skipped
}
