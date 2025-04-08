package custom_fields

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type customFieldTemplateService struct {
	log    logger.Logger
	neo4j  *neoRepo.Repositories
	events *events.EventsService
}

func NewCustomFieldTemplateService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService) interfaces.CustomFieldTemplateService {
	return &customFieldTemplateService{
		log:    log,
		neo4j:  neo4j,
		events: events,
	}
}

func (s *customFieldTemplateService) GetAll(ctx context.Context) (*neo4jentity.CustomFieldTemplateEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CustomFieldTemplateService.GetAll")
	defer spans.Finish()

	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	dbNodes, err := s.neo4j.CustomFieldTemplateReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	customFieldTemplateEntities := neo4jentity.CustomFieldTemplateEntities{}
	for _, dbNode := range dbNodes {
		customFieldTemplateEntities = append(customFieldTemplateEntities, *neo4jmapper.MapDbNodeToCustomFieldTemplateEntity(dbNode))
	}
	return &customFieldTemplateEntities, nil
}

func (s *customFieldTemplateService) GetById(ctx context.Context, customFieldTemplateId string) (*neo4jentity.CustomFieldTemplateEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CustomFieldTemplateService.GetById")
	defer spans.Finish()

	spans.TagEntity(customFieldTemplateId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	dbNode, err := s.neo4j.CustomFieldTemplateReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), customFieldTemplateId)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if dbNode == nil {
		err = errors.New("custom field template not found")
		spans.TraceError(err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToCustomFieldTemplateEntity(dbNode), nil
}

func (s *customFieldTemplateService) Save(ctx context.Context, id *string, input neo4jrepository.CustomFieldTemplateSaveFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CustomFieldTemplateService.Save")
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
	customFieldTemplateId := ""

	if id == nil || *id == "" {
		createFlow = true
		spans.LogKV("flow", "create")
		customFieldTemplateId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelCustomFieldTemplate)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}
	} else {
		spans.LogKV("flow", "update")
		customFieldTemplateId = *id

		// validate custom field template exists
		_, err = s.GetById(ctx, customFieldTemplateId)
		if err != nil {
			spans.TraceError(err)
			return "", err
		}
	}
	spans.TagEntity(customFieldTemplateId)

	if createFlow {
		// validate entity type is present and is valid when creating new custom field template
		if !supportedEntityTypeForCustomFieldTemplate(input.EntityType) {
			err = errors.New("entity type is missing or not supported")
			spans.TraceError(err)
			return "", err
		}
	}

	err = s.neo4j.CustomFieldTemplateWriteRepository.Save(ctx, tenant, customFieldTemplateId, input)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	if createFlow {
		err = s.events.Publisher.PublishFanoutEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.New_CreateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(input))
		if err != nil {
			spans.TraceError(errors.Wrap(err, "unable to publish message CreateCustomFieldTemplate"))
		}
	} else {
		err = s.events.Publisher.PublishFanoutEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.New_UpdateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(input))
		if err != nil {
			spans.TraceError(errors.Wrap(err, "unable to publish message UpdateCustomFieldTemplate"))
		}
	}

	return customFieldTemplateId, nil
}

func supportedEntityTypeForCustomFieldTemplate(entityType model.EntityType) bool {
	return entityType == model.ORGANIZATION || entityType == model.OPPORTUNITY || entityType == model.CONTACT || entityType == model.LOG_ENTRY
}

func (s *customFieldTemplateService) Delete(ctx context.Context, customFieldTemplateId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CustomFieldTemplateService.Delete")
	defer spans.Finish()

	spans.TagEntity(customFieldTemplateId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// validate custom field template exists
	_, err = s.GetById(ctx, customFieldTemplateId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = s.neo4j.CustomFieldTemplateWriteRepository.Delete(ctx, common.GetTenantFromContext(ctx), customFieldTemplateId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = s.events.Publisher.PublishFanoutEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.Delete{})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message"))
	}

	return nil
}
