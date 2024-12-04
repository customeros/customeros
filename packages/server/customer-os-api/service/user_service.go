package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	jobrolepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/job_role"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type UserService interface {
	GetAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model2.SortBy) (*utils.Pagination, error)
	AddRole(ctx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error)
	AddRoleInTenant(ctx context.Context, userId string, tenant string, role model.Role) (*neo4jentity.UserEntity, error)
	RemoveRole(ctx context.Context, userId string, role model.Role) (*neo4jentity.UserEntity, error)
	RemoveRoleInTenant(ctx context.Context, userId string, tenant string, role model.Role) (*neo4jentity.UserEntity, error)
	ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool

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

func (s *userService) GetAll(parentCtx context.Context, page, limit int, filter *model.Filter, sortBy []*model2.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	paginatedResult := utils.Pagination{
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

	dbNodesWithTotalCount, err := s.repositories.Neo4jRepositories.UserReadRepository.GetPaginatedCustomerUsers(
		ctx,
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

func (s *userService) addPlayerDbRelationshipToUser(relationship dbtype.Relationship, userEntity *neo4jentity.UserEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	userEntity.DefaultForPlayer = utils.GetBoolPropOrFalse(props, "default")
}
