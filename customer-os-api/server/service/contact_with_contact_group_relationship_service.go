package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ContactWithContactGroupRelationshipService interface {
	AddContactToGroup(contactId, groupId string) (bool, error)
	RemoveContactFromGroup(contactId, groupId string) (bool, error)
}

type contactWithContactGroupRelationshipService struct {
	driver *neo4j.Driver
}

func NewContactWithContactGroupRelationshipService(driver *neo4j.Driver) ContactWithContactGroupRelationshipService {
	return &contactWithContactGroupRelationshipService{
		driver: driver,
	}
}

func (s *contactWithContactGroupRelationshipService) AddContactToGroup(contactId, groupId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId}), (g:ContactGroup {id:$groupId})
			MERGE (c)-[:BELONGS_TO]->(g)
			MERGE (g)-[:CONTAINS]->(c)
			`,
			map[string]interface{}{
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

func (s *contactWithContactGroupRelationshipService) RemoveContactFromGroup(contactId, groupId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId}), (g:ContactGroup {id:$groupId})
			MATCH (c)-[r1:BELONGS_TO]->(g)
			MATCH (g)-[r2:CONTAINS]->(c)
            DELETE r1, r2
			`,
			map[string]interface{}{
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
