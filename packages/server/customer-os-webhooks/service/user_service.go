package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
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
	var appSource = utils.StringFirstNonEmpty(userInput.AppSource, constants.AppSourceCustomerOsWebhooks)
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
		userId = utils.NewUUIDIfEmpty(userId)
		userInput.Id = userId
		span.LogFields(log.String("userId", userId))

		// Create or update user
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
			return s.grpcClients.UserClient.UpsertUser(ctx, &userpb.UpsertUserGrpcRequest{
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
				Bot:             userInput.Bot,
				SourceFields: &commonpb.SourceFields{
					Source:    userInput.ExternalSystem,
					AppSource: appSource,
				},
				ExternalSystemFields: &commonpb.ExternalSystemFields{
					ExternalSystemId: userInput.ExternalSystem,
					ExternalId:       userInput.ExternalId,
					ExternalUrl:      userInput.ExternalUrl,
					ExternalIdSecond: userInput.ExternalIdSecond,
					ExternalSource:   userInput.ExternalSourceEntity,
					SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
				},
			})
		})
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertUser"))
			reason = fmt.Sprintf("failed sending event to upsert user with external reference %s for tenant %s :%s", userInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for user to be created in neo4j
		if !failedSync && !matchingUserExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				found, findErr := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, userId, neo4jutil.NodeLabelUser)
				if found && findErr == nil {
					break
				}
				time.Sleep(utils.BackOffExponentialDelay(i))
			}
		}
	}
	if !failedSync && userInput.HasEmail() {
		// Create or update email
		emailId, err := s.services.EmailService.CreateEmail(ctx, userInput.Email, userInput.ExternalSystem, appSource)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("Failed to create email address for user %s: %s", userId, err.Error())
			s.log.Error(reason)
		}
		// Link email to user
		if !failedSync {
			_, err = CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
				return s.grpcClients.UserClient.LinkEmailToUser(ctx, &userpb.LinkEmailToUserGrpcRequest{
					Tenant:    common.GetTenantFromContext(ctx),
					UserId:    userId,
					EmailId:   emailId,
					Primary:   true,
					AppSource: appSource,
				})
			})
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err, log.String("grpcMethod", "LinkEmailToUser"))
				reason = fmt.Sprintf("Failed to link email address %s for user %s: %s", userInput.Email, userId, err.Error())
				s.log.Error(reason)
			}
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
				_, err = CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
					return s.grpcClients.UserClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
						Tenant:        common.GetTenantFromContext(ctx),
						UserId:        userId,
						PhoneNumberId: phoneNumberId,
						Primary:       phoneNumberDtls.Primary,
						Label:         phoneNumberDtls.Label,
						AppSource:     appSource,
					})
				})
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

func (s *userService) mapDbNodeToUserEntity(dbNode dbtype.Node) *entity.UserEntity {
	props := utils.GetPropsFromNode(dbNode)
	return &entity.UserEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Roles:           utils.GetListStringPropOrEmpty(props, "roles"),
		Internal:        utils.GetBoolPropOrFalse(props, "internal"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
	}
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
