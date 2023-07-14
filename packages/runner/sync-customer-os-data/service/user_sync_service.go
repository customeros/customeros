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
	"strings"
	"time"
)

type userSyncService struct {
	repositories *repository.Repositories
}

func NewDefaultUserSyncService(repositories *repository.Repositories) SyncService {
	return &userSyncService{
		repositories: repositories,
	}
}

func (s *userSyncService) Sync(ctx context.Context, sourceService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		users := sourceService.GetDataForSync(common.USERS, batchSize, runId)
		if len(users) == 0 {
			logrus.Debugf("no users found for sync from %s for tenant %s", sourceService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d users from %s for tenant %s", len(users), sourceService.SourceId(), tenant)

		for _, v := range users {
			var failedSync = false
			var reason string
			userInput := v.(entity.UserData)
			userInput.Normalize()

			if userInput.Skip {
				if err := sourceService.MarkProcessed(userInput.SyncId, runId, true, true, "Input user marked for skip"); err != nil {
					failed++
					continue
				}
				skipped++
				continue
			}

			userInput.Email = strings.ToLower(userInput.Email)

			userId, err := s.repositories.UserRepository.GetMatchedUserId(ctx, tenant, userInput)
			if err != nil {
				failedSync = true
				reason = fmt.Sprintf("failed finding existing matched user with external reference %v for tenant %v :%v", userInput.ExternalId, tenant, err)
				logrus.Errorf(reason)
			}

			// Create new user id if not found
			if len(userId) == 0 {
				userUuid, _ := uuid.NewRandom()
				userId = userUuid.String()
			}
			userInput.Id = userId

			if !failedSync {
				err = s.repositories.UserRepository.MergeUser(ctx, tenant, syncDate, userInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merging user with external reference %v for tenant %v :%v", userInput.ExternalId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if userInput.HasEmail() && !failedSync {
				err = s.repositories.UserRepository.MergeEmail(ctx, tenant, userInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merging email for user with id %v for tenant %v :%v", userId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			if userInput.HasPhoneNumber() && !failedSync {
				err = s.repositories.UserRepository.MergePhoneNumber(ctx, tenant, userInput)
				if err != nil {
					failedSync = true
					reason = fmt.Sprintf("failed merging phone number for user with id %v for tenant %v :%v", userId, tenant, err)
					logrus.Errorf(reason)
				}
			}

			logrus.Debugf("successfully merged user with id %v for tenant %v from %v", userId, tenant, sourceService.SourceId())
			if err := sourceService.MarkProcessed(userInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
	return completed, failed, skipped
}
