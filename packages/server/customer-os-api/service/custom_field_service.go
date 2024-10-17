package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

// TODO alexb deprecate and remove
type CustomFieldService interface {
	MergeAndUpdateCustomFieldsForContact(ctx context.Context, contactId string, customFields *entity.CustomFieldEntities) error
	MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error)
	DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error)
	GetCustomFields(ctx context.Context, obj *model.CustomFieldEntityType) (*entity.CustomFieldEntities, error)
}

type customFieldService struct {
	log        logger.Logger
	repository *repository.Repositories
}

func NewCustomFieldService(log logger.Logger, repository *repository.Repositories) CustomFieldService {
	return &customFieldService{
		log:        log,
		repository: repository,
	}
}

func (s *customFieldService) getDriver() neo4j.DriverWithContext {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *customFieldService) MergeAndUpdateCustomFieldsForContact(ctx context.Context, contactId string, customFields *entity.CustomFieldEntities) error {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		if customFields != nil {
			for _, customField := range *customFields {
				if customField.Id == nil {
					dbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(ctx, tx, tenant, contactId, customField)
					if err != nil {
						return nil, err
					}
					if customField.TemplateId != nil {
						var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
						err := s.repository.CustomFieldRepository.LinkWithCustomFieldTemplateForContactInTx(ctx, tx, fieldId, contactId, *customField.TemplateId)
						if err != nil {
							return nil, err
						}
					}
				} else {
					_, err := s.repository.CustomFieldRepository.UpdateForContactInTx(ctx, tx, tenant, contactId, customField)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		return nil, nil
	})

	return err
}

func (s *customFieldService) GetCustomFields(ctx context.Context, entityType *model.CustomFieldEntityType) (*entity.CustomFieldEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbRecords, err := s.repository.CustomFieldRepository.GetCustomFields(ctx, session, common.GetContext(ctx).Tenant, entityType)
	if err != nil {
		return nil, err
	}

	customFieldEntities := entity.CustomFieldEntities{}

	for _, dbRecord := range dbRecords {
		customFieldEntity := s.mapDbNodeToCustomFieldEntity(dbRecord.Values[0].(dbtype.Node))
		customFieldEntities = append(customFieldEntities, *customFieldEntity)
	}

	return &customFieldEntities, nil
}

func (s *customFieldService) MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	customFieldNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		customFieldDbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		if entity.TemplateId != nil {
			var fieldId = utils.GetPropsFromNode(*customFieldDbNode)["id"].(string)
			if err = s.repository.CustomFieldRepository.LinkWithCustomFieldTemplateForContactInTx(ctx, tx, fieldId, contactId, *entity.TemplateId); err != nil {
				return nil, err
			}
		}
		return customFieldDbNode, err
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToCustomFieldEntity(*customFieldNode.(*dbtype.Node)), nil
}

func (s *customFieldService) UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	customFieldDbNode, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.repository.CustomFieldRepository.UpdateForContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
	})

	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToCustomFieldEntity(*customFieldDbNode.(*dbtype.Node)), nil
}

func (s *customFieldService) DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)
	err := s.repository.CustomFieldRepository.DeleteByNameFromContact(ctx, session, common.GetContext(ctx).Tenant, contactId, fieldName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)
	err := s.repository.CustomFieldRepository.DeleteByIdFromContact(ctx, session, common.GetContext(ctx).Tenant, contactId, fieldId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) mapDbNodeToCustomFieldEntity(node dbtype.Node) *entity.CustomFieldEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.CustomFieldEntity{
		Id:            utils.StringPtr(utils.GetStringPropOrEmpty(props, "id")),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		DataType:      utils.GetStringPropOrEmpty(props, "datatype"),
		Source:        neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Value: model.AnyTypeValue{
			Str:   utils.GetStringPropOrNil(props, entity.CustomFieldTextProperty.String()),
			Time:  utils.GetTimePropOrNil(props, entity.CustomFieldTimeProperty.String()),
			Int:   utils.GetInt64PropOrNil(props, entity.CustomFieldIntProperty.String()),
			Float: utils.GetFloatPropOrNil(props, entity.CustomFieldFloatProperty.String()),
			Bool:  utils.GetBoolPropOrNil(props, entity.CustomFieldBoolProperty.String()),
		},
	}
	return &result
}
