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
	span, ctx := tracing.StartTracerSpan(ctx, "UserSyncService.Sync")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	completed, failed, skipped := 0, 0, 0
	for {
		users := dataService.GetDataForSync(ctx, common.USERS, batchSize, runId)
		if len(users) == 0 {
			s.log.Debugf("no users found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		s.log.Infof("syncing %d users from %s for tenant %s", len(users), dataService.SourceId(), tenant)

		for _, v := range users {
			s.syncUser(ctx, v.(entity.UserData), dataService, syncDate, tenant, runId, &completed, &failed, &skipped)
		}
		if len(users) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *userSyncService) syncUser(ctx context.Context, userInput entity.UserData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserSyncService.syncUser")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	userInput.Normalize()

	if userInput.Skip {
		if err := dataService.MarkProcessed(ctx, userInput.SyncId, runId, true, true, "Input user marked for skip"); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	userInput.Email = strings.ToLower(userInput.Email)

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

	if !failedSync {
		err = s.repositories.UserRepository.MergeUser(ctx, tenant, syncDate, userInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging user with external reference %v for tenant %v :%v", userInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if userInput.HasEmail() && !failedSync {
		err = s.repositories.UserRepository.MergeEmail(ctx, tenant, userInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging email for user with id %v for tenant %v :%v", userId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if userInput.HasPhoneNumber() && !failedSync {
		err = s.repositories.UserRepository.MergePhoneNumber(ctx, tenant, userInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging phone number for user with id %v for tenant %v :%v", userId, tenant, err)
			s.log.Errorf(reason)
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
