package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CustomFieldTemplateService interface {
	GetAll(ctx context.Context) (*neo4jentity.CustomFieldTemplateEntities, error)
	GetById(ctx context.Context, customFieldTemplateId string) (*neo4jentity.CustomFieldTemplateEntity, error)
	Save(ctx context.Context, id *string, input neo4jrepository.CustomFieldTemplateSaveFields) (string, error)
	Delete(ctx context.Context, customFieldTemplateId string) error
}

type customFieldTemplateService struct {
	log      logger.Logger
	services *Services
}

func NewCustomFieldTemplateService(log logger.Logger, services *Services) CustomFieldTemplateService {
	return &customFieldTemplateService{
		log:      log,
		services: services,
	}
}

func (s *customFieldTemplateService) GetAll(ctx context.Context) (*neo4jentity.CustomFieldTemplateEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	dbNodes, err := s.services.Neo4jRepositories.CustomFieldTemplateReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	customFieldTemplateEntities := neo4jentity.CustomFieldTemplateEntities{}
	for _, dbNode := range dbNodes {
		customFieldTemplateEntities = append(customFieldTemplateEntities, *neo4jmapper.MapDbNodeToCustomFieldTemplateEntity(dbNode))
	}
	return &customFieldTemplateEntities, nil
}

func (s *customFieldTemplateService) GetById(ctx context.Context, customFieldTemplateId string) (*neo4jentity.CustomFieldTemplateEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, customFieldTemplateId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	dbNode, err := s.services.Neo4jRepositories.CustomFieldTemplateReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), customFieldTemplateId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if dbNode == nil {
		err = errors.New("custom field template not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToCustomFieldTemplateEntity(dbNode), nil
}

func (s *customFieldTemplateService) Save(ctx context.Context, id *string, input neo4jrepository.CustomFieldTemplateSaveFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateService.Save")
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
	customFieldTemplateId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
		customFieldTemplateId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelCustomFieldTemplate)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		span.LogKV("flow", "update")
		customFieldTemplateId = *id

		// validate custom field template exists
		_, err = s.GetById(ctx, customFieldTemplateId)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, customFieldTemplateId)

	if createFlow {
		// validate entity type is present and is valid when creating new custom field template
		if !supportedEntityTypeForCustomFieldTemplate(input.EntityType) {
			err = errors.New("entity type is missing or not supported")
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	err = s.services.Neo4jRepositories.CustomFieldTemplateWriteRepository.Save(ctx, tenant, customFieldTemplateId, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		err = s.services.RabbitMQService.PublishEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.New_CreateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(input))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateCustomFieldTemplate"))
		}
	} else {
		err = s.services.RabbitMQService.PublishEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.New_UpdateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(input))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateCustomFieldTemplate"))
		}
	}

	return customFieldTemplateId, nil
}

func supportedEntityTypeForCustomFieldTemplate(entityType model.EntityType) bool {
	return entityType == model.ORGANIZATION || entityType == model.OPPORTUNITY || entityType == model.CONTACT || entityType == model.LOG_ENTRY
}

func (s *customFieldTemplateService) Delete(ctx context.Context, customFieldTemplateId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateService.Delete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, customFieldTemplateId)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// validate custom field template exists
	_, err = s.GetById(ctx, customFieldTemplateId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.Neo4jRepositories.CustomFieldTemplateWriteRepository.Delete(ctx, common.GetTenantFromContext(ctx), customFieldTemplateId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, customFieldTemplateId, model.CUSTOM_FIELD_TEMPLATE, dto.Delete{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message"))
	}

	return nil
}
