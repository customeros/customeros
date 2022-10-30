package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
)

type ContactWithContactGroupRelationshipService interface {
	AddContactToGroup(ctx context.Context, contactId, groupId string) (bool, error)
	RemoveContactFromGroup(ctx context.Context, contactId, groupId string) (bool, error)
}

type contactWithContactGroupRelationshipService struct {
	driver *neo4j.Driver
}

func NewContactWithContactGroupRelationshipService(driver *neo4j.Driver) ContactWithContactGroupRelationshipService {
	return &contactWithContactGroupRelationshipService{
		driver: driver,
	}
}

func (s *contactWithContactGroupRelationshipService) AddContactToGroup(ctx context.Context, contactId, groupId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH 	(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
					(g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (c)-[:BELONGS_TO_GROUP]->(g)
			MERGE (g)-[:CONTAINS_CONTACT]->(c)
			`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"groupId":   groupId,
			})
		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *contactWithContactGroupRelationshipService) RemoveContactFromGroup(ctx context.Context, contactId, groupId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH 	(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
					(g:ContactGroup {id:$groupId})-[:GROUP_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MATCH (c)-[r1:BELONGS_TO_GROUP]->(g)
			MATCH (g)-[r2:CONTAINS_CONTACT]->(c)
            DELETE r1, r2
			`,
			map[string]interface{}{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"groupId":   groupId,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}
