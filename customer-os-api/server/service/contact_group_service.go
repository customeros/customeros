package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactGroupService interface {
	Create(ctx context.Context, contactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error)
	Delete(ctx context.Context, id string) (bool, error)
	FindAll(ctx context.Context) (*entity.ContactGroupEntities, error)
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupEntities, error)
}

type contactGroupService struct {
	driver *neo4j.Driver
}

func NewContactGroupService(driver *neo4j.Driver) ContactGroupService {
	return &contactGroupService{
		driver: driver,
	}
}

func (s *contactGroupService) Create(ctx context.Context, newContactGroup *entity.ContactGroupEntity) (*entity.ContactGroupEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (g:ContactGroup {
				  id: randomUUID(),
				  name: $name})-[:GROUP_BELONGS_TO_TENANT]->(t)
			RETURN g`,
			map[string]interface{}{
				"name":   newContactGroup.Name,
				"tenant": common.GetContext(ctx).Tenant,
			})

		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}
	return mapDbNodeToContactGroup(queryResult.(dbtype.Node)), nil
}

func (s *contactGroupService) Delete(ctx context.Context, id string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
            DETACH DELETE g
			`,
			map[string]interface{}{
				"groupId": id,
				"tenant":  common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactGroupService) FindAll(ctx context.Context) (*entity.ContactGroupEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (:Tenant {name:$tenant})<-[:GROUP_BELONGS_TO_TENANT]-(g:ContactGroup) 
				RETURN g
				ORDER BY g.name`,
			map[string]interface{}{
				"tenant": common.GetContext(ctx).Tenant,
			})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	contactGroups := entity.ContactGroupEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := mapDbNodeToContactGroup(dbRecord.Values[0].(dbtype.Node))
		contactGroups = append(contactGroups, *contactGroup)
	}

	return &contactGroups, nil
}

func (s *contactGroupService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$id})-[:BELONGS_TO_GROUP]->(g:ContactGroup)-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
				RETURN g 
				ORDER BY g.name`,
			map[string]interface{}{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	contactGroups := entity.ContactGroupEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := mapDbNodeToContactGroup(dbRecord.Values[0].(dbtype.Node))
		contactGroups = append(contactGroups, *contactGroup)
	}

	return &contactGroups, nil
}

func mapDbNodeToContactGroup(dbContactGroupNode dbtype.Node) *entity.ContactGroupEntity {
	props := utils.GetPropsFromNode(dbContactGroupNode)
	contactGroup := entity.ContactGroupEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &contactGroup
}
