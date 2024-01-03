package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type AnalysisRepository interface {
	LinkWithDescribesXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, analysisId string) error
	GetDescribesForAnalysis(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetDescribedByForXX(ctx context.Context, tenant string, ids []string, linkedWith LinkedWith) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAnalysis entity.AnalysisEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error)
}

type analysisRepository struct {
	driver *neo4j.DriverWithContext
}

func NewAnalysisRepository(driver *neo4j.DriverWithContext) AnalysisRepository {
	return &analysisRepository{
		driver: driver,
	}
}

func (r *analysisRepository) LinkWithDescribesXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, linkedWith LinkedWith, entityId, analysisId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalysisRepository.LinkWithDescribesXXInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (d:%s_%s {id:$entityId}) `, linkedWith, tenant)
	query += fmt.Sprintf(`MATCH (a:Analysis_%s {id:$analysisId}) `, tenant)
	query += `MERGE (a)-[r:DESCRIBES]->(d) `
	query += `return r `

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"entityId":   entityId,
			"analysisId": analysisId,
		})
	span.LogFields(log.String("query", query))
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *analysisRepository) GetDescribesForAnalysis(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalysisRepository.GetDescribesForAnalysis")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (a:Analysis_%s)-[DESCRIBES]->(d) " +
		" WHERE a.id IN $ids " +
		" RETURN d, a.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	span.LogFields(log.String("query", query))
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *analysisRepository) GetDescribedByForXX(ctx context.Context, tenant string, ids []string, linkedWith LinkedWith) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalysisRepository.GetDescribedByForXX")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (a:Analysis_%s)-[DESCRIBES]->(d:%s_%s) " +
		" WHERE d.id IN $ids " +
		" RETURN a, d.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, linkedWith, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	span.LogFields(log.String("query", query))
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *analysisRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newAnalysis entity.AnalysisEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AnalysisRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var createdAt time.Time
	createdAt = utils.Now()
	if newAnalysis.CreatedAt != nil {
		createdAt = *newAnalysis.CreatedAt
	}

	query := "MERGE (a:Analysis_%s {id:randomUUID()}) ON CREATE SET " +
		" a:Analysis, " +
		" a.source=$source, " +
		" a.createdAt=$createdAt, " +
		" a.analysisType=$analysisType, " +
		" a.content=$content, " +
		" a.contentType=$contentType, " +
		" a.sourceOfTruth=$sourceOfTruth, " +
		" a.appSource=$appSource " +
		" RETURN a"

	span.LogFields(log.String("query", query))

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"source":        source,
			"createdAt":     createdAt,
			"analysisType":  newAnalysis.AnalysisType,
			"content":       newAnalysis.Content,
			"contentType":   newAnalysis.ContentType,
			"sourceOfTruth": sourceOfTruth,
			"appSource":     newAnalysis.AppSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
