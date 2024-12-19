package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

type UserService interface {
	SyncUsers(ctx context.Context, users []model.UserData) (SyncResult, error)
	GetIdForReferencedUser(ctx context.Context, tenant, externalSystemId string, user model.ReferencedUser) (string, error)
}

type userService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewUserService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) UserService {
	return &userService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.UserSyncConcurrency,
	}
}

func (s *userService) SyncUsers(ctx context.Context, users []model.UserData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.SyncUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate user input before syncing
	for _, user := range users {
		if user.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(user.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", user.ExternalSystem))
			return SyncResult{}, errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all users concurrently
	for _, userData := range users {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(userData model.UserData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncUser(ctx, syncMutex, userData, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(userData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), users[0].ExternalSystem,
		users[0].AppSource, "user", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *userService) syncUser(ctx context.Context, syncMutex *sync.Mutex, userInput model.UserData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.syncUser")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, userInput.ExternalSystem)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "userInput", userInput)

	var tenant = common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""
	userInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, userInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", userInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if user sync should be skipped
	if userInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(userInput.SkipReason)
	}

	// Lock user and email creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if user already exists
	userId, err := s.repositories.UserRepository.GetMatchedUserId(ctx, tenant, userInput.ExternalSystem, userInput.ExternalId, userInput.Email)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched user with external reference %s for tenant %s :%s", userInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}

	if !failedSync {
		matchingUserExists := userId != ""
		span.LogFields(log.Bool("found matching user", matchingUserExists))

		// Create new user id if not found
		var inputUserId *string = nil
		if userId != "" {
			inputUserId = &userId
		}

		// Create or update user
		userFields := data_fields.UserFields{
			FirstName:       utils.StringPtr(userInput.FirstName),
			LastName:        utils.StringPtr(userInput.LastName),
			Name:            utils.StringPtr(userInput.Name),
			CreatedAt:       userInput.CreatedAt,
			Internal:        utils.BoolPtr(false),
			Source:          utils.StringPtr(userInput.ExternalSystem),
			Bot:             utils.BoolPtr(userInput.Bot),
			Timezone:        utils.StringPtr(userInput.Timezone),
			ProfilePhotoUrl: utils.StringPtr(userInput.ProfilePhotoUrl),
			ExternalSystem: &neo4jmodel.ExternalSystem{
				ExternalSystemId: userInput.ExternalSystem,
				ExternalId:       userInput.ExternalId,
				ExternalUrl:      userInput.ExternalUrl,
				ExternalIdSecond: userInput.ExternalIdSecond,
				ExternalSource:   userInput.ExternalSourceEntity,
				SyncDate:         &syncDate,
			},
		}
		userId, err = s.services.CommonServices.UserService.Save(ctx, nil, inputUserId, userFields)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed to save user with external reference %s for tenant %s :%s", userInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		userInput.Id = userId
		span.LogFields(log.String("userId", userId))
	}
	if !failedSync && userInput.HasEmail() {
		_, err = s.services.CommonServices.EmailService.Merge(ctx, nil, tenant,
			commonservice.EmailFields{
				Email:     userInput.Email,
				AppSource: userInput.AppSource,
				Source:    neo4jentity.DecodeDataSource(userInput.ExternalSystem),
				Primary:   true,
			},
			&commonservice.LinkWith{
				Type: commonmodel.USER,
				Id:   userId,
			})
		if err != nil {
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed to create and link email address %s with user %s: %s", userInput.Email, userId, err.Error())
			failedSync = true
		}
	}

	if !failedSync && userInput.HasPhoneNumbers() {
		for _, phoneNumberDtls := range userInput.PhoneNumbers {
			// Create or update phone number
			phoneNumberId, err := s.services.PhoneNumberService.CreatePhoneNumber(ctx, phoneNumberDtls.Number, userInput.ExternalSystem, userInput.AppSource)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to create phone number for user %s: %s", userId, err.Error())
				s.log.Error(reason)
			}
			// Link phone number to user
			if !failedSync {

				err := s.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithUser(ctx, tenant, userId, phoneNumberId, phoneNumberDtls.Label, phoneNumberDtls.Primary)
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err, log.String("grpcMethod", "LinkPhoneNumberToUser"))
					reason = fmt.Sprintf("Failed to link phone number %s for user %s: %s", phoneNumberDtls.Number, userId, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}

func (s *userService) GetIdForReferencedUser(ctx context.Context, tenant, externalSystemId string, user model.ReferencedUser) (string, error) {
	if !user.Available() {
		return "", nil
	}

	if user.ReferencedById() {
		return s.repositories.UserRepository.GetUserIdById(ctx, tenant, user.Id)
	} else if user.ReferencedByExternalId() {
		return s.repositories.UserRepository.GetUserIdByExternalId(ctx, tenant, user.ExternalId, externalSystemId)
	} else if user.ReferencedByExternalIdSecond() {
		return s.repositories.UserRepository.GetUserIdByExternalIdSecond(ctx, tenant, user.ExternalIdSecond, externalSystemId)
	}
	return "", nil
}
