package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type MarkdownEventService interface {
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input data_fields.MarkdownEventFields) (string, error)
}

type markdownEventService struct {
	log      logger.Logger
	services *Services
}

func NewMarkdownEventService(log logger.Logger, services *Services) MarkdownEventService {
	return &markdownEventService{
		log:      log,
		services: services,
	}
}

func (s *markdownEventService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input data_fields.MarkdownEventFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkdownEventService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	markdownEventId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if input.CreatedAt == nil || input.CreatedAt.IsZero() {
			input.CreatedAt = utils.NowPtr()
		}
		if input.Source == nil {
			input.Source = utils.ToPtr(neo4jentity.DataSourceOpenline)
		}
		if utils.IfNotNilString(input.AppSource) == "" {
			input.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}

		// check mandatory fields
		if utils.IfNotNilString(input.OrganizationId) == "" {
			err = errors.New("organizationId is required")
			tracing.TraceErr(span, err)
			return "", err
		}
		// validate organization exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, utils.IfNotNilString(input.OrganizationId), model.NodeLabelOrganization)
		if err != nil || !exists {
			err = errors.New("organization not found")
			tracing.TraceErr(span, err)
			return "", err
		}

		markdownEventId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelMarkdownEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		return "", errors.New("update not supported")
	}
	tracing.TagEntity(span, markdownEventId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		if createFlow {
			innerErr := s.services.Neo4jRepositories.MarkdownEventWriteRepository.CreateInTx(ctx, txWithPostCommit.Tx, tenant, markdownEventId, input)
			if innerErr != nil {
				s.log.Errorf("Error while saving markdown event %s: %s", markdownEventId, err.Error())
				return nil, innerErr
			}
		}

		if input.ExternalSystemAvailable() {
			innerErr := s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, markdownEventId, model.NodeLabelMarkdownEvent, *input.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link markdown event %s with external system %s: %s", markdownEventId, input.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				// historify markdown event
				err = s.services.RabbitMQService.PublishEvent(ctx, markdownEventId, model.MARKDOWN_EVENT, dto.CreateMarkdownEvent{input})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContact"))
				}

				// send event completed for organization for refresh
				utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), *input.OrganizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		span.LogFields(log.Bool("response.logEntryCreated", true))
	}
	return markdownEventId, nil
}
