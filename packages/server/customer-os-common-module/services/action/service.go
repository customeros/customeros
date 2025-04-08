package action

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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

func (s *actionService) SetOrganizationService(org interfaces.OrganizationService) {
	s.org = org
}

func (s *actionService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *actionService) GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*neo4j_entity.ActionEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ActionService.GetActionsForNodes")
	defer spans.Finish()

	spans.LogKV("entityType", entityType.String(), "ids", ids)

	records, err := s.neo4j.ActionReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	var data neo4j_entity.ActionEntities
	for _, v := range records {
		action := mapper.MapDbNodeToActionEntity(v.Node)
		action.DataloaderKey = v.LinkedNodeId
		data = append(data, *action)
	}

	return &data, nil
}

func (s *actionService) CreateActionForOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, actionFields data_fields.ActionFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ContractService.CreateActionForOrganization")
	defer spans.Finish()

	spans.TagEntity(organizationId)
	spans.LogObjectAsJson("actionFields", actionFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// set default values
	if actionFields.AppSource == nil {
		actionFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}
	if actionFields.Source == nil {
		actionFields.Source = utils.StringPtr(neo4j_entity.DataSourceOpenline.String())
	}
	if actionFields.CreatedAt == nil {
		actionFields.CreatedAt = utils.NowPtr()
	}
	if actionFields.ActionType == nil {
		actionFields.ActionType = utils.ToPtr(enum.ActionGeneric)
	}

	actionId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelAction)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// validate organization exists
		err = s.org.ValidateOrganizationExists(ctx, txWithPostCommit.Tx, organizationId)
		if err != nil {
			return nil, err
		}
		// create action
		err = s.neo4j.ActionWriteRepository.CreateV2(ctx, txWithPostCommit.Tx, tenant, actionId, organizationId, model.ORGANIZATION, actionFields)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishFanoutEvent(ctx, organizationId, model.ORGANIZATION, dto.AddActionToOrganization{ActionFields: actionFields})
			if err != nil {
				spans.TraceError(err)
			}
			s.events.Publisher.PublishNotification(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			return nil
		})

		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return actionId, nil
}
