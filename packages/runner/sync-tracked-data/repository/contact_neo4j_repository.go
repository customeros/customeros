package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"time"
)

type ContactRepository interface {
	GetOrCreateContactId(tenant, email, firstName, lastName, application string) (string, error)
}

type contactRepository struct {
	driver *neo4j.Driver
}

func NewContactRepository(driver *neo4j.Driver) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetOrCreateContactId(tenant, email, firstName, lastName, application string) (string, error) {
	session := (*r.driver).NewSession(
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	record, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(fmt.Sprintf(
			" MATCH (t:Tenant {name:$tenant}) "+
				" MERGE (e:Email {email: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET "+
				"				e.id=randomUUID(), "+
				"				e.createdAt=$now, "+
				"				e.updatedAt=$now, "+
				"				e.source=$source, "+
				"				e.sourceOfTruth=$sourceOfTruth, "+
				"				e.appSource=$appSource, "+
				"				e:%s "+
				" WITH DISTINCT t, e "+
				" MERGE (e)<-[rel:HAS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) "+
				" ON CREATE SET rel.primary=true, "+
				"				c.id=randomUUID(), "+
				"				c.firstName=$firstName, "+
				"				c.lastName=$lastName, "+
				"				c.createdAt=$now, "+
				"				c.updatedAt=$now, "+
				"				c.source=$source, "+
				"				c.sourceOfTruth=$sourceOfTruth, "+
				"				c.appSource=$appSource, "+
				"               c:%s "+
				" RETURN c.id", "Email_"+tenant, "Contact_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"email":         email,
				"firstName":     firstName,
				"lastName":      lastName,
				"source":        "openline",
				"sourceOfTruth": "openline",
				"appSource":     application,
				"now":           time.Now().UTC(),
			})
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})

	return record.(*db.Record).Values[0].(string), err
}
