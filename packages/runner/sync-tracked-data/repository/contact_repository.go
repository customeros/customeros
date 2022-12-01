package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

type ContactRepository interface {
	GetOrCreateContactId(tenant, email, firstName, lastName string) (string, error)
}

type contactRepository struct {
	driver *neo4j.Driver
}

func NewContactRepository(driver *neo4j.Driver) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetOrCreateContactId(tenant, email, firstName, lastName string) (string, error) {
	session := (*r.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	record, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(`
			MERGE (t:Tenant {name:$tenant})
			MERGE (e:Email {email: $email})<-[r:EMAILED_AT]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t)
            ON CREATE SET r.primary=true, e.id=randomUUID(), c.id=randomUUID(), c.firstName=$firstName, c.lastName=$lastName, c.createdAt=datetime({timezone: 'UTC'})
			RETURN c.id`,
			map[string]interface{}{
				"tenant":    tenant,
				"email":     email,
				"firstName": firstName,
				"lastName":  lastName,
			})
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})

	return record.(*db.Record).Values[0].(string), err
}
