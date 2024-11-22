package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
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

type LogEntryService interface {
	Save(ctx context.Context, id *string, logEntryFields data_fields.LogEntryFields) (string, error)
}

type logEntryService struct {
	log      logger.Logger
	services *Services
}

func NewLogEntryService(log logger.Logger, services *Services) LogEntryService {
	return &logEntryService{
		log:      log,
		services: services,
	}
}

func (s *logEntryService) Save(ctx context.Context, id *string, logEntryFields data_fields.LogEntryFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "logEntryFields", logEntryFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	logEntryId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if logEntryFields.CreatedAt == nil || logEntryFields.CreatedAt.IsZero() {
			logEntryFields.CreatedAt = utils.NowPtr()
		}
		if logEntryFields.StartedAt == nil || logEntryFields.StartedAt.IsZero() {
			logEntryFields.StartedAt = utils.NowPtr()
		}
		if utils.IfNotNilString(logEntryFields.Source) == "" {
			logEntryFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(logEntryFields.AppSource) == "" {
			logEntryFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}

		// check mandatory fields
		if logEntryFields.OrganizationId == nil || *logEntryFields.OrganizationId == "" {
			err = errors.New("organizationId is required")
			tracing.TraceErr(span, err)
			return "", err
		}
		// validate organization exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *logEntryFields.OrganizationId, model.NodeLabelOrganization)
		if err != nil || !exists {
			err = errors.New("organization not found")
			tracing.TraceErr(span, err)
			return "", err
		}

		logEntryId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelLogEntry)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		logEntryId = *id

		// validate log entry exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, logEntryId, model.NodeLabelLogEntry)
		if err != nil || !exists {
			err = errors.New("log entry not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, logEntryId)

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		if createFlow {
			innerErr := s.services.Neo4jRepositories.LogEntryWriteRepository.CreateInTx(ctx, &tx, tenant, logEntryId, logEntryFields)
			if innerErr != nil {
				s.log.Errorf("Error while saving log entry %s: %s", logEntryId, err.Error())
				return nil, innerErr
			}
		} else {
			innerErr := s.services.Neo4jRepositories.LogEntryWriteRepository.UpdateInTx(ctx, &tx, tenant, logEntryId, logEntryFields)
			if innerErr != nil {
				s.log.Errorf("Error while updating log entry %s: %s", logEntryId, err.Error())
				return nil, innerErr
			}
		}
		if logEntryFields.ExternalSystemAvailable() {
			innerErr := s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, &tx, tenant, logEntryId, model.NodeLabelLogEntry, *logEntryFields.ExternalSystem)
			if err != nil {
				s.log.Errorf("Error while link log entry %s with external system %s: %s", logEntryId, logEntryFields.ExternalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// send events
	if createFlow {
		err = s.services.OrganizationService.RequestRefreshLastTouchpoint(ctx, *logEntryFields.OrganizationId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to request refresh last touchpoint"))
		}
		err = s.services.RabbitMQService.PublishEvent(ctx, logEntryId, model.LOG_ENTRY, dto.CreateLogEntry{logEntryFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContact"))
		}
		utils.EventCompleted(ctx, tenant, model.LOG_ENTRY.String(), logEntryId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())
	} else {
		err = s.services.RabbitMQService.PublishEvent(ctx, logEntryId, model.LOG_ENTRY, dto.UpdateLogEntry{logEntryFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContact"))
		}
		if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.LOG_ENTRY.String(), logEntryId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	if createFlow {
		span.LogFields(log.Bool("response.logEntryCreated", true))
	} else {
		span.LogFields(log.Bool("response.logEntryUpdated", true))
	}
	return logEntryId, nil
}
