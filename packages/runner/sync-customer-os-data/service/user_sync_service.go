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
	"strings"
	"sync"
	"time"
)

type userSyncService struct {
	repositories *repository.Repositories
	log          logger.Logger
}

func NewDefaultUserSyncService(repositories *repository.Repositories, log logger.Logger) SyncService {
	return &userSyncService{
		repositories: repositories,
		log:          log,
	}
}

func (s *userSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	userSyncMutex := &sync.Mutex{}

	for {
		users := dataService.GetDataForSync(ctx, common.USERS, batchSize, runId)
		if len(users) == 0 {
			break
		}
		s.log.Infof("syncing %d users from %s for tenant %s", len(users), dataService.SourceId(), tenant)

		var wg sync.WaitGroup
		wg.Add(len(users))

		// Channel to collect results
		results := make(chan result, len(users))
		done := make(chan struct{})

		for _, v := range users {
			v := v
			go func(user entity.UserData) {
				defer wg.Done()
				var comp, fail, skip int
				s.syncUser(ctx, userSyncMutex, v.(entity.UserData), dataService, syncDate, tenant, runId, &comp, &fail, &skip)
				results <- result{comp, fail, skip}
			}(v.(entity.UserData))
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

		// Collect results
		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}
		if len(users) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *userSyncService) syncUser(ctx context.Context, userSyncMutex *sync.Mutex, userInput entity.UserData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserSyncService.syncUser")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	userInput.Normalize()

	if userInput.Skip {
		if err := dataService.MarkProcessed(ctx, userInput.SyncId, runId, true, true, userInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	userInput.Email = strings.ToLower(userInput.Email)

	userSyncMutex.Lock()
	userId, err := s.repositories.UserRepository.GetMatchedUserId(ctx, tenant, userInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched user with external reference %v for tenant %v :%v", userInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	// Create new user id if not found
	if len(userId) == 0 {
		userUuid, _ := uuid.NewRandom()
		userId = userUuid.String()
	}
	userInput.Id = userId
	span.LogFields(log.String("userId", userId))

	if !failedSync {
		err = s.repositories.UserRepository.MergeUser(ctx, tenant, syncDate, userInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging user with external reference %v for tenant %v :%v", userInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}
	userSyncMutex.Unlock()

	if userInput.HasEmail() && !failedSync {
		err = s.repositories.UserRepository.MergeEmail(ctx, tenant, userInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging email for user with id %v for tenant %v :%v", userId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if userInput.HasPhoneNumbers() && !failedSync {
		for _, phoneNumber := range userInput.PhoneNumbers {
			err = s.repositories.UserRepository.MergePhoneNumber(ctx, tenant, userInput, phoneNumber)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed merging phone number for user with id %v for tenant %v :%v", userId, tenant, err)
				s.log.Errorf(reason)
			}
		}
	}

	s.log.Debugf("successfully merged user with id %v for tenant %v from %v", userId, tenant, dataService.SourceId())
	if err := dataService.MarkProcessed(ctx, userInput.SyncId, runId, failedSync == false, false, reason); err != nil {
		tracing.TraceErr(span, err)
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
