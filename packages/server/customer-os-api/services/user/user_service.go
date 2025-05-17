package api_user

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	api_filters "github.com/customeros/customeros/packages/server/customer-os-api/services/filters"
	api_sort "github.com/customeros/customeros/packages/server/customer-os-api/services/sort"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type userService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

type CustomerAddJobRoleData struct {
	UserId        string
	JobRoleEntity *neo4jentity.JobRoleEntity
}

func NewUserService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.UserService {
	return &userService{
		log:          log,
		repositories: repositories,
	}
}

func (s *userService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *userService) ContainsRole(parentCtx context.Context, allowedRoles []model.Role) bool {
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.ContainsRole")
	defer spans.Finish()

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
	spans, ctx := telemetry.StartServiceSpan(parentCtx, "UserService.CanAddRemoveRole")
	defer spans.Finish()

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

func (s *userService) GetAll(ctx context.Context, page, limit int, filter *model.Filter, sortBy []*model2.SortBy) (*utils.Pagination, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UserService.GetAll")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	paginatedResult := utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := api_sort.BuildSort(sortBy, reflect.TypeOf(neo4jentity.UserEntity{}))
	if err != nil {
		return nil, err
	}
	cypherFilter, err := api_filters.BuildFilter(filter, reflect.TypeOf(neo4jentity.UserEntity{}))
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
