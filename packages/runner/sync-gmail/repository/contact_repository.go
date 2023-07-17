package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, tenant, firstName, lastName, source, appSource string, date time.Time) (*dbtype.Node, error)
}

type contactRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContactRepository(driver *neo4j.DriverWithContext) ContactRepository {
	return &contactRepository{
		driver: driver,
	}
}

func (r *contactRepository) CreateContact(ctx context.Context, tenant, firstName, lastName, source, appSource string, date time.Time) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (p:Contact:Contact_%s {id:randomUUID()}) 
		 		SET 	p.firstName = $firstName,
						p.lastName = $lastName,	
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.syncedWithEventStore = false
				MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(p) return c`, tenant)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"firstName":     firstName,
				"lastName":      lastName,
				"tenant":        tenant,
				"source":        source,
				"sourceOfTruth": source,
				"appSource":     appSource,
				"createdAt":     date,
				"updatedAt":     date,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
