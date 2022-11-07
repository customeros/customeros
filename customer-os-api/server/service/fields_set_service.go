package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type FieldsSetService interface {
	MergeFieldsSetToContact(ctx context.Context, contactId string, input *entity.FieldsSetEntity) (*entity.FieldsSetEntity, error)
}

type fieldsSetService struct {
	driver *neo4j.Driver
}

func NewFieldsSetService(driver *neo4j.Driver) FieldsSetService {
	return &fieldsSetService{
		driver: driver,
	}
}

func (s *fieldsSetService) MergeFieldsSetToContact(ctx context.Context, contactId string, input *entity.FieldsSetEntity) (*entity.FieldsSetEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:FieldsSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c)
            ON CREATE SET f.type=$type, f.id=randomUUID(), r.added=datetime({timezone: 'UTC'})
			RETURN f, r`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"name":      input.Name,
				"type":      input.Type,
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

	var entity = s.mapDbNodeToFieldsSetEntity((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node))
	s.addDbRelationshipToEntity((queryResult.([]*db.Record))[0].Values[1].(dbtype.Relationship), entity)
	return entity, nil
}

func (s *fieldsSetService) mapDbNodeToFieldsSetEntity(node dbtype.Node) *entity.FieldsSetEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.FieldsSetEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return &result
}

func (s *fieldsSetService) addDbRelationshipToEntity(relationship dbtype.Relationship, fieldsSetEntity *entity.FieldsSetEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	fieldsSetEntity.Added = utils.GetTimePropOrNow(props, "added")
}
