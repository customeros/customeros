package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UserService interface {
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

	_, err = s.services.EmailService.Merge(ctx, neo4jentity.EmailEntity{
		Email:     input.EmailInput.Email,
		AppSource: input.EmailInput.AppSource,
	}, &LinkWith{
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
