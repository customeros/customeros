package markdown_event

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"

	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type markdownEventService struct {
	log    logger.Logger
	neo4j  *neo4j_repository.Repositories
	events *events.EventsService
}

func NewMarkdownEventService(log logger.Logger, neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.MarkdownEventService {
	return &markdownEventService{
		log:    log,
		neo4j:  neo4j,
		events: events,
	}
}

func (s *markdownEventService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input data_fields.MarkdownEventFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MarkdownEventService.Save")
	defer spans.Finish()

	spans.LogObjectAsJson("input", input)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	markdownEventId := ""

	if id == nil || *id == "" {
		createFlow = true
		spans.LogKV("flow", "create")

		// prepare missing fields
		if input.CreatedAt == nil || input.CreatedAt.IsZero() {
			input.CreatedAt = utils.NowPtr()
		}
		if input.Source == nil {
			input.Source = utils.ToPtr(enum.SourceCustomerOS)
		}
		if utils.IfNotNilString(input.AppSource) == "" {
			input.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}

		// check mandatory fields
		if utils.IfNotNilString(input.OrganizationId) == "" {
			err = errors.New("organizationId is required")
			spans.TraceError(err)
			return "", err
		}
		// validate organization exists
		exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, utils.IfNotNilString(input.OrganizationId), model.NodeLabelOrganization)
		if err != nil || !exists {
			err = errors.New("organization not found")
			spans.TraceError(err)
			return "", err
		}

		markdownEventId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelMarkdownEvent)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}
	} else {
		return "", errors.New("update not supported")
	}
	spans.TagEntity(markdownEventId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			innerErr := s.neo4j.MarkdownEventWriteRepository.CreateInTx(ctx, txWithPostCommit.Tx, tenant, markdownEventId, input)
			if innerErr != nil {
				s.log.Errorf("Error while saving markdown event %s: %s", markdownEventId, err.Error())
				return nil, innerErr
			}
		}

		if input.ExternalSystemAvailable() {
			innerErr := s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, markdownEventId, model.NodeLabelMarkdownEvent, *input.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link markdown event %s with external system %s: %s", markdownEventId, input.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// send events
			if createFlow {
				// historify markdown event
				err = s.events.Publisher.PublishFanoutEvent(ctx, markdownEventId, model.MARKDOWN_EVENT, input)
				if err != nil {
					spans.TraceError(errors.Wrap(err, "unable to publish message CreateContact"))
				}

				// send event completed for organization for refresh
				s.events.Publisher.PublishNotification(ctx, tenant, *input.OrganizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	if createFlow {
		spans.LogFields(log.Bool("response.logEntryCreated", true))
	}
	return markdownEventId, nil
}
