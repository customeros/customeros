package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UserService interface {
	GetById(ctx context.Context, userId string) (*neo4jentity.UserEntity, error)
	GetAllUsersForTenant(ctx context.Context, tenant string) ([]*neo4jentity.UserEntity, error)
	FindUserByEmail(parentCtx context.Context, email string) (*neo4jentity.UserEntity, error)

	Create(ctx context.Context, input UserCreateData) (*string, error)
}

type userService struct {
	services *Services
}

func NewUserService(service *Services) UserService {
	return &userService{
		services: service,
	}
}

type UserCreateData struct {
	UserInput   neo4jentity.UserEntity
	EmailInput  neo4jentity.EmailEntity
	PlayerInput neo4jentity.PlayerEntity
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

func (s *userService) Create(ctx context.Context, input UserCreateData) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.Create")
	defer span.Finish()

	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.Object("input", input))

	tenant := common.GetTenantFromContext(ctx)

	var err error

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	input.UserInput.Id, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelUser)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.services.Neo4jRepositories.UserWriteRepository.CreateUser(ctx, input.UserInput)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	_, err = s.services.EmailService.Merge(ctx, tenant,
		EmailFields{
			Email:     input.EmailInput.Email,
			AppSource: input.UserInput.AppSource,
		},
		&LinkWith{
			Type:         commonModel.USER,
			Id:           input.UserInput.Id,
			Relationship: "HAS",
		})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.services.Neo4jRepositories.PlayerWriteRepository.Merge(ctx, input.UserInput.Id, input.PlayerInput)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &input.UserInput.Id, nil
}
