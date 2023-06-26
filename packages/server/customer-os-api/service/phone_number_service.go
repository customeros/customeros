package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PhoneNumberService interface {
	MergePhoneNumberTo(ctx context.Context, entityType entity.EntityType, entityId string, inputEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error)
	UpdatePhoneNumberFor(ctx context.Context, entityType entity.EntityType, entityId string, inputEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error)
	DetachFromEntityByPhoneNumber(ctx context.Context, entityType entity.EntityType, entityId, phoneNumber string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, phoneNumberId string) (bool, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, ids []string) (*entity.PhoneNumberEntities, error)

	mapDbNodeToPhoneNumberEntity(node dbtype.Node) *entity.PhoneNumberEntity
}

type phoneNumberService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) PhoneNumberService {
	return &phoneNumberService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *phoneNumberService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *phoneNumberService) GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, ids []string) (*entity.PhoneNumberEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetAllForEntityTypeByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("ids", ids))

	phoneNumbers, err := s.repositories.PhoneNumberRepository.GetAllForIds(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	phoneNumberEntities := make(entity.PhoneNumberEntities, 0, len(phoneNumbers))
	for _, v := range phoneNumbers {
		emailEntity := s.mapDbNodeToPhoneNumberEntity(*v.Node)
		s.addDbRelationshipToPhoneNumberEntity(*v.Relationship, emailEntity)
		emailEntity.DataloaderKey = v.LinkedNodeId
		phoneNumberEntities = append(phoneNumberEntities, *emailEntity)
	}
	return &phoneNumberEntities, nil
}

func (s *phoneNumberService) MergePhoneNumberTo(ctx context.Context, entityType entity.EntityType, entityId string, inputEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.MergePhoneNumberTo")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	var err error
	var phoneNumberNode *dbtype.Node
	var phoneNumberRelationship *dbtype.Relationship

	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		phoneNumberNode, phoneNumberRelationship, err = s.repositories.PhoneNumberRepository.MergePhoneNumberToInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, *inputEntity)
		if err != nil {
			return nil, err
		}
		phoneNumberId := utils.GetPropsFromNode(*phoneNumberNode)["id"].(string)
		if inputEntity.Primary == true {
			err := s.repositories.PhoneNumberRepository.SetOtherPhoneNumbersNonPrimaryInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, phoneNumberId)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, entityId)
	}

	var phoneNumberEntity = s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode)
	s.addDbRelationshipToPhoneNumberEntity(*phoneNumberRelationship, phoneNumberEntity)
	return phoneNumberEntity, nil
}

func (s *phoneNumberService) UpdatePhoneNumberFor(ctx context.Context, entityType entity.EntityType, entityId string, inputEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.UpdatePhoneNumberFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	var err error
	var phoneNumberNode *dbtype.Node
	var phoneNumberRelationship *dbtype.Relationship
	var detachCurrentPhoneNumber = false
	var phoneNumberExists = false

	currentPhoneNumberNode, err := s.repositories.PhoneNumberRepository.GetByIdAndRelatedEntity(ctx, entityType, common.GetTenantFromContext(ctx), inputEntity.Id, entityId)
	if err != nil {
		return nil, err
	}
	currentE164 := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*currentPhoneNumberNode), "e164")
	currentRawPhoneNumber := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*currentPhoneNumberNode), "rawPhoneNumber")
	if currentRawPhoneNumber == "" {
		phoneNumberExists, err = s.repositories.PhoneNumberRepository.Exists(ctx, common.GetContext(ctx).Tenant, inputEntity.RawPhoneNumber)
		if err != nil {
			return nil, err
		}
	}

	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		if len(inputEntity.RawPhoneNumber) == 0 || inputEntity.RawPhoneNumber == currentE164 || inputEntity.RawPhoneNumber == currentRawPhoneNumber ||
			(currentRawPhoneNumber == "" && !phoneNumberExists) {
			phoneNumberNode, phoneNumberRelationship, err = s.repositories.PhoneNumberRepository.UpdatePhoneNumberForInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, *inputEntity)
			if err != nil {
				return nil, err
			}
			phoneNumberId := utils.GetPropsFromNode(*phoneNumberNode)["id"].(string)
			if inputEntity.Primary == true {
				err := s.repositories.PhoneNumberRepository.SetOtherPhoneNumbersNonPrimaryInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, phoneNumberId)
				if err != nil {
					return nil, err
				}
			}
		} else {
			phoneNumberNode, phoneNumberRelationship, err = s.repositories.PhoneNumberRepository.MergePhoneNumberToInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, *inputEntity)
			if err != nil {
				return nil, err
			}
			phoneNumberId := utils.GetPropsFromNode(*phoneNumberNode)["id"].(string)
			if inputEntity.Primary == true {
				err := s.repositories.PhoneNumberRepository.SetOtherPhoneNumbersNonPrimaryInTx(ctx, tx, common.GetTenantFromContext(ctx), entityType, entityId, phoneNumberId)
				if err != nil {
					return nil, err
				}
			}
			detachCurrentPhoneNumber = true
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	if detachCurrentPhoneNumber {
		_, err = s.DetachFromEntityById(ctx, entityType, entityId, inputEntity.Id)
	}

	var phoneNumberEntity = s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode)
	s.addDbRelationshipToPhoneNumberEntity(*phoneNumberRelationship, phoneNumberEntity)
	return phoneNumberEntity, nil
}

func (s *phoneNumberService) DetachFromEntityByPhoneNumber(ctx context.Context, entityType entity.EntityType, entityId, phoneNumber string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	err := s.repositories.PhoneNumberRepository.RemoveRelationship(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumber)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *phoneNumberService) DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, phoneNumberId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId), log.String("phoneNumberId", phoneNumberId))

	err := s.repositories.PhoneNumberRepository.RemoveRelationshipById(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumberId)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpointSync(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *phoneNumberService) mapDbNodeToPhoneNumberEntity(node dbtype.Node) *entity.PhoneNumberEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.PhoneNumberEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		E164:           utils.GetStringPropOrEmpty(props, "e164"),
		RawPhoneNumber: utils.GetStringPropOrEmpty(props, "rawPhoneNumber"),
		Validated:      utils.GetBoolPropOrFalse(props, "validated"),
		Source:         entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}

func (s *phoneNumberService) addDbRelationshipToPhoneNumberEntity(relationship dbtype.Relationship, phoneNumberEntity *entity.PhoneNumberEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	phoneNumberEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
	phoneNumberEntity.Label = utils.GetStringPropOrEmpty(props, "label")
}
