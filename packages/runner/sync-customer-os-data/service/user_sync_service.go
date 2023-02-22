package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type UserSyncService interface {
	SyncUsers(dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type userSyncService struct {
	repositories *repository.Repositories
}

func NewUserSyncService(repositories *repository.Repositories) UserSyncService {
	return &userSyncService{
		repositories: repositories,
	}
}

func (s *userSyncService) SyncUsers(dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		users := dataService.GetUsersForSync(batchSize, runId)
		if len(users) == 0 {
			logrus.Debugf("no users found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d users from %s for tenant %s", len(users), dataService.SourceId(), tenant)

		for _, v := range users {
			var failedSync = false

			userId, err := s.repositories.UserRepository.GetMatchedUserId(tenant, *v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched user with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}
			// Create new user id if not found
			if len(userId) == 0 {
				userUuid, _ := uuid.NewRandom()
				userId = userUuid.String()
			}
			v.Id = userId

			err = s.repositories.UserRepository.MergeUser(tenant, syncDate, *v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed merging user with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			if len(v.Email) > 0 {
				err = s.repositories.UserRepository.MergeEmail(tenant, *v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merging email for user with id %v for tenant %v :%v", userId, tenant, err)
				}
			}

			if len(v.PhoneNumber) > 0 {
				err = s.repositories.UserRepository.MergePhoneNumber(tenant, *v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merging phone number for user with id %v for tenant %v :%v", userId, tenant, err)
				}
			}

			logrus.Debugf("successfully merged user with id %v for tenant %v from %v", userId, tenant, dataService.SourceId())
			if err := dataService.MarkUserProcessed(v.ExternalSyncId, runId, failedSync == false); err != nil {
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
	return completed, failed
}
