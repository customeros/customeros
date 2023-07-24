package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type AnalysisRepository interface {
	SummaryExistsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (bool, error)
	CreateSummaryForEmail(ctx context.Context, tenant, interactionEventId, summary, source, appSource string, createdAt time.Time) (*dbtype.Node, error)
}

type analysisRepository struct {
	driver *neo4j.DriverWithContext
}

func NewAnalysisRepository(driver *neo4j.DriverWithContext) AnalysisRepository {
	return &analysisRepository{
		driver: driver,
	}
}

func (r *analysisRepository) SummaryExistsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (bool, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})-[:DESCRIBES]->(a:Analysis_%s) RETURN a`, tenant, tenant)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return result != nil, nil
}

func (r *analysisRepository) CreateSummaryForEmail(ctx context.Context, tenant, interactionEventId, summary, source, appSource string, createdAt time.Time) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId}) `, tenant)
	query += fmt.Sprintf(`MERGE (i)<-[r:DESCRIBES]-(a:Analysis_%s{id:randomUUID()}) `, tenant)
	query += fmt.Sprintf("ON CREATE SET " +
		" a:Analysis, " +
		" a.createdAt=$createdAt, " +
		" a.analysisType=$analysisType, " +
		" a.content=$content, " +
		" a.contentType=$contentType, " +
		" a.source=$source, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a")

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
				"createdAt":          createdAt,
				"analysisType":       "summary",
				"content":            summary,
				"contentType":        "text",
				"source":             source,
				"sourceOfTruth":      source,
				"appSource":          appSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
