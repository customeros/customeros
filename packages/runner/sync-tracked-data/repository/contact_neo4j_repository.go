package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"golang.org/x/net/context"
	"time"
)

type ContactRepository interface {
	GetOrCreateContactByEmail(ctx context.Context, tenant, email, firstName, lastName, application string) (string, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) GetOrCreateContactByEmail(ctx context.Context, tenant, email, firstName, lastName, application string) (string, error) {
	session := (*r.driver).NewSession(ctx,
		neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
			BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close(ctx)

	record, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			" MATCH (t:Tenant {name:$tenant}) "+
				" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) "+
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
		record, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record, nil
	})

	return record.(*db.Record).Values[0].(string), err
}
