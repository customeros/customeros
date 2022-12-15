package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type ContactRepository interface {
	MergeContact(tenant string, syncDate time.Time, contact entity.ContactEntity) (string, error)
}

type contactRepository struct {
	driver *neo4j.Driver
}

func NewContactRepository(driver *neo4j.Driver) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) MergeContact(tenant string, syncDate time.Time, contact entity.ContactEntity) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				MERGE (c)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e)
				ON CREATE SET r.externalId=$externalId, c.id=randomUUID(), 
								c.firstName=$firstName, c.lastName=$lastName, r.syncDate=$syncDate, c.readonly=$readonly,
								c.createdAt=$createdAt
				ON MATCH SET 	c.firstName=$firstName, c.lastName=$lastName, r.syncDate=$syncDate, c.readonly=$readonly
				WITH c, t
				MERGE (c)-[:CONTACT_BELONGS_TO_TENANT]->(t)
				RETURN c.id`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": contact.ExternalSystem,
				"externalId":     contact.ExternalId,
				"firstName":      contact.FirstName,
				"lastName":       contact.LastName,
				"syncDate":       syncDate,
				"readonly":       contact.Readonly,
				"createdAt":      contact.CreatedAt,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(string), nil
}
