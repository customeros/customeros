package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type UserService interface {
	GetById(ctx context.Context, userId string) (*neo4jentity.UserEntity, error)
	GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error)
	FindUserByEmail(parentCtx context.Context, email string) (*neo4jentity.UserEntity, error)
	CreateUser(ctx context.Context, userEntity neo4jentity.UserEntity) (string, error)
	CreateTestUser(ctx context.Context, firstName, lastName string) (string, error)
}

type userService struct {
	services *Services
}

func NewUserService(service *Services) UserService {
	return &userService{
		services: service,
	}
}

func (s *userService) GetById(parentCtx context.Context, userId string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	node, err := s.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, common.GetContext(ctx).Tenant, userId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return mapper.MapDbNodeToUserEntity(node), nil
}

func (s *userService) GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error) {
	nodes, err := s.services.Neo4jRepositories.UserReadRepository.GetAllForTenant(ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsersForTenant: %w", err)
	}

	users := make([]*neo4jentity.UserEntity, len(nodes))

	for i, node := range nodes {
		users[i] = mapper.MapDbNodeToUserEntity(node)
	}

	return users, nil
}

func (s *userService) FindUserByEmail(parentCtx context.Context, email string) (*neo4jentity.UserEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserService.FindFirstUserWithRolesByEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, email)
	if err != nil {
		return nil, err
	}

	if userDbNode == nil {
		return nil, nil
	}

	return mapper.MapDbNodeToUserEntity(userDbNode), nil
}

func (s *userService) CreateUser(ctx context.Context, userEntity neo4jentity.UserEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.CreateUser")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	userId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelUser)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	userEntity.Id = userId
	if userEntity.AppSource == "" {
		userEntity.AppSource = common.GetAppSourceFromContext(ctx)
	}
	err = s.services.Neo4jRepositories.UserWriteRepository.CreateUser(ctx, userEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return userId, nil
}

func (s *userService) CreateTestUser(ctx context.Context, firstName, lastName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.CreateTestUser")
	defer span.Finish()
	span.LogKV("firstName", firstName, "lastName", lastName)

	return s.CreateUser(ctx, neo4jentity.UserEntity{
		FirstName: firstName,
		LastName:  lastName,
		Test:      true,
	})
}
