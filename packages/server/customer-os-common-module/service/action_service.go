package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ActionService interface {
	GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*neo4jentity.ActionEntities, error)
	CreateActionForOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, actionFields data_fields.ActionFields) (string, error)
}

type actionService struct {
	log      logger.Logger
	services *Services
}

func NewActionService(log logger.Logger, services *Services) ActionService {
	return &actionService{
		log:      log,
		services: services,
	}
}

func (s *actionService) GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*neo4jentity.ActionEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("ids", ids))

	records, err := s.services.Neo4jRepositories.ActionReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	var data neo4jentity.ActionEntities
	for _, v := range records {
		action := neo4jmapper.MapDbNodeToActionEntity(v.Node)
		action.DataloaderKey = v.LinkedNodeId
		data = append(data, *action)
	}

	return &data, nil
}

func (s *actionService) CreateActionForOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, actionFields data_fields.ActionFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.CreateActionForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, organizationId)
	tracing.LogObjectAsJson(span, "actionFields", actionFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate organization exists
	err = s.services.OrganizationService.ValidateOrganizationExists(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// set default values
	if actionFields.AppSource == nil {
		actionFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}
	if actionFields.Source == nil {
		actionFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
	}
	if actionFields.CreatedAt == nil {
		actionFields.CreatedAt = utils.NowPtr()
	}
	if actionFields.ActionType == nil {
		actionFields.ActionType = utils.ToPtr(enum.ActionGeneric)
	}

	actionId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelAction)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		// create action
		err = s.services.Neo4jRepositories.ActionWriteRepository.CreateV2(ctx, txWithPostCommit.Tx, tenant, actionId, organizationId, model.ORGANIZATION, actionFields)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.AddActionToOrganization{actionFields})
			if err != nil {
				tracing.TraceErr(span, err)
			}
			s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return actionId, nil
}
