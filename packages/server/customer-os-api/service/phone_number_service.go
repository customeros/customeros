package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type PhoneNumberService interface {
	GetAllForContact(ctx context.Context, contactId string) (*entity.PhoneNumberEntities, error)
	MergePhoneNumberToContact(ctx context.Context, id string, toEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error)
	UpdatePhoneNumberInContact(ctx context.Context, id string, toEntity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error)
	RemoveFromContactByE164(ctx context.Context, contactId, e164 string) (bool, error)
	RemoveFromContactById(ctx context.Context, contactId, phoneNumberId string) (bool, error)
}

type phoneNumberService struct {
	repositories *repository.Repositories
}

func NewPhoneNumberService(repositories *repository.Repositories) PhoneNumberService {
	return &phoneNumberService{
		repositories: repositories,
	}
}

func (s *phoneNumberService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *phoneNumberService) GetAllForContact(ctx context.Context, contactId string) (*entity.PhoneNumberEntities, error) {
	queryResult, err := s.repositories.PhoneNumberRepository.GetAllForContact(ctx, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	phoneNumberEntities := entity.PhoneNumberEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		phoneNumberEntity := s.mapDbNodeToPhoneNumberEntity(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToPhoneNumberEntity(dbRecord.Values[1].(dbtype.Relationship), phoneNumberEntity)
		phoneNumberEntities = append(phoneNumberEntities, *phoneNumberEntity)
	}

	return &phoneNumberEntities, nil
}

func (s *phoneNumberService) MergePhoneNumberToContact(ctx context.Context, contactId string, entity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	var err error
	var phoneNumberNode *dbtype.Node
	var phoneNumberRelationship *dbtype.Relationship

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		phoneNumberNode, phoneNumberRelationship, err = s.repositories.PhoneNumberRepository.MergePhoneNumberToContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		phoneNumberId := utils.GetPropsFromNode(*phoneNumberNode)["id"].(string)
		if entity.Primary == true {
			err := s.repositories.PhoneNumberRepository.SetOtherContactPhoneNumbersNonPrimaryInTx(ctx, tx, common.GetTenantFromContext(ctx), contactId, phoneNumberId)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	var phoneNumberEntity = s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode)
	s.addDbRelationshipToPhoneNumberEntity(*phoneNumberRelationship, phoneNumberEntity)
	return phoneNumberEntity, nil
}

func (s *phoneNumberService) UpdatePhoneNumberInContact(ctx context.Context, contactId string, entity *entity.PhoneNumberEntity) (*entity.PhoneNumberEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	var err error
	var phoneNumberNode *dbtype.Node
	var phoneNumberRelationship *dbtype.Relationship

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		phoneNumberNode, phoneNumberRelationship, err = s.repositories.PhoneNumberRepository.UpdatePhoneNumberByContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		phoneNumberId := utils.GetPropsFromNode(*phoneNumberNode)["id"].(string)
		if entity.Primary == true {
			err := s.repositories.PhoneNumberRepository.SetOtherContactPhoneNumbersNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, phoneNumberId)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var phoneNumberEntity = s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode)
	s.addDbRelationshipToPhoneNumberEntity(*phoneNumberRelationship, phoneNumberEntity)
	return phoneNumberEntity, nil
}

func (s *phoneNumberService) RemoveFromContactByE164(ctx context.Context, contactId, e164 string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[rel:HAS]->(p:PhoneNumber {e164:$e164})
            DELETE rel`,
			map[string]interface{}{
				"id":     contactId,
				"e164":   e164,
				"tenant": common.GetTenantFromContext(ctx),
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *phoneNumberService) RemoveFromContactById(ctx context.Context, contactId, phoneNumberId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[rel:HAS]->(p:PhoneNumber {id:$phoneNumberId})
            DELETE rel`,
			map[string]interface{}{
				"contactId":     contactId,
				"phoneNumberId": phoneNumberId,
				"tenant":        common.GetTenantFromContext(ctx),
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
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
