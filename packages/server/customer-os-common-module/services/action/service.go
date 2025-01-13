package action

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type actionService struct {
	log    logger.Logger
	neo4j  *neo4jRepository.Repositories
	events *events.EventsService
	org    interfaces.OrganizationService
}

func NewActionService(log logger.Logger, neo4j *neo4jRepository.Repositories, events *events.EventsService, org interfaces.OrganizationService) interfaces.ActionService {
	return &actionService{
		log:    log,
		neo4j:  neo4j,
		events: events,
		org:    org,
	}
}

func (s *actionService) GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*entity.ActionEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionService.GetActionsForNodes")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("ids", ids))

	records, err := s.neo4j.ActionReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	var data entity.ActionEntities
	for _, v := range records {
		action := mapper.MapDbNodeToActionEntity(v.Node)
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
	err = s.org.ValidateOrganizationExists(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// set default values
	if actionFields.AppSource == nil {
		actionFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}
	if actionFields.Source == nil {
		actionFields.Source = utils.StringPtr(entity.DataSourceOpenline.String())
	}
	if actionFields.CreatedAt == nil {
		actionFields.CreatedAt = utils.NowPtr()
	}
	if actionFields.ActionType == nil {
		actionFields.ActionType = utils.ToPtr(enum.ActionGeneric)
	}

	actionId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelAction)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// create action
		err = s.neo4j.ActionWriteRepository.CreateV2(ctx, txWithPostCommit.Tx, tenant, actionId, organizationId, model.ORGANIZATION, actionFields)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.AddActionToOrganization{actionFields})
			if err != nil {
				tracing.TraceErr(span, err)
			}
			s.events.Publisher.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
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
