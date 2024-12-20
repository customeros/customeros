package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

type PhoneNumberService interface {
	Merge(ctx context.Context, phoneNumber string, source neo4jentity.DataSource) (string, error)
	UpdatePhoneNumberFor(ctx context.Context, entityType commonModel.EntityType, entityId string, phoneId string, label *string, primary *bool) error
	DetachFromEntityByPhoneNumber(ctx context.Context, entityType commonModel.EntityType, entityId, phoneNumber string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType commonModel.EntityType, entityId, phoneNumberId string) (bool, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, ids []string) (*neo4jentity.PhoneNumberEntities, error)
	GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.PhoneNumberEntity, error)
}

type phoneNumberService struct {
	services *Services
}

func NewPhoneNumberService(services *Services) PhoneNumberService {
	return &phoneNumberService{
		services: services,
	}
}

func (s *phoneNumberService) Merge(ctx context.Context, phoneNumber string, source neo4jentity.DataSource) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.CreatePhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	tenant := common.GetTenantFromContext(ctx)

	var phoneNumberId string
	phoneNumber = strings.TrimSpace(phoneNumber)

	existingPhoneNumber, err := s.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if existingPhoneNumber == nil {
		phoneNumberId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelPhoneNumber)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		data := neo4jrepository.PhoneNumberCreateFields{
			RawPhoneNumber: phoneNumber,
			SourceFields: neo4jmodel.SourceFields{
				AppSource: common.GetAppSourceFromContext(ctx),
				Source:    source.String(),
			},
			CreatedAt: utils.Now(),
		}

		err = s.services.Neo4jRepositories.PhoneNumberWriteRepository.CreatePhoneNumber(ctx, tenant, phoneNumberId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	} else {
		phoneNumberId = existingPhoneNumber.Id
	}

	return phoneNumberId, nil
}

func (s *phoneNumberService) GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, ids []string) (*neo4jentity.PhoneNumberEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetAllForEntityTypeByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("ids", ids))

	phoneNumbers, err := s.services.Neo4jRepositories.PhoneNumberReadRepository.GetAllForLinkedEntityIds(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	phoneNumberEntities := make(neo4jentity.PhoneNumberEntities, 0, len(phoneNumbers))
	for _, v := range phoneNumbers {
		phoneNumberEntity := neo4jmapper.MapDbNodeToPhoneNumberEntity(v.Node)
		s.addDbRelationshipToPhoneNumberEntity(*v.Relationship, phoneNumberEntity)
		phoneNumberEntity.DataloaderKey = v.LinkedNodeId
		phoneNumberEntities = append(phoneNumberEntities, *phoneNumberEntity)
	}
	return &phoneNumberEntities, nil
}

func (s *phoneNumberService) UpdatePhoneNumberFor(ctx context.Context, entityType commonModel.EntityType, entityId string, phoneId string, label *string, primary *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.UpdatePhoneNumberFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", string(entityType)), log.String("entityId", entityId))

	tenant := common.GetTenantFromContext(ctx)

	phoneNumberEntity, err := s.GetById(ctx, phoneId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	existsById, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, entityId, entityType.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if !existsById {
		err = errors.New("Entity not found")
		tracing.TraceErr(span, err)
		return err
	}

	if entityType == commonModel.ORGANIZATION {
		err = s.services.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithOrganization(ctx, tenant, entityId, phoneNumberEntity.Id, utils.IfNotNilString(label), utils.IfNotNilBool(primary))
	} else if entityType == commonModel.CONTACT {
		err = s.services.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithContact(ctx, tenant, entityId, phoneNumberEntity.Id, utils.IfNotNilString(label), utils.IfNotNilBool(primary))
	} else if entityType == commonModel.USER {
		err = s.services.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithUser(ctx, tenant, entityId, phoneNumberEntity.Id, utils.IfNotNilString(label), utils.IfNotNilBool(primary))
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, entityId, entityType, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *phoneNumberService) DetachFromEntityByPhoneNumber(ctx context.Context, entityType commonModel.EntityType, entityId, phoneNumber string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	err := s.services.Neo4jRepositories.PhoneNumberWriteRepository.RemoveRelationship(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumber)

	//TODO: Update last touchpoint
	//if entityType == commonModel.ORGANIZATION {
	//	s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	//} else if entityType == commonModel.CONTACT {
	//	s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	//}

	return err == nil, err
}

func (s *phoneNumberService) DetachFromEntityById(ctx context.Context, entityType commonModel.EntityType, entityId, phoneNumberId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId), log.String("phoneNumberId", phoneNumberId))

	err := s.services.Neo4jRepositories.PhoneNumberWriteRepository.RemoveRelationshipById(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumberId)

	//TODO: Update last touchpoint
	//if entityType == commonModel.ORGANIZATION {
	//	s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	//} else if entityType == commonModel.CONTACT {
	//	s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	//}

	return err == nil, err
}

func (s *phoneNumberService) GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	phoneNumberNode, err := s.services.Neo4jRepositories.PhoneNumberReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), phoneNumberId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberNode), nil
}

func (s *phoneNumberService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	phoneNumberNode, err := s.services.Neo4jRepositories.PhoneNumberReadRepository.GetByPhoneNumber(ctx, common.GetTenantFromContext(ctx), phoneNumber)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberNode), nil
}

func (s *phoneNumberService) addDbRelationshipToPhoneNumberEntity(relationship dbtype.Relationship, phoneNumberEntity *neo4jentity.PhoneNumberEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	phoneNumberEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
	phoneNumberEntity.Label = utils.GetStringPropOrEmpty(props, "label")
}
