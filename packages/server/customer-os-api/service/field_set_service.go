package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type FieldSetService interface {
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.FieldSetEntities, error)
	MergeFieldSetToContact(ctx context.Context, contactId string, input *entity.FieldSetEntity, fieldSetDefinitionId *string) (*entity.FieldSetEntity, error)
	UpdateFieldSetInContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error)
	DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error)
	getDriver() neo4j.Driver
}

type fieldSetService struct {
	repository *repository.RepositoryContainer
}

func NewFieldSetService(repository *repository.RepositoryContainer) FieldSetService {
	return &fieldSetService{
		repository: repository,
	}
}

func (s *fieldSetService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *fieldSetService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.FieldSetEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet) 
				RETURN s, r`,
			map[string]any{
				"contactId": contact.ID,
				"tenant":    common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	fieldSetEntities := entity.FieldSetEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		fieldSetEntity := s.mapDbNodeToFieldSetEntity(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToEntity(dbRecord.Values[1].(dbtype.Relationship), fieldSetEntity)
		fieldSetEntities = append(fieldSetEntities, *fieldSetEntity)
	}

	return &fieldSetEntities, nil
}

func (s *fieldSetService) MergeFieldSetToContact(ctx context.Context, contactId string, input *entity.FieldSetEntity, fieldSetDefinitionId *string) (*entity.FieldSetEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c)
            ON CREATE SET f.id=randomUUID(), r.added=datetime({timezone: 'UTC'})
			RETURN f, r`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"name":      input.Name,
			})
		records, err := txResult.Collect()
		if err != nil {
			return nil, err
		}
		if fieldSetDefinitionId != nil {
			var fieldSetId = utils.GetPropsFromNode(records[0].Values[0].(dbtype.Node))["id"].(string)
			err := s.repository.FieldSetRepository.LinkWithFieldSetDefinitionInTx(common.GetContext(ctx).Tenant, fieldSetId, *fieldSetDefinitionId, tx)
			if err != nil {
				return nil, err
			}
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node))
	s.addDbRelationshipToEntity((queryResult.([]*db.Record))[0].Values[1].(dbtype.Relationship), fieldSetEntity)
	return fieldSetEntity, nil
}

func (s *fieldSetService) UpdateFieldSetInContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error) {
	session := (s.getDriver()).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name
			RETURN s, r`,
			map[string]interface{}{
				"tenant":     common.GetContext(ctx).Tenant,
				"contactId":  contactId,
				"fieldSetId": input.Id,
				"name":       input.Name,
			})
		records, err := txResult.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node))
	s.addDbRelationshipToEntity((queryResult.([]*db.Record))[0].Values[1].(dbtype.Relationship), fieldSetEntity)
	return fieldSetEntity, nil
}

func (s *fieldSetService) DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
				  (s)-[:HAS_PROPERTY]->(f:CustomField)
            DETACH DELETE f, s`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"tenant":     common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *fieldSetService) mapDbNodeToFieldSetEntity(node dbtype.Node) *entity.FieldSetEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.FieldSetEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &result
}

func (s *fieldSetService) addDbRelationshipToEntity(relationship dbtype.Relationship, fieldSetEntity *entity.FieldSetEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	fieldSetEntity.Added = utils.GetTimePropOrNow(props, "added")
}
