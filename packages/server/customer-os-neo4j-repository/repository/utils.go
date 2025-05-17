package neo4j_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

func LogAndExecuteWriteQuery(ctx context.Context, driver neo4j.DriverWithContext, cypher string, params map[string]any, spans *telemetry.Spans) error {
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func LogAndExecuteWriteQueryInTx(ctx context.Context, tx *neo4j.ManagedTransaction, driver *neo4j.DriverWithContext, database, cypher string, params map[string]any, spans *telemetry.Spans) error {
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, driver, database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
	}
	return err
}
