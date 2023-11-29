package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"log"
)

func ExecuteWriteQuery(ctx context.Context, driver *neo4j.DriverWithContext, cypher string, params map[string]any) {
	_, err := utils.ExecuteQuery(ctx, *driver, "neo4j", cypher, params)
	if err != nil {
		log.Fatalf("Error executing query %s", cypher)
		return
	}
}

func ExecuteReadQueryWithSingleReturn(ctx context.Context, driver *neo4j.DriverWithContext, query string, params map[string]any) any {
	session := utils.NewNeo4jReadSession(ctx, *driver, utils.WithDatabaseName("neo4j"))
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, err := tx.Run(ctx, query, params)
		if err != nil {
			log.Fatalf("Error executing query %s", query)
		}
		return record.Single(ctx)
	})
	if err != nil {
		log.Fatalf("Error executing query %s", query)
	}
	return queryResult
}

func ExecuteReadQueryWithCollectionReturn(ctx context.Context, driver *neo4j.DriverWithContext, query string, params map[string]any) []*db.Record {
	session := utils.NewNeo4jReadSession(ctx, *driver, utils.WithDatabaseName("neo4j"))
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx, query, params)
		if err != nil {
			log.Fatalf("Error executing query %s", query)
		}
		return records.Collect(ctx)
	})
	if err != nil {
		log.Fatalf("Error executing query %s", query)
	}
	return queryResult.([]*db.Record)
}
