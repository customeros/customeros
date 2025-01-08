package service

import (
	"context"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type UserService interface {
	GetAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model2.SortBy) (*utils.Pagination, error)
	ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool
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
