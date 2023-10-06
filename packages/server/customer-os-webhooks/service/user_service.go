package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	usergrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

const maxWorkersUserSync = 4

type UserService interface {
	SyncUsers(ctx context.Context, users []model.UserData) error
	mapDbNodeToUserEntity(dbNode dbtype.Node) *entity.UserEntity
}

type userService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewUserService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) UserService {
	return &userService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *userService) SyncUsers(ctx context.Context, users []model.UserData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.SyncUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		return errors.ErrTenantNotValid
	}

	// pre-validate user input before syncing
	for _, user := range users {
		if user.ExternalSystem == "" {
			return errors.ErrMissingExternalSystem
		}
		if !entity.IsValidDataSource(strings.ToLower(user.ExternalSystem)) {
			return errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, maxWorkersUserSync)

	syncMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all users concurrently
	for _, userData := range users {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue with Slack sync
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

			result := s.syncUser(ctx, syncMutex, userData, syncDate, common.GetTenantFromContext(ctx))
			statuses = append(statuses, result)
		}(userData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), users[0].ExternalSystem,
		users[0].AppSource, "user", syncDate, statuses)

	return nil
}

func (s *userService) syncUser(ctx context.Context, syncMutex *sync.Mutex, userInput model.UserData, syncDate time.Time, tenant string) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.syncUser")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", userInput.ExternalSystem), log.Object("userInput", userInput), log.String("tenant", tenant))

	var failedSync = false
	var reason = ""
	userInput.Normalize()

	// TODO: Merge external system, should be cached and moved to external system service
	err := s.repositories.ExternalSystemRepository.MergeExternalSystem(ctx, tenant, userInput.ExternalSystem, userInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", userInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		return NewFailedSyncStatus(reason)
	}

	// Check if user sync should be skipped
	if userInput.Skip {
		span.LogFields(log.Bool("skippedSync", true))
		return NewSkippedSyncStatus(userInput.SkipReason)
	}

	// Lock user and email creation
	syncMutex.Lock()
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
		userId = utils.NewUUIDIfEmpty(userId)
		userInput.Id = userId
		span.LogFields(log.String("userId", userId))

		// Create or update user
		_, err = s.grpcClients.UserClient.UpsertUser(ctx, &usergrpc.UpsertUserGrpcRequest{
			Tenant:          tenant,
			Id:              userId,
			LoggedInUserId:  "",
			FirstName:       userInput.FirstName,
			LastName:        userInput.LastName,
			Name:            userInput.Name,
			CreatedAt:       utils.ConvertTimeToTimestampPtr(userInput.CreatedAt),
			UpdatedAt:       utils.ConvertTimeToTimestampPtr(userInput.UpdatedAt),
			Internal:        false,
			ProfilePhotoUrl: userInput.ProfilePhotoUrl,
			Timezone:        userInput.Timezone,
			SourceFields: &commongrpc.SourceFields{
				Source:    userInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(userInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commongrpc.ExternalSystemFields{
				ExternalSystemId: userInput.ExternalSystem,
				ExternalId:       userInput.ExternalId,
				ExternalUrl:      userInput.ExternalUrl,
				ExternalIdSecond: utils.StringFirstNonEmpty(userInput.ExternalOwnerId, userInput.ExternalIdSecond),
				ExternalSource:   userInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed sending event to upsert user with external reference %s for tenant %s :%s", userInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for user to be created in neo4j
		if !failedSync && !matchingUserExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				user, findErr := s.repositories.UserRepository.GetById(ctx, tenant, userId)
				if user != nil && findErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
			}
		}
	}
	if !failedSync && userInput.HasEmail() {
		// Create or update email
		emailId, err := s.services.EmailService.CreateEmail(ctx, userInput.Email, userInput.ExternalSystem, userInput.AppSource)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed to create email address for user %s: %s", userId, err.Error())
			s.log.Error(reason)
		}
		// Link email to user
		if emailId != "" {
			_, err = s.grpcClients.UserClient.LinkEmailToUser(ctx, &usergrpc.LinkEmailToUserGrpcRequest{
				Tenant:  common.GetTenantFromContext(ctx),
				UserId:  userId,
				EmailId: emailId,
				Primary: true,
			})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("Failed to link email address %s for user %s: %s", userInput.Email, userId, err.Error())
				s.log.Error(reason)
			}
		}
	}
	syncMutex.Unlock()

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
			if phoneNumberId != "" {
				_, err = s.grpcClients.UserClient.LinkPhoneNumberToUser(ctx, &usergrpc.LinkPhoneNumberToUserGrpcRequest{
					Tenant:        common.GetTenantFromContext(ctx),
					UserId:        userId,
					PhoneNumberId: phoneNumberId,
					Primary:       phoneNumberDtls.Primary,
					Label:         phoneNumberDtls.Label,
				})
				if err != nil {
					failedSync = true
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("Failed to link phone number %s for user %s: %s", phoneNumberDtls.Number, userId, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		return NewFailedSyncStatus(reason)
	}
	return NewSuccessfulSyncStatus()
}

func (s *userService) mapDbNodeToUserEntity(dbNode dbtype.Node) *entity.UserEntity {
	props := utils.GetPropsFromNode(dbNode)
	return &entity.UserEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Roles:           utils.GetListStringPropOrEmpty(props, "roles"),
		Internal:        utils.GetBoolPropOrFalse(props, "internal"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
	}
}
