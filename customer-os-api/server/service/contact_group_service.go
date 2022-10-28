package service

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ContactGroupService interface {
	Create(ctx context.Context, contactGroup *entity.ContactGroupNode) (*entity.ContactGroupNode, error)
	FindAll(ctx context.Context) (*entity.ContactGroupNodes, error)
	Delete(ctx context.Context, id string) (bool, error)
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupNodes, error)
}

type contactGroupService struct {
	driver *neo4j.Driver
}

func NewContactGroupService(driver *neo4j.Driver) ContactGroupService {
	return &contactGroupService{
		driver: driver,
	}
}

func (s *contactGroupService) Create(ctx context.Context, newContactGroup *entity.ContactGroupNode) (*entity.ContactGroupNode, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (t:Tenant {name:$tenant})
			CREATE (g:ContactGroup {
				  id: randomUUID(),
				  name: $name})-[:BELONGS_TO]->(t)
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
	contactGroup := entity.ContactGroupNode{}
	mapstructure.Decode(utils.GetPropsFromNode(queryResult.(dbtype.Node)), &contactGroup)

	return &contactGroup, nil
}

func (s *contactGroupService) FindAll(ctx context.Context) (*entity.ContactGroupNodes, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (g:ContactGroup)--(:Tenant {name:$tenant}) RETURN g`, map[string]interface{}{
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

	contactGroups := entity.ContactGroupNodes{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := entity.ContactGroupNode{}
		mapstructure.Decode(utils.GetPropsFromNode(dbRecord.Values[0].(dbtype.Node)), &contactGroup)
		contactGroups = append(contactGroups, contactGroup)
	}

	return &contactGroups, nil
}

func (s *contactGroupService) Delete(ctx context.Context, id string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact), (g:ContactGroup {id:$groupId})-[r0:BELONGS_TO]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (c)-[r1:BELONGS_TO]->(g)
			OPTIONAL MATCH (g)-[r2:CONTAINS]->(c)
            DELETE r0, r1, r2, g
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

func (s *contactGroupService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.ContactGroupNodes, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (:Tenant {name:$tenant})<--(g:ContactGroup)<--(c:Contact {id:$id}) RETURN g`,
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

	contactGroups := entity.ContactGroupNodes{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		contactGroup := entity.ContactGroupNode{}
		mapstructure.Decode(utils.GetPropsFromNode(dbRecord.Values[0].(dbtype.Node)), &contactGroup)
		contactGroups = append(contactGroups, contactGroup)
	}

	return &contactGroups, nil
}
