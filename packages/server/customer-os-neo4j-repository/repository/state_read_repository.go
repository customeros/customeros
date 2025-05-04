package neo4j_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type StateReadRepository interface {
	GetStatesByCountryId(ctx context.Context, countryId string) ([]*dbtype.Node, error)
}

type stateReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewStateReadRepository(driver *neo4j.DriverWithContext, database string) StateReadRepository {
	return &stateReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *stateReadRepository) GetStatesByCountryId(ctx context.Context, countryId string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "StateReadRepository.GetStatesByCountryId")
	defer spans.Finish()
	spans.LogKV("countryId", countryId)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"countryId": countryId,
		}

		queryResult, err := tx.Run(ctx, "MATCH (s:Period)-[:BELONGS_TO_COUNTRY]->(c:Country { id: $countryId }) RETURN s", params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range dbRecords.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}
