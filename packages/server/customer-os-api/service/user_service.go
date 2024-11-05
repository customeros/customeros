package service

import (
	"context"
	"fmt"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	Create(ctx context.Context, UserEntity neo4jentity.UserEntity) (string, error)
	Update(ctx context.Context, userId, firstName, lastName string, name, timezone, profilePhotoURL *string) (*neo4jentity.UserEntity, error)
	GetAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)

	IsOwner(ctx context.Context, id string) (*bool, error)
	GetContactOwner(ctx context.Context, contactId string) (*neo4jentity.UserEntity, error)
	GetNoteCreator(ctx context.Context, noteId string) (*neo4jentity.UserEntity, error)
	GetUsersConnectedForContacts(ctx context.Context, contactIds []string) (*neo4jentity.UserEntities, error)
	GetUsersForEmails(ctx context.Context, emailIds []string) (*neo4jentity.UserEntities, error)
	GetUsersForPhoneNumbers(ctx context.Context, phoneNumberIds []string) (*neo4jentity.UserEntities, error)
	GetUsersForPlayers(ctx context.Context, playerIds []string) (*neo4jentity.UserEntities, error)
	GetUserOwnersForOrganizations(ctx context.Context, organizationIDs []string) (*neo4jentity.UserEntities, error)
	GetUserOwnersForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForServiceLineItems(ctx context.Context, serviceLineItemIds []string) (*neo4jentity.UserEntities, error)
	GetUsersWithMailboxes(ctx context.Context) (*neo4jentity.UserEntities, error)
	GetUserCreatorsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.UserEntities, error)
	GetUserAuthorsForLogEntries(ctx context.Context, logEntryIDs []string) (*neo4jentity.UserEntities, error)
	GetUserAuthorsForComments(ctx context.Context, commentIds []string) (*neo4jentity.UserEntities, error)
	GetUserForFlowSenders(ctx context.Context, flowSenderIds []string) (*neo4jentity.UserEntities, error)
	GetUsers(ctx context.Context, userIds []string) (*neo4jentity.UserEntities, error)
	GetDistinctOrganizationOwners(ctx context.Context) (*neo4jentity.UserEntities, error)
	GetReminderOwner(ctx context.Context, reminderId string) (*neo4jentity.UserEntity, error)
	AddRole(ctx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error)
	AddRoleInTenant(ctx context.Context, userId string, tenant string, role model.Role) (*neo4jentity.UserEntity, error)
	RemoveRole(ctx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error)
	RemoveRoleInTenant(ctx context.Context, userId string, tenant string, role model.Role) (*neo4jentity.UserEntity, error)
	ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool
	GetContractOwner(ctx context.Context, contractId string) (*neo4jentity.UserEntity, error)

	addPlayerDbRelationshipToUser(relationship dbtype.Relationship, userEntity *neo4jentity.UserEntity)

	CustomerAddJobRole(ctx context.Context, entity *CustomerAddJobRoleData) (*model.CustomerUser, error)
}

type userService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

type CustomerAddJobRoleData struct {
	UserId        string
	JobRoleEntity *neo4jentity.JobRoleEntity
}

func NewUserService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) UserService {
	return &userService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *userService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *userService) Update(ctx context.Context, userId, firstName, lastName string, name, timezone, profilePhotoURL *string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId))

	if userId != common.GetContext(ctx).UserId {
		if !s.ContainsRole(ctx, []model.Role{model.RoleAdmin, model.RolePlatformOwner, model.RoleOwner}) {
			return nil, fmt.Errorf("user can not update other user")
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.UpsertUser(ctx, &userpb.UpsertUserGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			Id:             userId,
			SourceFields: &commonpb.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: constants.AppSourceCustomerOsApi,
			},
			FirstName:       firstName,
			LastName:        lastName,
			Name:            utils.IfNotNilString(name),
			Timezone:        utils.IfNotNilString(timezone),
			ProfilePhotoUrl: utils.IfNotNilString(profilePhotoURL),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return s.services.CommonServices.UserService.GetById(ctx, userId)
}

func (s *userService) ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.ContainsRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	myRoles := common.GetRolesFromContext(ctx)
	for _, allowedRole := range allowedRoles {
		for _, myRole := range myRoles {
			if myRole == allowedRole.String() {
				return true
			}
		}
	}
	return false
}

func (s *userService) CanAddRemoveRole(parentCtx context.Context, role model.Role) bool {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.CanAddRemoveRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	switch role {
	case model.RoleAdmin:
		return false // this role is a special endpoint and can not be given to a user
	case model.RoleOwner:
		return s.ContainsRole(ctx, []model.Role{model.RoleAdmin, model.RolePlatformOwner, model.RoleOwner})
	case model.RolePlatformOwner:
		return s.ContainsRole(ctx, []model.Role{model.RoleAdmin, model.RolePlatformOwner})
	case model.RoleUser:
		return s.ContainsRole(ctx, []model.Role{model.RoleAdmin, model.RolePlatformOwner, model.RoleOwner})
	default:
		s.log.Errorf("unknown role: %s", role)
		return false
	}
}

func (s *userService) AddRole(parentCtx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.AddRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("role", string(role)))

	if !s.CanAddRemoveRole(ctx, role) {
		return nil, fmt.Errorf("logged-in user can not add role: %s", role)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.AddRole(ctx, &userpb.AddRoleGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			UserId:         userId,
			Role:           mapper.MapRoleToEntity(role),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.services.CommonServices.UserService.GetById(ctx, userId)
}

func (s *userService) CustomerAddJobRole(ctx context.Context, entity *CustomerAddJobRoleData) (*model.CustomerUser, error) {
	result := &model.CustomerUser{}

	jobRoleCreate := &jobrolepb.CreateJobRoleGrpcRequest{
		Tenant:      common.GetTenantFromContext(ctx),
		JobTitle:    entity.JobRoleEntity.JobTitle,
		Description: entity.JobRoleEntity.Description,
		Primary:     &entity.JobRoleEntity.Primary,
		StartedAt:   timestamppb.New(utils.IfNotNilTimeWithDefault(entity.JobRoleEntity.StartedAt, utils.Now())),
		EndedAt:     timestamppb.New(utils.IfNotNilTimeWithDefault(entity.JobRoleEntity.EndedAt, utils.Now())),
		AppSource:   entity.JobRoleEntity.AppSource,
		Source:      string(entity.JobRoleEntity.Source),
		CreatedAt:   timestamppb.New(entity.JobRoleEntity.CreatedAt),
	}

	contextWithTimeout, cancel := utils.GetLongLivedContext(ctx)
	defer cancel()

	jobRole, err := utils.CallEventsPlatformGRPCWithRetry[*jobrolepb.JobRoleIdGrpcResponse](func() (*jobrolepb.JobRoleIdGrpcResponse, error) {
		return s.grpcClients.JobRoleClient.CreateJobRole(contextWithTimeout, jobRoleCreate)
	})
	if err != nil {
		s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
		return nil, err
	}

	result.JobRole = &model.CustomerJobRole{
		ID: jobRole.Id,
	}
	user, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.LinkJobRoleToUser(contextWithTimeout, &userpb.LinkJobRoleToUserGrpcRequest{
			UserId:    entity.UserId,
			JobRoleId: jobRole.Id,
			Tenant:    common.GetTenantFromContext(ctx),
			AppSource: utils.StringFirstNonEmpty(entity.JobRoleEntity.AppSource, constants.AppSourceCustomerOsApi),
		})
	})
	if err != nil {
		s.log.Errorf("(%s) Failed to call method: {%v}", utils.GetFunctionName(), err.Error())
		return nil, err
	}
	result.ID = user.Id
	return result, nil
}

func (s *userService) AddRoleInTenant(parentCtx context.Context, userId, tenant string, role model.Role) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.AddRoleInTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("role", string(role)))

	if !s.CanAddRemoveRole(ctx, role) {
		return nil, fmt.Errorf("logged-in user can not add role: %s", role)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.AddRole(ctx, &userpb.AddRoleGrpcRequest{
			Tenant:         tenant,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			UserId:         userId,
			Role:           mapper.MapRoleToEntity(role),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.services.CommonServices.UserService.GetById(ctx, userId)
}

func (s *userService) RemoveRole(parentCtx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.RemoveRole")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("role", string(role)))

	if !s.CanAddRemoveRole(ctx, role) {
		return nil, fmt.Errorf("logged-in user can not remove role: %s", role)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.RemoveRole(ctx, &userpb.RemoveRoleGrpcRequest{
			Tenant:         common.GetTenantFromContext(ctx),
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			UserId:         userId,
			Role:           mapper.MapRoleToEntity(role),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.services.CommonServices.UserService.GetById(ctx, userId)
}

func (s *userService) RemoveRoleInTenant(parentCtx context.Context, userId string, tenant string, role model.Role) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.RemoveRoleInTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("role", string(role)))

	if !s.CanAddRemoveRole(ctx, role) {
		return nil, fmt.Errorf("logged-in user can not remove role: %s", role)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.RemoveRole(ctx, &userpb.RemoveRoleGrpcRequest{
			Tenant:         tenant,
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			UserId:         userId,
			Role:           mapper.MapRoleToEntity(role),
			AppSource:      constants.AppSourceCustomerOsApi,
		})
	})
	if err != nil {
		return nil, err
	}

	return s.services.CommonServices.UserService.GetById(ctx, userId)
}

func (s *userService) GetAll(parentCtx context.Context, page, limit int, filter *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(neo4jentity.UserEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := buildFilter(filter, reflect.TypeOf(neo4jentity.UserEntity{}))
	if err != nil {
		return nil, err
	}

	dbNodesWithTotalCount, err := s.repositories.UserRepository.GetPaginatedCustomerUsers(
		ctx,
		session,
		common.GetContext(ctx).Tenant,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherFilter,
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	users := make(neo4jentity.UserEntities, 0, len(dbNodesWithTotalCount.Nodes))
	for _, v := range dbNodesWithTotalCount.Nodes {
		entity := neo4jmapper.MapDbNodeToUserEntity(v)
		users = append(users, *entity)
	}
	paginatedResult.SetRows(&users)
	return &paginatedResult, nil
}

func (s *userService) IsOwner(parentCtx context.Context, userId string) (*bool, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.IsOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	isOwner, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repositories.UserRepository.IsOwner(ctx, tx, common.GetContext(ctx).Tenant, userId)
	})
	if err != nil {
		return nil, err
	}
	return isOwner.(*bool), nil
}

func (s *userService) GetContactOwner(parentCtx context.Context, contactId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetContactOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	ownerDbNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repositories.UserRepository.GetOwnerForContact(ctx, tx, common.GetContext(ctx).Tenant, contactId)
	})
	if err != nil {
		return nil, err
	} else if ownerDbNode.(*dbtype.Node) == nil {
		return nil, nil
	} else {
		return neo4jmapper.MapDbNodeToUserEntity(ownerDbNode.(*dbtype.Node)), nil
	}
}

func (s *userService) GetNoteCreator(parentCtx context.Context, noteId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetNoteCreator")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	userDbNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repositories.UserRepository.GetCreatorForNote(ctx, tx, common.GetContext(ctx).Tenant, noteId)
	})
	if err != nil {
		return nil, err
	} else if userDbNode.(*dbtype.Node) == nil {
		return nil, nil
	} else {
		return neo4jmapper.MapDbNodeToUserEntity(userDbNode.(*dbtype.Node)), nil
	}
}

func (s *userService) GetUsersConnectedForContacts(ctx context.Context, contactIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUsersConnectedForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.repositories.UserRepository.GetUsersConnectedForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsersForEmails(parentCtx context.Context, emailIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsersForEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.repositories.UserRepository.GetAllForEmails(ctx, common.GetTenantFromContext(ctx), emailIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsersForPhoneNumbers(parentCtx context.Context, phoneNumberIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsersForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.repositories.UserRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), phoneNumberIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsersForPlayers(parentCtx context.Context, playerIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsersForPlayers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	users, err := s.repositories.Neo4jRepositories.PlayerReadRepository.GetUsersForPlayer(ctx, playerIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		s.addPlayerDbRelationshipToUser(*v.Relationship, userEntity)
		userEntity.Tenant = v.Tenant
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserOwnersForOrganizations(parentCtx context.Context, organizationIDs []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserOwnersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	users, err := s.repositories.Neo4jRepositories.UserReadRepository.GetAllOwnersForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserOwnersForOpportunities(parentCtx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserOwnersForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.repositories.UserRepository.GetAllOwnersForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserCreatorsForOpportunities(parentCtx context.Context, opportunityIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunityIds", opportunityIds))

	users, err := s.repositories.UserRepository.GetAllCreatorsForOpportunities(ctx, common.GetTenantFromContext(ctx), opportunityIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}
func (s *userService) GetUserCreatorsForServiceLineItems(parentCtx context.Context, serviceLineItemIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("serviceLineItemIds", serviceLineItemIds))

	users, err := s.repositories.UserRepository.GetAllCreatorsForServiceLineItems(ctx, common.GetTenantFromContext(ctx), serviceLineItemIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}
func (s *userService) GetUsersWithMailboxes(ctx context.Context) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUsersWithMailboxes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	mailboxes, err := s.repositories.PostgresRepositories.TenantSettingsMailboxRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create a map to store unique usernames
	userEmailMap := make(map[string]struct{})
	for _, mailbox := range mailboxes {
		// Add the username to the map (map ensures uniqueness)
		userEmailMap[mailbox.Username] = struct{}{}
	}

	entities := make(neo4jentity.UserEntities, 0, len(userEmailMap))
	for userEmail := range userEmailMap {
		userNode, err := s.repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, userEmail)
		if err != nil {
			return nil, err
		}

		if userNode != nil {
			userEntity := neo4jmapper.MapDbNodeToUserEntity(userNode)
			entities = append(entities, *userEntity)
		}
	}

	return &entities, nil
}
func (s *userService) GetUserCreatorsForContracts(parentCtx context.Context, contractIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserCreatorsForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractIds", contractIds))

	users, err := s.repositories.UserRepository.GetAllCreatorsForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserAuthorsForLogEntries(parentCtx context.Context, logEntryIDs []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUserAuthorsForLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("logEntryIDs", logEntryIDs))

	users, err := s.repositories.UserRepository.GetAllAuthorsForLogEntries(ctx, common.GetTenantFromContext(ctx), logEntryIDs)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserAuthorsForComments(ctx context.Context, commentIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserAuthorsForComments")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("commentIds", commentIds))

	users, err := s.repositories.UserRepository.GetAllAuthorsForComments(ctx, common.GetTenantFromContext(ctx), commentIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUserForFlowSenders(ctx context.Context, flowSenderIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserForFlowSenders")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("flowSenderIds", flowSenderIds))

	users, err := s.repositories.UserRepository.GetAllSendersForFlowSenders(ctx, common.GetTenantFromContext(ctx), flowSenderIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(users))
	for _, v := range users {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(v.Node)
		userEntity.DataloaderKey = v.LinkedNodeId
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetUsers(parentCtx context.Context, userIds []string) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetUsers")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userIds", userIds))

	userDbNodes, err := s.repositories.UserRepository.GetUsers(ctx, common.GetTenantFromContext(ctx), userIds)
	if err != nil {
		return nil, err
	}
	userEntities := make(neo4jentity.UserEntities, 0, len(userDbNodes))
	for _, dbNode := range userDbNodes {
		userEntity := neo4jmapper.MapDbNodeToUserEntity(dbNode)
		userEntities = append(userEntities, *userEntity)
	}
	return &userEntities, nil
}

func (s *userService) GetDistinctOrganizationOwners(parentCtx context.Context) (*neo4jentity.UserEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetDistinctOrganizationOwners")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.repositories.UserRepository.GetDistinctOrganizationOwners(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}

	userEntities := make(neo4jentity.UserEntities, 0, len(dbNodes))
	for _, dbNode := range dbNodes {
		entity := neo4jmapper.MapDbNodeToUserEntity(dbNode)
		userEntities = append(userEntities, *entity)
	}
	return &userEntities, nil
}

func (s *userService) Create(ctx context.Context, userEntity neo4jentity.UserEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("userEntity", userEntity))

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
		return s.grpcClients.UserClient.UpsertUser(ctx, &userpb.UpsertUserGrpcRequest{
			Tenant: common.GetTenantFromContext(ctx),
			SourceFields: &commonpb.SourceFields{
				Source:    string(userEntity.Source),
				AppSource: userEntity.AppSource,
			},
			LoggedInUserId:  common.GetUserIdFromContext(ctx),
			FirstName:       userEntity.FirstName,
			LastName:        userEntity.LastName,
			Name:            userEntity.Name,
			Internal:        false,
			Bot:             false,
			ProfilePhotoUrl: userEntity.ProfilePhotoUrl,
			Timezone:        userEntity.Timezone,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelUser, span)
	return response.Id, nil
}

func (s *userService) addPlayerDbRelationshipToUser(relationship dbtype.Relationship, userEntity *neo4jentity.UserEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	userEntity.DefaultForPlayer = utils.GetBoolPropOrFalse(props, "default")
}

func (s *userService) GetContractOwner(parentCtx context.Context, contractId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetContractOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	ownerDbNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repositories.UserRepository.GetOwnerForContract(ctx, tx, common.GetContext(ctx).Tenant, contractId)
	})
	if err != nil {
		return nil, err
	} else if ownerDbNode.(*dbtype.Node) == nil {
		return nil, nil
	} else {
		return neo4jmapper.MapDbNodeToUserEntity(ownerDbNode.(*dbtype.Node)), nil
	}
}

func (s *userService) GetReminderOwner(ctx context.Context, reminderId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetReminderOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	ownerDbNode, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repositories.UserRepository.GetOwnerForReminder(ctx, tx, common.GetContext(ctx).Tenant, reminderId)
	})
	if err != nil {
		return nil, err
	} else if ownerDbNode.(*dbtype.Node) == nil {
		return nil, nil
	} else {
		return neo4jmapper.MapDbNodeToUserEntity(ownerDbNode.(*dbtype.Node)), nil
	}
}
